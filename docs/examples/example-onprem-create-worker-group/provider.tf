# Configure provider
#
# The provider block is empty. Configure authentication using
# environment variables or a credentials file as described in
# https://docs.cribl.io/cribl-as-code/terraform-auth/#terraform-auth-on-prem.

terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = ">= 1.20.138"
    }
  }
}

provider "criblio" {
}
