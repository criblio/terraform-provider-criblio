resource "criblio_workspace_system_settings" "example" {
  upgrade_settings = {
    disable_automatic_upgrade = true
  }
}

# To enable automatic upgrades, use disable_automatic_upgrade = false and
# automatic_upgrade_check_period = "24h" (or another supported duration).
