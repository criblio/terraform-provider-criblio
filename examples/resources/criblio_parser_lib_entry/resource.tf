resource "criblio_parser_lib_entry" "my_parserlibentry" {
  description = "...my_description..."
  group_id    = "myExistingGroupId"
  id          = "myUniqueParserIdToCRUD"
  lib         = "...my_lib..."
  tags        = "...my_tags..."
  type        = "delim"
}