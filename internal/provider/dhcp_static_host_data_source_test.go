package provider

import (
	"terraform-provider-dnsmasq/internal/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDhcpStaticHostDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() { setupDhcpStaticHostDataSourceTest(t) },
				Config:    providerConfig + testAccDhcpStaticHostDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dnsmasq_dhcp_static_host.test", "mac_address", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("data.dnsmasq_dhcp_static_host.test", "ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr("data.dnsmasq_dhcp_static_host.test", "hostname", "example"),
				),
			},
		},
		CheckDestroy: teardownDhcpStaticHostDataSourceTest,
	})
}

const testAccDhcpStaticHostDataSourceConfig = `
data "dnsmasq_dhcp_static_host" "test" {
  mac_address = "00:11:22:33:44:55"
}
`

func setupDhcpStaticHostDataSourceTest(t *testing.T) {
	dnsmasq := client.New(apiUrl, "")
	_, err := dnsmasq.CreateStaticDhcpHost(client.StaticDhcpHost{MacAddress: "00:11:22:33:44:55", IPAddress: "1.2.3.4", HostName: "example"})
	if err != nil {
		t.Error(err)
	}
}

func teardownDhcpStaticHostDataSourceTest(*terraform.State) error {
	dnsmasq := client.New(apiUrl, "")
	_, err := dnsmasq.DeleteStaticDhcpHost("00:11:22:33:44:55")
	return err
}
