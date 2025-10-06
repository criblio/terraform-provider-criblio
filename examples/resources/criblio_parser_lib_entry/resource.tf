resource "criblio_parser_lib_entry" "my_parserlibentry" {
  description = "Parser for CSV log format"
  group_id    = "Cribl"
  id          = "parser1"
  lib         = "custom"
  tags        = "csv,logs"
  type        = "csv"
}