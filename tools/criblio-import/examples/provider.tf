terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Authentication via environment variables or ~/.cribl/credentials
}
