resource "criblio_global_var" "my_globalvar" {
  description = "This is a test var."
  group_id    = "myExistingGroupId"
  id          = "test_var"
  lib         = "custom"
  tags        = "test"
  type        = "number"
  value       = 100
}