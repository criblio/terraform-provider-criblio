package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitor(t *testing.T) {
	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	id := "tf_monitor_" + suffix
	resourceName := "criblio_monitor.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			// Create
			{
				Config: monitorConfig(id, "Terraform Monitor Test", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform Monitor Test"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "product", "stream"),
				),
			},
			// Update: rename + disable
			{
				Config: monitorConfig(id, "Terraform Monitor Updated", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform Monitor Updated"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			// No-op plan (idempotency)
			{
				Config:   monitorConfig(id, "Terraform Monitor Updated", false),
				PlanOnly: true,
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     id,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorConfig(id, name string, enabled bool) string {
	return fmt.Sprintf(`
resource "criblio_monitor" "test" {
  id                      = %q
  name                    = %q
  enabled                 = %v
  product                 = "stream"
  firing_after            = 300
  ok_after                = 60
  schedule_interval_seconds = 60
}
`, id, name, enabled)
}
