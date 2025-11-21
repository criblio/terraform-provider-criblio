package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabaseConnection(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "id", "my_databaseconnection"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "description", "MySQL database connection example"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "auth_type", "connectionString"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "database_type", "mysql"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "tags", "test"),
						resource.TestCheckResourceAttr("criblio_database_connection.my_databaseconnection", "connection_string", "mysql://user:password@localhost:3306/mydb"),
					),
				},
			},
		})
	})
}
