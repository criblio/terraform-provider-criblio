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
      index       = 1
    }
  }

  routes = [
    {
      name        = "my_route_1"
      pipeline    = "main"
      group_id    = "mygroup"
      description = "Route events through the grouped main pipeline"
    },
    {
      name        = "my_route_2"
      pipeline    = "main"
      description = "Route events through the default main pipeline"
    }
  ]
}
