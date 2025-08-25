# Cribl Lake Destination Configuration
resource "criblio_destination" "cribl_lake" {
  id       = "cribl-lake-11"
  group_id = "default"

  output_cribl_lake = {
    id                                = "cribl-lake-12"
    type                              = "cribl_lake"
    description                       = "Cribl Lake destination"
    dest_path                         = "default_logs"
  }
}