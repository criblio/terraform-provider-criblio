terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Configuration options
  #server_url ="https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/main/app/api/v1"
}

resource "criblio_pack" "my_pack" {
  id           = "pack-from-file"
  group_id     = "default"
  description  = "Pack from file"
  disabled     = true
  display_name = "Pack from file"
  filename     = "cribl-palo-alto-networks-source-1.0.0.crbl"
  version      = "1.0.0"

}

# Output the pack details to see the read-only attributes
output "pack_details" {
  value = criblio_pack.my_pack
}
