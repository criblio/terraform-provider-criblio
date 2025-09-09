resource "criblio_collector" "my_collector" {
  collector = {
    conf = {
      auth_type                 = "secret"
      authentication            = "oauthSecret"
      aws_api_key               = "...my_aws_api_key..."
      aws_authentication_method = "secret"
      aws_secret                = "...my_aws_secret..."
      aws_secret_key            = "...my_aws_secret_key..."
      bucket                    = "...my_bucket..."
      collect_method            = "other"
      collect_url               = "...my_collect_url..."
      connection_id             = "...my_connection_id..."
      connection_string         = "...my_connection_string..."
      container_name            = "...my_container_name..."
      credentials_secret        = "...my_credentials_secret..."
      dataset                   = "...my_dataset..."
      disable_time_filter       = true
      earliest                  = "...my_earliest..."
      endpoint                  = "...my_endpoint..."
      extractors = [
        {
          # ...
        }
      ]
      handle_escaped_chars        = true
      latest                      = "...my_latest..."
      max_batch_size              = 9
      output_mode                 = "json"
      password                    = "...my_password..."
      path                        = "...my_path..."
      query                       = "...my_query..."
      query_validation_enabled    = false
      recurse                     = false
      region                      = "...my_region..."
      reject_unauthorized         = true
      search                      = "...my_search..."
      search_head                 = "...my_search_head..."
      service_account_credentials = "...my_service_account_credentials..."
      storage_account_name        = "...my_storage_account_name..."
      timeout                     = 61
      token                       = "...my_token..."
      token_secret                = "...my_token_secret..."
      use_round_robin_dns         = false
      username                    = "...my_username..."
    }
    destructive = true
    encoding    = "...my_encoding..."
    type        = "rest"
  }
  environment             = "...my_environment..."
  group_id                = "default"
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
    stale_channel_flush_ms = 3239544.85
    throttle_rate_per_sec  = "...my_throttle_rate_per_sec..."
    type                   = "collection"
  }
  remove_fields = [
    "..."
  ]
  resume_on_boot = true
  schedule = {
    cron_schedule       = "...my_cron_schedule..."
    enabled             = true
    max_concurrent_runs = 3.71
    run = {
      earliest                 = 4.02
      expression               = "...my_expression..."
      job_timeout              = "...my_job_timeout..."
      latest                   = 5.4
      log_level                = "silly"
      max_task_reschedule      = 7.19
      max_task_size            = "...my_max_task_size..."
      min_task_size            = "...my_min_task_size..."
      mode                     = "...my_mode..."
      reschedule_dropped_tasks = false
      time_range_type          = "...my_time_range_type..."
    }
    skippable = false
  }
  streamtags = [
    "..."
  ]
  ttl             = "...my_ttl..."
  worker_affinity = false
}