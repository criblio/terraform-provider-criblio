terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

# Configure authentication using the CRIBL_* environment variables.
provider "criblio" {}

resource "criblio_monitor_shared_config" "example" {
  id   = "platform-priority"
  name = "Platform priority"
  config = jsonencode({
    priority = { value = "P3" }
  })
}
