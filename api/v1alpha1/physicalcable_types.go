/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PhysicalCableSpec defines the desired state of PhysicalCable
type PhysicalCableSpec struct {
	SourceInterface string `json:"sourceInterface"`
	TargetPort      string `json:"targetPort"`
}

// PhysicalCableStatus defines the observed state of PhysicalCable
type PhysicalCableStatus struct {
	State string `json:"state,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=cable
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Source Interface",type=string,JSONPath=`.spec.sourceInterface`
// +kubebuilder:printcolumn:name="Target Port",type=string,JSONPath=`.spec.targetPort`
// +kubebuilder:printcolumn:name="Datacenter",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/datacenter`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/availability-zone`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// PhysicalCable is the Schema for the physicalcables API
type PhysicalCable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalCableSpec   `json:"spec,omitempty"`
	Status PhysicalCableStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PhysicalCableList contains a list of PhysicalCable
type PhysicalCableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalCable `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalCable{}, &PhysicalCableList{})
}
