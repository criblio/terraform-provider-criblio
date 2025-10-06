resource "criblio_commit" "my_commit" {
  effective = true
  files = [
    "config/system.yml",
    "routes/default.yml",
  ]
  group   = "ops"
  message = "Update system settings and routes"
}