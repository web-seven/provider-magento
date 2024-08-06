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

// CategoriesParameters are the configurable fields of a Categories.
type CategoriesParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// CategoriesObservation are the observable fields of a Categories.
type CategoriesObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A CategoriesSpec defines the desired state of a Categories.
type CategoriesSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       CategoriesParameters `json:"forProvider"`
}

// A CategoriesStatus represents the observed state of a Categories.
type CategoriesStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          CategoriesObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Categories is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Categories struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CategoriesSpec   `json:"spec"`
	Status CategoriesStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CategoriesList contains a list of Categories
type CategoriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Categories `json:"items"`
}

// Categories type metadata.
var (
	CategoriesKind             = reflect.TypeOf(Categories{}).Name()
	CategoriesGroupKind        = schema.GroupKind{Group: Group, Kind: CategoriesKind}.String()
	CategoriesKindAPIVersion   = CategoriesKind + "." + SchemeGroupVersion.String()
	CategoriesGroupVersionKind = SchemeGroupVersion.WithKind(CategoriesKind)
)

func init() {
	SchemeBuilder.Register(&Categories{}, &CategoriesList{})
}
