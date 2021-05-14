/*
Copyright 2021.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlgoCodingSpec defines the desired state of AlgoCoding
type AlgoCodingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AlgoCoding. Edit algocoding_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Number of Replicas (or pods)"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Replicas int32 `json:"replicas,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Web application image URL"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	WebImage string `json:"webImage,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="DB image URL associated with the Web application"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DBImage string `json:"dbImage,omitempty"`
}

// AlgoCodingStatus defines the observed state of AlgoCoding
type AlgoCodingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodNames []string `json:"podNames"`

	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Pod Size- here",xDescriptors="urn:alm:descriptor:com.tectonic.ui:Text"
	// the number of pod created
	Size int32 `json:"size,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AlgoCoding is the Schema for the algocodings API
type AlgoCoding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlgoCodingSpec   `json:"spec,omitempty"`
	Status AlgoCodingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlgoCodingList contains a list of AlgoCoding
type AlgoCodingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlgoCoding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlgoCoding{}, &AlgoCodingList{})
}
