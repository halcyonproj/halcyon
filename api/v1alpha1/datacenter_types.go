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

// DataCenterSpec defines the desired state of DataCenter
type DataCenterSpec struct {
	// +kubebuilder:validation:Enum="private";"colocation"
	Type     string             `json:"type"`
	Provider string             `json:"provider"`
	Location DataCenterLocation `json:"location"`
}

type DataCenterLocation struct {
	Address              string `json:"address,omitempty"`
	AvailabilityZoneName string `json:"availabilityZoneName"`
}

// DataCenterStatus defines the observed state of DataCenter
type DataCenterStatus struct {
	Conditions []DataCenterCondition `json:"conditions"`
}

// DataCenterCondition defines conditions of the DataCenter
type DataCenterCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastHeartbeatTime  string `json:"lastHeartbeatTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
	Message            string `json:"message"`
	Reason             string `json:"reason"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=dc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.spec.location.availabilityZoneName`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// DataCenter is the Schema for the datacenters API
type DataCenter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataCenterSpec   `json:"spec,omitempty"`
	Status DataCenterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DataCenterList contains a list of DataCenter
type DataCenterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataCenter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataCenter{}, &DataCenterList{})
}
