resource "criblio_collector" "my_collector" {
  group_id = "default"
  id       = "collector-2"
  input_collector_azure_blob = {
    collector = {
      conf = {
        auth_type         = "manual"
        connection_string = "DefaultEndpointsProtocol=https;AccountName=mystorageaccount;AccountKey=abc123;EndpointSuffix=core.windows.net"
        container_name    = "my-container"
        extractors = [
          {
            expression = "__input.azureBlob.key.split('/')[1]"
            key        = "account_id"
            pattern    = ".*.json"
          }
        ]
        max_batch_size       = 500
        path                 = "container/logs/2025/10/"
        recurse              = true
        storage_account_name = "mystorageaccount"
      }
      type = "azureblob"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_cribl_lake = {
    collector = {
      conf = {
        dataset = "my-dataset"
      }
      type = "cribllake"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_database = {
    collector = {
      conf = {
        connection_id            = "myDatabaseConnection"
        query                    = "SELECT * FROM logs WHERE severity = 'ERROR'"
        query_validation_enabled = true
      }
      type = "database"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_gcs = {
    collector = {
      conf = {
        auth_type = "manual"
        bucket    = "my-gcs-bucket"
        extractors = [
          {
            expression = "__input.gcs.key.split('/')[1]"
            key        = "account_id"
            pattern    = ".*.csv"
          }
        ]
        max_batch_size              = 200
        path                        = "logs/2025/10/"
        recurse                     = true
        service_account_credentials = "***REDACTED***"
      }
      type = "gcs"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_health_check = {
    collector = {
      conf = {
        authentication      = "none"
        collect_method      = "get"
        collect_url         = "https://api.example.com/health"
        credentials_secret  = "healthCredSecret"
        password            = "healthPassword"
        reject_unauthorized = true
        timeout             = 10
        username            = "healthUser"
      }
      type = "healthcheck"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_rest = {
    collector = {
      conf = {
        auth_header_expr = "Bearer ${token}"
        auth_header_key  = "Authorization"
        auth_request_headers = [
          {
            name  = "Content-Type"
            value = "application/json"
          }
        ]
        auth_request_params = [
          {
            name  = "param1"
            value = "value1"
          }
        ]
        authentication           = "basic"
        capture_headers          = true
        client_secret_param_name = "client_secret"
        collect_method           = "get"
        collect_request_headers = [
          {
            name  = "Accept"
            value = "application/json"
          }
        ]
        collect_request_params = [
          {
            name  = "limit"
            value = "100"
          }
        ]
        collect_url         = "https://api.example.com/data"
        credentials_secret  = "restCredentialsSecret"
        decode_url          = true
        disable_time_filter = false
        discovery = {
          discover_body       = "{\"discover\":true}"
          discover_data_field = "results"
          discover_method     = "get"
          discover_request_headers = [
            {
              name  = "Authorization"
              value = "Bearer token"
            }
          ]
          discover_request_params = [
            {
              # ...
            }
          ]
          discover_type        = "http"
          discover_url         = "api.example.com/discover"
          enable_discover_code = true
          format_result_code   = "200"
          item_list = [
            "item1",
            "item2",
          ]
          manual_discover_result = "{\"result\":\"true\"}"
          pagination = {
            attribute = [
              "records",
            ]
            last_page_expr     = "$.data.isLastPage"
            limit              = 100
            limit_field        = "limit"
            max_pages          = 10
            offset             = 0
            offset_field       = "offset"
            page_field         = "page"
            size               = 50
            size_field         = "pageSize"
            total_record_field = "totalRecords"
            type               = "offset"
            zero_indexed       = true
          }
        }
        login_body = "{\"username\":\"user\",\"password\":\"pass\"}"
        login_url  = "https://api.example.com/login"
        pagination = {
          attribute = [
            "records",
          ]
          last_page_expr     = "$.data.isLastPage"
          limit              = 100
          limit_field        = "limit"
          max_pages          = 10
          offset             = 0
          offset_field       = "offset"
          page_field         = "page"
          size               = 50
          size_field         = "pageSize"
          total_record_field = "totalRecords"
          type               = "offset"
          zero_indexed       = true
        }
        password            = "restPassword"
        reject_unauthorized = true
        retry_rules = {
          codes = [
            429,
            500,
          ]
          enable_header         = true
          interval              = 500
          limit                 = 3
          max_interval_ms       = 10000
          multiplier            = 1.5
          retry_connect_reset   = false
          retry_connect_timeout = true
          retry_header_name     = "Retry-After"
          type                  = "backoff"
        }
        safe_headers = [
          "Content-Type",
          "Authorization",
        ]
        scheduling = {
          state_tracking = {
            # ...
          }
        }
        timeout              = 30
        token                = "restBearerToken"
        token_resp_attribute = "access_token"
        token_secret         = "restBearerTokenSecret"
        use_round_robin_dns  = false
        username             = "restUser"
      }
      type = "rest"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_s3 = {
    collector = {
      conf = {
        aws_api_key               = "AKIAIOSFODNN7EXAMPLE"
        aws_authentication_method = "auto"
        aws_secret                = "awsSecretPairId"
        aws_secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        bucket                    = "my-s3-bucket"
        extractors = [
          {
            expression = "__input.s3.key.split('/')[1]"
            key        = "account_id"
            pattern    = ".*.gz"
          }
        ]
        max_batch_size = 100
        path           = "logs/2025/10/"
        recurse        = true
        region         = "us-east-1"
      }
      type = "s3"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_script = {
    collector = {
      conf = {
        collect_script  = "echo Hello Collect Script"
        discover_script = "echo Hello Discovery Script"
        shell           = "/bin/bash"
      }
      type = "script"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
  input_collector_splunk = {
    collector = {
      conf = {
        authentication       = "token"
        credentials_secret   = "splunkCredentialsSecret"
        disable_time_filter  = false
        earliest             = "-24h"
        endpoint             = "/services/search/jobs"
        handle_escaped_chars = true
        latest               = "now"
        output_mode          = "json"
        password             = "splunkPassword"
        reject_unauthorized  = true
        search               = "index=main sourcetype=syslog | stats count"
        search_head          = "https://splunk.example.com"
        timeout              = 60
        token                = "splunkBearerToken"
        token_secret         = "splunkBearerTokenSecret"
        use_round_robin_dns  = false
        username             = "splunkUser"
      }
      type = "splunk"
    }
    environment             = "production"
    id                      = "myInputCollectorJobId"
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "rule1",
        "rule2",
      ]
      metadata = [
        {
          name  = "sourceType"
          value = "`value_expression`"
        }
      ]
      output   = "defaultDestination"
      pipeline = "defaultPipeline"
      preprocess = {
        args = [
          "--flag",
          "value",
        ]
        command  = "cat"
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 20000
      throttle_rate_per_sec  = "42 MB"
      type                   = "collection"
    }
    remove_fields = [
      "field1",
      "field2",
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "0 * * * *"
      enabled             = true
      max_concurrent_runs = 2
      resume_missed       = true
      run = {
        earliest                 = 0
        expression               = "true"
        job_timeout              = "30m"
        latest                   = 10
        log_level                = "debug"
        max_task_reschedule      = 3
        max_task_size            = "10GB"
        min_task_size            = "1GB"
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "merge(state,newState)"
          state_update_expression = "state = state + 1"
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "tag1",
      "tag2",
    ]
    ttl             = "2h"
    worker_affinity = false
  }
}