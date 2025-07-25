terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "determined-gian-gkh6kzw"
  workspace_id = "main"
}

resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  apihttp_provider = {
    authentication_method = "none"
    available_endpoints = [
      {
        data_field = "test"
        headers = [
          {
            name  = "test"
            value = "test"
          }
        ]
        method = "POST"
        name   = "test"
        url    = "http://localhost"
      }
    ]
    description = "test"
    id          = "test_http"
    type        = "api_http"
  }
}

resource "criblio_search_dataset_provider" "my_elastic_provider" {
  api_elastic_search_provider = {
    description = "test"
    endpoint    = "http://localhost"
    id          = "test_elastic"
    password    = "test"
    type        = "api_elasticsearch"
    username    = "test"
  }
}
resource "criblio_search_dataset_provider" "my_s3_provider" {
  s3_provider = {
    assume_role_arn           = "arn:aws::123456789012:role/test"
    assume_role_external_id   = "test"
    aws_api_key               = "test"
    aws_authentication_method = "auto"
    aws_secret_key            = "test"
    bucket                    = "test_bucket"
    bucket_path_suggestion    = "test"
    description               = "test"
    enable_abac_tagging       = true
    enable_assume_role        = false
    endpoint                  = "https://s3.us-east-1.amazonaws.com"
    id                        = "test_s3"
    region                    = "us-east-1"
    reject_unauthorized       = true
    reuse_connections         = false
    session_token             = "test"
    signature_version         = "v4"
    type                      = "s3"
  }
}
