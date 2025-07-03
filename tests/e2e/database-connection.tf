resource "criblio_database_connection" "my_databaseconnection" {
  # Required fields according to OpenAPI schema
  auth_type     = "connectionString"
  database_type = "mysql"
  description   = "MySQL database connection example"
  id            = "my_databaseconnection"
  group_id      = "default"

  # Optional fields
  config_obj         = "test"
  connection_string  = "mysql://user:password@localhost:3306/mydb"
  connection_timeout = 1000
  password           = "test"
  request_timeout    = 60
  tags               = "test"
  user               = "test"
}

output "database_connection" {
  value = criblio_database_connection.my_databaseconnection
}
