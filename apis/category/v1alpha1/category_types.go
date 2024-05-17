/*
Copyright 2024 Web Seven license.

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

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// CustomAttributes type
type CustomAttributes struct {
	AttributeCode string `json:"attribute_code"`
	Value         string `json:"value"`
}

// CategoryParameters are the configurable fields of a Category.
type CategoryParameters struct {
	Name             string             `json:"name,omitempty"`
	IsActive         bool               `json:"isActive,omitempty"`
	Position         int                `json:"position,omitempty"`
	Level            int                `json:"level,omitempty"`
	Children         string             `json:"children,omitempty"`
	CreatedAt        string             `json:"createdAt,omitempty"`
	UpdatedAt        string             `json:"updatedAt,omitempty"`
	Path             string             `json:"path,omitempty"`
	AvailableSortBy  []string           `json:"availableSortBy,omitempty"`
	IncludeInMenu    bool               `json:"includeInMenu,omitempty"`
	CustomAttributes []CustomAttributes `json:"customAttributes,omitempty"`
	ParentID         int                `json:"parentId,omitempty"`
}

// CategoryObservation are the observable fields of a Category.
type CategoryObservation struct {
	ID           int    `json:"id,omitempty"`
	ParentID     int    `json:"parentId,omitempty"`
	Name         string `json:"name,omitempty"`
	IsActive     bool   `json:"isActive,omitempty"`
	Position     int    `json:"position,omitempty"`
	Level        int    `json:"level,omitempty"`
	ProductCount int    `json:"productCount,omitempty"`
}

// A CategorySpec defines the desired state of a Category.
type CategorySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       CategoryParameters `json:"forProvider"`
}

// A CategoryStatus represents the observed state of a Category.
type CategoryStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          CategoryObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Category is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,magento}
type Category struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CategorySpec   `json:"spec"`
	Status CategoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CategoryList contains a list of Category
type CategoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Category `json:"items"`
}

// Category type metadata.
var (
	CategoryKind             = reflect.TypeOf(Category{}).Name()
	CategoryGroupKind        = schema.GroupKind{Group: Group, Kind: CategoryKind}.String()
	CategoryKindAPIVersion   = CategoryKind + "." + SchemeGroupVersion.String()
	CategoryGroupVersionKind = SchemeGroupVersion.WithKind(CategoryKind)
)

func init() {
	SchemeBuilder.Register(&Category{}, &CategoryList{})
}
