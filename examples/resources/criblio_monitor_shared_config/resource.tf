resource "criblio_monitor_shared_config" "example" {
  id   = "platform-priority"
  name = "Platform priority"
  config = jsonencode({
    priority = { value = "P3" }
  })
}
