resource "criblio_global_var" "my_globalvar" {
  args = [
    {
      name = "val"
      type = "number"
    }
  ]
  description = "This is a test var."
  group_id    = "default"
  id          = "test_var"
  lib         = "custom"
  tags        = "test"
  type        = "number"
  value       = 100
}