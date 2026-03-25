resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  apihttp = {
    authentication_method = "none"
    available_endpoints = [
      {
        data_field = "results"
        headers = [
          {
            name  = "Content-Type"
            value = "application/json"
          }
        ]
        method = "POST"
        name   = "search"
        url    = "http://localhost:8080/api/search"
      }
    ]
    description = "Example API HTTP search dataset provider"
    id          = "example_api_http"
    type        = "api_http"
  }
}

resource "criblio_search_dataset_provider" "my_elastic_provider" {
  api_elasticsearch = {
    description = "Example Elasticsearch provider"
    endpoint    = "https://localhost:9200"
    id          = "example_elasticsearch"
    password    = "changeme"
    type        = "api_elasticsearch"
    username    = "elastic"
  }
}

resource "criblio_search_dataset_provider" "my_s3_provider" {
  s3 = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/example-role"
    assume_role_external_id   = "example-external-id"
    aws_api_key               = "AKIAEXAMPLE"
    aws_authentication_method = "auto"
    aws_secret_key            = "example-secret"
    bucket                    = "example-bucket"
    bucket_path_suggestion    = "logs/"
    description               = "Example S3 search dataset provider"
    enable_abac_tagging       = true
    enable_assume_role        = false
    endpoint                  = "https://s3.us-east-1.amazonaws.com"
    id                        = "example_s3"
    region                    = "us-east-1"
    reject_unauthorized       = true
    reuse_connections         = false
    session_token             = "example-session-token"
    signature_version         = "v4"
    type                      = "s3"
  }
}
