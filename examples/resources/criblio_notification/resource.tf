resource "criblio_notification" "my_notification" {
  condition = "high-volume"
  group     = "default"
  id        = "high-volume-http-src"
  conf = {
    name                 = "in_http"
    time_window          = "60s"
    data_volume          = "1GB"
    notify_on_resolution = true
  }

  mode = "direct"
  template_target_pairs = [{
    target_id   = "system_email"
    template_id = "default-email"
  }]
  targets = ["system_email"]
}
