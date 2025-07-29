resource "criblio_lookup_file" "my_lookupfile_1" {
  content     = "column1, column2, column3, column4"
  description = "my_description"
  group_id    = "default_search"
  id          = "test_id_lookupfile_1.csv"
  mode        = "memory"
  tags        = "my_tags"
}
