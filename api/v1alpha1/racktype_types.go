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

// RackTypeSpec defines the desired state of RackType
type RackTypeSpec struct {
	Manufacturer string             `json:"manufacturer"`
	Model        string             `json:"model"`
	FormFactor   RackTypeFormFactor `json:"formFactor"`
	Units        RackTypeUnits      `json:"units"`
}

type RackTypeFormFactor struct {
	// +kubebuilder:validation:Enum=cabinet;frame
	Enclosure string `json:"enclosure"`
	// +kubebuilder:validation:Enum=free-standing;mounted
	Position string `json:"position"`
	// +kubebuilder:validation:Enum=2;4
	Posts int `json:"posts"`
}

type RackTypeUnits struct {
	Start int `json:"start"`
	End   int `json:"end"`
	// +kubebuilder:validation:Enum=top-down;bottom-up
	Order string `json:"order"`
}

// RackTypeStatus defines the observed state of RackType
type RackTypeStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Manufacturer",type=string,JSONPath=`.spec.manufacturer`
// +kubebuilder:printcolumn:name="Model",type=string,JSONPath=`.spec.model`
// +kubebuilder:printcolumn:name="Enclosure",type=string,JSONPath=`.spec.formFactor.enclosure`

// RackType is the Schema for the racktypes API
type RackType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RackTypeSpec   `json:"spec,omitempty"`
	Status RackTypeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RackTypeList contains a list of RackType
type RackTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RackType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RackType{}, &RackTypeList{})
}
