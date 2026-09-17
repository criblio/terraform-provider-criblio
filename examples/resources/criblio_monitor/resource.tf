resource "criblio_monitor" "example" {
  id         = "cpu-high"
  name       = "CPU High"
  dataset_id = "metrics"
  enabled    = true
  type       = "threshold"

  query = {
    A = {
      mode   = "promql"
      promql = "sum(rate(cpu_usage[5m]))"
    }
  }

  expr    = []
  silence = []

  firing_condition = {
    fire_delay  = 60
    clear_delay = 60
  }

  firing_rule = {
    label = "A"
    threshold = [{
      severity      = "critical"
      limit         = 90
      included_tags = []
      excluded_tags = []
    }]
  }

  aetos_metadata = {}

  monitor_policy_notification_config = {
    enabled = false
    type    = "policy"
    config = [{
      id                    = "default-policy"
      order                 = 0
      conditions            = []
      template_target_pairs = []
    }]
  }

  priority_scalar_value = {
    value = "P2"
  }

  team_scalar_value = {
    value = "ops"
  }
}
