// Package categories provides functions for interacting with Magento categories.
package categories

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/web-seven/provider-magento/apis/category/v1alpha1"

	magento "github.com/web-seven/provider-magento/internal/client"
)

// Client is a struct that embeds the magento.Client struct.
type Client struct {
	magento.Client
}

// categoriesPath is the path for the Magento categories API.
const categoriesPath = "/rest/default/V1/categories/"

// GetCategoryByID retrieves a category by its ID.
func GetCategoryByID(c *magento.Client, id string) (*v1alpha1.Category, error) {
	queryParams := map[string]string{
		"searchCriteria[filterGroups][0][filters][0][field]": "entity_id",
		"searchCriteria[filterGroups][0][filters][0][value]": id,
	}
	query := url.Values{}
	for key, value := range queryParams {
		query.Add(key, value)
	}
	resp, err := c.Create().R().SetQueryParams(queryParams).Get(categoriesPath + "list")

	if err != nil {
		return nil, err
	}

	var searchResults struct {
		Items []v1alpha1.CategoryParameters `json:"items"`
	}
	err = json.Unmarshal(resp.Body(), &searchResults)
	if err != nil {
		return nil, err
	}

	if len(searchResults.Items) > 0 {
		category := v1alpha1.Category{
			Spec: v1alpha1.CategorySpec{
				ForProvider: searchResults.Items[0],
			},
		}
		return &category, nil
	} else {
		return nil, errors.New("category not found")
	}
}

// CreateCategory creates a new category.
func CreateCategory(c *magento.Client, category *v1alpha1.Category) (*v1alpha1.CategoryObservation, *resty.Response, error) {
	requestBody := map[string]interface{}{
		"category": category.Spec.ForProvider,
	}

	resp, err := c.Create().R().SetHeader("Content-Type", "application/json").SetBody(requestBody).Post(categoriesPath)
	if err != nil {
		return nil, nil, err
	}

	var categoryObservation *v1alpha1.CategoryObservation
	err = json.Unmarshal(resp.Body(), &categoryObservation)
	if err != nil {
		return nil, nil, err
	}

	return categoryObservation, resp, nil
}

// UpdateCategoryByID updates a category by its ID.
func UpdateCategoryByID(c *magento.Client, id int, observed *v1alpha1.Category) error {
	requestBody := map[string]interface{}{
		"category": observed.Spec.ForProvider,
	}

	_, err := c.Create().R().SetBody(requestBody).Put(categoriesPath + strconv.Itoa(id))
	if err != nil {
		return err
	}

	return nil
}

// DeleteCategoryByID deletes a category by its ID.
func DeleteCategoryByID(c *magento.Client, id int) error {
	_, err := c.Create().R().Delete(categoriesPath + strconv.Itoa(id))
	return err
}

// IsUpToDate checks if the observed category is up to date with the desired category.
func IsUpToDate(observed *v1alpha1.Category, desired *v1alpha1.Category) (bool, error) {
	if observed == nil || desired == nil {
		return false, errors.New("observed or desired category is nil")
	}

	if observed.Spec.ForProvider.Name != desired.Spec.ForProvider.Name ||
		observed.Spec.ForProvider.IncludeInMenu != desired.Spec.ForProvider.IncludeInMenu ||
		observed.Spec.ForProvider.ParentID != desired.Spec.ForProvider.ParentID ||
		observed.Spec.ForProvider.Position != desired.Spec.ForProvider.Position ||
		observed.Spec.ForProvider.Level != desired.Spec.ForProvider.Level {
		return false, nil
	}

	return true, nil
}
