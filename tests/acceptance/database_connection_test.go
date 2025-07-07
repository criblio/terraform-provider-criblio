package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestDatabaseConnection(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: dbConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "id", "my_databaseconnection"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "description", "MySQL database connection example"),
					),
				},
				{
					Config: dbConfig,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	})
}

var dbConfig = `
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

data "criblio_database_connection" "my_databaseconnection" {
  database_type = "default"
  group_id      = "default"
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}
`
