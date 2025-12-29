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

// PhysicalPortSpec defines the desired state of PhysicalPort
type PhysicalPortSpec struct{}

// PhysicalPortStatus defines the observed state of PhysicalPort
type PhysicalPortStatus struct {
	State               string `json:"state,omitempty"`
	Duplex              string `json:"duplex,omitempty"`
	AutoNegotiated      bool   `json:"autoNegotiated,omitempty"`
	NegotiatedSpeed     int    `json:"negotiatedSpeed,omitempty"`
	MaxSpeed            int    `json:"maxSpeed,omitempty"`
	ObservedMAC         string `json:"observedMAC,omitempty"`
	AggregationEnabled  bool   `json:"aggregationEnabled"`
	AggregationProtocol string `json:"aggregationProtocol,omitempty"`
	IsUplink            bool   `json:"isUplink"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=port
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Observed MAC",type=string,JSONPath=`.status.observedMAC`
// +kubebuilder:printcolumn:name="Duplex",type=string,priority=1,JSONPath=`.status.duplex`
// +kubebuilder:printcolumn:name="Link Speed",type=number,priority=1,JSONPath=`.status.negotiatedSpeed`
// +kubebuilder:printcolumn:name="Uplink",type=string,priority=1,JSONPath=`.status.isUplink`
// +kubebuilder:printcolumn:name="Aggregating",type=string,priority=1,JSONPath=`.status.aggregationEnabled`
// +kubebuilder:printcolumn:name="Aggregation Protocol",type=string,priority=1,JSONPath=`.status.aggregationProtocol`
// +kubebuilder:printcolumn:name="Switch",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/switch`
// +kubebuilder:printcolumn:name="Rack",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/rack`
// +kubebuilder:printcolumn:name="Datacenter",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/datacenter`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/availability-zone`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// PhysicalPort is the Schema for the physicalports API
type PhysicalPort struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalPortSpec   `json:"spec,omitempty"`
	Status PhysicalPortStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PhysicalPortList contains a list of PhysicalPort
type PhysicalPortList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalPort `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalPort{}, &PhysicalPortList{})
}
