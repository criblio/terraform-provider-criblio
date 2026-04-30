resource "criblio_custom_banner" "my_custombanner" {
  enabled      = true
  message      = "This is the banner message to be displayed"
  theme        = "purple"
  type         = "custom"
  link         = "https://thisisarealwebsite.com"
  link_display = "This flavor text redirects to link"
}
