resource "criblio_pack_vars" "my_packvars" {
  description = "This is a test var."
  group_id    = "...my_group_id..."
  id          = "test_var"
  lib         = "custom"
  pack        = "...my_pack..."
  tags        = "test"
  type        = "number"
  value       = 100
}