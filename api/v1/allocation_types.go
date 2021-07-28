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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AllocationSpec defines the desired state of Allocation
type AllocationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image is the allocation pod image
	Image string `json:"image,omitempty"`
	// Replicas is the allocation pod replicas
	Replicas int32 `json:"replicas,omitempty"`
	// NodeName is the allocation pod nodeName
	NodeName string `json:"nodeName,omitempty"`
	// MsName is the allocation msName
	MsName string `json:"msName,omitempty"`
	// Interval is the allocation interval (*second)
	Interval int `json:"interval,omitempty"`
	// MaxNum is the allocation max numbers
	MaxNum int `json:"maxNum,omitempty"`
}

// AllocationStatus defines the observed state of Allocation
type AllocationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Total is the allocation total deployment numbers
	Total int `json:"total,omitempty"`
	// LastScheduleTime is the allocation last schedule time
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Allocation is the Schema for the allocations API
type Allocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AllocationSpec   `json:"spec,omitempty"`
	Status AllocationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AllocationList contains a list of Allocation
type AllocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Allocation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Allocation{}, &AllocationList{})
}
