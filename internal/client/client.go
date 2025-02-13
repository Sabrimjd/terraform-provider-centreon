package client

import (
	"bytes"
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
	ActiveCheckEnabled        int              `json:"active_check_enabled"`  // 0=disabled, 1=enabled
	PassiveCheckEnabled       int              `json:"passive_check_enabled"` // 0=disabled, 1=enabled
	NotificationEnabled       int              `json:"notification_enabled"`  // 0=disabled, 1=enabled
	NotificationOptions       int              `json:"notification_options"`
	NotificationInterval      int              `json:"notification_interval"`
	NotificationTimeperiodID  int              `json:"notification_timeperiod_id"`
	AddInheritedContactGroup  bool             `json:"add_inherited_contact_group"`
	AddInheritedContact       bool             `json:"add_inherited_contact"`
	FirstNotificationDelay    int              `json:"first_notification_delay"`
	RecoveryNotificationDelay int              `json:"recovery_notification_delay"`
	AcknowledgementTimeout    int              `json:"acknowledgement_timeout"`
	FreshnessChecked          int              `json:"freshness_checked"` // 0=disabled, 1=enabled
	FreshnessThreshold        int              `json:"freshness_threshold"`
	FlapDetectionEnabled      int              `json:"flap_detection_enabled"` // 0=disabled, 1=enabled
	LowFlapThreshold          int              `json:"low_flap_threshold"`
	HighFlapThreshold         int              `json:"high_flap_threshold"`
	EventHandlerEnabled       int              `json:"event_handler_enabled"` // 0=disabled, 1=enabled
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

func (c *Client) CreateHost(host *CreateHostRequest) error {
	url := fmt.Sprintf("%s/configuration/hosts", c.BaseURL)

	jsonData, err := json.Marshal(host)
	if err != nil {
		return fmt.Errorf("error marshaling host data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AUTH-TOKEN", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("API request failed: %v", errorResponse)
	}

	return nil
}

func (c *Client) UpdateHost(host *CreateHostRequest) error {
	// First, get the host ID using the name
	hosts, err := c.GetHosts(1, 1, fmt.Sprintf("{\"name\":\"%s\"}", host.Name))
	if err != nil {
		return fmt.Errorf("error getting host ID: %v", err)
	}
	if len(hosts.Result) == 0 {
		return fmt.Errorf("host not found: %s", host.Name)
	}

	hostID := hosts.Result[0].ID
	url := fmt.Sprintf("%s/configuration/hosts/%d", c.BaseURL, hostID)

	jsonData, err := json.Marshal(host)
	if err != nil {
		return fmt.Errorf("error marshaling host data: %v", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AUTH-TOKEN", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Accept both 200 OK and 204 No Content as success responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("API request failed: %v", errorResponse)
	}

	return nil
}

func (c *Client) DeleteHost(name string) error {
	// First, get the host ID using the name
	hosts, err := c.GetHosts(1, 1, fmt.Sprintf("{\"name\":\"%s\"}", name))
	if err != nil {
		return fmt.Errorf("error getting host ID: %v", err)
	}
	if len(hosts.Result) == 0 {
		return fmt.Errorf("host not found: %s", name)
	}

	hostID := hosts.Result[0].ID
	url := fmt.Sprintf("%s/configuration/hosts/%d", c.BaseURL, hostID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-AUTH-TOKEN", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Accept both 200 OK and 204 No Content as success responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("API request failed: %v", errorResponse)
	}

	return nil
}

func (c *Client) GetMonitoringServers(limit int, page int, search string) (*MonitoringServersResponse, error) {
	url := fmt.Sprintf("%s/configuration/monitoring-servers?limit=%d&page=%d&search=%s",
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

	var response MonitoringServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetHostGroups(limit int, page int, search string) (*HostGroupsResponse, error) {
	url := fmt.Sprintf("%s/monitoring/hostgroups?limit=%d&page=%d&search=%s",
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

	var response HostGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetHostTemplates(limit int, page int, search string) (*HostTemplatesResponse, error) {
	url := fmt.Sprintf("%s/configuration/hosts/templates?limit=%d&page=%d&search=%s",
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

	var response HostTemplatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
