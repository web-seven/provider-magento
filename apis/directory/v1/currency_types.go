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

// CurrencyParameters are the configurable fields of a Currency.
type CurrencyParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// CurrencyObservation are the observable fields of a Currency.
type CurrencyObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A CurrencySpec defines the desired state of a Currency.
type CurrencySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       CurrencyParameters `json:"forProvider"`
}

// A CurrencyStatus represents the observed state of a Currency.
type CurrencyStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          CurrencyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Currency is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Currency struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CurrencySpec   `json:"spec"`
	Status CurrencyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CurrencyList contains a list of Currency
type CurrencyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Currency `json:"items"`
}

// Currency type metadata.
var (
	CurrencyKind             = reflect.TypeOf(Currency{}).Name()
	CurrencyGroupKind        = schema.GroupKind{Group: Group, Kind: CurrencyKind}.String()
	CurrencyKindAPIVersion   = CurrencyKind + "." + SchemeGroupVersion.String()
	CurrencyGroupVersionKind = SchemeGroupVersion.WithKind(CurrencyKind)
)

func init() {
	SchemeBuilder.Register(&Currency{}, &CurrencyList{})
}
