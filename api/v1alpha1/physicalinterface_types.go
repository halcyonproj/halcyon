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

// PhysicalInterfaceSpec defines the desired state of PhysicalInterface
type PhysicalInterfaceSpec struct {
	MAC    string                  `json:"mac,omitempty"`
	Device PhysicalInterfaceDevice `json:"device,omitempty"`
}

type PhysicalInterfaceDevice struct {
	Type    string `json:"type,omitempty"`
	RefName string `json:"refName,omitempty"`
}

// PhysicalInterfaceStatus defines the observed state of PhysicalInterface
type PhysicalInterfaceStatus struct {
	IPAddress  string `json:"ipAddress,omitempty"`
	MACAddress string `json:"macAddress,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="IP Address",type=string,JSONPath=`.status.ipAddress`
// +kubebuilder:printcolumn:name="MAC Address",type=string,JSONPath=`.status.macAddress`
// +kubebuilder:printcolumn:name="Compute",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/compute`
// +kubebuilder:printcolumn:name="Datacenter",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/datacenter`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/availability-zone`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// PhysicalInterface is the Schema for the physicalinterfaces API
type PhysicalInterface struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalInterfaceSpec   `json:"spec,omitempty"`
	Status PhysicalInterfaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PhysicalInterfaceList contains a list of PhysicalInterface
type PhysicalInterfaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalInterface `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalInterface{}, &PhysicalInterfaceList{})
}
