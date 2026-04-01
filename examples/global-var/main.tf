resource "criblio_global_var" "my_globalvar" {
  description = "test variable"
  group_id    = "default"
  id          = "sample_globalvar_1"
  lib         = "test"
  tags        = "test"
  type        = "number"
  value       = 100
}

output "global_var" {
  value = criblio_global_var.my_globalvar
}
