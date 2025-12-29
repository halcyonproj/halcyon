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

// RackSpec defines the desired state of Rack
type RackSpec struct {
	SerialNumber string       `json:"serialNumber"`
	RackTypeName string       `json:"rackTypeName"`
	Location     RackLocation `json:"location"`
}

type RackLocation struct {
	DataCenterName string `json:"dataCenterName"`
}

// RackStatus defines the observed state of Rack
type RackStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Rack Type",type=string,JSONPath=`.spec.rackTypeName`
// +kubebuilder:printcolumn:name="Serial Number",type=string,priority=1,JSONPath=`.spec.serialNumber`
// +kubebuilder:printcolumn:name="Datacenter",type=string,JSONPath=`.spec.location.dataCenterName`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,priority=1,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/availability-zone`
// +kubebuilder:printcolumn:name="Region",priority=1,type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// Rack is the Schema for the racks API
type Rack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RackSpec   `json:"spec,omitempty"`
	Status RackStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RackList contains a list of Rack
type RackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rack{}, &RackList{})
}
