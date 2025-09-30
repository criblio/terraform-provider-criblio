data "criblio_pack_vars" "my_packvars" {
  group_id = "myExistingGroupId"
  pack     = "myExistingPackId"
  with     = "{ '$ref': '#/components/schemas/ComplexModel' }"
}