package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGrok(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is 'prohibited by current license'")
	}

	resourceName := "criblio_grok.my_grok[0]"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "id", "test_grok"),
					),
				},
				{
					Config: grokUpdatedConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "tags", "updated"),
					),
				},
				{
					Config:   grokUpdatedConfig,
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default","id":"test_grok"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

const grokUpdatedConfig = `resource "criblio_grok" "my_grok" {
  count = 1

  group_id = "default"
  id       = "test_grok"
  tags     = "updated"
  content  = <<-EOT
TESTWORD [a-zA-Z]+
EOT
}
`
