resource "criblio_pack_vars" "my_packvars" {
  description = "This is a test var."
  group_id    = "myExistingGroupId"
  id          = "test_var"
  lib         = "custom"
  pack        = "myExistingPackId"
  tags        = "test"
  type        = "number"
  value       = 100
}