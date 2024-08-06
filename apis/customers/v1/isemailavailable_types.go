/*
Copyright 2022 The Crossplane Authors.

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

package v1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// IsemailavailableParameters are the configurable fields of a Isemailavailable.
type IsemailavailableParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// IsemailavailableObservation are the observable fields of a Isemailavailable.
type IsemailavailableObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A IsemailavailableSpec defines the desired state of a Isemailavailable.
type IsemailavailableSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       IsemailavailableParameters `json:"forProvider"`
}

// A IsemailavailableStatus represents the observed state of a Isemailavailable.
type IsemailavailableStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          IsemailavailableObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Isemailavailable is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Isemailavailable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IsemailavailableSpec   `json:"spec"`
	Status IsemailavailableStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IsemailavailableList contains a list of Isemailavailable
type IsemailavailableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Isemailavailable `json:"items"`
}

// Isemailavailable type metadata.
var (
	IsemailavailableKind             = reflect.TypeOf(Isemailavailable{}).Name()
	IsemailavailableGroupKind        = schema.GroupKind{Group: Group, Kind: IsemailavailableKind}.String()
	IsemailavailableKindAPIVersion   = IsemailavailableKind + "." + SchemeGroupVersion.String()
	IsemailavailableGroupVersionKind = SchemeGroupVersion.WithKind(IsemailavailableKind)
)

func init() {
	SchemeBuilder.Register(&Isemailavailable{}, &IsemailavailableList{})
}
