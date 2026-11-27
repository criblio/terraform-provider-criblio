package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchMonitor(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("System Insights monitors require Cribl.Cloud")
	}
	example, err := os.ReadFile("../../examples/resources/criblio_search_monitor/resource.tf")
	if err != nil {
		t.Fatal(err)
	}
	id := acctest.RandomWithPrefix("tf-insights")
	const address = "criblio_search_monitor.example"
	initial := strings.Replace(string(example), "worker-events-in", id, 1)
	updated := strings.Replace(initial, "Worker events in", "Updated monitor", 1)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{Config: initial + fmt.Sprintf(`
data "criblio_search_monitor" "one" { id = %s.id }
data "criblio_search_monitors" "all" { depends_on = [%s] }
output "found" { value = contains(data.criblio_search_monitors.all.items[*].id, data.criblio_search_monitor.one.id) }
`, address, address), Check: resource.TestCheckOutput("found", "true")},
			{Config: updated, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(address, "name", "Updated monitor"),
				resource.TestCheckResourceAttr(address, "is_default", "false"),
			)},
			{ResourceName: address, ImportState: true, ImportStateId: id, ImportStateVerify: true},
			{Config: updated, PlanOnly: true},
		},
	})
}
