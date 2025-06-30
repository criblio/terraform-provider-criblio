terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Configuration options
}

data "criblio_config_version" "my_configversion" {
  id         = "syslog-workers"
  depends_on = [criblio_commit.my_commit]
}

resource "criblio_commit" "my_commit" {
  effective = true
  group     = "syslog-workers"
  message   = "test"
}

resource "criblio_deploy" "my_deploy" {
  id      = "syslog-workers"
  version = data.criblio_config_version.my_configversion.items[0]
}

output "deploy" {
  value = criblio_deploy.my_deploy
}

output "config_version" {
  value = data.criblio_config_version.my_configversion
}

output "commit" {
  value = criblio_commit.my_commit
}
