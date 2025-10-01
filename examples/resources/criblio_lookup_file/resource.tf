resource "criblio_lookup_file" "my_lookupfile" {
  content     = "...my_content..."
  description = "...my_description..."
  group_id    = "myExistingGroupId"
  id          = "myNewLookupIdToCRUD"
  mode        = "disk"
  tags        = "...my_tags..."
}