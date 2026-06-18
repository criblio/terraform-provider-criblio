package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSchema(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_schema.my_schema"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.StaticDirectory("testdata/TestSchemas/plan-diff"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "test schema"),
						resource.TestCheckResourceAttr(resourceName, "id", "my_schema"),
					),
				},
				{
					Config: schemaUpdatedConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "test schema updated"),
					),
				},
				{
					Config:   schemaUpdatedConfig,
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default","id":"my_schema"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

const schemaUpdatedConfig = `resource "criblio_schema" "my_schema" {
  description = "test schema updated"
  group_id    = "default"
  id          = "my_schema"
  schema      = <<-EOT
{
  "$id": "https://example.com/person.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Person",
  "type": "object",
  "required": ["firstName", "lastName", "age"],
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be greater than zero, less than 42.",
      "type": "integer",
      "minimum": 0,
      "maximum": 42
    }
  }
}
EOT
}
`
