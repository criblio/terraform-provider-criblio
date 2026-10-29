terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}


variable "monitor_id" {
  type        = string
  description = "Unique ID for the monitor."
  default     = "cpu-high-example"
}

variable "dataset_id" {
  type        = string
  description = "ID of an existing metrics dataset."
  default     = "metrics"
}

variable "promql" {
  type        = string
  description = "PromQL expression returning CPU usage as a percentage."
  default     = "cpu_usage_percent"
}

variable "threshold" {
  type        = number
  description = "CPU usage percentage above which the monitor fires."
  default     = 80
}

# Configure the provider using your usual Cribl credentials and workspace.
resource "criblio_monitor" "cpu_high" {
  id         = var.monitor_id
  name       = "CPU High"
  dataset_id = var.dataset_id
  enabled    = true
  type       = "threshold"

  query = {
    A = {
      mode       = "promql"
      dataset_id = var.dataset_id
      promql     = var.promql
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
      limit           = var.threshold
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

data "criblio_monitor" "cpu_high" {
  id = criblio_monitor.cpu_high.id
}

output "monitor_id" {
  value = data.criblio_monitor.cpu_high.id
}
