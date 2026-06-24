data "criblio_config_version" "default" {
  id = "default"
}

resource "criblio_deploy" "my_deploy" {
  id      = "default"
  version = data.criblio_config_version.default.items[0]
}
