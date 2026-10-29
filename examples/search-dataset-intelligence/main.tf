terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

variable "dataset_id" {
  type        = string
  description = "ID of an existing Search dataset."
  default     = "cribl_edge_logs"
}

# AI and agentic Search must be enabled, with AI consent accepted.
resource "criblio_search_dataset_intelligence" "example" {
  dataset_id = var.dataset_id
}

output "intelligence_status" {
  value = criblio_search_dataset_intelligence.example.status
}

output "intelligence_content" {
  value = criblio_search_dataset_intelligence.example.content
}
