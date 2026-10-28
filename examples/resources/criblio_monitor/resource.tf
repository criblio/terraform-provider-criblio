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
  enabled    = true
  type       = "outlier"

  query = {
    A = {
      mode   = "promql"
      promql = "avg by (pod) (http_request_duration_seconds)"
    }
  }

  expr = []

  # Peer-relative (pooled MAD): flag pods deviating from their service cohort. algorithm defaults to "mad".
  outlier_config = {
    mode                 = "instant"
    cohort_labels        = ["service"]
    threshold_direction  = "either"
    evaluation_window_ms = 1800000
    pct                  = 20
    min_coverage         = 80
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
    enabled = true
    config = [{
      id                    = "critical-routing"
      order                 = 1
      conditions            = []
      template_target_pairs = [{
        target_id   = "nt-sns-prod"
        template_id = "default-template"
      }]
    }]
  }
}
