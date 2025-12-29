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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PhysicalComputeSpec defines the desired state of PhysicalCompute
type PhysicalComputeSpec struct {
	Location   PhysicalComputeLocation     `json:"location"`
	Interfaces []PhysicalComputeInterfaces `json:"interfaces,omitempty"`
	BMC        PhysicalComputeBMC          `json:"bmc"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="";"on";"off";"unknown"
	UID          string `json:"UID,omitempty"`
	ComputeClass string `json:"computeClass"`
}

type PhysicalComputeBMC struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="";"on";"off";"unknown"
	PowerState string                       `json:"powerState"`
	Connection PhysicalComputeBMCConnection `json:"connection"`
}

type PhysicalComputeBMCConnection struct {
	// +kubebuilder:validation:Enum=ipmi;redfish
	Type          string                          `json:"type"`
	InterfaceName string                          `json:"interfaceName"`
	CredentialRef PhysicalComputeBMCCredentialRef `json:"credentialRef"`
}

type PhysicalComputeBMCCredentialRef struct {
	Name string `json:"name"`
}

type PhysicalComputeLocation struct {
	RackName string `json:"rackName"`
}

type PhysicalComputeInterfaces struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
}

// PhysicalComputeStatus defines the observed state of PhysicalCompute
type PhysicalComputeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Enum=on;off;unknown
	UID        string                           `json:"UID,omitempty"`
	BMC        PhysicalComputeBMCStatus         `json:"bmc"`
	DeviceInfo PhysicalComputeDeviceInfoStatus  `json:"deviceInfo"`
	Resources  PhysicalComputeResourcesStatus   `json:"resources"`
	Conditions []PhysicalComputeConditionStatus `json:"conditions,omitempty"`
}

type PhysicalComputeDeviceInfoStatus struct {
	Board   PhysicalComputeBoardStatus   `json:"board"`
	Chassis PhysicalComputeChassisStatus `json:"chassis"`
	Product PhysicalComputeProductStatus `json:"product"`
}

type PhysicalComputeResourcesStatus struct {
	Capacity    corev1.ResourceList `json:"capacity,omitempty"`
	Allocatable corev1.ResourceList `json:"allocatable,omitempty"`
}

type PhysicalComputeBMCStatus struct {
	GUID        string `json:"GUID,omitempty"`
	PowerPolicy string `json:"powerPolicy,omitempty"`
	// +kubebuilder:validation:Enum=on;off;unknown
	PowerState       string `json:"powerState,omitempty"`
	BIOSVersion      string `json:"BIOSVersion,omitempty"`
	BIOSBuildDate    string `json:"BIOSBuildDate,omitempty"`
	IPMIVersion      string `json:"IPMIVersion,omitempty"`
	FirmwareVersion  string `json:"firmwareVersion,omitempty"`
	ManufacturerID   string `json:"manufacturerID,omitempty"`
	ManufacturerName string `json:"manufacturerName,omitempty"`
}

type PhysicalComputeConditionStatus struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	// +kubebuilder:validation:Format="date-time"
	LastTransitionTime string `json:"lastTransitionTime"`
}

type PhysicalComputeBoardStatus struct {
	Serial           string `json:"serial,omitempty"`
	PartNumber       string `json:"partNumber,omitempty"`
	ManufacturerName string `json:"manufacturerName,omitempty"`
}

type PhysicalComputeChassisStatus struct {
	Serial     string `json:"serial,omitempty"`
	PartNumber string `json:"partNumber,omitempty"`
}

type PhysicalComputeProductStatus struct {
	Serial           string `json:"serial,omitempty"`
	PartNumber       string `json:"partNumber,omitempty"`
	ManufacturerName string `json:"manufacturerName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=compute
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Power State",type=string,JSONPath=`.status.bmc.powerState`
// +kubebuilder:printcolumn:name="UID",type=string,JSONPath=`.status.UID`
// +kubebuilder:printcolumn:name="Manufacturer",type=string,JSONPath=`.status.deviceInfo.product.manufacturerName`
// +kubebuilder:printcolumn:name="Serial Number",type=string,JSONPath=`.status.deviceInfo.product.serial`
// +kubebuilder:printcolumn:name="Rack",type=string,JSONPath=`.spec.location.rackName`
// +kubebuilder:printcolumn:name="Datacenter",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/datacenter`
// +kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/availability-zone`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.metadata.labels.topology\.halcyonproj\.dev/region`

// PhysicalCompute is the Schema for the physicalcomputes API
type PhysicalCompute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalComputeSpec   `json:"spec,omitempty"`
	Status PhysicalComputeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PhysicalComputeList contains a list of PhysicalCompute
type PhysicalComputeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalCompute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalCompute{}, &PhysicalComputeList{})
}
