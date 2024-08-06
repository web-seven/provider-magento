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

// SearchParameters are the configurable fields of a Search.
type SearchParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// SearchObservation are the observable fields of a Search.
type SearchObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A SearchSpec defines the desired state of a Search.
type SearchSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SearchParameters `json:"forProvider"`
}

// A SearchStatus represents the observed state of a Search.
type SearchStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SearchObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Search is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Search struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SearchSpec   `json:"spec"`
	Status SearchStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SearchList contains a list of Search
type SearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Search `json:"items"`
}

// Search type metadata.
var (
	SearchKind             = reflect.TypeOf(Search{}).Name()
	SearchGroupKind        = schema.GroupKind{Group: Group, Kind: SearchKind}.String()
	SearchKindAPIVersion   = SearchKind + "." + SchemeGroupVersion.String()
	SearchGroupVersionKind = SchemeGroupVersion.WithKind(SearchKind)
)

func init() {
	SchemeBuilder.Register(&Search{}, &SearchList{})
}
