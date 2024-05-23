package magento

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

// GetResourceByID retrieves a resource by its ID at specified api endpoint.
func GetResourceByID(c *Client, id string) (map[string]interface{}, error) {
	if id == "" {
		return nil, errors.New("resource with ID" + id + " in " + c.Path + " not found")
	}
	resp, _ := c.Create().R().Get(c.Path + "/" + id)
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("resource in " + c.Path + " not found")
	}

	var resource map[string]interface{}
	err := json.Unmarshal(resp.Body(), &resource)
	if err != nil {
		return nil, errors.New("failed to unmarshal response body")
	}

	return resource, nil
}

// CreateResource creates a new resource at specified api endpoint.
func CreateResource(c *Client, observed map[string]interface{}) (map[string]interface{}, *resty.Response, error) {

	requestBody := map[string]interface{}{
		strings.ToLower(observed["kind"].(string)): observed["spec"].(map[string]interface{})["forProvider"],
	}
	resp, err := c.Create().R().SetHeader("Content-Type", "application/json").SetBody(requestBody).Post(c.Path)
	if err != nil {
		return nil, nil, err
	}

	var desired map[string]interface{}
	err = json.Unmarshal(resp.Body(), &desired)
	if err != nil {
		return nil, nil, err
	}

	return desired, resp, nil
}

// UpdateResourceByID updates a resource by its ID at specified api endpoint.
func UpdateResourceByID(c *Client, id string, observed map[string]interface{}) error {
	requestBody := map[string]interface{}{
		strings.ToLower(observed["kind"].(string)): observed["spec"].(map[string]interface{})["forProvider"],
	}
	_, err := c.Create().R().SetBody(requestBody).Put(c.Path + "/" + id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteResourceByID deletes a resource by its ID at specified api endpoint.
func DeleteResourceByID(c *Client, id string) error {
	_, err := c.Create().R().Delete(c.Path + "/" + id)
	return err
}

// IsUpToDate checks if the observed resource is up to date with the desired resource.
func IsUpToDate(observed map[string]interface{}, desired map[string]interface{}) (bool, error) {

	if observed == nil || desired == nil {
		return false, errors.New("observed or desired resource is nil")
	}
	if observed["spec"].(map[string]interface{})["forProvider"].(map[string]interface{})["name"] != desired["name"] {
		return false, nil
	}

	return true, nil
}
