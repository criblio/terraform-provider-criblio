package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMonitorSilence(t *testing.T) {
	monitoringPreCheck(t)
	id := "tf-mute-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	// The backend list defaults end to now, so verify an active silence.
	begins := time.Now().Add(-time.Minute).UnixMilli()
	monitoringLifecycle(t, "criblio_monitor_silence", id, monitorSilence(id, "Maintenance", begins), monitorSilence(id, "Maintenance updated", begins), "name", "Maintenance updated", nil,
		resource.TestCheckResourceAttr("criblio_monitor.example_threshold", "id", id+"-monitor"),
		resource.TestCheckResourceAttrPair("criblio_monitor_silence.test", "monitor_id", "criblio_monitor.example_threshold", "id"),
		resource.TestCheckResourceAttr("criblio_monitor_silence.test", "duration_type", "relative"),
		resource.TestCheckResourceAttr("criblio_monitor_silence.test", "duration_value", "1"),
		resource.TestCheckResourceAttr("criblio_monitor_silence.test", "duration_unit", "hours"),
		resource.TestCheckTypeSetElemNestedAttrs("data.criblio_monitor_silences.all", "items.*", map[string]string{
			"id":         id,
			"monitor_id": id + "-monitor",
		}),
	)
}

func monitorSilence(id, name string, begins int64) string {
	monitor := strings.ReplaceAll(monitorThreshold(80), "acc_test_threshold_monitor", id+"-monitor")
	return monitor + fmt.Sprintf(`
resource "criblio_monitor_silence" "test" {
 id = %q
 name = %q
 monitor_id = criblio_monitor.example_threshold.id
 begins_at = %d
 expires_at = %d
 duration_type = "relative"
 duration_value = 1
 duration_unit = "hours"
 tags = jsonencode({ host = ["worker-01", "worker-02"] })
}`, id, name, begins, begins+int64(time.Hour/time.Millisecond)) + monitoringDataSources("criblio_monitor_silence", "criblio_monitor_silences")
}
