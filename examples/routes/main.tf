resource "criblio_routes" "my_routes" {
  group_id = "default"

  comments = [
    {
      comment = "Evaluate grouped routes first"
      id      = "intro"
      index   = 0
    }
  ]

  groups = {
    mygroup = {
      name        = "firstgroup"
      description = "Group of related routes"
      disabled    = false
      index       = 1
    }
  }

  routes = [
    {
      name     = "my_route_1"
      pipeline = "main"
      group_id = "mygroup"
    },
    {
      name     = "my_route_2"
      pipeline = "main"
    }
  ]
}
