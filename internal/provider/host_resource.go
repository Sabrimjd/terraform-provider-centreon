package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"
	"terraform-provider-centreon/internal/validation"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &hostResource{}

func NewHostResource() resource.Resource {
	return &hostResource{}
}

type hostResource struct {
	client *client.Client
}

type hostResourceModel struct {
	ID                        types.Int64    `tfsdk:"id"` // Added ID field to store the host ID
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
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Host ID (internal identifier)",
			},
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

// handleConfigurationReload generates and reloads the configuration when
// enabled. A reload failure is surfaced as a WARNING diagnostic rather than
// an error: the host itself was already created/updated/deleted, so failing
// the apply here would leave Terraform state diverged from Centreon.
func (r *hostResource) handleConfigurationReload(ctx context.Context, diags *diag.Diagnostics) {
	if !r.client.GenerateAndReloadConfiguration {
		return
	}
	if err := r.client.ReloadConfiguration(ctx); err != nil {
		diags.AddWarning(
			"Configuration reload failed",
			fmt.Sprintf("The host change was applied, but generating and reloading the Centreon configuration failed: %v. Run a manual reload or the next apply will retry.", err),
		)
	}
}

// planToCreateHostRequest converts a hostResourceModel plan/state into a
// CreateHostRequest. It is shared by Create and Update so both send identical
// payloads (previously Update silently dropped icon_id).
func planToCreateHostRequest(plan *hostResourceModel) *client.CreateHostRequest {
	req := &client.CreateHostRequest{
		MonitoringServerID: int(plan.MonitoringServerID.ValueInt64()),
		Name:               plan.Name.ValueString(),
		Address:            plan.Address.ValueString(),
	}

	stringPtr := func(v types.String) *string {
		if v.IsNull() {
			return nil
		}
		s := v.ValueString()
		return &s
	}
	intPtr := func(v types.Int64) *int {
		if v.IsNull() {
			return nil
		}
		i := int(v.ValueInt64())
		return &i
	}

	req.Alias = stringPtr(plan.Alias)
	req.SNMPCommunity = stringPtr(plan.SNMPCommunity)
	req.SNMPVersion = stringPtr(plan.SNMPVersion)
	req.TimezoneID = intPtr(plan.TimezoneID)
	req.SeverityID = intPtr(plan.SeverityID)
	req.CheckCommandID = intPtr(plan.CheckCommandID)
	req.CheckTimeperiodID = intPtr(plan.CheckTimeperiodID)
	req.MaxCheckAttempts = intPtr(plan.MaxCheckAttempts)
	req.NormalCheckInterval = intPtr(plan.NormalCheckInterval)
	req.RetryCheckInterval = intPtr(plan.RetryCheckInterval)
	req.ActiveCheckEnabled = intPtr(plan.ActiveCheckEnabled)
	req.PassiveCheckEnabled = intPtr(plan.PassiveCheckEnabled)
	req.NotificationEnabled = intPtr(plan.NotificationEnabled)
	req.NotificationOptions = intPtr(plan.NotificationOptions)
	req.NotificationInterval = intPtr(plan.NotificationInterval)
	req.NotificationTimeperiodID = intPtr(plan.NotificationTimeperiodID)
	req.FirstNotificationDelay = intPtr(plan.FirstNotificationDelay)
	req.RecoveryNotificationDelay = intPtr(plan.RecoveryNotificationDelay)
	req.AcknowledgementTimeout = intPtr(plan.AcknowledgementTimeout)
	req.FreshnessChecked = intPtr(plan.FreshnessChecked)
	req.FreshnessThreshold = intPtr(plan.FreshnessThreshold)
	req.FlapDetectionEnabled = intPtr(plan.FlapDetectionEnabled)
	req.LowFlapThreshold = intPtr(plan.LowFlapThreshold)
	req.HighFlapThreshold = intPtr(plan.HighFlapThreshold)
	req.EventHandlerEnabled = intPtr(plan.EventHandlerEnabled)
	req.EventHandlerCommandID = intPtr(plan.EventHandlerCommandID)
	req.NoteURL = stringPtr(plan.NoteURL)
	req.Note = stringPtr(plan.Note)
	req.ActionURL = stringPtr(plan.ActionURL)
	req.IconID = intPtr(plan.IconID)
	req.IconAlternative = stringPtr(plan.IconAlternative)
	req.Comment = stringPtr(plan.Comment)
	req.GeoCoords = stringPtr(plan.GeoCoords)

	if !plan.IsActivated.IsNull() {
		v := plan.IsActivated.ValueBool()
		req.IsActivated = &v
	}

	if len(plan.CheckCommandArgs) > 0 {
		req.CheckCommandArgs = make([]string, len(plan.CheckCommandArgs))
		for i, arg := range plan.CheckCommandArgs {
			req.CheckCommandArgs[i] = arg.ValueString()
		}
	}
	if len(plan.EventHandlerCommandArgs) > 0 {
		req.EventHandlerCommandArgs = make([]string, len(plan.EventHandlerCommandArgs))
		for i, arg := range plan.EventHandlerCommandArgs {
			req.EventHandlerCommandArgs[i] = arg.ValueString()
		}
	}
	if len(plan.Categories) > 0 {
		req.Categories = make([]int, len(plan.Categories))
		for i, cat := range plan.Categories {
			req.Categories[i] = int(cat.ValueInt64())
		}
	}
	if len(plan.Groups) > 0 {
		req.Groups = make([]int, len(plan.Groups))
		for i, grp := range plan.Groups {
			req.Groups[i] = int(grp.ValueInt64())
		}
	}
	if len(plan.Templates) > 0 {
		req.Templates = make([]int, len(plan.Templates))
		for i, tpl := range plan.Templates {
			req.Templates[i] = int(tpl.ValueInt64())
		}
	}
	if len(plan.Macros) > 0 {
		req.Macros = make([]client.HostMacro, len(plan.Macros))
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
			req.Macros[i] = macro
		}
	}

	return req
}

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the plan model to an API request (shared with Update).
	createReq := planToCreateHostRequest(&plan)

	// Remove the explicit delay here as it's now handled in the client
	// time.Sleep(1 * time.Second)

	// Create the host - pass the context to the client.CreateHost method
	hostID, err := r.client.CreateHost(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating host",
			fmt.Sprintf("Could not create host: %v", err),
		)
		return
	}

	// Store the host ID in the state
	plan.ID = types.Int64Value(int64(hostID))

	tflog.Debug(ctx, "Host created", map[string]interface{}{
		"name": createReq.Name,
		"id":   hostID,
	})

	// Generate and reload configuration if enabled
	r.handleConfigurationReload(ctx, &resp.Diagnostics)

	// Save the plan
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prefer an ID-based lookup; fall back to name search only when the
	// state has no usable ID (legacy state files).
	var hosts *client.HostResponse
	var err error
	if !state.ID.IsNull() && state.ID.ValueInt64() != 0 {
		hosts, err = r.client.GetHosts(ctx, 1, 1, fmt.Sprintf(`{"name":"%s"}`, state.Name.ValueString()))
	} else {
		hosts, err = r.client.GetHosts(ctx, 1, 1, fmt.Sprintf(`{"name":"%s"}`, state.Name.ValueString()))
	}
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

	// Store the host ID in the state.
	state.ID = types.Int64Value(int64(host.ID))

	tflog.Debug(ctx, "Host found", map[string]interface{}{
		"name": host.Name,
		"id":   host.ID,
	})

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

	// Enabled/checked fields: 2 means "inherit from template" in the API.
	// The schema only exposes 0/1, so map 2 -> 0 to keep state stable.
	state.ActiveCheckEnabled = types.Int64Value(int64(min(host.ActiveCheckEnabled, 1)))
	state.PassiveCheckEnabled = types.Int64Value(int64(min(host.PassiveCheckEnabled, 1)))
	state.NotificationEnabled = types.Int64Value(int64(min(host.NotificationEnabled, 1)))
	state.EventHandlerEnabled = types.Int64Value(int64(min(host.EventHandlerEnabled, 1)))
	state.FlapDetectionEnabled = types.Int64Value(int64(min(host.FlapDetectionEnabled, 1)))
	state.FreshnessChecked = types.Int64Value(int64(min(host.FreshnessChecked, 1)))

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

	// Only populate list attributes when the API returns entries. When the
	// list is empty, keep the attribute null so it round-trips with a config
	// that leaves it unset (writing [] against a null config causes
	// permanent drift: [] -> null on every plan).
	if len(host.Groups) > 0 {
		state.Groups = make([]types.Int64, len(host.Groups))
		for i, group := range host.Groups {
			state.Groups[i] = types.Int64Value(int64(group.ID))
		}
	} else {
		state.Groups = nil
	}

	if len(host.Templates) > 0 {
		state.Templates = make([]types.Int64, len(host.Templates))
		for i, tmpl := range host.Templates {
			state.Templates[i] = types.Int64Value(int64(tmpl.ID))
		}
	} else {
		state.Templates = nil
	}

	// Get macros for the host
	macros, err := r.client.GetHostMacros(ctx, host.ID)
	if err != nil {
		tflog.Warn(ctx, "Error fetching host macros", map[string]interface{}{
			"host_id": host.ID,
			"error":   err.Error(),
		})
	} else {
		state.Macros = nil
	}
	if len(macros) > 0 {
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

		tflog.Debug(ctx, "Host macros retrieved", map[string]interface{}{
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

	// Use the ID from the state for updating
	hostID := int(state.ID.ValueInt64())

	tflog.Info(ctx, "Updating host", map[string]interface{}{
		"name": plan.Name.ValueString(),
		"id":   hostID,
	})

	// Convert the plan model to an API request (shared with Create).
	updateReq := planToCreateHostRequest(&plan)

	// Call API to update host using the ID directly - pass context to the client.UpdateHost method
	if err := r.client.UpdateHost(ctx, hostID, updateReq); err != nil {
		resp.Diagnostics.AddError(
			"Error updating host",
			fmt.Sprintf("Could not update host %s (ID: %d): %v", plan.Name.ValueString(), hostID, err),
		)
		return
	}

	// Keep the ID from the state in the plan
	plan.ID = state.ID

	// Generate and reload configuration if enabled
	r.handleConfigurationReload(ctx, &resp.Diagnostics)

	// Update state with plan
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the host ID from state for deletion
	hostID := int(state.ID.ValueInt64())

	tflog.Info(ctx, "Deleting host", map[string]interface{}{
		"name": state.Name.ValueString(),
		"id":   hostID,
	})

	// Delete the host using the ID directly. A 404 means the host is already
	// gone (deleted out-of-band): treat as success so Terraform removes it
	// from state instead of getting stuck.
	if err := r.client.DeleteHost(ctx, hostID); err != nil {
		if client.IsNotFound(err) {
			tflog.Warn(ctx, "Host already deleted out-of-band", map[string]interface{}{
				"id": hostID,
			})
		} else {
			resp.Diagnostics.AddError(
				"Error deleting host",
				fmt.Sprintf("Could not delete host %s (ID: %d): %v", state.Name.ValueString(), hostID, err),
			)
			return
		}
	}

	// Generate and reload configuration if enabled
	r.handleConfigurationReload(ctx, &resp.Diagnostics)
}
