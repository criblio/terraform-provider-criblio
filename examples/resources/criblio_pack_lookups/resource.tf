resource "criblio_pack_lookups" "my_packlookups" {
  content     = "US,United States\nCA,Canada"
  description = "Country code to name lookup"
  group_id    = "myExistingGroupId"
  id          = "countries.csv"
  mode        = "memory"
  pack        = "myExistingPackId"
  tags        = "geo,reference"
}