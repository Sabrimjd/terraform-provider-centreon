package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL    string
	APIKey     string
	Protocol   string
	Server     string
	Port       string
	APIVersion string
	HTTPClient *http.Client
}

type PlatformInfo struct {
	IsInstalled         bool `json:"is_installed"`
	HasUpgradeAvailable bool `json:"has_upgrade_available"`
}

func NewClient(protocol, server, port, apiVersion, apiKey string) *Client {
	return &Client{
		Protocol:   protocol,
		Server:     server,
		Port:       port,
		APIVersion: apiVersion,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
		BaseURL:    fmt.Sprintf("%s://%s:%s/centreon/api/%s", protocol, server, port, apiVersion),
	}
}

func (c *Client) GetPlatformInfo() (*PlatformInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/platform/installation/status", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-AUTH-TOKEN", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var platformInfo PlatformInfo
	if err := json.NewDecoder(resp.Body).Decode(&platformInfo); err != nil {
		return nil, err
	}

	return &platformInfo, nil
}