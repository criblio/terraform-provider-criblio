resource "criblio_collector" "splunk_access_log_collector" {
  group_id = "default"
  id       = "splunk-demo-collector"
  input_collector_splunk = {
    collector = {
      type = "splunk"
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
    }
    id                      = "splunk-demo-collector"
    environment             = "demo"
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
      type = "rest"
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

        discovery = {
          discover_type = "none"
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

resource "criblio_collector" "rest_api_collector_discovery_http" {
  group_id = "default"
  id       = "rest-api-demo-collector_discovery_http"
  input_collector_rest = {
    collector = {
      type = "rest"
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
        discovery = {
          discover_type   = "http"
          discover_method = "get"
          # API validates discoverUrl as a jsExpression; use a string literal, not a bare URL
          discover_url = "'https://discover.test.example.com/v1'"
          discover_request_params = [
            {
              name  = "filter"
              value = "'active=true'"
            },
            {
              name  = "limit"
              value = "'100'"
            }
          ]
          discover_request_headers = [
            {
              name  = "X-Custom-Header"
              value = "'test-value'"
            }
          ]
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
    }
    environment             = "demo"
    id                      = "rest-api-demo-collector_discovery_http"
    ignore_group_jobs_limit = false
  }
}

resource "criblio_collector" "rest_conf_update_test" {
  group_id = "default"
  id       = "rest-conf-update-test"
  input_collector_rest = {
    collector = {
      type = "rest"
      conf = {
        authentication = "none"
        discovery      = { discover_type = "none" }
        # Change this value after initial apply to test Bug 2 fix
        timeout             = 60
        collect_method      = "get"
        collect_url         = "'https://api.test.example.com/data'"
        reject_unauthorized = false
      }
    }
    environment             = "demo"
    id                      = "rest-conf-update-test"
    ignore_group_jobs_limit = false
  }
}

resource "criblio_collector" "rest_api_collector_discovery_json" {
  group_id = "default"
  id       = "rest-api-demo-collector_discovery_json"
  input_collector_rest = {
    collector = {
      type = "rest"
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
        discovery = {
          discover_type          = "json"
          manual_discover_result = "{\"result\":\"true\"}"
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
    }
    environment             = "demo"
    id                      = "rest-api-demo-collector_discovery_json"
    ignore_group_jobs_limit = false
  }
}

resource "criblio_collector" "rest_api_collector_discovery_list" {
  group_id = "default"
  id       = "rest-api-demo-collector_discovery_list"
  input_collector_rest = {
    collector = {
      type = "rest"
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
        discovery = {
          discover_type = "list"
          item_list     = ["foo", "bar"]
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
    }
    environment             = "demo"
    id                      = "rest-api-demo-collector_discovery_list"
    ignore_group_jobs_limit = false
  }
}

# ---------------------------------------------------------------------------
# Extra REST examples (pagination, OAuth, post body). All hostnames are
# intentionally fake (test.example, *.test.example) — replace for real use.
# Cribl stores collectUrl, discoverUrl, loginUrl, and most header `value`
# fields as **JavaScript expressions**; bare `https://...` is invalid. Use
# string-literal form, e.g. 'https://host/path' in TF: "'https://host/path'"
# ---------------------------------------------------------------------------

resource "criblio_collector" "rest_okta_system_log_events" {
  group_id = "default"
  id       = "example-okta-system-log-events"
  input_collector_rest = {
    id                      = "example-okta-system-log-events"
    environment             = "demo"
    ignore_group_jobs_limit = false
    remove_fields           = []
    resume_on_boot          = false
    streamtags              = ["okta", "example"]
    ttl                     = "4h"
    worker_affinity         = false
    collector = {
      type = "rest"
      conf = {
        discovery      = { discover_type = "none" }
        collect_method = "get"
        # collectUrl is a jsExpression in the API; wrap the URL in a string literal
        collect_url         = "'https://idp.test.example.com/api/v1/system/log'"
        authentication      = "none"
        timeout             = 0
        use_round_robin_dns = false
        disable_time_filter = false
        reject_unauthorized = true
        capture_headers     = false
        safe_headers        = []
        pagination = {
          type                    = "response_header_link"
          next_relation_attribute = "next"
          cur_relation_attribute  = "self"
          max_pages               = 0
        }
        retry_rules = {
          type            = "backoff"
          interval        = 1000
          limit           = 6
          multiplier      = 2
          max_interval_ms = 30000
          codes           = [429, 503]
          enable_header   = false
        }
        # Cribl JS: use $${ in Terraform to emit a literal ${ for the API
        collect_request_params = [
          {
            name  = "since"
            value = <<-E
`$${new Date((earliest * 1000 || Date.now() - 7*24*60*60*1000)).toISOString()}`
E
          },
          {
            name  = "until"
            value = <<-E
`$${new Date((latest * 1000 || Date.now())).toISOString()}`
E
          },
        ]
        collect_request_headers = [
          {
            name = "Authorization"
            # header value is also a jsExpression
            value = "'SSWS test-api-token-replace-me'"
          },
        ]
      }
    }
    input = {
      breaker_rulesets       = ["cribl", "default"]
      metadata               = []
      output                 = "default"
      pipeline               = "main"
      preprocess             = { disabled = true }
      send_to_routes         = true
      stale_channel_flush_ms = 10000
      throttle_rate_per_sec  = "0"
      type                   = "collection"
    }
    schedule = {
      cron_schedule       = "*/5 * * * *"
      enabled             = true
      max_concurrent_runs = 1
      resume_missed       = true
      skippable           = false
      run = {
        expression               = "true"
        job_timeout              = "0"
        log_level                = "info"
        max_task_reschedule      = 1
        max_task_size            = "10MB"
        min_task_size            = "1MB"
        mode                     = "run"
        reschedule_dropped_tasks = true
        time_range_type          = "relative"
        timestamp_timezone       = "UTC"
        earliest                 = -10
        latest                   = -5
      }
    }
  }
}

# Temporarily disabled: on-prem Cribl does not support the lib/jobs API used on
# destroy (e.g. m/default/lib/jobs/example-google-workspace-admin-reports), so
# post-test terraform destroy fails in on-prem integration runs. Re-enable for
# Cloud or when local UI/API matches.
# resource "criblio_collector" "rest_google_workspace_admin_reports" {
#   group_id = "default"
#   id       = "example-google-workspace-admin-reports"
#   input_collector_rest = {
#     id                      = "example-google-workspace-admin-reports"
#     environment             = "demo"
#     ignore_group_jobs_limit = false
#     remove_fields           = []
#     resume_on_boot          = false
#     streamtags              = ["google", "gws", "example"]
#     ttl                     = "4h"
#     worker_affinity         = false
#     collector = {
#       type = "rest"
#       conf = {
#         discovery = {
#           discover_type = "list"
#           item_list     = ["admin", "login", "drive", "token"]
#         }
#         collect_method = "get"
#         # Cribl discovery id: literal ${id} in the stored URL
#         collect_url                 = "'https://admin.reports.test.example.com/v1/activity/users/all/applications/$${id}'"
#         authentication              = "google_oauth"
#         subject                     = "admin-audit@example.com"
#         service_account_credentials = "eyJ0eXBlIjoidGVzdCJ9" # base64 JSON placeholder; replace in real use
#         scopes = [
#           "https://test.example.com/oauth/scope/reports.readonly", # test scope; use real product scopes in production
#         ]
#         timeout             = 0
#         use_round_robin_dns = false
#         disable_time_filter = false
#         reject_unauthorized = true
#         capture_headers     = false
#         safe_headers        = []
#         decode_url          = true
#         pagination = {
#           type           = "response_body"
#           max_pages      = 50
#           attribute      = ["nextPageToken"]
#           last_page_expr = "nextPageToken === null"
#         }
#         retry_rules = {
#           type              = "backoff"
#           interval          = 1000
#           limit             = 5
#           multiplier        = 2
#           max_interval_ms   = 20000
#           codes             = [429, 503]
#           enable_header     = true
#           retry_header_name = "retry-after"
#         }
#         collect_request_params = [
#           { name = "maxResults", value = "100" },
#           {
#             name  = "startTime"
#             value = <<-E
# `$${new Date((earliest * 1000 || Date.now() - 7*24*60*60*1000)).toISOString()}`
# E
#           },
#           {
#             name  = "endTime"
#             value = <<-E
# `$${new Date((latest * 1000 || Date.now())).toISOString()}`
# E
#           },
#           {
#             name  = "pageToken"
#             value = <<-E
# `$${__e && __e.nextPageToken != null ? __e.nextPageToken : undefined}`
# E
#           },
#         ]
#       }
#     }
#     input = {
#       breaker_rulesets       = ["cribl", "default"]
#       metadata               = []
#       output                 = "default"
#       pipeline               = "main"
#       preprocess             = { disabled = true }
#       send_to_routes         = true
#       stale_channel_flush_ms = 10000
#       throttle_rate_per_sec  = "0"
#       type                   = "collection"
#     }
#     schedule = {
#       cron_schedule       = "*/30 * * * *"
#       enabled             = true
#       max_concurrent_runs = 1
#       resume_missed       = true
#       skippable           = false
#       run = {
#         expression               = "true"
#         job_timeout              = "60m"
#         log_level                = "info"
#         max_task_reschedule      = 1
#         max_task_size            = "10MB"
#         min_task_size            = "1MB"
#         mode                     = "run"
#         reschedule_dropped_tasks = true
#         time_range_type          = "relative"
#         earliest                 = -35
#         latest                   = -5
#         state_tracking = {
#           enabled                 = true
#           state_update_expression = "__timestampExtracted !== false && {latestTime: (state.latestTime || 0) > _time ? state.latestTime : _time}"
#           state_merge_expression  = "(prevState.latestTime || 0) > newState.latestTime ? prevState : newState"
#         }
#       }
#     }
#   }
# }

resource "criblio_collector" "rest_crowdstrike_combined_alerts" {
  group_id = "default"
  id       = "example-crowdstrike-combined-alerts"
  input_collector_rest = {
    id                      = "example-crowdstrike-combined-alerts"
    environment             = "demo"
    ignore_group_jobs_limit = false
    remove_fields           = []
    resume_on_boot          = false
    streamtags              = ["crowdstrike", "example"]
    ttl                     = "4h"
    worker_affinity         = false
    collector = {
      type = "rest"
      conf = {
        discovery      = { discover_type = "none" }
        collect_method = "post_with_body"
        # For production, use a backtick-wrapped Cribl expression; static JSON is enough for local validation
        collect_url          = "'https://falcon.test.example.com/alerts/combined/alerts/v1'"
        collect_body         = "'{\"filter\":\"\",\"sort\":\"created_timestamp.asc\",\"limit\":1000}'"
        authentication       = "login"
        username             = "cs-client-id"
        password             = "cs-client-secret"
        login_url            = "'https://falcon.test.example.com/oauth2/token'"
        login_body           = join("", ["`", "client_id=$${username}", "&client_secret=$${password}", "`"])
        token_resp_attribute = "access_token"
        auth_header_key      = "Authorization"
        auth_header_expr     = join("", ["`", "Bearer $${token}", "`"])
        auth_request_headers = [
          { name = "Content-Type", value = "application/x-www-form-urlencoded" },
          { name = "accept", value = "application/json" },
        ]
        collect_request_headers = [
          { name = "Content-Type", value = "application/json" },
          { name = "Accept", value = "application/json" },
        ]
        timeout             = 0
        use_round_robin_dns = false
        reject_unauthorized = false
        capture_headers     = false
        safe_headers        = []
        pagination = {
          type           = "response_body"
          max_pages      = 250
          attribute      = ["meta", "pagination", "after"]
          last_page_expr = "meta && meta.pagination && (meta.pagination.after === null || meta.pagination.after === undefined || meta.pagination.after === '')"
        }
        retry_rules = {
          type            = "backoff"
          interval        = 1000
          limit           = 5
          multiplier      = 2
          max_interval_ms = 20000
          codes           = [429, 503]
          enable_header   = false
        }
      }
    }
    input = {
      breaker_rulesets       = ["cribl", "default"]
      metadata               = []
      output                 = "default"
      pipeline               = "main"
      preprocess             = { disabled = true }
      send_to_routes         = true
      stale_channel_flush_ms = 10000
      throttle_rate_per_sec  = "0"
      type                   = "collection"
    }
    schedule = {
      cron_schedule       = "*/5 * * * *"
      enabled             = true
      max_concurrent_runs = 1
      resume_missed       = true
      skippable           = false
      run = {
        expression               = "true"
        job_timeout              = "0"
        log_level                = "info"
        max_task_reschedule      = 1
        max_task_size            = "10MB"
        min_task_size            = "1MB"
        mode                     = "run"
        reschedule_dropped_tasks = true
        time_range_type          = "relative"
        timestamp_timezone       = "UTC"
        earliest                 = -10
        latest                   = -5
      }
    }
  }
}

resource "criblio_collector" "rest_mode_audit_logs" {
  group_id = "default"
  id       = "example-mode-audit-logs"
  input_collector_rest = {
    id                      = "example-mode-audit-logs"
    environment             = "demo"
    ignore_group_jobs_limit = false
    remove_fields           = []
    resume_on_boot          = false
    streamtags              = ["mode", "example"]
    ttl                     = "4h"
    worker_affinity         = false
    collector = {
      type = "rest"
      conf = {
        discovery           = { discover_type = "none" }
        collect_method      = "get"
        collect_url         = "'https://analytics.test.example.com/api/v1/audit_logs'"
        authentication      = "none"
        timeout             = 0
        use_round_robin_dns = false
        disable_time_filter = true
        reject_unauthorized = false
        capture_headers     = false
        safe_headers        = []
        pagination = {
          type           = "response_body"
          max_pages      = 1000
          attribute      = ["metadata", "next_token"]
          last_page_expr = "metadata && (metadata.next_token === null || metadata.next_token === undefined || metadata.next_token === '')"
        }
        retry_rules = {
          type            = "backoff"
          interval        = 1000
          limit           = 5
          multiplier      = 2
          max_interval_ms = 20000
          codes           = [429, 503]
          enable_header   = false
        }
        collect_request_params = [
          {
            name  = "start_timestamp"
            value = <<-E
`$${new Date(Math.floor(Date.now() / 60000) * 60000 - 10 * 60 * 1000).toISOString()}`
E
          },
          {
            name  = "end_timestamp"
            value = <<-E
`$${new Date(Math.floor(Date.now() / 60000) * 60000 - 5 * 60 * 1000).toISOString()}`
E
          },
          {
            name  = "next_token"
            value = <<-E
`$${__e && __e.metadata && __e.metadata.next_token ? __e.metadata.next_token : undefined}`
E
          },
        ]
        collect_request_headers = [
          { name = "Authorization", value = "'Basic dGVzdDp0ZXN0'" },
          { name = "Content-Type", value = "'application/json'" },
          { name = "Accept", value = "'application/json'" },
        ]
      }
    }
    input = {
      breaker_rulesets       = ["cribl", "default"]
      metadata               = []
      output                 = "default"
      pipeline               = "main"
      preprocess             = { disabled = true }
      send_to_routes         = true
      stale_channel_flush_ms = 10000
      throttle_rate_per_sec  = "0"
      type                   = "collection"
    }
    schedule = {
      cron_schedule       = "*/5 * * * *"
      enabled             = true
      max_concurrent_runs = 1
      resume_missed       = true
      skippable           = false
      run = {
        expression               = "true"
        job_timeout              = "0"
        log_level                = "info"
        max_task_reschedule      = 1
        max_task_size            = "10MB"
        min_task_size            = "1MB"
        mode                     = "run"
        reschedule_dropped_tasks = true
        time_range_type          = "relative"
        timestamp_timezone       = "UTC"
        earliest                 = -10
        latest                   = -5
      }
    }
  }
}

/*
# Not a valid collector type for on-prem, so commented out.

resource "criblio_collector" "cribl_lake" {
  group_id = "default"
  id       = "cribl_logs_lake"
  input_collector_cribl_lake = {
    collector = {
      conf = {
        dataset = "cribl_logs"
      }
      type = "cribl_lake"
    }
    id                      = "cribl_logs_lake"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "Cribl Ruleset",
      ]
      metadata = [
        {
          name  = "__replayed"
          value = "true"
        },
      ]
      send_to_routes         = true
      stale_channel_flush_ms = 10000
      throttle_rate_per_sec  = "0"
      type                   = "collection"
    }
    remove_fields   = []
    resume_on_boot  = true
    streamtags      = []
    ttl             = "4h"
    worker_affinity = false
  }
}
*/

/*
# Script collector: enable when your environment supports it
resource "criblio_collector" "script_collector" {
  group_id = "default"
  id       = "script-demo-collector"
  input_collector_script = {
    collector = {
      type = "script"
      conf = {
        collect_script  = "echo 1"
        discover_script = "echo 1"
        shell           = "/bin/bash"
      }
      destructive = false
    }
    "id"                      = "script-demo-collector"
    "ignore_group_jobs_limit" = false
    "input" = {
      metadata = [
        {
          name  = "saas_domain"
          value = "'cribl.cloud'"
        },
        {
          name  = "type"
          value = "'organizations'"
        },
      ]
      output = "devnull"
      preprocess = {
        disabled = true
      }
      send_to_routes         = false
      stale_channel_flush_ms = 10000
      throttle_rate_per_sec  = "0"
      type                   = "collection"
      pipeline               = "main"
    }
    "resume_on_boot" = false
    "schedule" = {
      cron_schedule       = "*\/4 * * * *"
      enabled             = true
      max_concurrent_runs = 1
      resume_missed       = true
      run = {
        expression               = "true"
        job_timeout              = "0"
        log_level                = "info"
        max_task_reschedule      = 1
        max_task_size            = "10MB"
        min_task_size            = "1MB"
        mode                     = "run"
        reschedule_dropped_tasks = true
        time_range_type          = "relative"
        timestamp_timezone       = "UTC"
      }
      skippable = false
    }
    "ttl"             = "4h"
    "type"            = "collection"
    "worker_affinity" = false
  }
}
*/

