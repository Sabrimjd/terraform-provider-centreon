// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &configurationHostsDataSource{}

func NewConfigurationHostsDataSource() datasource.DataSource {
	return &configurationHostsDataSource{}
}

type configurationHostsDataSource struct {
	client *client.Client
}

type hostTemplateModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type monitoringServerModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type hostGroupModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type searchModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type hostModel struct {
	ID                     types.Int64            `tfsdk:"id"`
	Name                   types.String           `tfsdk:"name"`
	Alias                  types.String           `tfsdk:"alias"`
	Address                types.String           `tfsdk:"address"`
	MonitoringServer       monitoringServerModel  `tfsdk:"monitoring_server"`
	Templates              []hostTemplateModel    `tfsdk:"templates"`
	NormalCheckInterval    types.String           `tfsdk:"normal_check_interval"`
	RetryCheckInterval     types.String           `tfsdk:"retry_check_interval"`
	NotificationTimeperiod types.String           `tfsdk:"notification_timeperiod"`
	CheckTimeperiod        types.String           `tfsdk:"check_timeperiod"`
	Severity               types.String           `tfsdk:"severity"`
	Categories             []types.String         `tfsdk:"categories"`
	Groups                 []hostGroupModel       `tfsdk:"groups"`
	IsActivated            types.Bool             `tfsdk:"is_activated"`
}

type configurationHostsDataSourceModel struct {
	Limit  types.Int64   `tfsdk:"limit"`
	Page   types.Int64   `tfsdk:"page"`
	Search searchModel   `tfsdk:"search"`
	Hosts  []hostModel   `tfsdk:"hosts"`
	Id     types.String  `tfsdk:"id"`
}

func (d *configurationHostsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_hosts"
}

func (d *configurationHostsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Searches for Centreon hosts.",
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "Number of results to return",
				Required:    true,
			},
			"page": schema.Int64Attribute{
				Description: "Page number",
				Required:    true,
				},
				"search": schema.ObjectAttribute{
				Description: "Search criteria",
				Optional:    true,
				AttributeTypes: map[string]attr.Type{
					"name":  types.StringType,
					"value": types.StringType,
				},
			},
			"hosts": schema.ListNestedAttribute{
				Description: "List of hosts matching the search criteria",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"alias": schema.StringAttribute{
							Computed: true,
						},
						"address": schema.StringAttribute{
							Computed: true,
						},
						"monitoring_server": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"templates": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"groups": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"normal_check_interval": schema.StringAttribute{
							Computed: true,
						},
						"retry_check_interval": schema.StringAttribute{
							Computed: true,
						},
						"notification_timeperiod": schema.StringAttribute{
							Computed: true,
						},
						"check_timeperiod": schema.StringAttribute{
							Computed: true,
						},
						"severity": schema.StringAttribute{
							Computed: true,
						},
						"categories": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"is_activated": schema.BoolAttribute{
							Computed: true,
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

func (d *configurationHostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client, got: nil",
		)
		return
	}

	d.client = client
}

func (d *configurationHostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state configurationHostsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create search JSON
	searchQuery := "{}"
	if !state.Search.Name.IsNull() && !state.Search.Value.IsNull() {
		searchQuery = fmt.Sprintf("{\"%s\":\"%s\"}", 
			state.Search.Name.ValueString(), 
			state.Search.Value.ValueString())
	}

	hostResponse, err := d.client.GetHosts(
		int(state.Limit.ValueInt64()),
		int(state.Page.ValueInt64()),
		searchQuery,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hosts",
			err.Error(),
		)
		return
	}

	// Map response to model
	state.Hosts = make([]hostModel, len(hostResponse.Result))
	for i, host := range hostResponse.Result {
		templates := make([]hostTemplateModel, len(host.Templates))
		for j, tpl := range host.Templates {
			templates[j] = hostTemplateModel{
				ID:   types.Int64Value(int64(tpl.ID)),
				Name: types.StringValue(tpl.Name),
			}
		}

		groups := make([]hostGroupModel, len(host.Groups))
		for j, grp := range host.Groups {
			groups[j] = hostGroupModel{
				ID:   types.Int64Value(int64(grp.ID)),
				Name: types.StringValue(grp.Name),
			}
		}

		// Convert []string to []types.String for Categories
		categories := make([]types.String, len(host.Categories))
		for j, cat := range host.Categories {
			categories[j] = types.StringValue(cat)
		}

		state.Hosts[i] = hostModel{
			ID:      types.Int64Value(int64(host.ID)),
			Name:    types.StringValue(host.Name),
			Alias:   types.StringValue(host.Alias),
			Address: types.StringValue(host.Address),
			MonitoringServer: monitoringServerModel{
				ID:   types.Int64Value(int64(host.MonitoringServer.ID)),
				Name: types.StringValue(host.MonitoringServer.Name),
			},
			Templates:              templates,
			Groups:                 groups,
			NormalCheckInterval:    types.StringValue(host.NormalCheckInterval),
			RetryCheckInterval:     types.StringValue(host.RetryCheckInterval),
			NotificationTimeperiod: types.StringValue(host.NotificationTimeperiod),
			CheckTimeperiod:        types.StringValue(host.CheckTimeperiod),
			Severity:               types.StringValue(host.Severity),
			Categories:             categories,
			IsActivated:            types.BoolValue(host.IsActivated),
		}
	}

	state.Id = types.StringValue("configuration_hosts")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
