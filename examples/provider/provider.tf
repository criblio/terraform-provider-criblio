terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "1.23.46"
    }
  }
}

provider "criblio" {
  cloud_domain    = "..." # Optional - can use CRIBL_CLOUD_DOMAIN environment variable
  organization_id = "..." # Optional - can use CRIBL_ORGANIZATION_ID environment variable
  server_url      = "..." # Optional
  workspace_id    = "..." # Optional - can use CRIBL_WORKSPACE_ID environment variable

  # Optional TLS settings for servers using internal CAs or self-signed certificates:
  # ca_cert_file = "/path/to/company-ca.pem" # Or CRIBL_CA_CERT_FILE
  # insecure_skip_verify = true # Or CRIBL_INSECURE_SKIP_VERIFY; testing only
}
