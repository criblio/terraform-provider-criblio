resource "criblio_mapping_ruleset" "my_mappingruleset" {
  active = false
  conf = {
    functions = [
      {
        conf = {
          add = [
            {
              name  = "...my_name..."
              value = "...my_value..."
            }
          ]
        }
        description = "...my_description..."
        disabled    = true
        filter      = "true"
        final       = false
        group_id    = "...my_group_id..."
        id          = "...my_id..."
      }
    ]
  }
  id      = "my-mapping-ruleset-id"
  product = "stream"
}