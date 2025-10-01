resource "criblio_cribl_lake_house" "my_cribllakehouse" {
  description = "My Lakehouse description"
  id          = "MyLakehouseId"
  status      = "operational"
  tier_size   = "small"
}