package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &hostTemplatesDataSource{}

func NewHostTemplatesDataSource() datasource.DataSource {
	return &hostTemplatesDataSource{}
}

type hostTemplatesDataSource struct {
	client *client.Client
}

type hostTemplateDetail struct {
	ID                      types.Int64    `tfsdk:"id"`
	Name                    types.String   `tfsdk:"name"`
	Alias                   types.String   `tfsdk:"alias"`
	SNMPVersion            types.String   `tfsdk:"snmp_version"`
	TimezoneID             types.Int64    `tfsdk:"timezone_id"`
	SeverityID             types.Int64    `tfsdk:"severity_id"`
	CheckCommandID         types.Int64    `tfsdk:"check_command_id"`
	CheckCommandArgs       []types.String `tfsdk:"check_command_args"`
	CheckTimeperiodID      types.Int64    `tfsdk:"check_timeperiod_id"`
	MaxCheckAttempts       types.Int64    `tfsdk:"max_check_attempts"`
	NormalCheckInterval    types.Int64    `tfsdk:"normal_check_interval"`
	RetryCheckInterval     types.Int64    `tfsdk:"retry_check_interval"`
	ActiveCheckEnabled     types.Int64    `tfsdk:"active_check_enabled"`
	PassiveCheckEnabled    types.Int64    `tfsdk:"passive_check_enabled"`
	NotificationEnabled    types.Int64    `tfsdk:"notification_enabled"`
	NotificationOptions    types.Int64    `tfsdk:"notification_options"`
	NotificationInterval   types.Int64    `tfsdk:"notification_interval"`
	NotificationTimeperiodID types.Int64  `tfsdk:"notification_timeperiod_id"`
	AddInheritedContactGroup types.Bool   `tfsdk:"add_inherited_contact_group"`
	AddInheritedContact    types.Bool    `tfsdk:"add_inherited_contact"`
	FirstNotificationDelay types.Int64    `tfsdk:"first_notification_delay"`
	RecoveryNotificationDelay types.Int64 `tfsdk:"recovery_notification_delay"`
	AcknowledgementTimeout types.Int64    `tfsdk:"acknowledgement_timeout"`
	FreshnessChecked      types.Int64    `tfsdk:"freshness_checked"`
	FreshnessThreshold    types.Int64    `tfsdk:"freshness_threshold"`
	FlapDetectionEnabled  types.Int64    `tfsdk:"flap_detection_enabled"`
	LowFlapThreshold      types.Int64    `tfsdk:"low_flap_threshold"`
	HighFlapThreshold     types.Int64    `tfsdk:"high_flap_threshold"`
	EventHandlerEnabled   types.Int64    `tfsdk:"event_handler_enabled"`
	EventHandlerCommandID types.Int64    `tfsdk:"event_handler_command_id"`
	EventHandlerCommandArgs []types.String `tfsdk:"event_handler_command_args"`
	NoteURL               types.String   `tfsdk:"note_url"`
	Note                  types.String   `tfsdk:"note"`
	ActionURL             types.String   `tfsdk:"action_url"`
	IconID               types.Int64    `tfsdk:"icon_id"`
	IconAlternative      types.String   `tfsdk:"icon_alternative"`
	Comment              types.String   `tfsdk:"comment"`
	IsLocked             types.Bool     `tfsdk:"is_locked"`
}

type hostTemplatesDataSourceModel struct {
	Limit     types.Int64           `tfsdk:"limit"`
	Page      types.Int64           `tfsdk:"page"`
	Search    *searchModel          `tfsdk:"search"`
	Templates []hostTemplateDetail  `tfsdk:"templates"`
	Id        types.String         `tfsdk:"id"`
}

func (d *hostTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_templates"
}

func (d *hostTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of host templates.",
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
			"templates": schema.ListNestedAttribute{
				Description: "List of host templates",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Template ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Template name",
							Computed:    true,
						},
						"alias": schema.StringAttribute{
							Description: "Template alias",
							Computed:    true,
						},
						"snmp_version": schema.StringAttribute{
							Description: "SNMP version",
							Computed:    true,
						},
						"timezone_id": schema.Int64Attribute{
							Description: "Timezone ID",
							Computed:    true,
						},
						"severity_id": schema.Int64Attribute{
							Description: "Severity ID",
							Computed:    true,
						},
						"check_command_id": schema.Int64Attribute{
							Description: "Check command ID",
							Computed:    true,
						},
						"check_command_args": schema.ListAttribute{
							Description: "Check command arguments",
							Computed:    true,
							ElementType: types.StringType,
						},
						"check_timeperiod_id": schema.Int64Attribute{
							Description: "Check timeperiod ID",
							Computed:    true,
						},
						"max_check_attempts": schema.Int64Attribute{
							Description: "Maximum check attempts",
							Computed:    true,
						},
						"normal_check_interval": schema.Int64Attribute{
							Description: "Normal check interval",
							Computed:    true,
						},
						"retry_check_interval": schema.Int64Attribute{
							Description: "Retry check interval",
							Computed:    true,
						},
						"active_check_enabled": schema.Int64Attribute{
							Description: "Active check enabled",
							Computed:    true,
						},
						"passive_check_enabled": schema.Int64Attribute{
							Description: "Passive check enabled",
							Computed:    true,
						},
						"notification_enabled": schema.Int64Attribute{
							Description: "Notification enabled",
							Computed:    true,
						},
						"notification_options": schema.Int64Attribute{
							Description: "Notification options",
							Computed:    true,
						},
						"notification_interval": schema.Int64Attribute{
							Description: "Notification interval",
							Computed:    true,
						},
						"notification_timeperiod_id": schema.Int64Attribute{
							Description: "Notification timeperiod ID",
							Computed:    true,
						},
						"add_inherited_contact_group": schema.BoolAttribute{
							Description: "Add inherited contact group",
							Computed:    true,
						},
						"add_inherited_contact": schema.BoolAttribute{
							Description: "Add inherited contact",
							Computed:    true,
						},
						"first_notification_delay": schema.Int64Attribute{
							Description: "First notification delay",
							Computed:    true,
						},
						"recovery_notification_delay": schema.Int64Attribute{
							Description: "Recovery notification delay",
							Computed:    true,
						},
						"acknowledgement_timeout": schema.Int64Attribute{
							Description: "Acknowledgement timeout",
							Computed:    true,
						},
						"freshness_checked": schema.Int64Attribute{
							Description: "Freshness checked",
							Computed:    true,
						},
						"freshness_threshold": schema.Int64Attribute{
							Description: "Freshness threshold",
							Computed:    true,
						},
						"flap_detection_enabled": schema.Int64Attribute{
							Description: "Flap detection enabled",
							Computed:    true,
						},
						"low_flap_threshold": schema.Int64Attribute{
							Description: "Low flap threshold",
							Computed:    true,
						},
						"high_flap_threshold": schema.Int64Attribute{
							Description: "High flap threshold",
							Computed:    true,
						},
						"event_handler_enabled": schema.Int64Attribute{
							Description: "Event handler enabled",
							Computed:    true,
						},
						"event_handler_command_id": schema.Int64Attribute{
							Description: "Event handler command ID",
							Computed:    true,
						},
						"event_handler_command_args": schema.ListAttribute{
							Description: "Event handler command arguments",
							Computed:    true,
							ElementType: types.StringType,
						},
						"note_url": schema.StringAttribute{
							Description: "Note URL",
							Computed:    true,
						},
						"note": schema.StringAttribute{
							Description: "Note",
							Computed:    true,
						},
						"action_url": schema.StringAttribute{
							Description: "Action URL",
							Computed:    true,
						},
						"icon_id": schema.Int64Attribute{
							Description: "Icon ID",
							Computed:    true,
						},
						"icon_alternative": schema.StringAttribute{
							Description: "Icon alternative text",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "Comment",
							Computed:    true,
						},
						"is_locked": schema.BoolAttribute{
							Description: "Is locked",
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

func (d *hostTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *hostTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostTemplatesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize empty search if not provided
	if (state.Search == nil) {
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

	templatesResponse, err := d.client.GetHostTemplates(
		int(state.Limit.ValueInt64()),
		int(state.Page.ValueInt64()),
		searchQuery,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Host Templates",
			err.Error(),
		)
		return
	}

	// Map response to model
	state.Templates = make([]hostTemplateDetail, len(templatesResponse.Result))
	for i, template := range templatesResponse.Result {
		// Convert []string to []types.String for CheckCommandArgs
		checkCommandArgs := make([]types.String, len(template.CheckCommandArgs))
		for j, arg := range template.CheckCommandArgs {
			checkCommandArgs[j] = types.StringValue(arg)
		}

		// Convert []string to []types.String for EventHandlerCommandArgs
		eventHandlerCommandArgs := make([]types.String, len(template.EventHandlerCommandArgs))
		for j, arg := range template.EventHandlerCommandArgs {
			eventHandlerCommandArgs[j] = types.StringValue(arg)
		}

		// Convert pointer values safely
		var snmpVersion types.String
		if template.SNMPVersion != nil {
			snmpVersion = types.StringValue(*template.SNMPVersion)
		} else {
			snmpVersion = types.StringNull()
		}

		var timezoneID types.Int64
		if template.TimezoneID != nil {
			timezoneID = types.Int64Value(int64(*template.TimezoneID))
		} else {
			timezoneID = types.Int64Null()
		}

		var severityID types.Int64
		if template.SeverityID != nil {
			severityID = types.Int64Value(int64(*template.SeverityID))
		} else {
			severityID = types.Int64Null()
		}

		var checkCommandID types.Int64
		if template.CheckCommandID != nil {
			checkCommandID = types.Int64Value(int64(*template.CheckCommandID))
		} else {
			checkCommandID = types.Int64Null()
		}

		var checkTimeperiodID types.Int64
		if template.CheckTimeperiodID != nil {
			checkTimeperiodID = types.Int64Value(int64(*template.CheckTimeperiodID))
		} else {
			checkTimeperiodID = types.Int64Null()
		}

		var maxCheckAttempts types.Int64
		if template.MaxCheckAttempts != nil {
			maxCheckAttempts = types.Int64Value(int64(*template.MaxCheckAttempts))
		} else {
			maxCheckAttempts = types.Int64Null()
		}

		var normalCheckInterval types.Int64
		if template.NormalCheckInterval != nil {
			normalCheckInterval = types.Int64Value(int64(*template.NormalCheckInterval))
		} else {
			normalCheckInterval = types.Int64Null()
		}

		var retryCheckInterval types.Int64
		if template.RetryCheckInterval != nil {
			retryCheckInterval = types.Int64Value(int64(*template.RetryCheckInterval))
		} else {
			retryCheckInterval = types.Int64Null()
		}

		var notificationOptions types.Int64
		if template.NotificationOptions != nil {
			notificationOptions = types.Int64Value(int64(*template.NotificationOptions))
		} else {
			notificationOptions = types.Int64Null()
		}

		var notificationInterval types.Int64
		if template.NotificationInterval != nil {
			notificationInterval = types.Int64Value(int64(*template.NotificationInterval))
		} else {
			notificationInterval = types.Int64Null()
		}

		var notificationTimeperiodID types.Int64
		if template.NotificationTimeperiodID != nil {
			notificationTimeperiodID = types.Int64Value(int64(*template.NotificationTimeperiodID))
		} else {
			notificationTimeperiodID = types.Int64Null()
		}

		var firstNotificationDelay types.Int64
		if template.FirstNotificationDelay != nil {
			firstNotificationDelay = types.Int64Value(int64(*template.FirstNotificationDelay))
		} else {
			firstNotificationDelay = types.Int64Null()
		}

		var recoveryNotificationDelay types.Int64
		if template.RecoveryNotificationDelay != nil {
			recoveryNotificationDelay = types.Int64Value(int64(*template.RecoveryNotificationDelay))
		} else {
			recoveryNotificationDelay = types.Int64Null()
		}

		var acknowledgementTimeout types.Int64
		if template.AcknowledgementTimeout != nil {
			acknowledgementTimeout = types.Int64Value(int64(*template.AcknowledgementTimeout))
		} else {
			acknowledgementTimeout = types.Int64Null()
		}

		var freshnessThreshold types.Int64
		if template.FreshnessThreshold != nil {
			freshnessThreshold = types.Int64Value(int64(*template.FreshnessThreshold))
		} else {
			freshnessThreshold = types.Int64Null()
		}

		var lowFlapThreshold types.Int64
		if template.LowFlapThreshold != nil {
			lowFlapThreshold = types.Int64Value(int64(*template.LowFlapThreshold))
		} else {
			lowFlapThreshold = types.Int64Null()
		}

		var highFlapThreshold types.Int64
		if template.HighFlapThreshold != nil {
			highFlapThreshold = types.Int64Value(int64(*template.HighFlapThreshold))
		} else {
			highFlapThreshold = types.Int64Null()
		}

		var eventHandlerCommandID types.Int64
		if template.EventHandlerCommandID != nil {
			eventHandlerCommandID = types.Int64Value(int64(*template.EventHandlerCommandID))
		} else {
			eventHandlerCommandID = types.Int64Null()
		}

		var noteURL types.String
		if template.NoteURL != nil {
			noteURL = types.StringValue(*template.NoteURL)
		} else {
			noteURL = types.StringNull()
		}

		var note types.String
		if template.Note != nil {
			note = types.StringValue(*template.Note)
		} else {
			note = types.StringNull()
		}

		var actionURL types.String
		if template.ActionURL != nil {
			actionURL = types.StringValue(*template.ActionURL)
		} else {
			actionURL = types.StringNull()
		}

		var iconID types.Int64
		if template.IconID != nil {
			iconID = types.Int64Value(int64(*template.IconID))
		} else {
			iconID = types.Int64Null()
		}

		var iconAlternative types.String
		if template.IconAlternative != nil {
			iconAlternative = types.StringValue(*template.IconAlternative)
		} else {
			iconAlternative = types.StringNull()
		}

		state.Templates[i] = hostTemplateDetail{
			ID:                      types.Int64Value(int64(template.ID)),
			Name:                    types.StringValue(template.Name),
			Alias:                   types.StringValue(template.Alias),
			SNMPVersion:            snmpVersion,
			TimezoneID:             timezoneID,
			SeverityID:             severityID,
			CheckCommandID:         checkCommandID,
			CheckCommandArgs:       checkCommandArgs,
			CheckTimeperiodID:      checkTimeperiodID,
			MaxCheckAttempts:       maxCheckAttempts,
			NormalCheckInterval:    normalCheckInterval,
			RetryCheckInterval:     retryCheckInterval,
			ActiveCheckEnabled:     types.Int64Value(int64(template.ActiveCheckEnabled)),
			PassiveCheckEnabled:    types.Int64Value(int64(template.PassiveCheckEnabled)),
			NotificationEnabled:    types.Int64Value(int64(template.NotificationEnabled)),
			NotificationOptions:    notificationOptions,
			NotificationInterval:   notificationInterval,
			NotificationTimeperiodID: notificationTimeperiodID,
			AddInheritedContactGroup: types.BoolValue(template.AddInheritedContactGroup),
			AddInheritedContact:    types.BoolValue(template.AddInheritedContact),
			FirstNotificationDelay: firstNotificationDelay,
			RecoveryNotificationDelay: recoveryNotificationDelay,
			AcknowledgementTimeout: acknowledgementTimeout,
			FreshnessChecked:      types.Int64Value(int64(template.FreshnessChecked)),
			FreshnessThreshold:    freshnessThreshold,
			FlapDetectionEnabled:  types.Int64Value(int64(template.FlapDetectionEnabled)),
			LowFlapThreshold:      lowFlapThreshold,
			HighFlapThreshold:     highFlapThreshold,
			EventHandlerEnabled:   types.Int64Value(int64(template.EventHandlerEnabled)),
			EventHandlerCommandID: eventHandlerCommandID,
			EventHandlerCommandArgs: eventHandlerCommandArgs,
			NoteURL:               noteURL,
			Note:                  note,
			ActionURL:             actionURL,
			IconID:               iconID,
			IconAlternative:      iconAlternative,
			Comment:              types.StringValue(template.Comment),
			IsLocked:             types.BoolValue(template.IsLocked),
		}
	}

	state.Id = types.StringValue("host_templates")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
