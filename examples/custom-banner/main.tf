resource "criblio_custom_banner" "my_custom_banner" {
  enabled = true
  message = "Scheduled maintenance window: Saturday 2am-4am UTC"
  theme   = "purple"
  type    = "custom"

  link         = "https://status.example.com"
  link_display = "View status page"
}

data "criblio_custom_banner" "maintenance" {
  depends_on = [criblio_custom_banner.my_custom_banner]
}

