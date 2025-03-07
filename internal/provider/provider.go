package provider

import (
	"context"
	"terraform-provider-centreon/internal/client"
	"terraform-provider-centreon/internal/logging"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &centreonProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &centreonProvider{
			version: version,
		}
	}
}

// centreonProvider is the provider implementation.
type centreonProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type centreonProviderModel struct {
	Protocol                       types.String `tfsdk:"protocol"`
	Server                         types.String `tfsdk:"server"`
	Port                           types.String `tfsdk:"port"`
	APIVersion                     types.String `tfsdk:"api_version"`
	APIKey                         types.String `tfsdk:"api_key"`
	GenerateAndReloadConfiguration types.Bool   `tfsdk:"generate_and_reload_configuration"`
}

// Metadata returns the provider type name.
func (p *centreonProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "centreon"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *centreonProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol to use for API calls (http or https)",
			},
			"server": schema.StringAttribute{
				Required:    true,
				Description: "Centreon server hostname (eg. 'centreon.example.com')",
			},
			"port": schema.StringAttribute{
				Required:    true,
				Description: "Centreon server port (eg. 80, 443)",
			},
			"api_version": schema.StringAttribute{
				Required:    true,
				Description: "API version to use (e.g., 'latest', 'v24.10')",
			},
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "API key for authentication",
			},
			"generate_and_reload_configuration": schema.BoolAttribute{
				Optional:    true,
				Description: "When true, automatically generates and reloads the configuration for all monitoring servers after applying changes",
			},
		},
	}
}

// Configure prepares a Centreon API client for data sources and resources.
func (p *centreonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Initialize file logger
	logCtx, err := logging.InitializeFileLogger(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to setup logging",
			err.Error(),
		)
		return
	}

	// Use the logging context for the rest of the configuration
	ctx = logCtx

	var config centreonProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Protocol.IsNull() || config.Server.IsNull() || config.Port.IsNull() ||
		config.APIVersion.IsNull() || config.APIKey.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"All provider configuration fields are required",
		)
		return
	}

	logging.Info(ctx, "Configuring Centreon client",
		map[string]interface{}{
			"server":      config.Server.ValueString(),
			"protocol":    config.Protocol.ValueString(),
			"port":        config.Port.ValueString(),
			"api_version": config.APIVersion.ValueString(),
		})

	client := client.NewClient(
		config.Protocol.ValueString(),
		config.Server.ValueString(),
		config.Port.ValueString(),
		config.APIVersion.ValueString(),
		config.APIKey.ValueString(),
	)
	client.GenerateAndReloadConfiguration = !config.GenerateAndReloadConfiguration.IsNull() && config.GenerateAndReloadConfiguration.ValueBool()

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *centreonProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPlatformInfoDataSource,
		NewHostsDataSource,
		NewMonitoringServersDataSource,
		NewHostGroupsDataSource,
		NewHostTemplatesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *centreonProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHostResource,
	}
}
