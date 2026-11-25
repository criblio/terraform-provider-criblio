package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const searchUsageGroupPath = "/m/default_search/search/usage-groups"

// searchUsageGroupMember is a user that already exists in the organization.
type searchUsageGroupMember struct {
	ID    string
	Email string
}

// Usage Groups are a Cribl.Cloud-only Search feature.
// CRIBL_TEST_USAGE_GROUP_USER_ID and CRIBL_TEST_USAGE_GROUP_USER_EMAIL name the
// designated test member used for assignment coverage; without them the
// lifecycle still runs, membership aside. The member's original Usage Group is
// restored on cleanup.
func searchUsageGroupPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Usage Groups are a Cribl.Cloud-only Search feature; the on-prem fixture does not serve this endpoint")
	}
	if os.Getenv("TF_ACC") == "" {
		return
	}
	// Usage Groups are a licensed Cribl.Cloud Search feature. Where the
	// endpoint is not served, or the credentials cannot reach it, skip rather
	// than failing the suite on an environment gap.
	client, err := searchUsageGroupAPIClient()
	if err != nil {
		t.Skipf("Usage Group acceptance coverage needs cloud credentials: %v", err)
	}
	if _, err := restclient.Get[[]map[string]any](context.Background(), client, searchUsageGroupPath); err != nil {
		t.Skipf("Usage Groups endpoint is not available in this environment: %v", err)
	}
}

func TestSearchUsageGroup(t *testing.T) {
	searchUsageGroupPreCheck(t)

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	groupID := "tf_usage_group_" + suffix

	member := searchUsageGroupTestMember(t)

	resourceName := "criblio_search_usage_group.test"
	created := searchUsageGroupCreatedConfig(groupID, member)
	updated := searchUsageGroupUpdatedConfig(groupID, member)

	// There is deliberately no out-of-band drift step here, unlike
	// guard_rule_test.go. x-terraform-refresh-from-api is not set on this
	// resource, so reads preserve configured inputs and an edit made outside
	// Terraform to `rules`, `description` or `users` does not surface as a diff.
	// That protects the clean plan after import, which the API's canonical rule
	// descriptions and the UI's string limits would otherwise dirty. Drift in
	// `enabled`, `users_count` and `coordinator_heap_memory_limit` is still
	// picked up, since those refresh unconditionally.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		// Delete is verified against the API rather than only through state,
		// since the endpoint refuses to delete built-in groups and a silent
		// no-op would otherwise look like success.
		CheckDestroy: testCheckSearchUsageGroupDestroyed(groupID),
		Steps: []resource.TestStep{
			{
				Config: created,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", groupID),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform acceptance Usage Group"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					// Numeric and percentage limits both round-trip through the
					// dynamic `limit` attribute.
					resource.TestCheckResourceAttr(resourceName, "rules.max_results_per_search.type", "maxResultsPerSearch"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_results_per_search.limit", "100000"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_concurrent_adhoc_searches_per_user.type", "maxConcurrentAdhocSearchesPerUser"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_concurrent_adhoc_searches_per_user.limit", "50%"),
					resource.TestCheckResourceAttrSet(resourceName, "users_count"),
					searchUsageGroupCheckMember(resourceName, member),
					// Single-group data source resolves the group we just created.
					resource.TestCheckResourceAttr("data.criblio_search_usage_group.by_id", "id", groupID),
					resource.TestCheckResourceAttr("data.criblio_search_usage_group.by_id", "description", "Terraform acceptance Usage Group"),
					// Group-list data source returns a populated collection.
					testCheckListDataSourceHasItems("data.criblio_search_usage_groups.all"),
				),
			},
			{
				Config:   created,
				PlanOnly: true,
			},
			{
				Config: updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform acceptance Usage Group (revised)"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_results_per_search.limit", "75000"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_concurrent_adhoc_searches_per_user.limit", "25%"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_timerange_width.type", "maxTimerangeWidth"),
					resource.TestCheckResourceAttr(resourceName, "rules.max_timerange_width.limit", "604800"),
					searchUsageGroupCheckMember(resourceName, member),
				),
			},
			{
				Config:   updated,
				PlanOnly: true,
			},
			// The import ID is the Usage Group ID; the group is always
			// default_search. Every attribute must come back from the API.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     groupID,
				ImportStateVerify: true,
			},
			{
				Config:   updated,
				PlanOnly: true,
			},
		},
	})
}

func searchUsageGroupCreatedConfig(groupID string, member *searchUsageGroupMember) string {
	return searchUsageGroupConfig(groupID, member, false)
}

func searchUsageGroupUpdatedConfig(groupID string, member *searchUsageGroupMember) string {
	return searchUsageGroupConfig(groupID, member, true)
}

// searchUsageGroupConfig renders the test configuration. The revised variant
// changes the group's own fields and every limit, and adds a rule, so the update
// path covers both edits and additions.
func searchUsageGroupConfig(groupID string, member *searchUsageGroupMember, revised bool) string {
	description := "Terraform acceptance Usage Group"
	enabled := "true"
	rules := `    max_results_per_search = {
      type  = "maxResultsPerSearch"
      limit = 100000
    }
    max_concurrent_adhoc_searches_per_user = {
      type  = "maxConcurrentAdhocSearchesPerUser"
      limit = "50%"
    }`
	if revised {
		description += " (revised)"
		enabled = "false"
		rules = `    max_results_per_search = {
      type  = "maxResultsPerSearch"
      limit = 75000
    }
    max_concurrent_adhoc_searches_per_user = {
      type  = "maxConcurrentAdhocSearchesPerUser"
      limit = "25%"
    }
    max_timerange_width = {
      type  = "maxTimerangeWidth"
      limit = 604800
    }`
	}

	users := ""
	if member != nil {
		users = fmt.Sprintf(`

  users = {
    %[1]q = {
      id    = %[1]q
      email = %[2]q
    }
  }`, member.ID, member.Email)
	}

	return fmt.Sprintf(`resource "criblio_search_usage_group" "test" {
  id          = %q
  description = %q
  enabled     = %s

  rules = {
%s
  }%s
}

data "criblio_search_usage_group" "by_id" {
  id = criblio_search_usage_group.test.id
}

data "criblio_search_usage_groups" "all" {
  depends_on = [criblio_search_usage_group.test]
}
`, groupID, description, enabled, rules, users)
}

// TestSearchUsageGroupImportPlansNoOp asserts the acceptance criterion that an
// import leaves a clean plan, using a plannable `import {}` block: the framework
// requires the planned action to be a no-op. That construct needs Terraform
// 1.5.0 or later, while the acceptance workflow pins 1.1.7, so it lives in its
// own test with a version gate instead of inside TestSearchUsageGroup, which
// would otherwise skip wholesale on the older CLI and lose all CI coverage.
// TestSearchUsageGroup still covers import on every version via
// ImportStateVerify followed by a PlanOnly step.
func TestSearchUsageGroupImportPlansNoOp(t *testing.T) {
	searchUsageGroupPreCheck(t)

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	groupID := "tf_usage_group_import_" + suffix
	resourceName := "criblio_search_usage_group.test"
	config := searchUsageGroupCreatedConfig(groupID, nil)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0),
		},
		CheckDestroy: testCheckSearchUsageGroupDestroyed(groupID),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(resourceName, "id", groupID),
			},
			{
				ResourceName:    resourceName,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   groupID,
				Config:          config,
			},
		},
	})
}

// searchUsageGroupCheckMember asserts the member landed in the group, and does
// nothing when membership is not under test.
func searchUsageGroupCheckMember(resourceName string, member *searchUsageGroupMember) resource.TestCheckFunc {
	if member == nil {
		return func(*terraform.State) error { return nil }
	}
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "users."+member.ID+".id", member.ID),
		resource.TestCheckResourceAttr(resourceName, "users."+member.ID+".email", member.Email),
		resource.TestCheckResourceAttr(resourceName, "users_count", "1"),
	)
}

// searchUsageGroupTestMember returns the designated test member for assignment
// coverage, following the pattern in TestSearchDatasetACL.
// CRIBL_TEST_USAGE_GROUP_USER_ID and CRIBL_TEST_USAGE_GROUP_USER_EMAIL name the
// member; both must be set, and the ID must be an existing organization member
// because the API silently drops IDs it does not know. Membership coverage is
// skipped when they are unset, so the rest of the lifecycle still runs.
//
// A user belongs to one Usage Group at a time, so assigning them here moves them
// out of whichever group holds them. The original group is captured up front and
// restored on cleanup, so the member's assignment is the same after the test as
// before it.
func searchUsageGroupTestMember(t *testing.T) *searchUsageGroupMember {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		return nil
	}
	id := os.Getenv("CRIBL_TEST_USAGE_GROUP_USER_ID")
	email := os.Getenv("CRIBL_TEST_USAGE_GROUP_USER_EMAIL")
	if id == "" || email == "" {
		t.Log("set CRIBL_TEST_USAGE_GROUP_USER_ID and CRIBL_TEST_USAGE_GROUP_USER_EMAIL to cover user assignment; skipping that coverage")
		return nil
	}

	client, err := searchUsageGroupAPIClient()
	if err != nil {
		t.Logf("skipping user assignment coverage, no API client: %v", err)
		return nil
	}

	origin, body := searchUsageGroupMembershipOf(t, client, id)
	if origin == "" {
		t.Logf("member %q is in no Usage Group; it will fall back to the catch-all group after the test", email)
	} else {
		t.Logf("member %q is currently in Usage Group %q; it will be restored there on cleanup", email, origin)
	}
	t.Cleanup(func() { searchUsageGroupRestoreMembership(t, client, id, origin, body) })

	return &searchUsageGroupMember{ID: id, Email: email}
}

// searchUsageGroupMembershipOf reports which Usage Group holds the user, with
// that group's full representation so it can be restored verbatim. The update
// endpoint replaces the whole object, so a partial body would drop its rules.
func searchUsageGroupMembershipOf(t *testing.T, client *restclient.Client, userID string) (string, map[string]any) {
	t.Helper()
	groups, err := restclient.Get[[]map[string]any](context.Background(), client, searchUsageGroupPath)
	if err != nil || groups == nil {
		t.Logf("could not read current Usage Group membership: %v", err)
		return "", nil
	}
	for _, group := range *groups {
		users, ok := group["users"].(map[string]any)
		if !ok {
			continue
		}
		if _, member := users[userID]; !member {
			continue
		}
		id, _ := group["id"].(string)
		return id, group
	}
	return "", nil
}

// searchUsageGroupRestoreMembership puts the member back in the group that held
// them. Members of the catch-all `default` group need no action: deleting the
// test group returns them there. Any other group is restored from the captured
// representation, because membership does not move back on its own.
func searchUsageGroupRestoreMembership(t *testing.T, client *restclient.Client, userID, origin string, body map[string]any) {
	t.Helper()
	if origin == "" || origin == "default" || body == nil {
		return
	}
	current, _ := searchUsageGroupMembershipOf(t, client, userID)
	if current == origin {
		return
	}
	path := fmt.Sprintf("%s/%s", searchUsageGroupPath, origin)
	delete(body, "usersCount")
	if _, err := restclient.Patch[map[string]any, map[string]any](context.Background(), client, path, body); err != nil {
		t.Errorf("restore member %q to Usage Group %q: %v; the member is now in %q", userID, origin, err, current)
		return
	}
	t.Logf("restored member %q to Usage Group %q", userID, origin)
}

// testCheckSearchUsageGroupDestroyed confirms the API no longer serves the group
// once Terraform has destroyed it.
func testCheckSearchUsageGroupDestroyed(groupID string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client, err := searchUsageGroupAPIClient()
		if err != nil {
			return err
		}
		path := fmt.Sprintf("%s/%s", searchUsageGroupPath, groupID)
		if _, err := restclient.Get[map[string]any](context.Background(), client, path); err != nil {
			if restclient.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("read %s after destroy: %w", path, err)
		}
		return fmt.Errorf("Usage Group %s still exists after destroy", groupID)
	}
}

// searchUsageGroupAPIClient builds a REST client from the same credentials the
// provider uses, applying the provider's defaults for workspace, org, and cloud
// domain.
func searchUsageGroupAPIClient() (*restclient.Client, error) {
	credentials, err := auth.GetCredentials()
	if err != nil || credentials == nil {
		return nil, fmt.Errorf("load Cribl credentials: %w", err)
	}
	if credentials.Workspace == "" {
		credentials.Workspace = "main"
	}
	if credentials.CloudDomain == "" {
		credentials.CloudDomain = "cribl.cloud"
	}
	if credentials.OrganizationID == "" {
		return nil, fmt.Errorf("CRIBL_ORGANIZATION_ID must be set to verify Usage Group deletion")
	}

	baseURL := auth.ConstructBaseURL(
		auth.ConstructBaseURLInput{BaseURL: "https://{workspaceId}-{organizationId}.{cloudDomain}"},
		credentials,
	)

	return restclient.New(restclient.Config{
		BaseURL:             baseURL,
		ProviderOrgID:       credentials.OrganizationID,
		ProviderWorkspaceID: credentials.Workspace,
		ProviderCloudDomain: credentials.CloudDomain,
		Credentials:         credentials,
		HTTPClient:          &http.Client{},
	}), nil
}
