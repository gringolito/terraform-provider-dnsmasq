# Create a Static DHCP lease reservation
resource "dnsmasq_dhcp_static_host" "example" {
  mac_address = "00:11:22:33:44:55"
  ip_address  = "1.2.3.4"
  hostname    = "example"
}
