// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

type HostTemplate struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MonitoringServer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type HostGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Host struct {
	ID                     int              `json:"id"`
	Name                   string           `json:"name"`
	Alias                  string           `json:"alias"`
	Address                string           `json:"address"`
	MonitoringServer       MonitoringServer `json:"monitoring_server"`
	Templates              []HostTemplate   `json:"templates"`
	NormalCheckInterval    string           `json:"normal_check_interval"`
	RetryCheckInterval     string           `json:"retry_check_interval"`
	NotificationTimeperiod string           `json:"notification_timeperiod"`
	CheckTimeperiod        string           `json:"check_timeperiod"`
	Severity               string           `json:"severity"`
	Categories             []string         `json:"categories"`
	Groups                 []HostGroup      `json:"groups"`
	IsActivated            bool             `json:"is_activated"`
}

type HostResponse struct {
	Result []Host `json:"result"`
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

func (c *Client) GetHosts(limit int, page int, search string) (*HostResponse, error) {
	url := fmt.Sprintf("%s/configuration/hosts?limit=%d&page=%d&search=%s",
		c.BaseURL, limit, page, search)

	req, err := http.NewRequest("GET", url, nil)
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

	var hostResponse HostResponse
	if err := json.NewDecoder(resp.Body).Decode(&hostResponse); err != nil {
		return nil, err
	}

	return &hostResponse, nil
}
