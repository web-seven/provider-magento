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

// StoreviewsParameters are the configurable fields of a Storeviews.
type StoreviewsParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// StoreviewsObservation are the observable fields of a Storeviews.
type StoreviewsObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A StoreviewsSpec defines the desired state of a Storeviews.
type StoreviewsSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       StoreviewsParameters `json:"forProvider"`
}

// A StoreviewsStatus represents the observed state of a Storeviews.
type StoreviewsStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          StoreviewsObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Storeviews is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Storeviews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreviewsSpec   `json:"spec"`
	Status StoreviewsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StoreviewsList contains a list of Storeviews
type StoreviewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Storeviews `json:"items"`
}

// Storeviews type metadata.
var (
	StoreviewsKind             = reflect.TypeOf(Storeviews{}).Name()
	StoreviewsGroupKind        = schema.GroupKind{Group: Group, Kind: StoreviewsKind}.String()
	StoreviewsKindAPIVersion   = StoreviewsKind + "." + SchemeGroupVersion.String()
	StoreviewsGroupVersionKind = SchemeGroupVersion.WithKind(StoreviewsKind)
)

func init() {
	SchemeBuilder.Register(&Storeviews{}, &StoreviewsList{})
}
