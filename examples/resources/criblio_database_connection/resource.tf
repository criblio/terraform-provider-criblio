resource "criblio_database_connection" "my_databaseconnection" {
  auth_type          = "password"
  config_obj         = "{\"ssl\":true}"
  connection_string  = "postgresql://db.example.com:5432/appdb"
  connection_timeout = 30
  database_type      = "postgres"
  description        = "Primary application database"
  group_id           = "Cribl"
  id                 = "db-prod-01"
  password           = "$${{secret:db_password}"
  request_timeout    = 60
  tags               = "prod,db"
  user               = "appuser"
}