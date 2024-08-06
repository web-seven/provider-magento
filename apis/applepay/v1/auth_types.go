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

// AuthParameters are the configurable fields of a Auth.
type AuthParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// AuthObservation are the observable fields of a Auth.
type AuthObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A AuthSpec defines the desired state of a Auth.
type AuthSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       AuthParameters `json:"forProvider"`
}

// A AuthStatus represents the observed state of a Auth.
type AuthStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          AuthObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Auth is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Auth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthSpec   `json:"spec"`
	Status AuthStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AuthList contains a list of Auth
type AuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Auth `json:"items"`
}

// Auth type metadata.
var (
	AuthKind             = reflect.TypeOf(Auth{}).Name()
	AuthGroupKind        = schema.GroupKind{Group: Group, Kind: AuthKind}.String()
	AuthKindAPIVersion   = AuthKind + "." + SchemeGroupVersion.String()
	AuthGroupVersionKind = SchemeGroupVersion.WithKind(AuthKind)
)

func init() {
	SchemeBuilder.Register(&Auth{}, &AuthList{})
}
