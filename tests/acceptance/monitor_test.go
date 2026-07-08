// Hand-written: do not regenerate (listed in .codegen-ignore).

package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitor(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping criblio_monitor for on-prem deployments: Aetos is cloud-only.")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	id := "tf_monitor_" + suffix
	resourceName := "criblio_monitor.test"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				// Step 1: Create
				{
					Config: monitorConfig(id, "Terraform Monitor Test", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", id),
						resource.TestCheckResourceAttr(resourceName, "name", "Terraform Monitor Test"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "type", "threshold"),

						// managed_by is computed and stamped by the backend when the
						// Terraform provider User-Agent is detected.
						resource.TestCheckResourceAttrSet(resourceName, "managed_by"),

						// Data source mirrors the resource.
						resource.TestCheckResourceAttrPair(
							"data.criblio_monitor.test", "id",
							resourceName, "id",
						),
						resource.TestCheckResourceAttrPair(
							"data.criblio_monitor.test", "managed_by",
							resourceName, "managed_by",
						),

						// List data source contains at least one item.
						testCheckListDataSourceHasItems("data.criblio_monitors.all"),
					),
				},

				// Step 2: Update — rename + disable
				{
					Config: monitorConfig(id, "Terraform Monitor Updated", false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "Terraform Monitor Updated"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					),
				},

				// Step 3: No-op plan (idempotency)
				{
					Config:   monitorConfig(id, "Terraform Monitor Updated", false),
					PlanOnly: true,
				},

				// Step 4: Import
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     id,
					ImportStateVerify: true,
					// managed_by is read-only; id is set via ImportStatePassthroughID.
					// All other fields must round-trip correctly.
				},
			},
		})
	})
}

// monitorConfig returns a Terraform configuration that creates a criblio_monitor
// resource plus a data source and list data source to exercise all three.
func monitorConfig(id, name string, enabled bool) string {
	return fmt.Sprintf(`
resource "criblio_monitor" "test" {
  id      = %q
  name    = %q
  enabled = %v
  type    = "threshold"

  priority = jsonencode({ value = "P3" })
  team     = jsonencode({ value = "tf-acceptance" })

  query = jsonencode({
    A = { mode = "promql", promql = "up" }
  })

  expr = jsonencode([])

  firing_condition = jsonencode({
    fire_delay  = 60
    clear_delay = 60
  })

  firing_rule = jsonencode({
    label     = "down"
    threshold = [
      {
        severity     = "critical"
        limit        = 0
        operator     = "lt"
        includedTags = []
        excludedTags = []
      }
    ]
  })

  metadata     = jsonencode({})
  notification = jsonencode({ enabled = false, type = "policy", config = [] })
  silence      = []
}

data "criblio_monitor" "test" {
  id         = criblio_monitor.test.id
  depends_on = [criblio_monitor.test]
}

data "criblio_monitors" "all" {
  depends_on = [criblio_monitor.test]
}
`, id, name, enabled)
}
