resource "criblio_parser_lib_entry" "test_parser_lib_entry_unique" {
  description           = "test_description"
  group_id              = "default_search"
  id                    = "test_parser_lib_entry_unique_123"
  lib                   = "custom"
  tags                  = "test_tags"
  type                  = "delim"
}
