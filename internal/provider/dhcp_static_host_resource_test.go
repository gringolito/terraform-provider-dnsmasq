package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDhcpStaticHostResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpStaticHostResourceConfig("00:11:22:33:44:55", "1.2.3.4", "example"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dnsmasq_dhcp_static_host.test", "mac_address", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("dnsmasq_dhcp_static_host.test", "ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr("dnsmasq_dhcp_static_host.test", "hostname", "example"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "dnsmasq_dhcp_static_host.test",
				ImportState:       true,
				ImportStateId:     "00:11:22:33:44:55",
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpStaticHostResourceConfig("00:11:22:33:44:55", "10.20.30.40", "new-example"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dnsmasq_dhcp_static_host.test", "mac_address", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("dnsmasq_dhcp_static_host.test", "ip_address", "10.20.30.40"),
					resource.TestCheckResourceAttr("dnsmasq_dhcp_static_host.test", "hostname", "new-example"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpStaticHostResourceConfig(macAddress string, ipAddress string, hostName string) string {
	return providerConfig + fmt.Sprintf(`
resource "dnsmasq_dhcp_static_host" "test" {
  mac_address = %q
  ip_address = %q
  hostname = %q
}
`, macAddress, ipAddress, hostName)
}
