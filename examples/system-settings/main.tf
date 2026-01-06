data "criblio_system_info" "my_systeminfo" {
}

output "system_settings" {
  value = data.criblio_system_info.my_systeminfo
}

output "version" {
  value = jsondecode(data.criblio_system_info.my_systeminfo.items[0].build["VERSION"])
}
