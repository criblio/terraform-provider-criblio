resource "criblio_pack_lookups" "my_packlookups" {
  content     = "...my_content..."
  description = "...my_description..."
  group_id    = "myExistingGroupId"
  id          = "myUniqueLookupIdToCRUD"
  mode        = "memory"
  pack        = "myExistingPackId"
  tags        = "...my_tags..."
}