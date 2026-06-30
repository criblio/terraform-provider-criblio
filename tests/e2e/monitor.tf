resource "criblio_monitor" "demo" {
  id      = "tf-demo-monitor"
  name    = "Terraform Demo Monitor"
  enabled = true
  product = "aetos"

  # Firing / recovery timing (seconds)
  firing_after             = 300
  ok_after                 = 60
  schedule_interval_seconds = 60

  # Plain PromQL string (simplest form of query)
  query = jsonencode("up")

  # Rules and params as JSON — use jsonencode() for structured values
  rules  = jsonencode([])
  params = jsonencode({})
}

data "criblio_monitor" "demo" {
  id = criblio_monitor.demo.id

  depends_on = [criblio_monitor.demo]
}

output "monitor_id" {
  value = criblio_monitor.demo.id
}

output "monitor_managed_by" {
  description = "Stamped by backend when speakeasy-sdk/terraform User-Agent is detected"
  value       = data.criblio_monitor.demo.id
}
