package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMonitor(t *testing.T) {
	id := "tf-monitor-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: monitorConfig(id, "created", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_monitor.test", "id", id),
					resource.TestCheckResourceAttr("criblio_monitor.test", "name", "Terraform monitor created"),
					resource.TestCheckResourceAttr("criblio_monitor.test", "managed_by", "terraform"),
					resource.TestCheckResourceAttrPair("data.criblio_monitor.test", "id", "criblio_monitor.test", "id"),
					testCheckListDataSourceHasItems("data.criblio_monitors.all"),
				),
			},
			{
				Config: monitorConfig(id, "updated", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_monitor.test", "name", "Terraform monitor updated"),
					resource.TestCheckResourceAttr("criblio_monitor.test", "enabled", "false"),
				),
			},
			{Config: monitorConfig(id, "updated", false), PlanOnly: true},
			{
				ResourceName:      "criblio_monitor.test",
				ImportState:       true,
				ImportStateId:     id,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorConfig(id, suffix string, enabled bool) string {
	return fmt.Sprintf(`
resource "criblio_monitor" "test" {
  id         = %[1]q
  name       = "Terraform monitor %[2]s"
  dataset_id = "metrics"
  enabled    = %[3]t
  type       = "threshold"

  query = {
    A = {
      mode   = "promql"
      promql = "sum(rate(cpu_usage[5m]))"
    }
  }
  expr    = []
  silence = []

  firing_condition = {
    fire_delay  = 60
    clear_delay = 60
  }
  firing_rule = {
    label = "A"
    threshold = [{
      severity      = "critical"
      limit         = 90
      included_tags = []
      excluded_tags = []
    }]
  }
  aetos_metadata = {}
  monitor_policy_notification_config = {
    enabled = false
    type    = "policy"
    config = [{
      id                    = "default-policy"
      order                 = 0
      conditions            = []
      template_target_pairs = []
    }]
  }
  priority_scalar_value = { value = "P2" }
  team_scalar_value     = { value = "ops" }
}

data "criblio_monitor" "test" {
  id         = criblio_monitor.test.id
  depends_on = [criblio_monitor.test]
}

data "criblio_monitors" "all" {
  depends_on = [criblio_monitor.test]
}
`, id, suffix, enabled)
}
