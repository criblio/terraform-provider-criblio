resource "criblio_parquet_schema" "my_parquetschema" {
  description = "OCSF 1.1.0 base schema"
  group_id    = "Cribl"
  id          = "ocsf_1_1_0"
  schema      = "{\"type\":\"object\",\"properties\":{\"_raw\":{\"type\":\"string\"}}}"
}