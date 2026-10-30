package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"
)

func TestMonitorProfile(t *testing.T) {
	monitoringPreCheck(t)
	id := "tf-profile-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	monitoringLifecycle(t, "criblio_monitor_profile", id, monitorProfile(id, "Platform"), monitorProfile(id, "Platform updated"), "name", "Platform updated", nil)
}

func monitorProfile(id, name string) string {
	return fmt.Sprintf(`resource "criblio_monitor_shared_config" "priority" {
 id = %q
 config = jsonencode({ priority = { value = "P3" } })
}
resource "criblio_monitor_profile" "test" {
 id = %q
 name = %q
 config_ids = [criblio_monitor_shared_config.priority.id]
}`, id+"-priority", id, name) + monitoringDataSources("criblio_monitor_profile", "criblio_monitor_profiles")
}
