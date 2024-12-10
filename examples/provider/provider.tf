# Configuration-based authentication
provider "dnsmasq" {
  api_url   = "http://localhost:6904"
  api_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}

# Create a Static DHCP lease reservation
resource "dnsmasq_dhcp_static_host" "example" {
  mac_address = "00:11:22:33:44:55"
  ip_address  = "1.2.3.4"
  hostname    = "example"
}
