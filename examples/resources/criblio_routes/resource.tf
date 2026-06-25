resource "criblio_routes" "my_routes" {
  group_id = "default"
  id       = "default"

  routes = [
    {
      description = "Route application errors to the main pipeline"
      filter      = "level == 'error'"
      final       = true
      name        = "Errors to main"
      pipeline    = "main"
    }
  ]
}
