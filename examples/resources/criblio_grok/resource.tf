resource "criblio_grok" "my_grok" {
  content  = "...my_content..."
  group_id = "myExistingGroupId"
  id       = "myGrokIdToCRUD"
  size     = 6.73
  tags     = "...my_tags..."
}