package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Defaults to default_logs and the user ID used in examples/search-dataset-acl.
// Always uses the existing test-team team.
// CRIBL_TEST_DATASET_ACL_DATASET_ID and CRIBL_TEST_DATASET_ACL_USER
// optionally override the dataset and user defaults.
// Configured ACL maps are authoritative; the test changes shares and revokes them on cleanup.
// CRIBL_TEST_DATASET_ACL_USER must be a Search member ID from GET /products/search/users
// (not an email), or an API-credential principal in the form <client-id>@clients.
func TestSearchDatasetACL(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Search dataset ACL acceptance tests require Cribl.Cloud")
	}
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	dataset := os.Getenv("CRIBL_TEST_DATASET_ACL_DATASET_ID")
	if dataset == "" {
		dataset = "default_logs"
	}
	user := os.Getenv("CRIBL_TEST_DATASET_ACL_USER")
	if user == "" {
		user = "google-oauth2|115928590796700585084"
	}
	const team = "test-team"
	config := func(policy string) string {
		return fmt.Sprintf(`resource "criblio_search_dataset_acl" "test" {
 dataset_id = %q
 users = { %q = %q }
 teams = { %q = %q }
}

data "criblio_search_dataset_acl" "test" {
 dataset_id = criblio_search_dataset_acl.test.dataset_id
 depends_on = [criblio_search_dataset_acl.test]
}`, dataset, user, policy, team, policy)
	}
	name := "criblio_search_dataset_acl.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{Config: config("DatasetRead"), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "users."+user, "DatasetRead"),
				resource.TestCheckResourceAttr(name, "teams."+team, "DatasetRead"),
				resource.TestCheckResourceAttr("data.criblio_search_dataset_acl.test", "users."+user, "DatasetRead"),
			)},
			{Config: config("DatasetMaintain"), Check: resource.TestCheckResourceAttr(name, "teams."+team, "DatasetMaintain")},
			{Config: config("DatasetMaintain"), PlanOnly: true},
			{ResourceName: name, ImportState: true, ImportStateId: dataset, ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "dataset_id"},
		},
	})
}
