resource "criblio_pack_lookups" "my_packlookups" {
  group_id = "...my_group_id..."
  id       = "...my_id..."
  lookup_file_input1 = {
    description = "...my_description..."
    file_info = {
      filename = "...my_filename..."
    }
    id   = "...my_id..."
    mode = "disk"
    size = 8.47
    tags = "...my_tags..."
  }
  lookup_file_input2 = {
    content     = "...my_content..."
    description = "...my_description..."
    id          = "...my_id..."
    mode        = "memory"
    size        = 5.09
    tags        = "...my_tags..."
  }
  pack = "...my_pack..."
}