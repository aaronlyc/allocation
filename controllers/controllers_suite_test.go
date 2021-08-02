package controllers

import (
	"context"
	"fmt"
	myappv1 "github.com/aaronlyc/allocation/api/v1"
	"github.com/aaronlyc/allocation/controllers/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Testing with Ginkgo", func() {
	testConfig := &myappv1.Allocation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-allocate",
			Namespace: "test-ns",
		},
		Spec: myappv1.AllocationSpec{
			Image:    "busybox:latest",
			Replicas: 1,
			NodeName: "test-node",
			MsName:   "test-ms",
			Interval: 10,
			MaxNum:   100,
		},
	}

	var c *mock.MockClient
	var ct *gomock.Controller

	BeforeEach(func() {
		var t gomock.TestReporter
		ct = gomock.NewController(t)
		c = mock.NewMockClient(ct)
	}, 10)

	AfterEach(func() {
		ct.Finish() // 断言 DB.Get() 方法是否被调用
	}, 10)

	Describe("test GetAllocate", func() {

		It("test no error", func() {
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj client.Object) error {
				p := obj.(*myappv1.Allocation)
				*p = *testConfig

				return nil
			})
			testAllocate := &AllocationReconciler{
				Log:    ctrl.Log,
				Client: c,
			}
			allocate, err := testAllocate.GetAllocate(context.TODO(), types.NamespacedName{})
			Expect(err).NotTo(HaveOccurred())
			Expect(allocate.Name).To(Equal("test-allocate"))
		})

		It("test has error", func() {
			testErr := fmt.Errorf("test has error")
			c.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ client.ObjectKey, _ client.Object) error {
				return testErr
			})
			testAllocate := &AllocationReconciler{
				Log:    ctrl.Log,
				Client: c,
			}
			_, err := testAllocate.GetAllocate(context.TODO(), types.NamespacedName{})
			Expect(err).To(Equal(testErr))
		})
	})

	Describe("test dealAllocate", func() {
		It("test deal allocate", func() {
			// create成功时
			c.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			testAllocate := &AllocationReconciler{
				Log:    ctrl.Log,
				Client: c,
			}
			err := testAllocate.dealAllocate(context.TODO(), testConfig)
			Expect(err).NotTo(HaveOccurred())

			// create失败时
			c.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("处理失败"))
			testAllocate2 := &AllocationReconciler{
				Log:    ctrl.Log,
				Client: c,
			}
			err = testAllocate2.dealAllocate(context.TODO(), testConfig)
			Expect(err).To(HaveOccurred())
		})
	})

})
