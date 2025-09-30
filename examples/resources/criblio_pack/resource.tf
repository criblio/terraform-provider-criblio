resource "criblio_pack" "my_pack" {
  allow_custom_functions = false
  author                 = "...my_author..."
  description            = "...my_description..."
  disabled               = false
  display_name           = "...my_display_name..."
  exports = [
    "..."
  ]
  filename               = "./myfile.json"
  force                  = true
  group_id               = "myExistingGroupId"
  id                     = "myUniquePackIdToCRUD"
  inputs                 = 8.61
  is_disabled            = false
  min_log_stream_version = "...my_min_log_stream_version..."
  outputs                = 6.97
  source                 = "myExistingSourceId"
  spec                   = "main"
  tags = {
    data_type = [
      "..."
    ]
    domain = [
      "..."
    ]
    streamtags = [
      "..."
    ]
    technology = [
      "..."
    ]
  }
  version = "...my_version..."
}