resource "criblio_workspace_system_settings" "main" {
  upgrade_settings = {
    disable_automatic_upgrade = true
  }
}

data "criblio_workspace_system_settings" "main" {
  depends_on = [criblio_workspace_system_settings.main]
}
