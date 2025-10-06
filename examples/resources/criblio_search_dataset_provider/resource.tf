resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  api_aws_provider = {
    account_configs = [
      {
        assume_role_arn         = "arn:aws:iam::123456789012:role/MyRole"
        assume_role_external_id = "external-id-123"
        aws_api_key             = "AKIAIOSFODNN7EXAMPLE"
        aws_secret_key          = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        name                    = "aws-account-1"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_azure_data_explorer_provider = {
    client_id     = "11111111-1111-1111-1111-111111111111"
    client_secret = "superSecretAzureKey"
    description   = "my generic provider description"
    id            = "myUniqueGenericProviderId"
    tenant_id     = "00000000-0000-0000-0000-000000000000"
    type          = "generic"
  }
  api_azure_provider = {
    account_configs = [
      {
        client_id     = "12345678-aaaa-bbbb-cccc-123456789abc"
        client_secret = "superSecret"
        name          = "azure-account-1"
        tenant_id     = "tenant-12345"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_elastic_search_provider = {
    description = "my generic provider description"
    endpoint    = "https://elasticsearch.example.com"
    id          = "myUniqueGenericProviderId"
    password    = "elastic_pass"
    type        = "generic"
    username    = "elastic_user"
  }
  api_gcp_provider = {
    account_configs = [
      {
        name                        = "gcp-account-1"
        service_account_credentials = "{...}"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_google_workspace_provider = {
    account_configs = [
      {
        name                        = "workspace-account-1"
        service_account_credentials = "{...}"
        subject                     = "admin@example.com"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_ms_graph_provider = {
    account_configs = [
      {
        client_id     = "client-456"
        client_secret = "secret-789"
        name          = "msgraph-account-1"
        tenant_id     = "tenant-123"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_okta_provider = {
    account_configs = [
      {
        api_token       = "00aBcDefGhIjKlMnOpQrStUvWxYz123456"
        domain_endpoint = "https://dev-123456.okta.com"
        name            = "okta-account-1"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_open_search_provider = {
    description = "my generic provider description"
    endpoint    = "https://opensearch.example.com"
    id          = "myUniqueGenericProviderId"
    password    = "opensearch_pass"
    type        = "generic"
    username    = "opensearch_user"
  }
  api_tailscale_provider = {
    account_configs = [
      {
        client_id     = "tailscale-client-id-123"
        client_secret = "tailscale-client-secret-abc"
        name          = "my-tailscale-account"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  api_zoom_provider = {
    account_configs = [
      {
        account_id    = "123456789"
        client_id     = "zoom-client-id-abc"
        client_secret = "zoom-client-secret-xyz"
        name          = "my-zoom-account"
      }
    ]
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  apihttp_provider = {
    authentication_method = "basic"
    available_endpoints = [
      {
        data_field = "users"
        headers = [
          {
            name  = "Authorization"
            value = "Bearer token123"
          }
        ]
        method = "GET"
        name   = "getUsers"
        url    = "https://api.example.com/users"
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
    authentication_method = "connection_string"
    client_id             = "00000000-0000-0000-0000-000000000000"
    client_secret         = "superSecretClientKey"
    connection_string     = "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=myKey;EndpointSuffix=core.windows.net"
    description           = "my generic provider description"
    id                    = "myUniqueGenericProviderId"
    location              = "eastus2"
    sas_configs = [
      {
        blob_sas_url   = "https://myaccount.blob.core.windows.net/my-container?sp=rl&st=2025-09-30T00:00:00Z&se=2025-10-01T00:00:00Z&spr=https&sv=2021-08-06&sr=c&sig=mysignature"
        container_name = "my-container"
      }
    ]
    storage_account_name = "myAccountName"
    tenant_id            = "11111111-1111-1111-1111-111111111111"
    type                 = "generic"
  }
  click_house_provider = {
    description = "my generic provider description"
    endpoint    = "https://clickhouse.example.com:8443"
    id          = "myUniqueGenericProviderId"
    password    = "click_password"
    type        = "generic"
    username    = "click_user"
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
    endpoint                    = "https://storage.googleapis.com"
    id                          = "myUniqueGenericProviderId"
    service_account_credentials = "{\n  \"type\": \"service_account\",\n  \"project_id\": \"my-gcp-project\",\n  \"private_key_id\": \"abcd1234efgh5678ijkl90\",\n  \"private_key\": \"-----BEGIN PRIVATE KEY-----\\nMIIEvgIBADANBgkqhkiG9...==\\n-----END PRIVATE KEY-----\\n\",\n  \"client_email\": \"my-service-account@my-gcp-project.iam.gserviceaccount.com\",\n  \"client_id\": \"123456789012345678901\",\n  \"auth_uri\": \"https://accounts.google.com/o/oauth2/auth\",\n  \"token_uri\": \"https://oauth2.googleapis.com/token\",\n  \"auth_provider_x509_cert_url\": \"https://www.googleapis.com/oauth2/v1/certs\",\n  \"client_x509_cert_url\": \"https://www.googleapis.com/robot/v1/metadata/x509/my-service-account%40my-gcp-project.iam.gserviceaccount.com\"\n}\n"
    type                        = "generic"
  }
  meta_provider = {
    description = "my generic provider description"
    id          = "myUniqueGenericProviderId"
    type        = "generic"
  }
  prometheus_provider = {
    auth_type       = "basic"
    description     = "my generic provider description"
    endpoint        = "https://prometheus.example.com"
    id              = "myUniqueGenericProviderId"
    max_concurrency = 10
    password        = "prom_pass"
    token           = "prometheusBearerToken123"
    type            = "generic"
    username        = "prom_user"
  }
  s3_provider = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/my-s3-role"
    assume_role_external_id   = "external-id-123"
    aws_api_key               = "AKIAEXAMPLE"
    aws_authentication_method = "auto"
    aws_secret_key            = "mySecretKey123"
    bucket                    = "my-s3-bucket"
    bucket_path_suggestion    = "https://s3.us-east-1.amazonaws.com/mybucket"
    description               = "my generic provider description"
    enable_abac_tagging       = true
    enable_assume_role        = true
    endpoint                  = "https://s3.us-east-1.amazonaws.com"
    id                        = "myUniqueGenericProviderId"
    region                    = "us-east-1"
    reject_unauthorized       = true
    reuse_connections         = true
    session_token             = "MyAWSSessionToken"
    signature_version         = "v4"
    type                      = "generic"
  }
  snowflake_provider = {
    account_identifier = "myorg-myaccount"
    description        = "my generic provider description"
    endpoint           = "https://myorg-myaccount.snowflakecomputing.com"
    id                 = "myUniqueGenericProviderId"
    max_concurrency    = 10
    passphrase         = "myPassphrase"
    priv_key           = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqh...\n-----END PRIVATE KEY-----"
    type               = "generic"
    username           = "snowflake_user"
  }
}