package provider

import (
	"context"
	"terraform-provider-centreon/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &platformInfoDataSource{}

func NewPlatformInfoDataSource() datasource.DataSource {
	return &platformInfoDataSource{}
}

type platformInfoDataSource struct{
	client *client.Client
}

type platformInfoDataSourceModel struct {
	IsInstalled         types.Bool `tfsdk:"is_installed"`
	HasUpgradeAvailable types.Bool `tfsdk:"has_upgrade_available"`
	Id                  types.String `tfsdk:"id"`
}

func (d *platformInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform_info"
}

func (d *platformInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches Centreon platform installation status information.",
		Attributes: map[string]schema.Attribute{
			"is_installed": schema.BoolAttribute{
				Description: "Indicates if Centreon is installed",
				Computed:    true,
			},
			"has_upgrade_available": schema.BoolAttribute{
				Description: "Indicates if an upgrade is available",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Placeholder identifier",
				Computed:    true,
			},
		},
	}
}

func (d *platformInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *platformInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state platformInfoDataSourceModel

	if d.client == nil {
		resp.Diagnostics.AddError(
			"Client not configured",
			"The client was not properly configured",
		)
		return
	}
	
	platformInfo, err := d.client.GetPlatformInfo()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Platform Info",
			err.Error(),
		)
		return
	}

	state.IsInstalled = types.BoolValue(platformInfo.IsInstalled)
	state.HasUpgradeAvailable = types.BoolValue(platformInfo.HasUpgradeAvailable)
	state.Id = types.StringValue("platform_info")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}