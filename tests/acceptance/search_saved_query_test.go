package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchSavedQuery(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Test criblio_search_saved_query resources
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery", "id", "test_saved_query"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery", "name", "test"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery", "description", "test"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery", "is_private", "true"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery_with_notifications", "id", "test_saved_query_with_notifications_2"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery_with_notifications", "name", "test_with_notifications"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery_with_notifications", "description", "test with notifications"),
						resource.TestCheckResourceAttr("criblio_search_saved_query.my_searchsavedquery_with_notifications", "is_private", "true"),

						// Test criblio_notification_target resource
						resource.TestCheckResourceAttr("criblio_notification_target.my_notificationtarget", "id", "test_slack_target_2"),

						// Test criblio_notification resource
						resource.TestCheckResourceAttr("criblio_notification.my_notification", "id", "test_notification"),
						resource.TestCheckResourceAttr("criblio_notification.my_notification", "condition", "search"),
						resource.TestCheckResourceAttr("criblio_notification.my_notification", "disabled", "false"),
						resource.TestCheckResourceAttr("criblio_notification.my_notification", "group", "default_search"),
					),
				},
			},
		})
	})
}
