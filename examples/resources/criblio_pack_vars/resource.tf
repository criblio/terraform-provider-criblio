resource "criblio_pack_vars" "my_packvars" {
  args = [
    {
      name = "val"
      type = "number"
    }
  ]
  description = "This is a test var."
  group_id    = "Cribl"
  id          = "test_var"
  lib         = "custom"
  pack        = "example-pack"
  tags        = "test"
  type        = "number"
  value       = 100
}