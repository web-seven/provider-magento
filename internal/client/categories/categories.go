package categories

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/web-seven/provider-magento/apis/category/v1alpha1"

	magento "github.com/web-seven/provider-magento/internal/client"
)

type Client struct {
	magento.Client
}

func GetCategoryByName(c *magento.Client, categoryName string) (*v1alpha1.Category, error) {
	path := fmt.Sprintf("%s/rest/all/V1/categories/list", c.BaseURL)

	queryParams := map[string]string{
		"searchCriteria[filterGroups][0][filters][0][field]": "name",
		"searchCriteria[filterGroups][0][filters][0][value]": categoryName,
	}

	query := url.Values{}
	for key, value := range queryParams {
		query.Add(key, value)
	}

	restyClient, err := c.CreateRestyClient()
	if err != nil {
		return nil, err
	}

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+c.AccessToken).
		SetQueryParams(queryParams).
		Get(path)
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

func CreateCategory(c *magento.Client, category *v1alpha1.Category) (*v1alpha1.Category, error) {
	path := c.BaseURL + "/rest/default/V1/categories"
	requestBody := map[string]interface{}{
		"category": category.Spec.ForProvider,
	}
	restyClient, err := c.CreateRestyClient()
	if err != nil {
		return nil, err
	}

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+c.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(path)

	if err != nil {
		return nil, errors.New("failed to create category: " + err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status() + ": " + string(resp.Body()))
	}

	var newCategory *v1alpha1.Category
	err = json.Unmarshal(resp.Body(), &newCategory)
	if err != nil {
		return nil, errors.New("failed to parse response: " + err.Error())
	}

	return newCategory, nil
}

func UpdateCategoryByName(c *magento.Client, categoryName string, observed *v1alpha1.Category) error {
	category, err := GetCategoryByName(c, categoryName)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/rest/all/V1/categories/%d", c.BaseURL, category.Spec.ForProvider.ID)

	restyClient, err := c.CreateRestyClient()
	if err != nil {
		return err
	}

	requestBody := map[string]interface{}{
		"category": observed.Spec.ForProvider,
	}

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+c.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Put(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("failed to update category")
	}

	return nil
}

func DeleteCategoryByName(c *magento.Client, categoryName string) error {
	category, err := GetCategoryByName(c, categoryName)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/rest/all/V1/categories/%d", c.BaseURL, category.Spec.ForProvider.ID)

	restyClient, err := c.CreateRestyClient()
	if err != nil {
		return err
	}

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+c.AccessToken).
		Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("failed to delete category")
	}

	return nil
}

func IsCategoryUpToDate(name string, observed *v1alpha1.Category, desired *v1alpha1.Category) (bool, error) {
	if observed == nil || desired == nil {
		return false, errors.New("observed or desired category is nil")
	}

	if observed.Spec.ForProvider.Name != desired.Spec.ForProvider.Name ||
		observed.Spec.ForProvider.IsActive != desired.Spec.ForProvider.IsActive ||
		observed.Spec.ForProvider.IncludeInMenu != desired.Spec.ForProvider.IncludeInMenu {
		return false, nil
	}

	return true, nil
}
