package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// monitorSilenceTagsEnv names the tag matchers to cover, as a JSON object such
// as {"host":["worker-01"]}.
const monitorSilenceTagsEnv = "CRIBL_TEST_MONITOR_SILENCE_TAGS"

func TestMonitorSilence(t *testing.T) {
	monitoringPreCheck(t)
	if os.Getenv(monitorSilenceTagsEnv) == "" {
		t.Logf("set %s to a JSON object of tag matchers to cover selective instance muting; skipping that coverage", monitorSilenceTagsEnv)
	}
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

// monitorSilenceTags renders the `tags` attribute from the designated
// environment variable, following the pattern in TestSearchDatasetACL. It is
// gated because the acceptance workspace's API rejects the attribute with
// `400 tags is not a supported field`, which failed this test at its first
// apply; the request body omits the key when the configuration does not set it.
func monitorSilenceTags() string {
	matchers := os.Getenv(monitorSilenceTagsEnv)
	if matchers == "" {
		return ""
	}
	return fmt.Sprintf("\n tags = %q", matchers)
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
 duration_unit = "hours"%s
}`, id, name, begins, begins+int64(time.Hour/time.Millisecond), monitorSilenceTags()) + monitoringDataSources("criblio_monitor_silence", "criblio_monitor_silences")
}
