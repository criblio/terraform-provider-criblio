package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"
)

func TestMonitorSharedConfig(t *testing.T) {
	monitoringPreCheck(t)
	id := "tf-shared-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	monitoringLifecycle(t, "criblio_monitor_shared_config", id, monitorSharedConfig(id, "P3"), monitorSharedConfig(id, "P2"), "name", "Priority P2", nil)
}

func monitorSharedConfig(id, priority string) string {
	return fmt.Sprintf(`resource "criblio_monitor_shared_config" "test" {
 id = %q
 name = %q
 config = jsonencode({ priority = { value = %q } })
}`, id, "Priority "+priority, priority) + monitoringDataSources("criblio_monitor_shared_config", "criblio_monitor_shared_configs")
}
