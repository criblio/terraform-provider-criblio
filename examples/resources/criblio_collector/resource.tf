resource "criblio_collector" "my_collector" {
  group_id = "default"
  id       = "...my_id..."
  input_collector_azure_blob = {
    collector = {
      conf = {
        auth_type         = "clientCert"
        connection_string = "...my_connection_string..."
        container_name    = "...my_container_name..."
        extractors = [
          {
            # ...
          }
        ]
        max_batch_size       = 7
        path                 = "...my_path..."
        recurse              = false
        storage_account_name = "...my_storage_account_name..."
      }
      type = "azureblob"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = false
      }
      send_to_routes         = true
      stale_channel_flush_ms = 15172601.08
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = true
      max_concurrent_runs = 3.7
      resume_missed       = true
      run = {
        earliest                 = 4.68
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 6.6
        log_level                = "debug"
        max_task_reschedule      = 5.24
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "run"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = false
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = true
  }
  input_collector_cribl_lake = {
    collector = {
      conf = {
        dataset = "...my_dataset..."
      }
      type = "cribllake"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 35996665.2
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = false
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = false
      max_concurrent_runs = 10.17
      resume_missed       = false
      run = {
        earliest                 = 3.52
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 3.89
        log_level                = "warn"
        max_task_reschedule      = 8.06
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = false
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = true
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = false
  }
  input_collector_database = {
    collector = {
      conf = {
        connection_id            = "...my_connection_id..."
        query                    = "...my_query..."
        query_validation_enabled = false
      }
      type = "database"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = false
      }
      send_to_routes         = false
      stale_channel_flush_ms = 13453157.01
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = false
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = false
      max_concurrent_runs = 7.04
      resume_missed       = false
      run = {
        earliest                 = 6.12
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 4.5
        log_level                = "info"
        max_task_reschedule      = 6.3
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "preview"
        reschedule_dropped_tasks = false
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = true
  }
  input_collector_gcs = {
    collector = {
      conf = {
        auth_type = "clientSecret"
        bucket    = "...my_bucket..."
        extractors = [
          {
            # ...
          }
        ]
        max_batch_size              = 5
        path                        = "...my_path..."
        recurse                     = false
        service_account_credentials = "...my_service_account_credentials..."
      }
      type = "gcs"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = false
      }
      send_to_routes         = false
      stale_channel_flush_ms = 32381408.37
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = false
      max_concurrent_runs = 9.94
      resume_missed       = true
      run = {
        earliest                 = 8.47
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 1.77
        log_level                = "warn"
        max_task_reschedule      = 10.09
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = false
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "relative"
        time_warning = {
          # ...
        }
      }
      skippable = true
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = true
  }
  input_collector_health_check = {
    collector = {
      conf = {
        authentication      = "tokenSecret"
        collect_method      = "other"
        collect_url         = "...my_collect_url..."
        credentials_secret  = "...my_credentials_secret..."
        password            = "...my_password..."
        reject_unauthorized = true
        timeout             = 926
        username            = "...my_username..."
      }
      type = "healthcheck"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = false
      }
      send_to_routes         = true
      stale_channel_flush_ms = 30099853.21
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = true
      max_concurrent_runs = 3.54
      resume_missed       = true
      run = {
        earliest                 = 1.52
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 1.43
        log_level                = "silly"
        max_task_reschedule      = 1.88
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "preview"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = false
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "absolute"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = true
  }
  input_collector_rest = {
    collector = {
      conf = {
        auth_header_expr = "...my_auth_header_expr..."
        auth_header_key  = "...my_auth_header_key..."
        auth_request_headers = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        auth_request_params = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        authentication           = "google_oauth"
        capture_headers          = true
        client_secret_param_name = "...my_client_secret_param_name..."
        collect_method           = "get"
        collect_request_headers = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        collect_request_params = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        collect_url         = "...my_collect_url..."
        credentials_secret  = "...my_credentials_secret..."
        decode_url          = false
        disable_time_filter = false
        discovery = {
          discover_body       = "...my_discover_body..."
          discover_data_field = "...my_discover_data_field..."
          discover_method     = "get"
          discover_request_headers = [
            {
              name  = "...my_name..."
              value = "...my_value..."
            }
          ]
          discover_request_params = [
            {
              # ...
            }
          ]
          discover_type        = "http"
          discover_url         = "...my_discover_url..."
          enable_discover_code = true
          format_result_code   = "...my_format_result_code..."
          item_list = [
            "..."
          ]
          pagination = {
            attribute = [
              "..."
            ]
            last_page_expr     = "...my_last_page_expr..."
            limit              = 2
            limit_field        = "...my_limit_field..."
            max_pages          = 2
            offset             = 5
            offset_field       = "...my_offset_field..."
            page_field         = "...my_page_field..."
            size               = 0
            size_field         = "...my_size_field..."
            total_record_field = "...my_total_record_field..."
            type               = "page"
            zero_indexed       = false
          }
        }
        login_body = "...my_login_body..."
        login_url  = "...my_login_url..."
        pagination = {
          attribute = [
            "..."
          ]
          last_page_expr     = "...my_last_page_expr..."
          limit              = 9
          limit_field        = "...my_limit_field..."
          max_pages          = 4
          offset             = 7
          offset_field       = "...my_offset_field..."
          page_field         = "...my_page_field..."
          size               = 8
          size_field         = "...my_size_field..."
          total_record_field = "...my_total_record_field..."
          type               = "none"
          zero_indexed       = false
        }
        password            = "...my_password..."
        reject_unauthorized = true
        retry_rules = {
          codes = [
            6
          ]
          enable_header         = false
          interval              = 5
          limit                 = 9
          max_interval_ms       = 8
          multiplier            = 8.98
          retry_connect_reset   = false
          retry_connect_timeout = true
          retry_header_name     = "...my_retry_header_name..."
          type                  = "backoff"
        }
        safe_headers = [
          "..."
        ]
        scheduling = {
          state_tracking = {
            # ...
          }
        }
        timeout              = 1766
        token                = "...my_token..."
        token_resp_attribute = "...my_token_resp_attribute..."
        token_secret         = "...my_token_secret..."
        use_round_robin_dns  = true
        username             = "...my_username..."
      }
      type = "rest"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = true
      }
      send_to_routes         = true
      stale_channel_flush_ms = 26906061.17
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = false
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = false
      max_concurrent_runs = 7.88
      resume_missed       = true
      run = {
        earliest                 = 7.89
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 8.93
        log_level                = "debug"
        max_task_reschedule      = 4.17
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "preview"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = false
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "absolute"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = false
  }
  input_collector_s3 = {
    collector = {
      conf = {
        aws_api_key               = "...my_aws_api_key..."
        aws_authentication_method = "auto"
        aws_secret                = "...my_aws_secret..."
        aws_secret_key            = "...my_aws_secret_key..."
        bucket                    = "...my_bucket..."
        extractors = [
          {
            # ...
          }
        ]
        max_batch_size = 6
        path           = "...my_path..."
        recurse        = false
        region         = "...my_region..."
      }
      type = "s3"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = true
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = true
      }
      send_to_routes         = false
      stale_channel_flush_ms = 8564102.73
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = false
      max_concurrent_runs = 10.2
      resume_missed       = false
      run = {
        earliest                 = 7.6
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 9.2
        log_level                = "info"
        max_task_reschedule      = 7.38
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "preview"
        reschedule_dropped_tasks = false
        state_tracking = {
          enabled                 = true
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "absolute"
        time_warning = {
          # ...
        }
      }
      skippable = true
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = true
  }
  input_collector_splunk = {
    collector = {
      conf = {
        authentication       = "basic"
        credentials_secret   = "...my_credentials_secret..."
        disable_time_filter  = false
        earliest             = "...my_earliest..."
        endpoint             = "...my_endpoint..."
        handle_escaped_chars = false
        latest               = "...my_latest..."
        output_mode          = "csv"
        password             = "...my_password..."
        reject_unauthorized  = false
        search               = "...my_search..."
        search_head          = "...my_search_head..."
        timeout              = 297
        token                = "...my_token..."
        token_secret         = "...my_token_secret..."
        use_round_robin_dns  = true
        username             = "...my_username..."
      }
      type = "splunk"
    }
    environment             = "...my_environment..."
    id                      = "...my_id..."
    ignore_group_jobs_limit = false
    input = {
      breaker_rulesets = [
        "..."
      ]
      metadata = [
        {
          name  = "...my_name..."
          value = "...my_value..."
        }
      ]
      output   = "...my_output..."
      pipeline = "...my_pipeline..."
      preprocess = {
        args = [
          "..."
        ]
        command  = "...my_command..."
        disabled = true
      }
      send_to_routes         = false
      stale_channel_flush_ms = 531943.62
      throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
      type                   = "collection"
    }
    remove_fields = [
      "..."
    ]
    resume_on_boot = true
    saved_state = {
      # ...
    }
    schedule = {
      cron_schedule       = "...my_cron_schedule..."
      enabled             = false
      max_concurrent_runs = 6.42
      resume_missed       = true
      run = {
        earliest                 = 7.9
        expression               = "...my_expression..."
        job_timeout              = "...my_job_timeout..."
        latest                   = 3.67
        log_level                = "error"
        max_task_reschedule      = 6.56
        max_task_size            = "...my_max_task_size..."
        min_task_size            = "...my_min_task_size..."
        mode                     = "list"
        reschedule_dropped_tasks = true
        state_tracking = {
          enabled                 = false
          state_merge_expression  = "...my_state_merge_expression..."
          state_update_expression = "...my_state_update_expression..."
        }
        time_range_type = "absolute"
        time_warning = {
          # ...
        }
      }
      skippable = false
    }
    streamtags = [
      "..."
    ]
    ttl             = "...my_ttl..."
    worker_affinity = false
  }
}