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

// StoregroupsParameters are the configurable fields of a Storegroups.
type StoregroupsParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// StoregroupsObservation are the observable fields of a Storegroups.
type StoregroupsObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A StoregroupsSpec defines the desired state of a Storegroups.
type StoregroupsSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       StoregroupsParameters `json:"forProvider"`
}

// A StoregroupsStatus represents the observed state of a Storegroups.
type StoregroupsStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          StoregroupsObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Storegroups is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Storegroups struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoregroupsSpec   `json:"spec"`
	Status StoregroupsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StoregroupsList contains a list of Storegroups
type StoregroupsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Storegroups `json:"items"`
}

// Storegroups type metadata.
var (
	StoregroupsKind             = reflect.TypeOf(Storegroups{}).Name()
	StoregroupsGroupKind        = schema.GroupKind{Group: Group, Kind: StoregroupsKind}.String()
	StoregroupsKindAPIVersion   = StoregroupsKind + "." + SchemeGroupVersion.String()
	StoregroupsGroupVersionKind = SchemeGroupVersion.WithKind(StoregroupsKind)
)

func init() {
	SchemeBuilder.Register(&Storegroups{}, &StoregroupsList{})
}
