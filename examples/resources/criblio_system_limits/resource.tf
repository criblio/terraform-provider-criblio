# Manage the lookup allowlist on the two Edge fleets from PLAT-17435.
# The provider reads existing limits and preserves unconfigured settings.
resource "criblio_system_limits" "lookup_paths" {
  for_each = toset(["stream-leaders", "stream-leaders-with-insights"])

  group_id             = each.value
  lookup_allowed_paths = ["/tmp"]
}
