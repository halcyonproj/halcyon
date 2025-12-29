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

// PhysicalSwitchSpec defines the desired state of PhysicalSwitch
type PhysicalSwitchSpec struct {
	Location    PhysicalSwitchLocation `json:"location"`
	Identifier  string                 `json:"identifier"`
	FabricClass string                 `json:"fabricClass"`
}

type PhysicalSwitchLocation struct {
	RackName string `json:"rackName"`
}

// PhysicalSwitchStatus defines the observed state of PhysicalSwitch
type PhysicalSwitchStatus struct {
	IPAddress    string `json:"ipAddress"`
	MACAddress   string `json:"macAddress"`
	DeviceID     string `json:"deviceID"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	Serial       string `json:"serial"`
	Version      string `json:"version"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=switch
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Manufacturer",type=string,JSONPath=`.status.manufacturer`
// +kubebuilder:printcolumn:name="Model",type=string,JSONPath=`.status.model`
// +kubebuilder:printcolumn:name="Identifier",type=string,JSONPath=`.spec.identifier`
// +kubebuilder:printcolumn:name="Rack",type=string,JSONPath=`.spec.location.rackName`
// +kubebuilder:printcolumn:name="Datacenter",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/datacenter`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/availability-zone`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// PhysicalSwitch is the Schema for the physicalswitches API
type PhysicalSwitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalSwitchSpec   `json:"spec,omitempty"`
	Status PhysicalSwitchStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PhysicalSwitchList contains a list of PhysicalSwitch
type PhysicalSwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalSwitch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalSwitch{}, &PhysicalSwitchList{})
}
