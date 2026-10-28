resource "criblio_monitor_config_profile" "example" {
  id          = "team-sre-defaults"
  name        = "Team SRE Defaults"
  description = "Shared monitor defaults inherited by SRE monitors"

  # IDs of criblio_shared_config resources this profile composes (empty = none yet).
  config_ids = []
}
