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

// ProductsParameters are the configurable fields of a Products.
type ProductsParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// ProductsObservation are the observable fields of a Products.
type ProductsObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A ProductsSpec defines the desired state of a Products.
type ProductsSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ProductsParameters `json:"forProvider"`
}

// A ProductsStatus represents the observed state of a Products.
type ProductsStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ProductsObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Products is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Products struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductsSpec   `json:"spec"`
	Status ProductsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProductsList contains a list of Products
type ProductsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Products `json:"items"`
}

// Products type metadata.
var (
	ProductsKind             = reflect.TypeOf(Products{}).Name()
	ProductsGroupKind        = schema.GroupKind{Group: Group, Kind: ProductsKind}.String()
	ProductsKindAPIVersion   = ProductsKind + "." + SchemeGroupVersion.String()
	ProductsGroupVersionKind = SchemeGroupVersion.WithKind(ProductsKind)
)

func init() {
	SchemeBuilder.Register(&Products{}, &ProductsList{})
}
