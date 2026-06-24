resource "criblio_commit" "my_commit" {
  effective = true
  group     = "default"
  message   = "Update Terraform-managed configuration"
}
