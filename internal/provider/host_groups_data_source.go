package provider

import (
	"context"
	"fmt"
	"terraform-provider-centreon/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &hostGroupsDataSource{}

func NewHostGroupsDataSource() datasource.DataSource {
	return &hostGroupsDataSource{}
}

type hostGroupsDataSource struct {
	client *client.Client
}

type hostGroupDetail struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type hostGroupsDataSourceModel struct {
	Limit  types.Int64        `tfsdk:"limit"`
	Page   types.Int64        `tfsdk:"page"`
	Search *searchModel       `tfsdk:"search"`
	Groups []hostGroupDetail `tfsdk:"groups"`
	Id     types.String      `tfsdk:"id"`
}

func (d *hostGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_groups"
}

func (d *hostGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of host groups.",
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
			"groups": schema.ListNestedAttribute{
				Description: "List of host groups",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Group ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Group name",
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

func (d *hostGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *hostGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostGroupsDataSourceModel

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

	groupsResponse, err := d.client.GetHostGroups(
		int(state.Limit.ValueInt64()),
		int(state.Page.ValueInt64()),
		searchQuery,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Host Groups",
			err.Error(),
		)
		return
	}

	// Map response to model
	state.Groups = make([]hostGroupDetail, len(groupsResponse.Result))
	for i, group := range groupsResponse.Result {
		state.Groups[i] = hostGroupDetail{
			ID:   types.Int64Value(int64(group.ID)),
			Name: types.StringValue(group.Name),
		}
	}

	state.Id = types.StringValue("host_groups")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
