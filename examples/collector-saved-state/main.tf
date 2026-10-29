terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}


variable "group_id" {
  type    = string
  default = "default"
}

variable "collector_id" {
  type    = string
  default = "terraform-saved-state-test"
}

variable "connection_id" {
  description = "An existing Cribl database connection in the selected worker group."
  type        = string
  default     = "terraform-saved-state-test"
}

variable "query" {
  description = "Plain SQL query; the example quotes it as a JavaScript string for Cribl. Its result must include the tracking column."
  type        = string
  default     = "SELECT id, name, email FROM users"
}

variable "tracking_column" {
  type    = string
  default = "id"
}

variable "saved_state" {
  description = "Nested checkpoint, legacy string map, or null to read runtime state from Cribl."
  type        = any
  default = {
    databaseCollectorState = {
      data = {
        id = {
          value = 3
        }
      }
    }
  }
}

resource "criblio_collector" "saved_state_test" {
  group_id = var.group_id
  id       = var.collector_id

  input_collector_database = {
    description = "Terraform saved_state compatibility test"
    collector = {
      type = "database"
      conf = {
        connection_id = var.connection_id
        # Cribl expects a JavaScript expression returning SQL.
        query = jsonencode(var.query)
      }
    }
    input = {
      type           = "collection"
      send_to_routes = true
    }
    saved_state = var.saved_state
    schedule = {
      # Run manually from Cribl when you want to generate a runtime checkpoint.
      enabled       = false
      cron_schedule = "0 * * * *"
      run = {
        mode        = "run"
        expression  = "true"
        job_timeout = "5m"
        state_tracking = {
          enabled         = true
          tracking_column = var.tracking_column
        }
      }
    }
  }
}

# Read the API checkpoint separately: configured resource state is preserved,
# whereas this data source shows what Cribl actually returns.
data "criblio_collector" "api" {
  group_id = criblio_collector.saved_state_test.group_id
  id       = criblio_collector.saved_state_test.id
}

output "resource_saved_state" {
  value = criblio_collector.saved_state_test.input_collector_database.saved_state
}

output "api_saved_state" {
  value = data.criblio_collector.api.input_collector_database.saved_state
}
