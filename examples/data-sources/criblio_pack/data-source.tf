data "criblio_pack" "my_pack" {
  disabled = false
  group_id = "Cribl"
  id       = "observability-pack"
  with     = "outputs,inputs"
}