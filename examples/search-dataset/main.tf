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

resource "criblio_search_dataset" "my_searchdataset" {
  apihttp_dataset = {
    description = "test"
    enabled_endpoints = [
      "http://localhost"
    ]
    id = "test_http_dataset"
    metadata = {
      created             = "2021-07-10T14:32:53.487Z"
      enable_acceleration = false
      modified            = "2022-07-12T22:13:15.898Z"
      tags = [
        "test"
      ]
    }
    provider_id = "test_http"
    type        = "api_http"
  }
}

resource "criblio_search_dataset" "my_elastic_search_dataset" {
  api_elastic_search_dataset = {
    description = "test"
    id          = "test_elastic_dataset"
    index       = "test"
    metadata = {
      created             = "2021-06-28T12:13:39.681Z"
      enable_acceleration = false
      modified            = "2022-04-18T23:19:36.636Z"
      tags = [
        "test"
      ]
    }
    provider_id     = "test_elastic"
    timestamp_field = "test"
    type            = "api_elasticsearch"
  }
}

resource "criblio_search_dataset" "my_s3_dataset" {
  s3_dataset = {
    auto_detect_region = false
    bucket             = "test_bucket"
    description        = "test"
    extra_paths = [
      {
        auto_detect_region = false
        bucket             = "test_bucket"
        filter             = "test"
        path               = "logs/*.log"
        region             = "us-east-1"
      }
    ]
    filter = "test"
    id     = "test_s3_dataset"
    metadata = {
      created             = "2021-10-01T19:20:31.326Z"
      enable_acceleration = false
      modified            = "2021-02-10T22:10:18.036Z"
      tags = [
        "test"
      ]
    }
    path                   = "logs/*.log"
    provider_id            = "S3"
    region                 = "us-east-1"
    storage_classes = [
      "STANDARD"
    ]
    type = "s3"
  }
}

/*
resource "criblio_search_dataset" "my_cribl_lake_dataset" {
  cribl_lake_dataset = {
    description = "test"
    id          = "test_cribl_lake_dataset"
    metadata = {
      created             = "2021-09-10T09:02:49.190Z"
      enable_acceleration = false
      modified            = "2021-08-05T01:30:06.408Z"
      tags = [
        "test"
      ]
    }
    provider_id = "cribl_lake"
    type        = "cribl_lake"
  }
}
*/
