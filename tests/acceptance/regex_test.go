package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRegex(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_regex.my_regex"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "test_regex_2"),
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "lib", "custom"),
						resource.TestCheckResourceAttr(resourceName, "tags", "test"),
						resource.TestCheckResourceAttr(resourceName, "id", "test_regex_2"),
					),
				},
				{
					Config: regexUpdatedConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "test_regex_updated"),
						resource.TestCheckResourceAttr(resourceName, "sample_data", "10.0.0.1"),
					),
				},
				{
					Config:   regexUpdatedConfig,
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default","id":"test_regex_2"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

const regexUpdatedConfig = `resource "criblio_regex" "my_regex" {
  description = "test_regex_updated"
  group_id    = "default"
  id          = "test_regex_2"
  lib         = "custom"
  regex       = "/\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b/"
  sample_data = "10.0.0.1"
  tags        = "test"
}
`
