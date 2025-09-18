data "criblio_instance_settings" "my_instancesettings" {
}

output "instance_settings" {
  value = data.criblio_instance_settings.my_instancesettings
}
