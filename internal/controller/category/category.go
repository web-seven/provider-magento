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

package category

import (
	"context"
	"net/http"
	"strconv"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	magento "github.com/web-seven/provider-magento/internal/client"
	"github.com/web-seven/provider-magento/internal/client/categories"

	"github.com/web-seven/provider-magento/apis/category/v1alpha1"
	apisv1alpha1 "github.com/web-seven/provider-magento/apis/v1alpha1"
	"github.com/web-seven/provider-magento/internal/features"
)

const (
	errNotCategory  = "managed resource is not a Category custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// A NoOpService does nothing.
type MagentoService struct {
	client *magento.Client
}

var (
	newMagentoService = func(creds []byte, baseURL string) (*MagentoService, error) {

		// Create a new Magento API client
		c := magento.NewClient(baseURL, string(creds))
		return &MagentoService{
			client: c,
		}, nil
	}
)

// Setup adds a controller that reconciles Category managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.CategoryGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.CategoryGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newMagentoService}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.Category{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte, baseURL string) (*MagentoService, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Category)
	if !ok {
		return nil, errors.New(errNotCategory)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data, pc.Spec.MagentoURL)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{service: svc}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service *MagentoService
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Category)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotCategory)
	}

	externalID := cr.GetAnnotations()["external-id"]
	desired, err := categories.GetCategoryByID(c.service.client, externalID)
	if err != nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, err
	}

	id, _ := strconv.Atoi(externalID)
	if desired != nil {
		cr.Status.SetConditions(xpv1.Available())
		cr.Status.AtProvider.Name = desired.Spec.ForProvider.Name
		cr.Status.AtProvider.ParentID = desired.Spec.ForProvider.ParentID
		cr.Status.AtProvider.IsActive = desired.Spec.ForProvider.IsActive
		cr.Status.AtProvider.Position = desired.Spec.ForProvider.Position
		cr.Status.AtProvider.Level = desired.Spec.ForProvider.Level
		cr.Status.AtProvider.ID = id
	}

	isUpToDate, _ := categories.IsUpToDate(cr, desired)
	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  isUpToDate,
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Category)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotCategory)
	}
	cr.SetConditions(xpv1.Creating())
	category, resp, err := categories.CreateCategory(c.service.client, cr)
	cr.SetAnnotations(map[string]string{"external-id": strconv.Itoa(category.ID)})
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

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Category)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotCategory)
	}
	err := categories.UpdateCategoryByID(c.service.client, cr.Status.AtProvider.ID, cr)

	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Category)
	if !ok {
		return errors.New(errNotCategory)
	}
	cr.SetConditions(xpv1.Deleting())
	err := categories.DeleteCategoryByID(c.service.client, cr.Status.AtProvider.ID)
	if err != nil {
		return err
	}

	return nil
}
