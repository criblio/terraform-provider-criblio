resource "criblio_parser_lib_entry" "my_parserlibentry" {
  additional_properties = "{ \"see\": \"documentation\" }"
  description           = "...my_description..."
  group_id              = "...my_group_id..."
  id                    = "...my_id..."
  lib                   = "...my_lib..."
  tags                  = "...my_tags..."
  type                  = "delim"
}