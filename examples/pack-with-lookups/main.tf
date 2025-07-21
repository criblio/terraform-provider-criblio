resource "criblio_pack_lookups" "my_packlookups" {
  group_id = "default"
  id       = "my_id"
  lookup_file_input1 = {
    content = "fooobar"
    description = "my_description"
    id   = "my_id"
    mode = "disk"
    size = 8.47
    tags = "my_tags"
  }
  pack = criblio_pack.my_pack.id
}

resource "criblio_pack" "my_pack" {
  id           = "pack-with-lookups"
  group_id     = "default"
  description  = "Pack with lookups"
  disabled     = true
  display_name = "Pack from source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"

}
