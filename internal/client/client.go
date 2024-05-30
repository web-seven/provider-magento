package magento

import (
	"github.com/go-resty/resty/v2"
)

// Client struct to hold the Magento API client configuration
type Client struct {
	BaseURL     string
	AccessToken string
	Path        string
}

// NewClient initializes a new Magento API client configuration
func NewClient(baseURL, accessToken string) *Client {
	return &Client{
		BaseURL:     baseURL,
		AccessToken: accessToken,
	}
}

// CreateRestyClient creates a new resty client with configuration from Client
func (c *Client) Create() *resty.Client {
	restyClient := resty.New()
	restyClient.SetBaseURL(c.BaseURL)
	restyClient.SetAuthToken(c.AccessToken)

	return restyClient
}
