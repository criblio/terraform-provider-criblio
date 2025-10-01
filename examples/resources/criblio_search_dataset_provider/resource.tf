resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  api_aws_provider = {
    account_configs = [
      {
        assume_role_arn         = "...my_assume_role_arn..."
        assume_role_external_id = "...my_assume_role_external_id..."
        aws_api_key             = "...my_aws_api_key..."
        aws_secret_key          = "...my_aws_secret_key..."
        name                    = "...my_name..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_azure_data_explorer_provider = {
    client_id     = "...my_client_id..."
    client_secret = "...my_client_secret..."
    description   = "my generic provider description"
    id            = "myUniqueGenericProviderId"
    tenant_id     = "...my_tenant_id..."
    type          = "generic"
  }
  api_azure_provider = {
    account_configs = [
      {
        client_id     = "...my_client_id..."
        client_secret = "...my_client_secret..."
        name          = "...my_name..."
        tenant_id     = "...my_tenant_id..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_elastic_search_provider = {
    description = "my generic provider description"
    endpoint    = "...my_endpoint..."
    id          = "myUniqueGenericProviderId"
    password    = "...my_password..."
    type        = "generic"
    username    = "...my_username..."
  }
  api_gcp_provider = {
    account_configs = [
      {
        name                        = "...my_name..."
        service_account_credentials = "...my_service_account_credentials..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_google_workspace_provider = {
    account_configs = [
      {
        name                        = "...my_name..."
        service_account_credentials = "...my_service_account_credentials..."
        subject                     = "...my_subject..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  apihttp_provider = {
    authentication_method = "none"
    available_endpoints = [
      {
        data_field = "...my_data_field..."
        headers = [
          {
            name  = "...my_name..."
            value = "...my_value..."
          }
        ]
        method = "POST"
        name   = "...my_name..."
        url    = "...my_url..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_ms_graph_provider = {
    account_configs = [
      {
        client_id     = "...my_client_id..."
        client_secret = "...my_client_secret..."
        name          = "...my_name..."
        tenant_id     = "...my_tenant_id..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_okta_provider = {
    account_configs = [
      {
        api_token       = "...my_api_token..."
        domain_endpoint = "...my_domain_endpoint..."
        name            = "...my_name..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_open_search_provider = {
    description = "my generic provider description"
    endpoint    = "...my_endpoint..."
    id          = "myUniqueGenericProviderId"
    password    = "...my_password..."
    type        = "generic"
    username    = "...my_username..."
  }
  api_tailscale_provider = {
    account_configs = [
      {
        client_id     = "...my_client_id..."
        client_secret = "...my_client_secret..."
        name          = "...my_name..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_zoom_provider = {
    account_configs = [
      {
        account_id    = "...my_account_id..."
        client_id     = "...my_client_id..."
        client_secret = "...my_client_secret..."
        name          = "...my_name..."
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  aws_security_lake_provider = {
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  azure_blob_provider = {
    authentication_method = "client_secret"
    client_id             = "...my_client_id..."
    client_secret         = "...my_client_secret..."
    connection_string     = "...my_connection_string..."
    description           = "my generic provider description"
    id                    = "myUniqueGenericProviderId"
    location              = "...my_location..."
    sas_configs = [
      {
        blob_sas_url   = "...my_blob_sas_url..."
        container_name = "...my_container_name..."
      }
    ]
    storage_account_name = "...my_storage_account_name..."
    tenant_id            = "...my_tenant_id..."
    type                 = "generic"
  }
  click_house_provider = {
    description = "my generic provider description"
    endpoint    = "...my_endpoint..."
    id          = "myUniqueGenericProviderId"
    password    = "...my_password..."
    type        = "generic"
    username    = "...my_username..."
  }
  cribl_leader_provider = {
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  edge_provider = {
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  gcs_provider = {
    description                 = "my generic provider description"
    endpoint                    = "...my_endpoint..."
    id                          = "myUniqueGenericProviderId"
    service_account_credentials = "...my_service_account_credentials..."
    type                        = "generic"
  }
  meta_provider = {
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  prometheus_provider = {
    auth_type       = "none"
    description     = "my generic provider description"
    endpoint        = "...my_endpoint..."
    id              = "myUniqueGenericProviderId"
    max_concurrency = 9.88
    password        = "...my_password..."
    token           = "...my_token..."
    type            = "generic"
    username        = "...my_username..."
  }
  s3_provider = {
    assume_role_arn           = "...my_assume_role_arn..."
    assume_role_external_id   = "...my_assume_role_external_id..."
    aws_api_key               = "...my_aws_api_key..."
    aws_authentication_method = "auto"
    aws_secret_key            = "...my_aws_secret_key..."
    bucket                    = "...my_bucket..."
    bucket_path_suggestion    = "...my_bucket_path_suggestion..."
    description               = "my generic provider description"
    enable_abac_tagging       = true
    enable_assume_role        = false
    endpoint                  = "...my_endpoint..."
    id                        = "myUniqueGenericProviderId"
    region                    = "...my_region..."
    reject_unauthorized       = true
    reuse_connections         = false
    session_token             = "...my_session_token..."
    signature_version         = "v4"
    type                      = "generic"
  }
  snowflake_provider = {
    account_identifier = "...my_account_identifier..."
    description        = "my generic provider description"
    endpoint           = "...my_endpoint..."
    id                 = "myUniqueGenericProviderId"
    max_concurrency    = 5
    passphrase         = "...my_passphrase..."
    priv_key           = "...my_priv_key..."
    type               = "generic"
    username           = "...my_username..."
  }
}