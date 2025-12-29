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

// FabricClassSpec defines the desired state of FabricClass
type FabricClassSpec struct {
	ConfigurationFrom FabricClassConfiguration `json:"configurationFrom"`
}

type FabricClassConfiguration struct {
	SecretRef FabricClassConfigurationSecretRef `json:"secretRef"`
}

type FabricClassConfigurationSecretRef struct {
	Name string `json:"name"`
}

// FabricClassStatus defines the observed state of FabricClass
type FabricClassStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FabricClass is the Schema for the fabricclasses API
type FabricClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FabricClassSpec   `json:"spec,omitempty"`
	Status FabricClassStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FabricClassList contains a list of FabricClass
type FabricClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FabricClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FabricClass{}, &FabricClassList{})
}
