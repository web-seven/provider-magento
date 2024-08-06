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

// WebsitesParameters are the configurable fields of a Websites.
type WebsitesParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// WebsitesObservation are the observable fields of a Websites.
type WebsitesObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A WebsitesSpec defines the desired state of a Websites.
type WebsitesSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       WebsitesParameters `json:"forProvider"`
}

// A WebsitesStatus represents the observed state of a Websites.
type WebsitesStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          WebsitesObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Websites is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,category={crossplane,managed,magento}
type Websites struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebsitesSpec   `json:"spec"`
	Status WebsitesStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebsitesList contains a list of Websites
type WebsitesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Websites `json:"items"`
}

// Websites type metadata.
var (
	WebsitesKind             = reflect.TypeOf(Websites{}).Name()
	WebsitesGroupKind        = schema.GroupKind{Group: Group, Kind: WebsitesKind}.String()
	WebsitesKindAPIVersion   = WebsitesKind + "." + SchemeGroupVersion.String()
	WebsitesGroupVersionKind = SchemeGroupVersion.WithKind(WebsitesKind)
)

func init() {
	SchemeBuilder.Register(&Websites{}, &WebsitesList{})
}
