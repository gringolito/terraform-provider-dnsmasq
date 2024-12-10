package provider

import (
	"context"
	"os"

	"terraform-provider-dnsmasq/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure dnsmasqProvider satisfies various provider interfaces.
var _ provider.Provider = &dnsmasqProvider{}

// dnsmasqProvider defines the provider implementation.
type dnsmasqProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// dnsmasqProviderModel describes the provider data model.
type dnsmasqProviderModel struct {
	URL   types.String `tfsdk:"api_url"`
	Token types.String `tfsdk:"api_token"`
}

func (p *dnsmasqProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dnsmasq"
	resp.Version = p.version
}

func (p *dnsmasqProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use the dnsmasq provider to manage dnsmasq resources using the dnsmasq-manager API (see: https://github.com/gringolito/dnsmasq-manager).",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "dnsmasq-manager API URL.",
				Required:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "dnsmasq-manager API JWT authentication token.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *dnsmasqProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config dnsmasqProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown dnsmasq-manager API URL",
			"The provider cannot create the dnsmasq client as there is an unknown configuration value for the dnsmasq-manager API URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DMM_API_URL environment variable.",
		)
	}
	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown dnsmasq-manager API token",
			"The provider cannot create the dnsmasq client as there is an unknown configuration value for the dnsmasq-manager API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DMM_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	url := os.Getenv("DMM_API_URL")
	token := os.Getenv("DMM_API_TOKEN")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing dnsmasq-manager API URL",
			"The provider cannot create the dnsmasq client as there is a missing or empty value for the dnsmasq-manager API URL. "+
				"Set the host value in the configuration or use the DMM_API_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := client.New(url, token)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *dnsmasqProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDhcpStaticHostResource,
	}
}

func (p *dnsmasqProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDhcpStaticHostDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &dnsmasqProvider{
			version: version,
		}
	}
}
