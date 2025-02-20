package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"
	"terraform-provider-centreon/internal/logging"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &hostsDataSource{}

func NewHostsDataSource() datasource.DataSource {
	return &hostsDataSource{}
}

type hostsDataSource struct {
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

type hostModel struct {
	ID                     types.Int64           `tfsdk:"id"`
	Name                   types.String          `tfsdk:"name"`
	Alias                  types.String          `tfsdk:"alias"`
	Address                types.String          `tfsdk:"address"`
	MonitoringServer       monitoringServerModel `tfsdk:"monitoring_server"`
	Templates              []hostTemplateModel   `tfsdk:"templates"`
	NormalCheckInterval    types.String          `tfsdk:"normal_check_interval"`
	RetryCheckInterval     types.String          `tfsdk:"retry_check_interval"`
	NotificationTimeperiod types.String          `tfsdk:"notification_timeperiod"`
	CheckTimeperiod        types.String          `tfsdk:"check_timeperiod"`
	Severity               types.String          `tfsdk:"severity"`
	Categories             []types.String        `tfsdk:"categories"`
	Groups                 []hostGroupModel      `tfsdk:"groups"`
	IsActivated            types.Bool            `tfsdk:"is_activated"`
}

type hostsDataSourceModel struct {
	Limit  types.Int64  `tfsdk:"limit"`
	Page   types.Int64  `tfsdk:"page"`
	Search searchModel  `tfsdk:"search"`
	Hosts  []hostModel  `tfsdk:"hosts"`
	Id     types.String `tfsdk:"id"`
}

func (d *hostsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosts"
}

func (d *hostsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *hostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		logging.Error(context.Background(), "No provider data available")
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client, got: nil",
		)
		logging.Error(context.Background(), "Invalid provider data type")
		return
	}

	d.client = client
	logging.Debug(context.Background(), "Hosts data source configured successfully")
}

func (d *hostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		logging.Error(ctx, "Failed to get configuration", map[string]interface{}{
			"error": resp.Diagnostics.Errors(),
		})
		return
	}

	// Create search JSON
	searchQuery := "{}"
	if !state.Search.Name.IsNull() && !state.Search.Value.IsNull() {
		searchQuery = fmt.Sprintf("{\"%s\":\"%s\"}",
			state.Search.Name.ValueString(),
			state.Search.Value.ValueString())
		logging.Debug(ctx, "Using search query", map[string]interface{}{
			"query": searchQuery,
		})
	}

	logging.Info(ctx, "Fetching hosts", map[string]interface{}{
		"limit":  state.Limit.ValueInt64(),
		"page":   state.Page.ValueInt64(),
		"search": searchQuery,
	})

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
		logging.Error(ctx, "Failed to fetch hosts", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	logging.Debug(ctx, "Successfully retrieved hosts", map[string]interface{}{
		"count": len(hostResponse.Result),
	})

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
			categories[j] = types.StringValue(fmt.Sprintf("%d", cat))
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
			NormalCheckInterval:    types.StringValue(fmt.Sprintf("%d", host.NormalCheckInterval)),
			RetryCheckInterval:     types.StringValue(fmt.Sprintf("%d", host.RetryCheckInterval)),
			NotificationTimeperiod: types.StringValue(fmt.Sprintf("%d", host.NotificationTimeperiodID)),
			CheckTimeperiod:        types.StringValue(fmt.Sprintf("%d", host.CheckTimeperiodID)),
			Severity:               types.StringValue(fmt.Sprintf("%d", host.SeverityID)),
			Categories:             categories,
			IsActivated:            types.BoolValue(host.IsActivated),
		}
	}

	state.Id = types.StringValue("hosts")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
