package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &monitoringServersDataSource{}

func NewMonitoringServersDataSource() datasource.DataSource {
	return &monitoringServersDataSource{}
}

type monitoringServersDataSource struct {
	client *client.Client
}

type monitoringServerDetail struct {
	ID                       types.Int64  `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Address                  types.String `tfsdk:"address"`
	IsLocalhost              types.Bool   `tfsdk:"is_localhost"`
	IsDefault                types.Bool   `tfsdk:"is_default"`
	SSHPort                  types.Int64  `tfsdk:"ssh_port"`
	LastRestart              types.String `tfsdk:"last_restart"`
	EngineStartCommand       types.String `tfsdk:"engine_start_command"`
	EngineStopCommand        types.String `tfsdk:"engine_stop_command"`
	EngineRestartCommand     types.String `tfsdk:"engine_restart_command"`
	EngineReloadCommand      types.String `tfsdk:"engine_reload_command"`
	NagiosBin                types.String `tfsdk:"nagios_bin"`
	NagiostatsBin            types.String `tfsdk:"nagiostats_bin"`
	BrokerReloadCommand      types.String `tfsdk:"broker_reload_command"`
	CentreonBrokerCfgPath    types.String `tfsdk:"centreonbroker_cfg_path"`
	CentreonBrokerModulePath types.String `tfsdk:"centreonbroker_module_path"`
	CentreonBrokerLogsPath   types.String `tfsdk:"centreonbroker_logs_path"`
	CentreonConnectorPath    types.String `tfsdk:"centreonconnector_path"`
	InitScriptCentreontrapd  types.String `tfsdk:"init_script_centreontrapd"`
	SnmpTrapdPathConf        types.String `tfsdk:"snmp_trapd_path_conf"`
	RemoteID                 types.Int64  `tfsdk:"remote_id"`
	RemoteServerUseAsProxy   types.Bool   `tfsdk:"remote_server_use_as_proxy"`
	IsUpdated                types.Bool   `tfsdk:"is_updated"`
	IsActivate               types.Bool   `tfsdk:"is_activate"`
}

type monitoringServersDataSourceModel struct {
	Limit   types.Int64              `tfsdk:"limit"`
	Page    types.Int64              `tfsdk:"page"`
	Search  *searchModel             `tfsdk:"search"`
	Servers []monitoringServerDetail `tfsdk:"servers"`
	Id      types.String             `tfsdk:"id"`
}

func (d *monitoringServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitoring_servers"
}

func (d *monitoringServersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of monitoring servers.",
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "Number of results to return",
				Required:    true,
			},
			"page": schema.Int64Attribute{
				Description: "Page number",
				Required:    true,
			},
			"search": schema.SingleNestedAttribute{
				Description: "Search criteria",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Field name to search",
						Optional:    true,
					},
					"value": schema.StringAttribute{
						Description: "Value to search for",
						Optional:    true,
					},
				},
			},
			"servers": schema.ListNestedAttribute{
				Description: "List of monitoring servers",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Server ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Server name",
							Computed:    true,
						},
						"address": schema.StringAttribute{
							Description: "Server address",
							Computed:    true,
						},
						"is_localhost": schema.BoolAttribute{
							Description: "Whether this is the localhost",
							Computed:    true,
						},
						"is_default": schema.BoolAttribute{
							Description: "Whether this is the default server",
							Computed:    true,
						},
						"ssh_port": schema.Int64Attribute{
							Description: "SSH port",
							Computed:    true,
						},
						"last_restart": schema.StringAttribute{
							Description: "Last restart time",
							Computed:    true,
						},
						"engine_start_command": schema.StringAttribute{
							Description: "Engine start command",
							Computed:    true,
						},
						"engine_stop_command": schema.StringAttribute{
							Description: "Engine stop command",
							Computed:    true,
						},
						"engine_restart_command": schema.StringAttribute{
							Description: "Engine restart command",
							Computed:    true,
						},
						"engine_reload_command": schema.StringAttribute{
							Description: "Engine reload command",
							Computed:    true,
						},
						"nagios_bin": schema.StringAttribute{
							Description: "Nagios binary path",
							Computed:    true,
						},
						"nagiostats_bin": schema.StringAttribute{
							Description: "Nagiostats binary path",
							Computed:    true,
						},
						"broker_reload_command": schema.StringAttribute{
							Description: "Broker reload command",
							Computed:    true,
						},
						"centreonbroker_cfg_path": schema.StringAttribute{
							Description: "Centreon broker config path",
							Computed:    true,
						},
						"centreonbroker_module_path": schema.StringAttribute{
							Description: "Centreon broker module path",
							Computed:    true,
						},
						"centreonbroker_logs_path": schema.StringAttribute{
							Description: "Centreon broker logs path",
							Computed:    true,
						},
						"centreonconnector_path": schema.StringAttribute{
							Description: "Centreon connector path",
							Computed:    true,
						},
						"init_script_centreontrapd": schema.StringAttribute{
							Description: "Init script for centreontrapd",
							Computed:    true,
						},
						"snmp_trapd_path_conf": schema.StringAttribute{
							Description: "SNMP trapd configuration path",
							Computed:    true,
						},
						"remote_id": schema.Int64Attribute{
							Description: "Remote ID",
							Computed:    true,
						},
						"remote_server_use_as_proxy": schema.BoolAttribute{
							Description: "Whether to use as proxy",
							Computed:    true,
						},
						"is_updated": schema.BoolAttribute{
							Description: "Whether the server is updated",
							Computed:    true,
						},
						"is_activate": schema.BoolAttribute{
							Description: "Whether the server is activated",
							Computed:    true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Description: "Placeholder identifier",
				Computed:    true,
			},
		},
	}
}

func (d *monitoringServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *monitoringServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state monitoringServersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize empty search if not provided
	if state.Search == nil {
		state.Search = &searchModel{
			Name:  types.StringNull(),
			Value: types.StringNull(),
		}
	}

	// Create search JSON
	searchQuery := "{}"
	if !state.Search.Name.IsNull() && !state.Search.Value.IsNull() {
		searchQuery = fmt.Sprintf("{\"%s\":\"%s\"}",
			state.Search.Name.ValueString(),
			state.Search.Value.ValueString())
	}

	serversResponse, err := d.client.GetMonitoringServers(
		int(state.Limit.ValueInt64()),
		int(state.Page.ValueInt64()),
		searchQuery,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Monitoring Servers",
			err.Error(),
		)
		return
	}

	// Map response to model
	state.Servers = make([]monitoringServerDetail, len(serversResponse.Result))
	for i, server := range serversResponse.Result {
		// Handle optional string pointer
		var centreonBrokerLogsPath types.String
		if server.CentreonBrokerLogsPath != nil {
			centreonBrokerLogsPath = types.StringValue(*server.CentreonBrokerLogsPath)
		} else {
			centreonBrokerLogsPath = types.StringNull()
		}

		// Handle optional int pointer
		var remoteID types.Int64
		if server.RemoteID != nil {
			remoteID = types.Int64Value(int64(*server.RemoteID))
		} else {
			remoteID = types.Int64Null()
		}

		state.Servers[i] = monitoringServerDetail{
			ID:                       types.Int64Value(int64(server.ID)),
			Name:                     types.StringValue(server.Name),
			Address:                  types.StringValue(server.Address),
			IsLocalhost:              types.BoolValue(server.IsLocalhost),
			IsDefault:                types.BoolValue(server.IsDefault),
			SSHPort:                  types.Int64Value(int64(server.SSHPort)),
			LastRestart:              types.StringValue(server.LastRestart),
			EngineStartCommand:       types.StringValue(server.EngineStartCommand),
			EngineStopCommand:        types.StringValue(server.EngineStopCommand),
			EngineRestartCommand:     types.StringValue(server.EngineRestartCommand),
			EngineReloadCommand:      types.StringValue(server.EngineReloadCommand),
			NagiosBin:                types.StringValue(server.NagiosBin),
			NagiostatsBin:            types.StringValue(server.NagiostatsBin),
			BrokerReloadCommand:      types.StringValue(server.BrokerReloadCommand),
			CentreonBrokerCfgPath:    types.StringValue(server.CentreonBrokerCfgPath),
			CentreonBrokerModulePath: types.StringValue(server.CentreonBrokerModulePath),
			CentreonBrokerLogsPath:   centreonBrokerLogsPath,
			CentreonConnectorPath:    types.StringValue(server.CentreonConnectorPath),
			InitScriptCentreontrapd:  types.StringValue(server.InitScriptCentreontrapd),
			SnmpTrapdPathConf:        types.StringValue(server.SnmpTrapdPathConf),
			RemoteID:                 remoteID,
			RemoteServerUseAsProxy:   types.BoolValue(server.RemoteServerUseAsProxy),
			IsUpdated:                types.BoolValue(server.IsUpdated),
			IsActivate:               types.BoolValue(server.IsActivate),
		}
	}

	state.Id = types.StringValue("monitoring_servers")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
