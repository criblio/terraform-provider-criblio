
resource "criblio_collector" "splunk_access_log_collector" {
  group_id = "default"
  id       = "splunk-demo-collector"
  input_collector_splunk = {
    collector = {
      conf = {
        authentication      = "token"
        disable_time_filter = false
        earliest            = "-24h@h"
        endpoint            = "https://demo.splunk.example.com:8089"
        latest              = "now"
        output_mode         = "json"
        reject_unauthorized = false
        search              = "index=main earliest=-1h@h"
        search_head         = "https://demo.splunk.example.com:8000"
        timeout             = 300
        token               = "your-splunk-token-here"
        use_round_robin_dns = false
        username            = "cribl-user"
      }
      type = "splunk"
    }
    environment             = "demo"
    id                      = "splunk-demo-collector"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "cribl",
        "default"
      ]
      metadata = [
        {
          name  = "source_type"
          value = "splunk_search"
        },
        {
          name  = "environment"
          value = "demo"
        }
      ]
      output   = "default"
      pipeline = "main"
      preprocess = {
        disabled = true
      }
      send_to_routes         = false
      stale_channel_flush_ms = 300000
      throttle_rate_per_sec  = "0"
      type                   = "collection"
    }
    remove_fields = [
      "_internal",
      "_raw"
    ]
    resume_on_boot = true
    schedule = {
      cron_schedule       = "0 */2 * * *"
      enabled             = true
      max_concurrent_runs = 1
      run = {
        earliest                 = 0
        expression               = "index=main earliest=-1h@h"
        job_timeout              = "30m"
        latest                   = 1
        log_level                = "info"
        max_task_reschedule      = 3
        max_task_size            = "1GB"
        min_task_size            = "1MB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        time_range_type          = "relative"
      }
      skippable = false
    }
    streamtags = [
      "splunk",
      "demo",
      "collector"
    ]
    ttl             = "4h"
    worker_affinity = false
  }
}

resource "criblio_collector" "rest_api_collector" {
  group_id = "default"
  id       = "rest-api-demo-collector"
  input_collector_rest = {
    collector = {
      conf = {
        # Mandatory REST collector parameters
        base_url     = "https://api.demo.example.com"
        auth_type    = "manual"
        bearer_token = ""
        oauth_config = {
          client_id     = ""
          client_secret = ""
          token_url     = ""
          scope         = ""
        }

        # Additional REST configuration
        authentication = "basic"
        body           = ""
        headers = {
          "Content-Type" = "application/json"
          "Accept"       = "application/json"
        }
        method              = "GET"
        password            = "demo-password"
        path                = "/api/v1/logs"
        reject_unauthorized = false
        timeout             = 30
        url                 = "https://api.demo.example.com/api/v1/logs"
        username            = "demo-user"
        collect_method      = "get"
        collect_url         = "api.demo.example.com/api/v1/logs"
      }
      type = "rest"
    }
    environment             = "demo"
    id                      = "rest-api-demo-collector"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "cribl",
        "default"
      ]
      metadata = [
        {
          name  = "source_type"
          value = "rest_api"
        },
        {
          name  = "api_endpoint"
          value = "api_v1_logs"
        }
      ]
      output   = "default"
      pipeline = "main"
      preprocess = {
        disabled = false
      }
      send_to_routes         = false
      stale_channel_flush_ms = 300000
      throttle_rate_per_sec  = "10"
      type                   = "collection"
    }
    remove_fields = [
      "_metadata",
      "_headers"
    ]
    resume_on_boot = true
    schedule = {
      cron_schedule       = "*/15 * * * *" # Every 15 minutes
      enabled             = true
      max_concurrent_runs = 1
      run = {
        earliest                 = 0
        expression               = ""
        job_timeout              = "10m"
        latest                   = 1
        log_level                = "info"
        max_task_reschedule      = 2
        max_task_size            = "500MB"
        min_task_size            = "1KB"
        mode                     = "list"
        reschedule_dropped_tasks = false
        time_range_type          = "relative"
      }
      skippable = true
    }
    streamtags = [
      "rest",
      "api",
      "demo",
      "collector"
    ]
    ttl             = "2h"
    worker_affinity = false
  }
}