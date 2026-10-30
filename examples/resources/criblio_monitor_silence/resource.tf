# Use an existing metrics dataset and a PromQL expression available in it.
resource "criblio_monitor" "cpu_high" {
  id         = "silence-example-cpu-high"
  name       = "CPU High"
  dataset_id = "metrics"
  enabled    = true
  type       = "threshold"
  managed_by = "terraform"
  silence    = []

  monitor_policy_notification_config = {
    type    = "policy"
    enabled = false
    config  = []
  }

  query = {
    A = {
      mode       = "promql"
      dataset_id = "metrics"
      promql     = "cpu_usage_percent"
    }
  }

  expr = []

  firing_condition = {
    fire_delay  = 300
    clear_delay = 60
  }

  firing_rule = {
    label = "A"
    threshold = [{
      severity        = "warning"
      limit           = 80
      operator        = "gt"
      included_tags   = []
      excluded_tags   = []
      times_triggered = 1
    }]
  }

  aetos_metadata        = {}
  priority_scalar_value = { value = "P3" }
  team_scalar_value     = { value = "" }
}

resource "criblio_monitor_silence" "example" {
  id          = "maintenance-window"
  monitor_id  = criblio_monitor.cpu_high.id
  name        = "Maintenance window"
  description = "Mute host alerts during maintenance"
  # Set a fixed current window before applying; values are Unix milliseconds.
  # Future windows can be excluded by the backend list endpoint.
  begins_at      = 1793325780159 # 2026-10-30 02:03:00.159 UTC.
  expires_at     = 1793329380159 # One hour after begins_at.
  duration_type  = "relative"
  duration_value = 1
  duration_unit  = "hours"
  tags = jsonencode({
    host = ["worker-01", "worker-02"]
  })
}
