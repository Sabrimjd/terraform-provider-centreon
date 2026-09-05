package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// defaultTimeout bounds every HTTP round trip so a hung Centreon API can
// never wedge a Terraform apply indefinitely.
const defaultTimeout = 30 * time.Second

// hostMutex serializes host mutations (create/update/delete). The Centreon
// API historically rejects concurrent writes to the configuration, so writes
// are kept sequential — but without artificial delays. Reads stay parallel.
var hostMutex sync.Mutex

type Client struct {
	BaseURL                        string
	APIKey                         string
	Protocol                       string
	Server                         string
	Port                           string
	APIVersion                     string
	HTTPClient                     *http.Client
	GenerateAndReloadConfiguration bool
}

type PlatformInfo struct {
	IsInstalled         bool `json:"is_installed"`
	HasUpgradeAvailable bool `json:"has_upgrade_available"`
}

type HostTemplate struct {
	ID                        int      `json:"id"`
	Name                      string   `json:"name"`
	Alias                     string   `json:"alias"`
	SNMPVersion               *string  `json:"snmp_version"`
	TimezoneID                *int     `json:"timezone_id"`
	SeverityID                *int     `json:"severity_id"`
	CheckCommandID            *int     `json:"check_command_id"`
	CheckCommandArgs          []string `json:"check_command_args"`
	CheckTimeperiodID         *int     `json:"check_timeperiod_id"`
	MaxCheckAttempts          *int     `json:"max_check_attempts"`
	NormalCheckInterval       *int     `json:"normal_check_interval"`
	RetryCheckInterval        *int     `json:"retry_check_interval"`
	ActiveCheckEnabled        int      `json:"active_check_enabled"`
	PassiveCheckEnabled       int      `json:"passive_check_enabled"`
	NotificationEnabled       int      `json:"notification_enabled"`
	NotificationOptions       *int     `json:"notification_options"`
	NotificationInterval      *int     `json:"notification_interval"`
	NotificationTimeperiodID  *int     `json:"notification_timeperiod_id"`
	AddInheritedContactGroup  bool     `json:"add_inherited_contact_group"`
	AddInheritedContact       bool     `json:"add_inherited_contact"`
	FirstNotificationDelay    *int     `json:"first_notification_delay"`
	RecoveryNotificationDelay *int     `json:"recovery_notification_delay"`
	AcknowledgementTimeout    *int     `json:"acknowledgement_timeout"`
	FreshnessChecked          int      `json:"freshness_checked"`
	FreshnessThreshold        *int     `json:"freshness_threshold"`
	FlapDetectionEnabled      int      `json:"flap_detection_enabled"`
	LowFlapThreshold          *int     `json:"low_flap_threshold"`
	HighFlapThreshold         *int     `json:"high_flap_threshold"`
	EventHandlerEnabled       int      `json:"event_handler_enabled"`
	EventHandlerCommandID     *int     `json:"event_handler_command_id"`
	EventHandlerCommandArgs   []string `json:"event_handler_command_args"`
	NoteURL                   *string  `json:"note_url"`
	Note                      *string  `json:"note"`
	ActionURL                 *string  `json:"action_url"`
	IconID                    *int     `json:"icon_id"`
	IconAlternative           *string  `json:"icon_alternative"`
	Comment                   string   `json:"comment"`
	IsLocked                  bool     `json:"is_locked"`
}

type MonitoringServer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MonitoringServerDetail struct {
	ID                       int     `json:"id"`
	Name                     string  `json:"name"`
	Address                  string  `json:"address"`
	IsLocalhost              bool    `json:"is_localhost"`
	IsDefault                bool    `json:"is_default"`
	SSHPort                  int     `json:"ssh_port"`
	LastRestart              string  `json:"last_restart"`
	EngineStartCommand       string  `json:"engine_start_command"`
	EngineStopCommand        string  `json:"engine_stop_command"`
	EngineRestartCommand     string  `json:"engine_restart_command"`
	EngineReloadCommand      string  `json:"engine_reload_command"`
	NagiosBin                string  `json:"nagios_bin"`
	NagiostatsBin            string  `json:"nagiostats_bin"`
	BrokerReloadCommand      string  `json:"broker_reload_command"`
	CentreonBrokerCfgPath    string  `json:"centreonbroker_cfg_path"`
	CentreonBrokerModulePath string  `json:"centreonbroker_module_path"`
	CentreonBrokerLogsPath   *string `json:"centreonbroker_logs_path"`
	CentreonConnectorPath    string  `json:"centreonconnector_path"`
	InitScriptCentreontrapd  string  `json:"init_script_centreontrapd"`
	SnmpTrapdPathConf        string  `json:"snmp_trapd_path_conf"`
	RemoteID                 *int    `json:"remote_id"`
	RemoteServerUseAsProxy   bool    `json:"remote_server_use_as_proxy"`
	IsUpdated                bool    `json:"is_updated"`
	IsActivate               bool    `json:"is_activate"`
}

type HostGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Host struct {
	ID                        int              `json:"id"`
	Name                      string           `json:"name"`
	Alias                     string           `json:"alias"`
	Address                   string           `json:"address"`
	MonitoringServer          MonitoringServer `json:"monitoring_server"`
	SNMPCommunity             string           `json:"snmp_community"`
	SNMPVersion               string           `json:"snmp_version"`
	TimezoneID                int              `json:"timezone_id"`
	SeverityID                int              `json:"severity_id"`
	CheckCommandID            int              `json:"check_command_id"`
	CheckCommandArgs          []string         `json:"check_command_args"`
	CheckTimeperiodID         int              `json:"check_timeperiod_id"`
	MaxCheckAttempts          int              `json:"max_check_attempts"`
	NormalCheckInterval       int              `json:"normal_check_interval"`
	RetryCheckInterval        int              `json:"retry_check_interval"`
	ActiveCheckEnabled        int              `json:"active_check_enabled"`  // 0=disabled, 1=enabled, 2=inherit/default
	PassiveCheckEnabled       int              `json:"passive_check_enabled"` // 0=disabled, 1=enabled, 2=inherit/default
	NotificationEnabled       int              `json:"notification_enabled"`  // 0=disabled, 1=enabled, 2=inherit/default
	NotificationOptions       int              `json:"notification_options"`
	NotificationInterval      int              `json:"notification_interval"`
	NotificationTimeperiodID  int              `json:"notification_timeperiod_id"`
	AddInheritedContactGroup  bool             `json:"add_inherited_contact_group"`
	AddInheritedContact       bool             `json:"add_inherited_contact"`
	FirstNotificationDelay    int              `json:"first_notification_delay"`
	RecoveryNotificationDelay int              `json:"recovery_notification_delay"`
	AcknowledgementTimeout    int              `json:"acknowledgement_timeout"`
	FreshnessChecked          int              `json:"freshness_checked"` // 0=disabled, 1=enabled, 2=inherit/default
	FreshnessThreshold        int              `json:"freshness_threshold"`
	FlapDetectionEnabled      int              `json:"flap_detection_enabled"` // 0=disabled, 1=enabled, 2=inherit/default
	LowFlapThreshold          int              `json:"low_flap_threshold"`
	HighFlapThreshold         int              `json:"high_flap_threshold"`
	EventHandlerEnabled       int              `json:"event_handler_enabled"` // 0=disabled, 1=enabled, 2=inherit/default
	EventHandlerCommandID     int              `json:"event_handler_command_id"`
	EventHandlerCommandArgs   []string         `json:"event_handler_command_args"`
	NoteURL                   string           `json:"note_url"`
	Note                      string           `json:"note"`
	ActionURL                 string           `json:"action_url"`
	IconID                    int              `json:"icon_id"`
	IconAlternative           string           `json:"icon_alternative"`
	Comment                   string           `json:"comment"`
	Categories                []int            `json:"categories"`
	Groups                    []HostGroup      `json:"groups"`
	Templates                 []HostTemplate   `json:"templates"`
	IsActivated               bool             `json:"is_activated"`
	GeoCoords                 string           `json:"geo_coords"`
}

type HostResponse struct {
	Result []Host `json:"result"`
}

type CreateHostRequest struct {
	MonitoringServerID        int         `json:"monitoring_server_id"`
	Name                      string      `json:"name"`
	Address                   string      `json:"address"`
	Alias                     *string     `json:"alias,omitempty"`
	SNMPCommunity             *string     `json:"snmp_community,omitempty"`
	SNMPVersion               *string     `json:"snmp_version,omitempty"`
	TimezoneID                *int        `json:"timezone_id,omitempty"`
	SeverityID                *int        `json:"severity_id,omitempty"`
	CheckCommandID            *int        `json:"check_command_id,omitempty"`
	CheckCommandArgs          []string    `json:"check_command_args,omitempty"`
	CheckTimeperiodID         *int        `json:"check_timeperiod_id,omitempty"`
	MaxCheckAttempts          *int        `json:"max_check_attempts,omitempty"`
	NormalCheckInterval       *int        `json:"normal_check_interval,omitempty"`
	RetryCheckInterval        *int        `json:"retry_check_interval,omitempty"`
	ActiveCheckEnabled        *int        `json:"active_check_enabled,omitempty"`
	PassiveCheckEnabled       *int        `json:"passive_check_enabled,omitempty"`
	NotificationEnabled       *int        `json:"notification_enabled,omitempty"`
	NotificationOptions       *int        `json:"notification_options,omitempty"`
	NotificationInterval      *int        `json:"notification_interval,omitempty"`
	NotificationTimeperiodID  *int        `json:"notification_timeperiod_id,omitempty"`
	AddInheritedContactGroup  *bool       `json:"add_inherited_contact_group,omitempty"`
	AddInheritedContact       *bool       `json:"add_inherited_contact,omitempty"`
	FirstNotificationDelay    *int        `json:"first_notification_delay,omitempty"`
	RecoveryNotificationDelay *int        `json:"recovery_notification_delay,omitempty"`
	AcknowledgementTimeout    *int        `json:"acknowledgement_timeout,omitempty"`
	FreshnessChecked          *int        `json:"freshness_checked,omitempty"`
	FreshnessThreshold        *int        `json:"freshness_threshold,omitempty"`
	FlapDetectionEnabled      *int        `json:"flap_detection_enabled,omitempty"`
	LowFlapThreshold          *int        `json:"low_flap_threshold,omitempty"`
	HighFlapThreshold         *int        `json:"high_flap_threshold,omitempty"`
	EventHandlerEnabled       *int        `json:"event_handler_enabled,omitempty"`
	EventHandlerCommandID     *int        `json:"event_handler_command_id,omitempty"`
	EventHandlerCommandArgs   []string    `json:"event_handler_command_args,omitempty"`
	NoteURL                   *string     `json:"note_url,omitempty"`
	Note                      *string     `json:"note,omitempty"`
	ActionURL                 *string     `json:"action_url,omitempty"`
	IconID                    *int        `json:"icon_id,omitempty"`
	IconAlternative           *string     `json:"icon_alternative,omitempty"`
	Comment                   *string     `json:"comment,omitempty"`
	IsActivated               *bool       `json:"is_activated,omitempty"`
	Categories                []int       `json:"categories,omitempty"`
	Groups                    []int       `json:"groups,omitempty"`
	Templates                 []int       `json:"templates,omitempty"`
	Macros                    []HostMacro `json:"macros,omitempty"`
	GeoCoords                 *string     `json:"geo_coords,omitempty"`
}

type HostMacro struct {
	Name        string  `json:"name"`
	Value       *string `json:"value"`
	IsPassword  bool    `json:"is_password"`
	Description *string `json:"description"`
}

type MonitoringServersResponse struct {
	Result []MonitoringServerDetail `json:"result"`
	Meta   Meta                     `json:"meta"`
}

type HostGroupsResponse struct {
	Result []HostGroup `json:"result"`
	Meta   Meta        `json:"meta"`
}

type HostTemplatesResponse struct {
	Result []HostTemplate `json:"result"`
	Meta   Meta           `json:"meta"`
}

type HostMacroResponse struct {
	Result []HostMacro `json:"result"`
}

type Meta struct {
	Page   int                    `json:"page"`
	Limit  int                    `json:"limit"`
	Search map[string]interface{} `json:"search"`
	SortBy map[string]interface{} `json:"sort_by"`
	Total  int                    `json:"total"`
}

func NewClient(protocol, server, port, apiVersion, apiKey string) *Client {
	return &Client{
		Protocol:   protocol,
		Server:     server,
		Port:       port,
		APIVersion: apiVersion,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		BaseURL:    fmt.Sprintf("%s://%s:%s/centreon/api/%s", protocol, server, port, apiVersion),
	}
}

// prettyPrintJSON formats a JSON string with indentation if it parses as
// JSON; otherwise it returns the input unchanged.
func prettyPrintJSON(input string) string {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 {
		return input
	}
	if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
		var parsed interface{}
		if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
			if pretty, err := json.MarshalIndent(parsed, "", "  "); err == nil {
				return string(pretty)
			}
		}
	}
	return input
}

// doRequest executes an HTTP request against the Centreon API. It never logs
// request or response bodies (they can contain SNMP communities and macro
// passwords) and returns non-2xx responses as *APIError.
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("X-AUTH-TOKEN", c.APIKey)

	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		tflog.Warn(ctx, "Centreon API request failed", map[string]interface{}{
			"method":   req.Method,
			"path":     req.URL.Path,
			"duration": time.Since(start).String(),
			"error":    err.Error(),
		})
		return nil, fmt.Errorf("error making request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		tflog.Warn(ctx, "Error reading Centreon API response body", map[string]interface{}{
			"status": resp.StatusCode,
			"error":  err.Error(),
		})
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))

	tflog.Debug(ctx, "Centreon API response", map[string]interface{}{
		"method":   req.Method,
		"path":     req.URL.Path,
		"status":   resp.StatusCode,
		"duration": time.Since(start).String(),
		"bytes":    len(body),
	})

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HandleAPIError(resp, body)
	}

	return resp, nil
}

func (c *Client) GetPlatformInfo(ctx context.Context) (*PlatformInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/platform/installation/status", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var platformInfo PlatformInfo
	if err := json.NewDecoder(resp.Body).Decode(&platformInfo); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &platformInfo, nil
}

// GetHosts lists hosts. The search parameter is a raw JSON object string
// (e.g. `{"name":"web-1"}`) and is URL-encoded as a query value.
func (c *Client) GetHosts(ctx context.Context, limit, page int, search string) (*HostResponse, error) {
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/configuration/hosts?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hostResponse HostResponse
	if err := json.NewDecoder(resp.Body).Decode(&hostResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &hostResponse, nil
}

// CreateHost creates a host and returns the new host ID.
func (c *Client) CreateHost(ctx context.Context, host *CreateHostRequest) (int, error) {
	// The Centreon API serializes configuration writes poorly; keep writes
	// sequential (no artificial delay, the mutex is enough).
	hostMutex.Lock()
	defer hostMutex.Unlock()

	payload, err := json.Marshal(host)
	if err != nil {
		return 0, fmt.Errorf("error marshaling host data: %w", err)
	}

	tflog.Debug(ctx, "Creating host", map[string]interface{}{
		"name": host.Name,
		"url":  c.BaseURL + "/configuration/hosts",
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/configuration/hosts", bytes.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		tflog.Warn(ctx, "Host creation failed", map[string]interface{}{
			"name":  host.Name,
			"error": err.Error(),
		})
		return 0, err
	}
	defer resp.Body.Close()

	// Centreon >= 23.04 returns the created host ID in the response body.
	var created struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err == nil && created.ID != 0 {
		tflog.Debug(ctx, "Host created", map[string]interface{}{
			"name": host.Name,
			"id":   created.ID,
		})
		return created.ID, nil
	}

	// Fallback for API versions that do not return the ID: resolve it by
	// searching for the host name.
	hosts, err := c.GetHosts(ctx, 1, 1, fmt.Sprintf(`{"name":"%s"}`, host.Name))
	if err != nil {
		return 0, fmt.Errorf("host was created but error getting its ID: %w", err)
	}
	if len(hosts.Result) == 0 {
		return 0, fmt.Errorf("host was created but could not retrieve its ID")
	}
	return hosts.Result[0].ID, nil
}

// UpdateHost updates a host by ID.
func (c *Client) UpdateHost(ctx context.Context, hostID int, host *CreateHostRequest) error {
	hostMutex.Lock()
	defer hostMutex.Unlock()

	payload, err := json.Marshal(host)
	if err != nil {
		return fmt.Errorf("error marshaling host data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s/configuration/hosts/%d", c.BaseURL, hostID), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// DeleteHost deletes a host by ID.
func (c *Client) DeleteHost(ctx context.Context, hostID int) error {
	hostMutex.Lock()
	defer hostMutex.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/configuration/hosts/%d", c.BaseURL, hostID), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetMonitoringServers(ctx context.Context, limit, page int, search string) (*MonitoringServersResponse, error) {
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/configuration/monitoring-servers?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response MonitoringServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &response, nil
}

func (c *Client) GetHostGroups(ctx context.Context, limit, page int, search string) (*HostGroupsResponse, error) {
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/monitoring/hostgroups?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response HostGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &response, nil
}

func (c *Client) GetHostTemplates(ctx context.Context, limit, page int, search string) (*HostTemplatesResponse, error) {
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/configuration/hosts/templates?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response HostTemplatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &response, nil
}

// ReloadConfiguration generates and reloads the configuration for all
// monitoring servers.
func (c *Client) ReloadConfiguration(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/configuration/monitoring-servers/generate-and-reload", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetHostMacros retrieves macros for a given host ID.
func (c *Client) GetHostMacros(ctx context.Context, hostID int) ([]HostMacro, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/configuration/hosts/%d/macros", c.BaseURL, hostID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to fetch host macros: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var macroResponse HostMacroResponse
	if err := json.NewDecoder(resp.Body).Decode(&macroResponse); err != nil {
		return nil, fmt.Errorf("error decoding host macros response: %w", err)
	}
	return macroResponse.Result, nil
}
