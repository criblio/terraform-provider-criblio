resource "criblio_schema" "my_schema" {
  description = "...my_description..."
  group_id    = "myExistingGroupId"
  id          = "myUniqueSchemaIdToCRUD"
  schema      = "...my_schema..."
}