resource "criblio_pack" "my_pack" {
  id           = "EdgeSupervisorHealth"
  group_id     = "default_fleet"
  description  = "EdgeSupervisorHealth - Monitors Supervisor process health on Search Deux instances"
  disabled     = false
  display_name = "EdgeSupervisorHealth"
  version      = "1.0.0"
  author       = "Observability Team"
}
