resource "criblio_monitor" "example" {
  id         = "cpu-high"
  name       = "CPU High"
  dataset_id = "metrics"
  managed_by = "terraform"
  silence    = []
  enabled    = true
  type       = "threshold"

  query = {
    A = {
      mode   = "promql"
      promql = "sum(rate(cpu_usage[5m]))"
    }
  }

  monitor_policy_notification_config = {
    type    = "policy"
    enabled = false
    config  = []
  }

  expr = []

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

  priority_scalar_value = {
    value = "P2"
  }

  team_scalar_value = {
    value = "ops"
  }
}

resource "criblio_monitor" "example_outlier" {
  id         = "latency-peer-outlier"
  name       = "Latency Peer Outlier"
  dataset_id = "metrics"
  managed_by = "terraform"
  silence    = []
  enabled    = true
  type       = "outlier"

  query = {
    A = {
      mode   = "promql"
      promql = "avg by (pod) (http_request_duration_seconds)"
    }
  }

  monitor_policy_notification_config = {
    type    = "policy"
    enabled = false
    config  = []
  }

  expr = []

  # Outlier creation is currently rejected by the backend.
  # These are the fields defined by the published upstream schema.
  outlier_config = {
    algorithm = "mad"
  }

  firing_condition = {
    fire_delay  = 300
    clear_delay = 60
  }

  firing_rule = {
    label = "A"
    threshold = [{
      severity      = "warning"
      limit         = 2
      operator      = "gt"
      included_tags = []
      excluded_tags = []
    }]
  }

  aetos_metadata = {}

  priority_scalar_value = {
    value = "P2"
  }

  team_scalar_value = {
    value = "ops"
  }
}

resource "criblio_monitor" "example_target_notification" {
  id         = "cpu-high-notify"
  name       = "CPU High (direct target notification)"
  dataset_id = "metrics"
  managed_by = "terraform"
  silence    = []
  enabled    = true
  type       = "threshold"

  query = {
    A = {
      mode   = "promql"
      promql = "sum(rate(cpu_usage[5m]))"
    }
  }

  expr = []

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

  priority_scalar_value = {
    value = "P2"
  }

  team_scalar_value = {
    value = "ops"
  }

  # Direct delivery: firings route straight to notification target/template pairs.
  monitor_target_notification_config = {
    enabled = true
    config = [{
      target_id   = "nt-sns-prod"
      template_id = "default-template"
    }]
  }
}

resource "criblio_monitor" "example_policy_notification" {
  id         = "latency-high-notify"
  name       = "Latency High (policy routing)"
  dataset_id = "metrics"
  managed_by = "terraform"
  silence    = []
  enabled    = true
  type       = "threshold"

  query = {
    A = {
      mode   = "promql"
      promql = "histogram_quantile(0.99, sum by (le) (rate(http_request_duration_seconds_bucket[5m])))"
    }
  }

  expr = []

  firing_condition = {
    fire_delay  = 60
    clear_delay = 60
  }

  firing_rule = {
    label = "A"
    threshold = [{
      severity      = "warning"
      limit         = 1
      included_tags = []
      excluded_tags = []
    }]
  }

  aetos_metadata = {}

  priority_scalar_value = {
    value = "P2"
  }

  team_scalar_value = {
    value = "ops"
  }

  # Policy routing: ordered policies decide which target/template pairs receive each firing.
  monitor_policy_notification_config = {
    type    = "policy"
    enabled = true
    config = [{
      id         = "critical-routing"
      order      = 1
      conditions = jsonencode([])
      template_target_pairs = [{
        target_id   = "nt-sns-prod"
        template_id = "default-template"
      }]
    }]
  }
}
