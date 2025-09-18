resource "criblio_routes" "my_routes" {
  group_id = "default"
  routes = [
    {
      name     = "my_route_1"
      pipeline = "main"
    },
    {
      name     = "my_route_2"
      pipeline = "main"
    }
  ]
}
