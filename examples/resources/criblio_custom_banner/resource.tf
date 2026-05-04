resource "criblio_custom_banner" "my_custombanner" {
  created = 1759154100
  custom_themes = [
    "themes",
  ]
  enabled           = true
  id                = "myUniqueBannerMessageId"
  invert_font_color = false
  link              = "https://thisisarealwebsite.com"
  link_display      = "This flavor text redirects to link"
  message           = "This is the banner message to be displayed"
  theme             = "purple"
  type              = "custom"
}