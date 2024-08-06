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

// TypesParameters are the configurable fields of a Types.
type TypesParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// TypesObservation are the observable fields of a Types.
type TypesObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A TypesSpec defines the desired state of a Types.
type TypesSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TypesParameters `json:"forProvider"`
}

// A TypesStatus represents the observed state of a Types.
type TypesStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          TypesObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Types is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Types struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TypesSpec   `json:"spec"`
	Status TypesStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TypesList contains a list of Types
type TypesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Types `json:"items"`
}

// Types type metadata.
var (
	TypesKind             = reflect.TypeOf(Types{}).Name()
	TypesGroupKind        = schema.GroupKind{Group: Group, Kind: TypesKind}.String()
	TypesKindAPIVersion   = TypesKind + "." + SchemeGroupVersion.String()
	TypesGroupVersionKind = SchemeGroupVersion.WithKind(TypesKind)
)

func init() {
	SchemeBuilder.Register(&Types{}, &TypesList{})
}
