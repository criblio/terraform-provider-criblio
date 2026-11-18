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
)

const searchUsageGroupPath = "/m/default_search/search/usage-groups"

// searchUsageGroupMember is a user that already exists in the organization.
type searchUsageGroupMember struct {
	ID    string
	Email string
}

func searchUsageGroupPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Usage Groups are a Cribl.Cloud-only Search feature; the on-prem fixture does not serve this endpoint")
	}
}

func TestSearchUsageGroup(t *testing.T) {
	searchUsageGroupPreCheck(t)

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	groupID := "tf_usage_group_" + suffix

	// The API stores only IDs of users that already exist in the organization and
	// silently drops anything else, so membership can be exercised only with a
	// real member. A user belongs to one Usage Group at a time, so this moves them
	// into the test group for the duration of the test; destroying the group
	// returns them to the catch-all `default` group. Set
	// CRIBL_ACC_SKIP_USAGE_GROUP_MEMBERS=1 to leave membership untested.
	member := searchUsageGroupExistingMember(t)

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
			// An import block over the same configuration has to plan as a no-op,
			// which is what a clean plan after import means.
			{
				ResourceName:    resourceName,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   groupID,
				Config:          updated,
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

// searchUsageGroupExistingMember returns a user that the organization already
// knows, discovered from whichever Usage Group currently holds one. It returns
// nil when membership testing is switched off or the organization reports no
// members, so the rest of the lifecycle still runs.
func searchUsageGroupExistingMember(t *testing.T) *searchUsageGroupMember {
	t.Helper()
	if os.Getenv("CRIBL_ACC_SKIP_USAGE_GROUP_MEMBERS") != "" {
		t.Log("CRIBL_ACC_SKIP_USAGE_GROUP_MEMBERS is set; skipping user assignment coverage")
		return nil
	}
	if os.Getenv("TF_ACC") == "" {
		return nil
	}

	client, err := searchUsageGroupAPIClient()
	if err != nil {
		t.Fatalf("build API client: %v", err)
	}
	type usageGroup struct {
		ID    string `json:"id"`
		Users map[string]struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"users"`
	}
	groups, err := restclient.Get[[]usageGroup](context.Background(), client, searchUsageGroupPath)
	if err != nil {
		t.Fatalf("list Usage Groups: %v", err)
	}
	if groups != nil {
		for _, group := range *groups {
			for _, user := range group.Users {
				if user.ID == "" || user.Email == "" {
					continue
				}
				t.Logf("using existing member %q (currently in Usage Group %q) for assignment coverage", user.Email, group.ID)
				return &searchUsageGroupMember{ID: user.ID, Email: user.Email}
			}
		}
	}
	t.Log("no existing Usage Group member found; skipping user assignment coverage")
	return nil
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
