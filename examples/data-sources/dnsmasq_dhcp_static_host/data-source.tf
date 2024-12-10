# Retrieve the Static DHCP lease reservation for host 00:11:22:33:44:55
data "dnsmasq_dhcp_static_host" "example" {
  mac_address = "00:11:22:33:44:55"
}

output "ip_address" {
  value = data.dnsmasq_dhcp_static_host.example.ip_address
}

output "hostname" {
  value = data.dnsmasq_dhcp_static_host.example.hostname
}
