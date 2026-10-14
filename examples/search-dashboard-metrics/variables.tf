variable "dashboard_id" {
  type        = string
  description = "Unique ID for a disposable test dashboard; do not use the live GoatBot ID."
  default     = "bot_metrics_provider_test"

  validation {
    condition     = can(regex("^[a-zA-Z0-9_]+$", var.dashboard_id))
    error_message = "Dashboard IDs may contain only letters, numbers, and underscores; hyphens are not supported by the API."
  }
}

variable "dashboard_name" {
  type    = string
  default = "Bot Metrics Provider Test"
}

variable "dataset_id" {
  type    = string
  default = "metrics"
}

variable "category" {
  type        = string
  description = "Existing dashboard category ID, or null to omit it."
  default     = null
}

variable "worker_query" {
  type    = string
  default = "sum by (issue, outcome) (goatbot_run_outcome_total)"
}

variable "worker_group_by" {
  type    = list(string)
  default = ["issue", "outcome"]
}

variable "step_strategy" {
  type = object({
    type          = string
    target_points = optional(number)
    step          = optional(number)
  })
  default = { type = "auto", target_points = 240 }
}
