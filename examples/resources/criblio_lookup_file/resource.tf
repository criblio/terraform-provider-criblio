resource "criblio_lookup_file" "my_lookupfile" {
  content     = "US,United States\nCA,Canada"
  description = "Country code to name lookup"
  group_id    = "Cribl"
  id          = "countries.csv"
  mode        = "memory"
  tags        = "geo,reference"
}