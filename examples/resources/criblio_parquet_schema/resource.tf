resource "criblio_parquet_schema" "my_parquetschema" {
  description = "...my_description..."
  group_id    = "myExistingGroupId"
  id          = "myUniqueParquetSchemaIdToCRUD"
  schema      = "...my_schema..."
}