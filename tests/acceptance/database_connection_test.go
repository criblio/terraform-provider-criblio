package tests

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabaseConnection(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		resourceName := "criblio_database_connection.my_databaseconnection"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: databaseConnectionConfig("MySQL database connection example", "test", 60),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "my_databaseconnection"),
						resource.TestCheckResourceAttr(resourceName, "description", "MySQL database connection example"),
						resource.TestCheckResourceAttr(resourceName, "auth_type", "connectionString"),
						resource.TestCheckResourceAttr(resourceName, "database_type", "mysql"),
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "tags", "test"),
						resource.TestCheckResourceAttr(resourceName, "connection_string", "mysql://user:password@8.8.8.8:3306/mydb"),
					),
				},
				{
					Config: databaseConnectionConfig("Updated MySQL database connection example", "updated", 90),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "Updated MySQL database connection example"),
						resource.TestCheckResourceAttr(resourceName, "tags", "updated"),
						resource.TestCheckResourceAttr(resourceName, "request_timeout", "90"),
					),
				},
				{
					Config:   databaseConnectionConfig("Updated MySQL database connection example", "updated", 90),
					PlanOnly: true,
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateId:           `{"group_id":"default","id":"my_databaseconnection"}`,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"config_obj", "connection_string", "creds_secrets", "password", "text_secret"},
				},
			},
		})
	})
}

func databaseConnectionConfig(description string, tags string, requestTimeout int) string {
	return `resource "criblio_database_connection" "my_databaseconnection" {
  auth_type          = "connectionString"
  config_obj         = "test"
  connection_string  = "mysql://user:password@8.8.8.8:3306/mydb"
  connection_timeout = 1000
  database_type      = "mysql"
  description        = "` + description + `"
  group_id           = "default"
  id                 = "my_databaseconnection"
  password           = "test"
  request_timeout    = ` + strconv.Itoa(requestTimeout) + `
  tags               = "` + tags + `"
  user               = "test"
}
`
}
