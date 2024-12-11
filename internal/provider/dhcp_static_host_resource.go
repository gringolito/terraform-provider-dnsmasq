package provider

import (
	"context"
	"fmt"
	"terraform-provider-dnsmasq/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DhcpStaticHostResource{}
	_ resource.ResourceWithConfigure   = &DhcpStaticHostResource{}
	_ resource.ResourceWithImportState = &DhcpStaticHostResource{}
)

func NewDhcpStaticHostResource() resource.Resource {
	return &DhcpStaticHostResource{}
}

// DhcpStaticHostResource defines the resource implementation.
type DhcpStaticHostResource struct {
	client client.Client
}

// DhcpStaticHostResourceModel describes the resource data model.
type DhcpStaticHostResourceModel struct {
	MacAddress types.String `tfsdk:"mac_address"`
	IPAddress  types.String `tfsdk:"ip_address"`
	HostName   types.String `tfsdk:"hostname"`
}

func (m *DhcpStaticHostResourceModel) toDnsmasq() client.StaticDhcpHost {
	return client.StaticDhcpHost{
		MacAddress: m.MacAddress.ValueString(),
		IPAddress:  m.IPAddress.ValueString(),
		HostName:   m.HostName.ValueString(),
	}
}

func (m *DhcpStaticHostResourceModel) fromDnsmasq(host *client.StaticDhcpHost) {
	m.MacAddress = types.StringValue(host.MacAddress)
	m.IPAddress = types.StringValue(host.IPAddress)
	m.HostName = types.StringValue(host.HostName)
}

func (r *DhcpStaticHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_static_host"
}

func (r *DhcpStaticHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a static DHCP lease reservation resource. This allow static DHCP lease reservations to be allocated, modified, and released.",

		Attributes: map[string]schema.Attribute{
			"mac_address": schema.StringAttribute{
				MarkdownDescription: "Host MAC address.",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP address to be assigned to the host on the static DHCP lease reservation.",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname to be assigned to the host on the static DHCP lease reservation.",
				Required:            true,
			},
		},
	}
}

func (r *DhcpStaticHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DhcpStaticHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan DhcpStaticHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := r.client.CreateStaticDhcpHost(plan.toDnsmasq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to create DHCP Static Host", err.Error())
		return
	}

	tflog.Trace(ctx, "created a DHCP static host resource")

	// Save data into Terraform state
	var state DhcpStaticHostResourceModel
	state.fromDnsmasq(host)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DhcpStaticHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state DhcpStaticHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := r.client.ReadStaticDhcpHost(state.MacAddress.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DHCP Static Host", err.Error())
		return
	}

	tflog.Trace(ctx, "read a DHCP static host resource")

	// Save updated data into Terraform state
	state.fromDnsmasq(host)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DhcpStaticHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan DhcpStaticHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := r.client.UpdateStaticDhcpHost(plan.toDnsmasq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update DHCP Static Host", err.Error())
		return
	}

	tflog.Trace(ctx, "updated a DHCP static host resource")

	// Save updated data into Terraform state
	var state DhcpStaticHostResourceModel
	state.fromDnsmasq(host)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DhcpStaticHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read Terraform prior state state into the model
	var state DhcpStaticHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteStaticDhcpHost(state.MacAddress.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete DHCP Static Host", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted a DHCP static host resource")
}

func (r *DhcpStaticHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("mac_address"), req, resp)
}