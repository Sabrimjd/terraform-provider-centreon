package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"

	"terraform-provider-centreon/internal/logging"
	"terraform-provider-centreon/internal/validation"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &hostResource{}

func NewHostResource() resource.Resource {
	return &hostResource{}
}

type hostResource struct {
	client *client.Client
}

type hostResourceModel struct {
	MonitoringServerID        types.Int64    `tfsdk:"monitoring_server_id"`
	Name                      types.String   `tfsdk:"name"`
	Address                   types.String   `tfsdk:"address"`
	Alias                     types.String   `tfsdk:"alias"`
	SNMPCommunity             types.String   `tfsdk:"snmp_community"`
	SNMPVersion               types.String   `tfsdk:"snmp_version"`
	TimezoneID                types.Int64    `tfsdk:"timezone_id"`
	SeverityID                types.Int64    `tfsdk:"severity_id"`
	CheckCommandID            types.Int64    `tfsdk:"check_command_id"`
	CheckCommandArgs          []types.String `tfsdk:"check_command_args"`
	CheckTimeperiodID         types.Int64    `tfsdk:"check_timeperiod_id"`
	MaxCheckAttempts          types.Int64    `tfsdk:"max_check_attempts"`
	NormalCheckInterval       types.Int64    `tfsdk:"normal_check_interval"`
	RetryCheckInterval        types.Int64    `tfsdk:"retry_check_interval"`
	ActiveCheckEnabled        types.Int64    `tfsdk:"active_check_enabled"`
	PassiveCheckEnabled       types.Int64    `tfsdk:"passive_check_enabled"`
	NotificationEnabled       types.Int64    `tfsdk:"notification_enabled"`
	NotificationOptions       types.Int64    `tfsdk:"notification_options"`
	NotificationInterval      types.Int64    `tfsdk:"notification_interval"`
	NotificationTimeperiodID  types.Int64    `tfsdk:"notification_timeperiod_id"`
	FirstNotificationDelay    types.Int64    `tfsdk:"first_notification_delay"`
	RecoveryNotificationDelay types.Int64    `tfsdk:"recovery_notification_delay"`
	AcknowledgementTimeout    types.Int64    `tfsdk:"acknowledgement_timeout"`
	FreshnessChecked          types.Int64    `tfsdk:"freshness_checked"`
	FreshnessThreshold        types.Int64    `tfsdk:"freshness_threshold"`
	FlapDetectionEnabled      types.Int64    `tfsdk:"flap_detection_enabled"`
	LowFlapThreshold          types.Int64    `tfsdk:"low_flap_threshold"`
	HighFlapThreshold         types.Int64    `tfsdk:"high_flap_threshold"`
	EventHandlerEnabled       types.Int64    `tfsdk:"event_handler_enabled"`
	EventHandlerCommandID     types.Int64    `tfsdk:"event_handler_command_id"`
	EventHandlerCommandArgs   []types.String `tfsdk:"event_handler_command_args"`
	NoteURL                   types.String   `tfsdk:"note_url"`
	Note                      types.String   `tfsdk:"note"`
	ActionURL                 types.String   `tfsdk:"action_url"`
	IconID                    types.Int64    `tfsdk:"icon_id"`
	IconAlternative           types.String   `tfsdk:"icon_alternative"`
	Comment                   types.String   `tfsdk:"comment"`
	IsActivated               types.Bool     `tfsdk:"is_activated"`
	Categories                []types.Int64  `tfsdk:"categories"`
	Groups                    []types.Int64  `tfsdk:"groups"`
	Templates                 []types.Int64  `tfsdk:"templates"`
	Macros                    []macroModel   `tfsdk:"macros"`
	GeoCoords                 types.String   `tfsdk:"geo_coords"`
}

type macroModel struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	IsPassword  types.Bool   `tfsdk:"is_password"`
	Description types.String `tfsdk:"description"`
}

func (r *hostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *hostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Centreon host.",
		Attributes: map[string]schema.Attribute{
			"monitoring_server_id": schema.Int64Attribute{
				Required:    true,
				Description: "ID of the host's monitoring server",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Host name",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "IP or domain of the host",
				Validators: []validator.String{
					validation.HostnameOrIPValidator{},
				},
			},
			"alias": schema.StringAttribute{
				Optional:    true,
				Description: "Host alias",
			},
			"snmp_community": schema.StringAttribute{
				Optional:    true,
				Description: "Community of the SNMP agent",
				Sensitive:   true,
			},
			"snmp_version": schema.StringAttribute{
				Optional:    true,
				Description: "Version of the SNMP agent (1, 2c, or 3)",
				Validators: []validator.String{
					validation.SNMPVersionValidator{},
				},
			},
			"geo_coords": schema.StringAttribute{
				Optional:    true,
				Description: "Geographic coordinates of the host (format: latitude,longitude)",
				Validators: []validator.String{
					validation.GeoCoordsValidator{},
				},
			},
			"notification_options": schema.Int64Attribute{
				Optional:    true,
				Description: "Notification options (sum of: 1=DOWN, 2=UNREACHABLE, 4=RECOVERY, 8=FLAPPING, 16=DOWNTIME_SCHEDULED)",
				Validators: []validator.Int64{
					validation.NotificationOptionsValidator{},
				},
			},
			"timezone_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Timezone ID",
			},
			"severity_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Severity ID",
			},
			"check_command_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Check command ID",
			},
			"check_command_args": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Check command arguments",
			},
			"check_timeperiod_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Check timeperiod ID",
			},
			"max_check_attempts": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of retry attempts for host checks",
			},
			"normal_check_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between normal checks",
			},
			"retry_check_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between retry checks",
			},
			"active_check_enabled": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether active checks are enabled (0=disabled, 1=enabled)",
				Default:     int64default.StaticInt64(0),
			},
			"passive_check_enabled": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether passive checks are enabled (0=disabled, 1=enabled)",
				Default:     int64default.StaticInt64(0),
			},
			"notification_enabled": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether notifications are enabled (0=disabled, 1=enabled)",
				Default:     int64default.StaticInt64(0),
			},
			"notification_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between notifications",
			},
			"notification_timeperiod_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Notification timeperiod ID",
			},
			"first_notification_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Delay before first notification",
			},
			"recovery_notification_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Delay before recovery notification",
			},
			"acknowledgement_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Acknowledgement timeout",
			},
			"freshness_checked": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether freshness is checked (0=disabled, 1=enabled)",
				Default:     int64default.StaticInt64(0),
			},
			"freshness_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "Freshness threshold in seconds",
			},
			"flap_detection_enabled": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether flap detection is enabled (0=disabled, 1=enabled)",
				Default:     int64default.StaticInt64(0),
			},
			"low_flap_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "Low flap threshold",
			},
			"high_flap_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "High flap threshold",
			},
			"event_handler_enabled": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether event handler is enabled (0=disabled, 1=enabled)",
				Default:     int64default.StaticInt64(0),
			},
			"event_handler_command_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Event handler command ID",
			},
			"event_handler_command_args": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Event handler command arguments",
			},
			"note_url": schema.StringAttribute{
				Optional:    true,
				Description: "URL with additional host information",
			},
			"note": schema.StringAttribute{
				Optional:    true,
				Description: "Additional notes about the host",
			},
			"action_url": schema.StringAttribute{
				Optional:    true,
				Description: "URL for additional host actions",
			},
			"icon_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Icon ID",
			},
			"icon_alternative": schema.StringAttribute{
				Optional:    true,
				Description: "Alternative text for icon",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comments about the host",
			},
			"categories": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of category IDs",
			},
			"groups": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of group IDs",
			},
			"templates": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of template IDs",
			},
			"macros": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Host macros",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Macro name",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "Macro value",
						},
						"is_password": schema.BoolAttribute{
							Required:    true,
							Description: "Whether the macro value is a password",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Macro description",
						},
					},
				},
			},
			"is_activated": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the host is activated",
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *hostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Helper function to handle configuration reload if enabled.
func (r *hostResource) handleConfigurationReload() error {
	if r.client.GenerateAndReloadConfiguration {
		if err := r.client.ReloadConfiguration(); err != nil {
			return fmt.Errorf("failed to generate and reload configuration: %v", err)
		}
	}
	return nil
}

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan model to a CreateHostRequest
	createReq := &client.CreateHostRequest{
		MonitoringServerID: int(plan.MonitoringServerID.ValueInt64()),
		Name:               plan.Name.ValueString(),
		Address:            plan.Address.ValueString(),
	}

	// Set optional fields
	if !plan.Alias.IsNull() {
		v := plan.Alias.ValueString()
		createReq.Alias = &v
	}
	if !plan.SNMPCommunity.IsNull() {
		v := plan.SNMPCommunity.ValueString()
		createReq.SNMPCommunity = &v
	}
	if !plan.SNMPVersion.IsNull() {
		v := plan.SNMPVersion.ValueString()
		createReq.SNMPVersion = &v
	}
	if !plan.TimezoneID.IsNull() {
		v := int(plan.TimezoneID.ValueInt64())
		createReq.TimezoneID = &v
	}
	if !plan.SeverityID.IsNull() {
		v := int(plan.SeverityID.ValueInt64())
		createReq.SeverityID = &v
	}
	if !plan.CheckCommandID.IsNull() {
		v := int(plan.CheckCommandID.ValueInt64())
		createReq.CheckCommandID = &v
	}
	if len(plan.CheckCommandArgs) > 0 {
		args := make([]string, len(plan.CheckCommandArgs))
		for i, arg := range plan.CheckCommandArgs {
			args[i] = arg.ValueString()
		}
		createReq.CheckCommandArgs = args
	}
	if !plan.CheckTimeperiodID.IsNull() {
		v := int(plan.CheckTimeperiodID.ValueInt64())
		createReq.CheckTimeperiodID = &v
	}
	if !plan.MaxCheckAttempts.IsNull() {
		v := int(plan.MaxCheckAttempts.ValueInt64())
		createReq.MaxCheckAttempts = &v
	}
	if !plan.NormalCheckInterval.IsNull() {
		v := int(plan.NormalCheckInterval.ValueInt64())
		createReq.NormalCheckInterval = &v
	}
	if !plan.RetryCheckInterval.IsNull() {
		v := int(plan.RetryCheckInterval.ValueInt64())
		createReq.RetryCheckInterval = &v
	}
	if !plan.ActiveCheckEnabled.IsNull() && plan.ActiveCheckEnabled.ValueInt64() != 0 {
		v := int(plan.ActiveCheckEnabled.ValueInt64())
		createReq.ActiveCheckEnabled = &v
	}
	if !plan.PassiveCheckEnabled.IsNull() && plan.PassiveCheckEnabled.ValueInt64() != 0 {
		v := int(plan.PassiveCheckEnabled.ValueInt64())
		createReq.PassiveCheckEnabled = &v
	}
	if !plan.NotificationEnabled.IsNull() && plan.NotificationEnabled.ValueInt64() != 0 {
		v := int(plan.NotificationEnabled.ValueInt64())
		createReq.NotificationEnabled = &v
	}
	if !plan.NotificationOptions.IsNull() {
		v := int(plan.NotificationOptions.ValueInt64())
		createReq.NotificationOptions = &v
	}
	if !plan.NotificationInterval.IsNull() {
		v := int(plan.NotificationInterval.ValueInt64())
		createReq.NotificationInterval = &v
	}
	if !plan.NotificationTimeperiodID.IsNull() {
		v := int(plan.NotificationTimeperiodID.ValueInt64())
		createReq.NotificationTimeperiodID = &v
	}
	if !plan.FirstNotificationDelay.IsNull() {
		v := int(plan.FirstNotificationDelay.ValueInt64())
		createReq.FirstNotificationDelay = &v
	}
	if !plan.RecoveryNotificationDelay.IsNull() {
		v := int(plan.RecoveryNotificationDelay.ValueInt64())
		createReq.RecoveryNotificationDelay = &v
	}
	if !plan.AcknowledgementTimeout.IsNull() {
		v := int(plan.AcknowledgementTimeout.ValueInt64())
		createReq.AcknowledgementTimeout = &v
	}
	if !plan.FreshnessChecked.IsNull() && plan.FreshnessChecked.ValueInt64() != 0 {
		v := int(plan.FreshnessChecked.ValueInt64())
		createReq.FreshnessChecked = &v
	}
	if !plan.FreshnessThreshold.IsNull() {
		v := int(plan.FreshnessThreshold.ValueInt64())
		createReq.FreshnessThreshold = &v
	}
	if !plan.FlapDetectionEnabled.IsNull() && plan.FlapDetectionEnabled.ValueInt64() != 0 {
		v := int(plan.FlapDetectionEnabled.ValueInt64())
		createReq.FlapDetectionEnabled = &v
	}
	if !plan.LowFlapThreshold.IsNull() {
		v := int(plan.LowFlapThreshold.ValueInt64())
		createReq.LowFlapThreshold = &v
	}
	if !plan.HighFlapThreshold.IsNull() {
		v := int(plan.HighFlapThreshold.ValueInt64())
		createReq.HighFlapThreshold = &v
	}
	if !plan.EventHandlerEnabled.IsNull() && plan.EventHandlerEnabled.ValueInt64() != 0 {
		v := int(plan.EventHandlerEnabled.ValueInt64())
		createReq.EventHandlerEnabled = &v
	}
	if !plan.EventHandlerCommandID.IsNull() {
		v := int(plan.EventHandlerCommandID.ValueInt64())
		createReq.EventHandlerCommandID = &v
	}
	if len(plan.EventHandlerCommandArgs) > 0 {
		args := make([]string, len(plan.EventHandlerCommandArgs))
		for i, arg := range plan.EventHandlerCommandArgs {
			args[i] = arg.ValueString()
		}
		createReq.EventHandlerCommandArgs = args
	}
	if !plan.NoteURL.IsNull() {
		v := plan.NoteURL.ValueString()
		createReq.NoteURL = &v
	}
	if !plan.Note.IsNull() {
		v := plan.Note.ValueString()
		createReq.Note = &v
	}
	if !plan.ActionURL.IsNull() {
		v := plan.ActionURL.ValueString()
		createReq.ActionURL = &v
	}
	if !plan.IconID.IsNull() {
		v := int(plan.IconID.ValueInt64())
		createReq.IconID = &v
	}
	if !plan.IconAlternative.IsNull() {
		v := plan.IconAlternative.ValueString()
		createReq.IconAlternative = &v
	}
	if !plan.Comment.IsNull() {
		v := plan.Comment.ValueString()
		createReq.Comment = &v
	}
	if !plan.IsActivated.IsNull() {
		v := plan.IsActivated.ValueBool()
		createReq.IsActivated = &v
	}
	if !plan.GeoCoords.IsNull() {
		v := plan.GeoCoords.ValueString()
		createReq.GeoCoords = &v
	}

	// Convert slice fields
	if len(plan.Categories) > 0 {
		categories := make([]int, len(plan.Categories))
		for i, cat := range plan.Categories {
			categories[i] = int(cat.ValueInt64())
		}
		createReq.Categories = categories
	}
	if len(plan.Groups) > 0 {
		groups := make([]int, len(plan.Groups))
		for i, grp := range plan.Groups {
			groups[i] = int(grp.ValueInt64())
		}
		createReq.Groups = groups
	}
	if len(plan.Templates) > 0 {
		templates := make([]int, len(plan.Templates))
		for i, tpl := range plan.Templates {
			templates[i] = int(tpl.ValueInt64())
		}
		createReq.Templates = templates
	}

	// Convert macros
	if len(plan.Macros) > 0 {
		createReq.Macros = make([]client.HostMacro, len(plan.Macros))
		for i, m := range plan.Macros {
			macro := client.HostMacro{
				Name:       m.Name.ValueString(),
				IsPassword: m.IsPassword.ValueBool(),
			}
			if !m.Value.IsNull() {
				v := m.Value.ValueString()
				macro.Value = &v
			}
			if !m.Description.IsNull() {
				v := m.Description.ValueString()
				macro.Description = &v
			}
			createReq.Macros[i] = macro
		}
	}

	// Create the host
	err := r.client.CreateHost(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating host",
			fmt.Sprintf("Could not create host: %v", err),
		)
		return
	}

	// Generate and reload configuration if enabled
	if err := r.handleConfigurationReload(); err != nil {
		resp.Diagnostics.AddError(
			"Error after creating host",
			err.Error(),
		)
		return
	}

	// Save the plan
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get host details from API
	hosts, err := r.client.GetHosts(1, 1, fmt.Sprintf("{\"name\":\"%s\"}", state.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host",
			fmt.Sprintf("Could not read host %s: %v", state.Name.ValueString(), err),
		)
		return
	}

	if len(hosts.Result) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	host := hosts.Result[0]

	// Update state with values from API, only if they differ from defaults
	state.Name = types.StringValue(host.Name)
	state.Address = types.StringValue(host.Address)
	state.Alias = types.StringValue(host.Alias)
	state.MonitoringServerID = types.Int64Value(int64(host.MonitoringServer.ID))

	// Only set if not empty/default
	if host.SNMPCommunity != "" {
		state.SNMPCommunity = types.StringValue(host.SNMPCommunity)
	}
	if host.SNMPVersion != "" {
		state.SNMPVersion = types.StringValue(host.SNMPVersion)
	}
	if host.TimezoneID != 0 {
		state.TimezoneID = types.Int64Value(int64(host.TimezoneID))
	}
	if host.SeverityID != 0 {
		state.SeverityID = types.Int64Value(int64(host.SeverityID))
	}
	if host.CheckCommandID != 0 {
		state.CheckCommandID = types.Int64Value(int64(host.CheckCommandID))
	}
	if host.CheckTimeperiodID != 0 {
		state.CheckTimeperiodID = types.Int64Value(int64(host.CheckTimeperiodID))
	}
	if host.MaxCheckAttempts != 0 {
		state.MaxCheckAttempts = types.Int64Value(int64(host.MaxCheckAttempts))
	}
	if host.NormalCheckInterval != 0 {
		state.NormalCheckInterval = types.Int64Value(int64(host.NormalCheckInterval))
	}
	if host.RetryCheckInterval != 0 {
		state.RetryCheckInterval = types.Int64Value(int64(host.RetryCheckInterval))
	}
	if host.NotificationOptions != 0 {
		state.NotificationOptions = types.Int64Value(int64(host.NotificationOptions))
	}
	if host.NotificationInterval != 0 {
		state.NotificationInterval = types.Int64Value(int64(host.NotificationInterval))
	}
	if host.NotificationTimeperiodID != 0 {
		state.NotificationTimeperiodID = types.Int64Value(int64(host.NotificationTimeperiodID))
	}

	if host.FirstNotificationDelay != 0 {
		state.FirstNotificationDelay = types.Int64Value(int64(host.FirstNotificationDelay))
	}
	if host.RecoveryNotificationDelay != 0 {
		state.RecoveryNotificationDelay = types.Int64Value(int64(host.RecoveryNotificationDelay))
	}
	if host.AcknowledgementTimeout != 0 {
		state.AcknowledgementTimeout = types.Int64Value(int64(host.AcknowledgementTimeout))
	}
	if host.FreshnessThreshold != 0 {
		state.FreshnessThreshold = types.Int64Value(int64(host.FreshnessThreshold))
	}
	if host.LowFlapThreshold != 0 {
		state.LowFlapThreshold = types.Int64Value(int64(host.LowFlapThreshold))
	}
	if host.HighFlapThreshold != 0 {
		state.HighFlapThreshold = types.Int64Value(int64(host.HighFlapThreshold))
	}
	if host.EventHandlerCommandID != 0 {
		state.EventHandlerCommandID = types.Int64Value(int64(host.EventHandlerCommandID))
	}
	if host.IconID != 0 {
		state.IconID = types.Int64Value(int64(host.IconID))
	}

	// Only set string fields if not empty
	if host.NoteURL != "" {
		state.NoteURL = types.StringValue(host.NoteURL)
	}
	if host.Note != "" {
		state.Note = types.StringValue(host.Note)
	}
	if host.ActionURL != "" {
		state.ActionURL = types.StringValue(host.ActionURL)
	}
	if host.IconAlternative != "" {
		state.IconAlternative = types.StringValue(host.IconAlternative)
	}
	if host.Comment != "" {
		state.Comment = types.StringValue(host.Comment)
	}
	if host.GeoCoords != "" {
		state.GeoCoords = types.StringValue(host.GeoCoords)
	}

	// Handle special fields with default value of 2
	if host.ActiveCheckEnabled != 2 {
		state.ActiveCheckEnabled = types.Int64Value(int64(host.ActiveCheckEnabled))
	}
	if host.PassiveCheckEnabled != 2 {
		state.PassiveCheckEnabled = types.Int64Value(int64(host.PassiveCheckEnabled))
	}
	if host.NotificationEnabled != 2 {
		state.NotificationEnabled = types.Int64Value(int64(host.NotificationEnabled))
	}
	if host.EventHandlerEnabled != 2 {
		state.EventHandlerEnabled = types.Int64Value(int64(host.EventHandlerEnabled))
	}
	if host.FlapDetectionEnabled != 2 {
		state.FlapDetectionEnabled = types.Int64Value(int64(host.FlapDetectionEnabled))
	}
	if host.FreshnessChecked != 2 {
		state.FreshnessChecked = types.Int64Value(int64(host.FreshnessChecked))
	}

	// Handle enabled/checked fields with default value of 0
	state.ActiveCheckEnabled = types.Int64Value(int64(host.ActiveCheckEnabled))
	state.PassiveCheckEnabled = types.Int64Value(int64(host.PassiveCheckEnabled))
	state.NotificationEnabled = types.Int64Value(int64(host.NotificationEnabled))
	state.EventHandlerEnabled = types.Int64Value(int64(host.EventHandlerEnabled))
	state.FlapDetectionEnabled = types.Int64Value(int64(host.FlapDetectionEnabled))
	state.FreshnessChecked = types.Int64Value(int64(host.FreshnessChecked))

	// Only set arrays if not empty
	if len(host.CheckCommandArgs) > 0 {
		state.CheckCommandArgs = make([]types.String, len(host.CheckCommandArgs))
		for i, arg := range host.CheckCommandArgs {
			state.CheckCommandArgs[i] = types.StringValue(arg)
		}
	}

	if len(host.EventHandlerCommandArgs) > 0 {
		state.EventHandlerCommandArgs = make([]types.String, len(host.EventHandlerCommandArgs))
		for i, arg := range host.EventHandlerCommandArgs {
			state.EventHandlerCommandArgs[i] = types.StringValue(arg)
		}
	}

	if len(host.Categories) > 0 {
		state.Categories = make([]types.Int64, len(host.Categories))
		for i, cat := range host.Categories {
			state.Categories[i] = types.Int64Value(int64(cat))
		}
	}

	// Always set these fields as they are required
	state.Groups = make([]types.Int64, len(host.Groups))
	for i, group := range host.Groups {
		state.Groups[i] = types.Int64Value(int64(group.ID))
	}

	state.Templates = make([]types.Int64, len(host.Templates))
	for i, tmpl := range host.Templates {
		state.Templates[i] = types.Int64Value(int64(tmpl.ID))
	}

	// Get macros for the host
	macros, err := r.client.GetHostMacros(host.ID)
	if err != nil {
		logging.Warn(ctx, "Error fetching host macros", map[string]interface{}{
			"host_id": host.ID,
			"error":   err.Error(),
		})
	} else if len(macros) > 0 {
		state.Macros = make([]macroModel, len(macros))
		for i, m := range macros {
			mac := macroModel{
				Name:       types.StringValue(m.Name),
				IsPassword: types.BoolValue(m.IsPassword),
			}

			// According to API docs, if is_password is true and value is null,
			// the value is considered unchanged. However, for our purposes,
			// we should retrieve the value if possible
			if m.Value != nil {
				mac.Value = types.StringValue(*m.Value)
			}

			if m.Description != nil {
				mac.Description = types.StringValue(*m.Description)
			}

			state.Macros[i] = mac
		}

		logging.Info(ctx, "Host macros retrieved", map[string]interface{}{
			"host":   host.Name,
			"count":  len(macros),
			"macros": macros,
		})
	}

	state.IsActivated = types.BoolValue(host.IsActivated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan hostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to compare
	var state hostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logging.Info(ctx, "Updating host", map[string]interface{}{
		"host": plan.Name.ValueString(),
	})

	// Create update request using the same structure as create
	updateReq := &client.CreateHostRequest{
		MonitoringServerID: int(plan.MonitoringServerID.ValueInt64()),
		Name:               plan.Name.ValueString(),
		Address:            plan.Address.ValueString(),
	}

	// Only include fields that are actually set in the plan
	if !plan.Alias.IsNull() {
		v := plan.Alias.ValueString()
		updateReq.Alias = &v
	}
	if !plan.SNMPCommunity.IsNull() {
		v := plan.SNMPCommunity.ValueString()
		updateReq.SNMPCommunity = &v
	}
	if !plan.SNMPVersion.IsNull() {
		v := plan.SNMPVersion.ValueString()
		updateReq.SNMPVersion = &v
	}
	if !plan.TimezoneID.IsNull() {
		v := int(plan.TimezoneID.ValueInt64())
		updateReq.TimezoneID = &v
	}
	if !plan.SeverityID.IsNull() {
		v := int(plan.SeverityID.ValueInt64())
		updateReq.SeverityID = &v
	}
	if !plan.CheckCommandID.IsNull() {
		v := int(plan.CheckCommandID.ValueInt64())
		updateReq.CheckCommandID = &v
	}
	if !plan.CheckTimeperiodID.IsNull() {
		v := int(plan.CheckTimeperiodID.ValueInt64())
		updateReq.CheckTimeperiodID = &v
	}
	if !plan.MaxCheckAttempts.IsNull() {
		v := int(plan.MaxCheckAttempts.ValueInt64())
		updateReq.MaxCheckAttempts = &v
	}
	if !plan.NormalCheckInterval.IsNull() {
		v := int(plan.NormalCheckInterval.ValueInt64())
		updateReq.NormalCheckInterval = &v
	}
	if !plan.RetryCheckInterval.IsNull() {
		v := int(plan.RetryCheckInterval.ValueInt64())
		updateReq.RetryCheckInterval = &v
	}
	if !plan.ActiveCheckEnabled.IsNull() {
		v := int(plan.ActiveCheckEnabled.ValueInt64())
		updateReq.ActiveCheckEnabled = &v
	}
	if !plan.PassiveCheckEnabled.IsNull() {
		v := int(plan.PassiveCheckEnabled.ValueInt64())
		updateReq.PassiveCheckEnabled = &v
	}
	if !plan.NotificationEnabled.IsNull() {
		v := int(plan.NotificationEnabled.ValueInt64())
		updateReq.NotificationEnabled = &v
	}
	if !plan.NotificationOptions.IsNull() {
		v := int(plan.NotificationOptions.ValueInt64())
		updateReq.NotificationOptions = &v
	}
	if !plan.NotificationInterval.IsNull() {
		v := int(plan.NotificationInterval.ValueInt64())
		updateReq.NotificationInterval = &v
	}
	if !plan.NotificationTimeperiodID.IsNull() {
		v := int(plan.NotificationTimeperiodID.ValueInt64())
		updateReq.NotificationTimeperiodID = &v
	}
	if !plan.FirstNotificationDelay.IsNull() {
		v := int(plan.FirstNotificationDelay.ValueInt64())
		updateReq.FirstNotificationDelay = &v
	}
	if !plan.RecoveryNotificationDelay.IsNull() {
		v := int(plan.RecoveryNotificationDelay.ValueInt64())
		updateReq.RecoveryNotificationDelay = &v
	}
	if !plan.AcknowledgementTimeout.IsNull() {
		v := int(plan.AcknowledgementTimeout.ValueInt64())
		updateReq.AcknowledgementTimeout = &v
	}
	if !plan.FreshnessChecked.IsNull() {
		v := int(plan.FreshnessChecked.ValueInt64())
		updateReq.FreshnessChecked = &v
	}
	if !plan.FreshnessThreshold.IsNull() {
		v := int(plan.FreshnessThreshold.ValueInt64())
		updateReq.FreshnessThreshold = &v
	}
	if !plan.FlapDetectionEnabled.IsNull() {
		v := int(plan.FlapDetectionEnabled.ValueInt64())
		updateReq.FlapDetectionEnabled = &v
	}
	if !plan.LowFlapThreshold.IsNull() {
		v := int(plan.LowFlapThreshold.ValueInt64())
		updateReq.LowFlapThreshold = &v
	}
	if !plan.HighFlapThreshold.IsNull() {
		v := int(plan.HighFlapThreshold.ValueInt64())
		updateReq.HighFlapThreshold = &v
	}
	if !plan.EventHandlerEnabled.IsNull() {
		v := int(plan.EventHandlerEnabled.ValueInt64())
		updateReq.EventHandlerEnabled = &v
	}
	if !plan.EventHandlerCommandID.IsNull() {
		v := int(plan.EventHandlerCommandID.ValueInt64())
		updateReq.EventHandlerCommandID = &v
	}

	// Only include non-empty arrays
	if len(plan.CheckCommandArgs) > 0 {
		args := make([]string, len(plan.CheckCommandArgs))
		for i, arg := range plan.CheckCommandArgs {
			args[i] = arg.ValueString()
		}
		updateReq.CheckCommandArgs = args
	}
	if len(plan.EventHandlerCommandArgs) > 0 {
		args := make([]string, len(plan.EventHandlerCommandArgs))
		for i, arg := range plan.EventHandlerCommandArgs {
			args[i] = arg.ValueString()
		}
		updateReq.EventHandlerCommandArgs = args
	}
	if len(plan.Categories) > 0 {
		categories := make([]int, len(plan.Categories))
		for i, cat := range plan.Categories {
			categories[i] = int(cat.ValueInt64())
		}
		updateReq.Categories = categories
	}

	// Always include groups and templates as they are required
	if len(plan.Groups) > 0 {
		groups := make([]int, len(plan.Groups))
		for i, grp := range plan.Groups {
			groups[i] = int(grp.ValueInt64())
		}
		updateReq.Groups = groups
	}
	if len(plan.Templates) > 0 {
		templates := make([]int, len(plan.Templates))
		for i, tpl := range plan.Templates {
			templates[i] = int(tpl.ValueInt64())
		}
		updateReq.Templates = templates
	}

	// Only include non-empty string fields
	if !plan.NoteURL.IsNull() {
		v := plan.NoteURL.ValueString()
		updateReq.NoteURL = &v
	}
	if !plan.Note.IsNull() {
		v := plan.Note.ValueString()
		updateReq.Note = &v
	}
	if !plan.ActionURL.IsNull() {
		v := plan.ActionURL.ValueString()
		updateReq.ActionURL = &v
	}
	if !plan.IconAlternative.IsNull() {
		v := plan.IconAlternative.ValueString()
		updateReq.IconAlternative = &v
	}
	if !plan.Comment.IsNull() {
		v := plan.Comment.ValueString()
		updateReq.Comment = &v
	}
	if !plan.GeoCoords.IsNull() {
		v := plan.GeoCoords.ValueString()
		updateReq.GeoCoords = &v
	}
	if !plan.IsActivated.IsNull() {
		v := plan.IsActivated.ValueBool()
		updateReq.IsActivated = &v
	}

	// Handle macros - always include them in the update to ensure
	// they are properly updated per the OpenAPI documentation
	if len(plan.Macros) > 0 {
		logging.Info(ctx, "Including macros in update", map[string]interface{}{
			"host":   plan.Name.ValueString(),
			"macros": len(plan.Macros),
		})

		updateReq.Macros = make([]client.HostMacro, len(plan.Macros))
		for i, m := range plan.Macros {
			macro := client.HostMacro{
				Name:       m.Name.ValueString(),
				IsPassword: m.IsPassword.ValueBool(),
			}

			// Always include the value when updating
			if !m.Value.IsNull() {
				v := m.Value.ValueString()
				macro.Value = &v
			}

			if !m.Description.IsNull() {
				v := m.Description.ValueString()
				macro.Description = &v
			}

			updateReq.Macros[i] = macro
		}
	}

	// Call API to update host
	if err := r.client.UpdateHost(updateReq); err != nil {
		resp.Diagnostics.AddError(
			"Error updating host",
			fmt.Sprintf("Could not update host %s: %v", plan.Name.ValueString(), err),
		)
		return
	}

	// Generate and reload configuration if enabled
	if err := r.handleConfigurationReload(); err != nil {
		resp.Diagnostics.AddError(
			"Error after updating host",
			err.Error(),
		)
		return
	}

	// Update state with plan
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the host using the client
	if err := r.client.DeleteHost(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting host",
			fmt.Sprintf("Could not delete host %s: %v", state.Name.ValueString(), err),
		)
		return
	}

	// Generate and reload configuration if enabled
	if err := r.handleConfigurationReload(); err != nil {
		resp.Diagnostics.AddError(
			"Error after deleting host",
			err.Error(),
		)
		return
	}
}
