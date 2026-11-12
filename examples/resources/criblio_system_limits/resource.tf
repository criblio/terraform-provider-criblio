# The provider reads existing limits and preserves unconfigured settings.
resource "criblio_system_limits" "default" {
  group_id       = "default"
  min_free_space = "2 GB"
}
