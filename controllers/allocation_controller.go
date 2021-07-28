/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	myappv1 "github.com/aaronlyc/allocation/api/v1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AllocationReconciler reconciles a Allocation object
type AllocationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.aaron.domain,resources=allocations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.aaron.domain,resources=allocations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.aaron.domain,resources=allocations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Allocation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AllocationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := r.Log.WithName("Reconcile")
	logs.V(8).Info("Begin Reconcile >>>> ")
	// your logic here

	allocation, err := r.GetAllocate(ctx, req.NamespacedName)
	if err != nil {
		logs.Error(err, "获取CR对象失败 <<<<")
		return ctrl.Result{}, err
	}
	logs.Info("Total数量", "total:", allocation.Status.Total)

	err = r.dealAllocate(ctx, allocation)
	if err != nil {
		logs.Error(err, "创建deployment失败 <<<<")
		return ctrl.Result{}, err
	}

	// 添加指标
	allocationMetrics.WithLabelValues(allocation.Spec.MsName, "total").Set(float64(allocation.Status.Total))
	allocationMetrics.WithLabelValues(allocation.Spec.MsName, "leave").Set(float64(allocation.Spec.MaxNum - allocation.Status.Total))

	// 更新
	err = r.Status().Update(ctx, allocation)
	if err != nil {
		logs.Error(err, "更新CR状态失败 <<<<")
	}

	logs.V(8).Info("Stop Reconcile <<<< ")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AllocationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.Allocation{}).
		Complete(r)
}

func (r *AllocationReconciler) GetAllocate(ctx context.Context, namespacedName types.NamespacedName) (*myappv1.Allocation, error) {
	logs := r.Log.WithName("getAllocate")
	allocation := &myappv1.Allocation{}
	err := r.Get(ctx, namespacedName, allocation)
	if err != nil {
		logs.Error(err, "获取CR对象失败 <<<<")
		return nil, err
	}
	return allocation, nil
}

func (r *AllocationReconciler) dealAllocate(ctx context.Context, allocation *myappv1.Allocation) error {
	logs := r.Log.WithName("getAllocate")
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", allocation.Spec.MsName, NewRandString()),
			Namespace: allocation.GetNamespace(),
			Labels: map[string]string{
				"app": allocation.Spec.MsName}},
		Spec: appsv1.DeploymentSpec{
			Replicas: &allocation.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"foo": "bar"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"foo": "bar"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
			},
		},
	}

	if allocation.Status.Total < allocation.Spec.MaxNum {
		//time.Sleep(time.Duration(allocation.Spec.Interval) * time.Second)
		err := r.Create(ctx, dep)
		if err != nil {
			logs.Error(err, "创建deployment失败 <<<<")
			return err
		}
		logs.Info("创建新的deployment成功 >>>>")
		allocation.Status.Total++
	}
	return nil
}
