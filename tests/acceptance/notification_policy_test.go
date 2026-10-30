package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"
)

func TestNotificationPolicy(t *testing.T) {
	monitoringPreCheck(t)
	id := "tf-policy-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	monitoringLifecycle(t, "criblio_notification_policy", id, notificationPolicy(id, 10), notificationPolicy(id, 20), "order", "20", nil)
}

func notificationPolicy(id string, order int) string {
	return notificationTargetConfig(id+"-target", "policy") + fmt.Sprintf(`
resource "criblio_notification_template" "policy" {
 id = %q
 type = "sns"
 content = "{{monitor_name}}: {{status}}"
}
resource "criblio_notification_policy" "test" {
 id = %q
 order = %d
 disabled = true
 final = true
 wait_to_group = 0
 group_by_labels = ["host"]
 conditions = jsonencode([[{key = "team", operator = "=", value = "platform"}]])
 template_target_pairs = [{
  target_id = criblio_notification_target.my_notificationtarget.id
  template_id = criblio_notification_template.policy.id
 }]
}`, id+"-template", id, order) + monitoringDataSources("criblio_notification_policy", "criblio_notification_policies")
}
