package provider

import (
	"context"
	"fmt"
	"terraform-provider-dnsmasq/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DhcpStaticHostDataSource{}

func NewDhcpStaticHostDataSource() datasource.DataSource {
	return &DhcpStaticHostDataSource{}
}

// DhcpStaticHostDataSource defines the data source implementation.
type DhcpStaticHostDataSource struct {
	client client.Client
}

// DhcpStaticHostDataSourceModel describes the data source data model.
type DhcpStaticHostDataSourceModel struct {
	MacAddress types.String `tfsdk:"mac_address"`
	IPAddress  types.String `tfsdk:"ip_address"`
	HostName   types.String `tfsdk:"hostname"`
	Id         types.String `tfsdk:"id"`
}

func (d *DhcpStaticHostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_static_host"
}

func (d *DhcpStaticHostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a static DHCP lease reservation.",

		Attributes: map[string]schema.Attribute{
			"mac_address": schema.StringAttribute{
				MarkdownDescription: "Host MAC address (id) to filter the search.",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP address assigned to the host on the static DHCP lease reservation.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname assigned to the host on the static DHCP lease reservation.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Host MAC address identifier.",
				Computed:            true,
			},
		},
	}
}

func (d *DhcpStaticHostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DhcpStaticHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read Terraform configuration data into the model
	var data DhcpStaticHostDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := d.client.ReadStaticDhcpHost(data.MacAddress.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DHCP Static Host", err.Error())
		return
	}

	tflog.Trace(ctx, "read a DHCP static host data source")

	// Save data into Terraform state
	data.Id = types.StringValue((host.MacAddress))
	data.IPAddress = types.StringValue(host.IPAddress)
	data.HostName = types.StringValue(host.HostName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
