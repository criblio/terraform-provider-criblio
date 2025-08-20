########################
### Search ###
########################



### Search Dataset Provider ###
################################
resource "criblio_search_dataset_provider" "sbx_s3_provider" {
  s3_provider = {
    id          = "sbx_s3provider" # Replace 'john' with your first name
    description = "Search Sandbox S3"
    type        = "s3"

    # Set authentication method to assumeRole
    aws_authentication_method = "auto"

    # AssumeRole configuration
    enable_assume_role      = true
    assume_role_arn         = "arn:aws:iam::498946618345:role/searchingformeaning"
    assume_role_external_id = "whathappensinvegas"

    # These can be empty when using AssumeRole
    aws_api_key    = ""
    aws_secret_key = ""

    # Standard S3 settings
    signature_version   = "v4"
    reuse_connections   = true
    reject_unauthorized = true
  }
}

################################
### Search Dataset ###
################################



resource "criblio_search_dataset" "sbx_weblogs" {
  s3_dataset = {
    id          = "sbx_weblogs" # Replace 'john' with your first name
    description = "Search Sandbox Web Logs"
    provider_id = criblio_search_dataset_provider.sbx_s3_provider.id
    bucket      = "everydayimsearching"

    # Escape the $ to prevent Terraform interpolation
    path = "data/$${dataSource}/$${_time:%Y}/$${_time:%m}/$${_time:%d}/$${_time:%H}"
    #     path = "everydayimsearching/"
    auto_detect_region = false
    region             = "us-east-1"

    # Remove or use a valid JS expression
    # filter = ""  # or remove this line entirely

    skip_event_time_filter = false

    storage_classes = [
      "STANDARD"
    ]

    metadata = {
      enable_acceleration = false # Set to false to avoid needing scanMode/earliest
      tags = [
        "weblogs",
        "sandbox"
      ]
    }

    type = "s3"
  }
}

################################
### Search Saved Query ###
################################



resource "criblio_search_saved_query" "my_searchsavedquery" {
  description = "terraform saved query using search sample dataset"
  earliest    = "-1h"
  id          = "search_sample_saved"
  is_private  = true
  latest      = "now"
  name        = "search_sample_saved"
  query       = "dataset=\"cribl_search_sample\" | limit 1000"
  schedule = {
    cron_schedule = "0 */6 * * *" # Run every 6 hours
    enabled       = true
    keep_last_n   = 7 # Keep last 7 results
    notifications = {
      disabled = false
    }
    tz = "UTC" # Or "America/New_York", "America/Los_Angeles", etc.
  }
}

resource "criblio_search_saved_query" "my_searchsavedquery_lake" {
  description = "terraform saved query using lake dataset"
  earliest    = "-1h"
  id          = "search_lake_saved"
  is_private  = true
  latest      = "now"
  name        = "search_lake_saved"
  query       = "dataset=\"${criblio_cribl_lake_dataset.my_cribllakedataset.id}\" | limit 1000"
  schedule = {
    cron_schedule = "0 */6 * * *" # Run every 6 hours
    enabled       = true
    keep_last_n   = 7 # Keep last 7 results
    notifications = {
      disabled = false
    }
    tz = "UTC" # Or "America/New_York", "America/Los_Angeles", etc.
  }
}
resource "criblio_pack" "my_pack" {
  id           = "cribl-search-aws-vpc-flow-logs"
  group_id     = "default_search"
  description  = "Pack from file"
  disabled     = true
  display_name = "terraform pack"
  filename     = "cribl-search-aws-vpc-flow-logs-0.1.1.crbl"
  version      = "1.0.0"

}