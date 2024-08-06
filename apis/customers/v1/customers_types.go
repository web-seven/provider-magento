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

// CustomersParameters are the configurable fields of a Customers.
type CustomersParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// CustomersObservation are the observable fields of a Customers.
type CustomersObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A CustomersSpec defines the desired state of a Customers.
type CustomersSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       CustomersParameters `json:"forProvider"`
}

// A CustomersStatus represents the observed state of a Customers.
type CustomersStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          CustomersObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Customers is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Customers struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomersSpec   `json:"spec"`
	Status CustomersStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomersList contains a list of Customers
type CustomersList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Customers `json:"items"`
}

// Customers type metadata.
var (
	CustomersKind             = reflect.TypeOf(Customers{}).Name()
	CustomersGroupKind        = schema.GroupKind{Group: Group, Kind: CustomersKind}.String()
	CustomersKindAPIVersion   = CustomersKind + "." + SchemeGroupVersion.String()
	CustomersGroupVersionKind = SchemeGroupVersion.WithKind(CustomersKind)
)

func init() {
	SchemeBuilder.Register(&Customers{}, &CustomersList{})
}
