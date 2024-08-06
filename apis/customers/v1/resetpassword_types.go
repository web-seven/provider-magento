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

// ResetpasswordParameters are the configurable fields of a Resetpassword.
type ResetpasswordParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// ResetpasswordObservation are the observable fields of a Resetpassword.
type ResetpasswordObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A ResetpasswordSpec defines the desired state of a Resetpassword.
type ResetpasswordSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ResetpasswordParameters `json:"forProvider"`
}

// A ResetpasswordStatus represents the observed state of a Resetpassword.
type ResetpasswordStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ResetpasswordObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Resetpassword is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Resetpassword struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResetpasswordSpec   `json:"spec"`
	Status ResetpasswordStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ResetpasswordList contains a list of Resetpassword
type ResetpasswordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resetpassword `json:"items"`
}

// Resetpassword type metadata.
var (
	ResetpasswordKind             = reflect.TypeOf(Resetpassword{}).Name()
	ResetpasswordGroupKind        = schema.GroupKind{Group: Group, Kind: ResetpasswordKind}.String()
	ResetpasswordKindAPIVersion   = ResetpasswordKind + "." + SchemeGroupVersion.String()
	ResetpasswordGroupVersionKind = SchemeGroupVersion.WithKind(ResetpasswordKind)
)

func init() {
	SchemeBuilder.Register(&Resetpassword{}, &ResetpasswordList{})
}
