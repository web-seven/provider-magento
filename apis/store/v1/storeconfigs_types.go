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

// StoreconfigsParameters are the configurable fields of a Storeconfigs.
type StoreconfigsParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// StoreconfigsObservation are the observable fields of a Storeconfigs.
type StoreconfigsObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A StoreconfigsSpec defines the desired state of a Storeconfigs.
type StoreconfigsSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       StoreconfigsParameters `json:"forProvider"`
}

// A StoreconfigsStatus represents the observed state of a Storeconfigs.
type StoreconfigsStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          StoreconfigsObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Storeconfigs is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Storeconfigs struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreconfigsSpec   `json:"spec"`
	Status StoreconfigsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StoreconfigsList contains a list of Storeconfigs
type StoreconfigsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Storeconfigs `json:"items"`
}

// Storeconfigs type metadata.
var (
	StoreconfigsKind             = reflect.TypeOf(Storeconfigs{}).Name()
	StoreconfigsGroupKind        = schema.GroupKind{Group: Group, Kind: StoreconfigsKind}.String()
	StoreconfigsKindAPIVersion   = StoreconfigsKind + "." + SchemeGroupVersion.String()
	StoreconfigsGroupVersionKind = SchemeGroupVersion.WithKind(StoreconfigsKind)
)

func init() {
	SchemeBuilder.Register(&Storeconfigs{}, &StoreconfigsList{})
}
