# Example 1: Search Saved Query from JSON file
module "saved_query_from_file" {
  source = "../../modules/json-config-search-saved-query"

  json_file_path = "./configs/saved_query.json"
}

# Example 2: Search Saved Query with inline JSON and overrides
module "saved_query_inline" {
  source = "../../modules/json-config-search-saved-query"

  json_content = jsonencode({
    id          = "inline_saved_query"
    name        = "Inline Query Example"
    description = "Created from inline JSON"
    query       = "dataset=\"cribl_edge_appscope_metrics\" | summarize count=count() by host"
    earliest    = "-24h"
    latest      = "now"
    isPrivate   = false
  })

  additional_config = {
    description = "Modified description via additional_config"
  }
}

# Example 3: Search Dashboard from JSON file
module "dashboard_from_file" {
  source = "../../modules/json-config-search-dashboard"

  json_file_path = "./configs/dashboard.json"
}

# Outputs showing all created search resources
output "search_resources_created" {
  value = {
    saved_queries = [
      {
        id   = module.saved_query_from_file.query_id
        name = module.saved_query_from_file.query_name
      },
      {
        id   = module.saved_query_inline.query_id
        name = module.saved_query_inline.query_name
      }
    ]
    dashboards = [
      {
        id   = module.dashboard_from_file.dashboard_id
        name = module.dashboard_from_file.dashboard_name
      }
    ]
  }
}
