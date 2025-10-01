data "criblio_pack" "my_pack" {
  disabled = true
  group_id = "myExistingGroupId"
  id       = "myUniquePackIdToCRUD"
  with     = "input1, input2, input3"
}