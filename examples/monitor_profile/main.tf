terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

# Configure authentication using the CRIBL_* environment variables.
provider "criblio" {}

resource "criblio_monitor_shared_config" "priority" {
  id     = "profile-priority"
  name   = "Profile priority"
  config = jsonencode({ priority = { value = "P3" } })
}

resource "criblio_monitor_profile" "example" {
  id          = "platform-profile"
  name        = "Platform profile"
  description = "Reusable platform monitor defaults"
  config_ids  = [criblio_monitor_shared_config.priority.id]
}
