package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestSchemas(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: schemaConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_schema.my_schema", "description", "test schema"),
						resource.TestCheckResourceAttr("criblio_schema.my_schema", "id", "my_schema"),
					),
				},
				{
					Config: schemaConfig,
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

var schemaConfig = `
resource "criblio_schema" "my_schema" {
  description           = "test schema"
  group_id              = "default"
  id                    = "my_schema"
  schema                = <<-EOT
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

output "schema" {
  value = criblio_schema.my_schema
}

data "criblio_schema" "my_schema" {
  group_id = "default"
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}
`
