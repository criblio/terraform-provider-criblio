resource "criblio_mapping_ruleset" "my_mappingruleset" {
  conf = {
    functions = [
      {
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'default'"
            }
          ]
        }
        description = "Default routing"
        disabled    = false
        filter      = "true"
        final       = true
        group_id    = "default"
        id          = "eval"
      }
    ]
  }
  id      = "default"
  product = "stream"
}
