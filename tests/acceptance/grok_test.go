package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGrok(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is 'prohibited by current license'")
	}

	resourceName := "criblio_grok.my_grok"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: grokConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "id", "test_grok"),
					),
				},
				{
					Config: grokUpdatedConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "content", "UPDATEDWORD [a-zA-Z]+\n"),
					),
				},
				{
					Config:   grokUpdatedConfig,
					PlanOnly: true,
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateId:           `{"group_id":"default","id":"test_grok"}`,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"tags"},
				},
			},
		})
	})
}

const grokUpdatedConfig = `resource "criblio_grok" "my_grok" {
  group_id = "default"
  id       = "test_grok"
  content  = <<-EOT
UPDATEDWORD [a-zA-Z]+
EOT
}
`

const grokConfig = `resource "criblio_grok" "my_grok" {
  group_id = "default"
  id       = "test_grok"
  content  = <<-EOT
TESTWORD [a-zA-Z]+
EOT
}
`
