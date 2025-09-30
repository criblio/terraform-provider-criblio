resource "criblio_database_connection" "my_databaseconnection" {
  auth_type          = "...my_auth_type..."
  config_obj         = "...my_config_obj..."
  connection_string  = "...my_connection_string..."
  connection_timeout = 9.28
  database_type      = "oracle"
  description        = "...my_description..."
  group_id           = "myExistingGroupId"
  id                 = "myUniqueDatabaseConnIdToCRUD"
  password           = "...my_password..."
  request_timeout    = 0.19
  tags               = "...my_tags..."
  user               = "...my_user..."
}