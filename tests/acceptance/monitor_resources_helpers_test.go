package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func monitoringPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("The on-prem acceptance fixture does not enable customer metrics and V2 Notifications")
	}
}

func monitoringLifecycle(t *testing.T, tfType, id, initial, updated, changedField, want string, ignore []string, checks ...resource.TestCheckFunc) {
	t.Helper()
	name := tfType + ".test"
	listType := map[string]string{"criblio_monitor_shared_config": "criblio_monitor_shared_configs", "criblio_monitor_profile": "criblio_monitor_profiles", "criblio_monitor_silence": "criblio_monitor_silences", "criblio_notification_policy": "criblio_notification_policies"}[tfType]
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{Config: initial, Check: resource.ComposeAggregateTestCheckFunc(
				append([]resource.TestCheckFunc{
					resource.TestCheckResourceAttr(name, "id", id),
					resource.TestCheckResourceAttr("data."+tfType+".by_id", "id", id),
					testCheckListDataSourceHasItems("data." + listType + ".all"),
				}, checks...)...,
			)},
			{Config: initial, PlanOnly: true},
			{Config: updated, Check: resource.ComposeAggregateTestCheckFunc(append([]resource.TestCheckFunc{
				resource.TestCheckResourceAttr(name, changedField, want),
			}, checks...)...)},
			{Config: updated, PlanOnly: true},
			{ResourceName: name, ImportState: true, ImportStateId: id, ImportStateVerify: true, ImportStateVerifyIgnore: ignore},
			{Config: updated, PlanOnly: true},
		},
	})
}

func monitoringDataSources(tfType, listType string) string {
	return fmt.Sprintf(`
data %q "by_id" { id = %s.test.id }
data %q "all" { depends_on = [%s.test] }
`, tfType, tfType, listType, tfType)
}
