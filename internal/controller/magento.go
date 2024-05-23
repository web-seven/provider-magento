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

package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	magento "github.com/web-seven/provider-magento/internal/client"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apisv1alpha1 "github.com/web-seven/provider-magento/apis/v1alpha1"
	"github.com/web-seven/provider-magento/internal/features"
)

const (
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient = "cannot create new Service"
	api          = "/rest/"
	id           = "external-id"
)

// MagentoService is a service that can connect to Magento API.
type MagentoService struct {
	client *magento.Client
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube                   client.Client
	usage                  resource.Tracker
	createMagentoServiceFn func(creds []byte, baseURL string) (*MagentoService, error)
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	kube client.Client
	// A 'client' used to connect to the external resource API. In practice this
	service *MagentoService
}

// newMagentoService creates a new MagentoService.
var (
	newMagentoService = func(creds []byte, baseURL string) (*MagentoService, error) {

		// Create a new Magento API client
		c := magento.NewClient(baseURL, string(creds))
		return &MagentoService{
			client: c,
		}, nil
	}
)

// isValidGVK returns true if the GroupVersionKind is a valid Magento API resource.
func isValidGVK(gvk schema.GroupVersionKind) bool {
	return gvk.Group == "magento.web7.md" && gvk.Version == "v1alpha1" &&
		!strings.Contains(gvk.Kind, "Options") && !strings.Contains(gvk.Kind, "List") &&
		!strings.Contains(gvk.Kind, "Event") && !strings.Contains(gvk.Kind, "Config")
}

// Setup adds a controller that reconciles managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName("magento.web7.md")
	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}
	scheme := mgr.GetScheme()
	gvks := mgr.GetScheme().AllKnownTypes()

	// Filter GroupVersionKind to include resources that are part of the Magento API
	var filteredGvks []schema.GroupVersionKind
	for gvk := range gvks {
		if isValidGVK(gvk) {
			filteredGvks = append(filteredGvks, gvk)
		}
	}

	// Create controllers for each filtered GroupVersionKind from schema
	for _, gvk := range filteredGvks {
		obj, err := scheme.New(gvk)
		if err != nil {
			return err
		}
		r := managed.NewReconciler(mgr,
			resource.ManagedKind(gvk),
			managed.WithExternalConnecter(&connector{
				kube:                   mgr.GetClient(),
				usage:                  resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
				createMagentoServiceFn: newMagentoService}),
			managed.WithLogger(o.Logger.WithValues("controller", name)),
			managed.WithPollInterval(o.PollInterval),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
			managed.WithConnectionPublishers(cps...))

		err = ctrl.NewControllerManagedBy(mgr).
			Named(name).
			WithOptions(o.ForControllerRuntime()).
			WithEventFilter(resource.DesiredStateChanged()).
			For(obj.(client.Object)).
			Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))

		if err != nil {
			return err
		}
	}

	return nil
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.createMagentoServiceFn(data, pc.Spec.MagentoURL)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}
	client := c.kube
	return &external{service: svc, kube: client}, nil
}

// getDesiredCRD returns the CustomResourceDefinition that matches the group and kind.
func getDesiredCRD(crds *v1.CustomResourceDefinitionList, group string, kind string) *v1.CustomResourceDefinition {
	for _, crd := range crds.Items {
		if crd.Spec.Group == group && crd.Spec.Names.Kind == kind {
			return &crd
		}
	}
	return nil
}

// Observe checks if the external resource exists and if it is up to date with the managed resource.
func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	err := v1.AddToScheme(c.kube.Scheme())
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	crds := &v1.CustomResourceDefinitionList{}
	_ = c.kube.List(ctx, crds)
	crd := getDesiredCRD(crds, mg.GetObjectKind().GroupVersionKind().Group, mg.GetObjectKind().GroupVersionKind().Kind)
	c.service.client.Path = api + strings.ToUpper(mg.GetObjectKind().GroupVersionKind().Version) + "/" + crd.Spec.Names.Plural

	externalID := mg.GetAnnotations()[id]
	desired, err := magento.GetResourceByID(c.service.client, externalID)
	if err != nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, err
	}

	observed, err := runtime.DefaultUnstructuredConverter.ToUnstructured(mg)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	if desired != nil {
		observed["status"].(map[string]interface{})["atProvider"].(map[string]interface{})["id"] = desired["id"]
		observed["status"].(map[string]interface{})["atProvider"].(map[string]interface{})["name"] = desired["name"]
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(observed, mg)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	mg.SetConditions(xpv1.Available())

	isUpToDate, _ := magento.IsUpToDate(observed, desired)

	return managed.ExternalObservation{

		ResourceExists:    true,
		ResourceUpToDate:  isUpToDate,
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

// Create a new resource at the external API.
func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	mg.SetConditions(xpv1.Creating())

	observed, err := runtime.DefaultUnstructuredConverter.ToUnstructured(mg)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	resource, resp, err := magento.CreateResource(c.service.client, observed)
	mg.SetAnnotations(map[string]string{id: fmt.Sprintf("%v", resource["id"])})

	if resp.StatusCode() != http.StatusOK {
		return managed.ExternalCreation{}, errors.New(resp.String())
	}
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

// Update the external resource to reflect the managed resource's desired state.
func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	observed, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(mg)
	externalID := mg.GetAnnotations()[id]

	err := magento.UpdateResourceByID(c.service.client, externalID, observed)

	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

// Delete the external resource.
func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	externalID := mg.GetAnnotations()[id]
	mg.SetConditions(xpv1.Deleting())
	err := magento.DeleteResourceByID(c.service.client, externalID)
	if err != nil {
		return err
	}

	return nil
}
