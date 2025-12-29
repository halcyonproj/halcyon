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

// AvailabilityZoneSpec defines the desired state of AvailabilityZone
type AvailabilityZoneSpec struct {
	Location AvailabilityZoneLocation `json:"location"`
}

type AvailabilityZoneLocation struct {
	RegionName string `json:"regionName"`
}

// AvailabilityZoneStatus defines the observed state of AvailabilityZone
type AvailabilityZoneStatus struct {
	Conditions []AvailabilityZoneCondition `json:"conditions"`
}

// AvailabilityZoneCondition defines conditions of the AvailabilityZone
type AvailabilityZoneCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastHeartbeatTime  string `json:"lastHeartbeatTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
	Message            string `json:"message"`
	Reason             string `json:"reason"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=az
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.spec.location.regionName`

// AvailabilityZone is the Schema for the availabilityzones API
type AvailabilityZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AvailabilityZoneSpec   `json:"spec,omitempty"`
	Status AvailabilityZoneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AvailabilityZoneList contains a list of AvailabilityZone
type AvailabilityZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AvailabilityZone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AvailabilityZone{}, &AvailabilityZoneList{})
}
