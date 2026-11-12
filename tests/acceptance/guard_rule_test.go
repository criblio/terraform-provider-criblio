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
)

const guardRuleGroupID = "default"

func guardRulePreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Guard Rules acceptance coverage is Cloud-only; the on-prem fixture is not provisioned with a Guard rule library")
	}
}

func TestGuardRule(t *testing.T) {
	guardRulePreCheck(t)

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	ruleID := "tf_guard_rule_" + suffix
	ruleset := "tf-guard-ruleset-" + suffix

	resourceName := "criblio_guard_rule.test"
	created := guardRuleCreatedConfig(ruleID, ruleset)
	updated := guardRuleUpdatedConfig(ruleID, ruleset)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: created,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", guardRuleGroupID),
					resource.TestCheckResourceAttr(resourceName, "id", ruleID),
					resource.TestCheckResourceAttr(resourceName, "description", "Internal employee badge identifiers"),
					resource.TestCheckResourceAttr(resourceName, "regex", `/\bBADGE-\d{8}\b/`),
					resource.TestCheckResourceAttr(resourceName, "lib", "custom"),
					resource.TestCheckResourceAttr(resourceName, "rulesets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rulesets.0", ruleset),
					resource.TestCheckResourceAttr(resourceName, "context_keywords.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "context_keywords.0.keyword", "badge"),
					resource.TestCheckResourceAttr(resourceName, "context_keywords.0.placement", "before"),
					// Single-rule data source resolves the rule we just created.
					resource.TestCheckResourceAttr("data.criblio_guard_rule.by_id", "id", ruleID),
					resource.TestCheckResourceAttr("data.criblio_guard_rule.by_id", "regex", `/\bBADGE-\d{8}\b/`),
					// Rule-list data source returns a populated collection.
					testCheckListDataSourceHasItems("data.criblio_guard_rules.all"),
				),
			},
			{
				Config:   created,
				PlanOnly: true,
			},
			{
				Config: updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Internal employee badge identifiers (revised)"),
					resource.TestCheckResourceAttr(resourceName, "regex", `/\bBADGE-\d{6,8}\b/`),
					resource.TestCheckResourceAttr(resourceName, "rulesets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "context_keywords.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "context_keywords.1.keyword", "employee"),
					resource.TestCheckResourceAttr(resourceName, "context_keywords.1.placement", "after"),
				),
			},
			{
				Config:   updated,
				PlanOnly: true,
			},
			// An edit made outside Terraform must surface as a diff on refresh.
			{
				PreConfig: func() {
					guardRulePatchOutOfBand(t, ruleID, ruleset)
				},
				Config:             updated,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Re-apply converges the drifted rule back to the declared configuration.
			{
				Config: updated,
				Check: resource.TestCheckResourceAttr(
					resourceName, "description", "Internal employee badge identifiers (revised)"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf(`{"group_id":%q,"id":%q}`, guardRuleGroupID, ruleID),
				ImportStateVerify: true,
			},
			// Import must leave a clean plan.
			{
				Config:   updated,
				PlanOnly: true,
			},
		},
	})
}

func guardRuleCreatedConfig(ruleID, ruleset string) string {
	return fmt.Sprintf(`resource "criblio_guard_rule" "test" {
  group_id    = %q
  id          = %q
  description = "Internal employee badge identifiers"
  regex       = "/\\bBADGE-\\d{8}\\b/"
  lib         = "custom"
  rulesets    = [%q]

  context_keywords = [
    {
      keyword   = "badge"
      placement = "before"
    },
  ]
}

data "criblio_guard_rule" "by_id" {
  group_id = %q
  id       = criblio_guard_rule.test.id
}

data "criblio_guard_rules" "all" {
  group_id   = %q
  depends_on = [criblio_guard_rule.test]
}
`, guardRuleGroupID, ruleID, ruleset, guardRuleGroupID, guardRuleGroupID)
}

func guardRuleUpdatedConfig(ruleID, ruleset string) string {
	return fmt.Sprintf(`resource "criblio_guard_rule" "test" {
  group_id    = %q
  id          = %q
  description = "Internal employee badge identifiers (revised)"
  regex       = "/\\bBADGE-\\d{6,8}\\b/"
  lib         = "custom"
  rulesets    = [%q, %q]

  context_keywords = [
    {
      keyword   = "badge"
      placement = "before"
    },
    {
      keyword   = "employee"
      placement = "after"
    },
  ]
}

data "criblio_guard_rule" "by_id" {
  group_id = %q
  id       = criblio_guard_rule.test.id
}

data "criblio_guard_rules" "all" {
  group_id   = %q
  depends_on = [criblio_guard_rule.test]
}
`, guardRuleGroupID, ruleID, ruleset, ruleset+"-extra", guardRuleGroupID, guardRuleGroupID)
}

// guardRulePatchOutOfBand edits the rule straight through the Control Plane API,
// bypassing Terraform, so the next refresh has real drift to detect.
func guardRulePatchOutOfBand(t *testing.T, ruleID, ruleset string) {
	t.Helper()

	client := guardRuleAPIClient(t)
	path := fmt.Sprintf("/m/%s/lib/sds-rules/%s", guardRuleGroupID, ruleID)
	body := map[string]any{
		"id":          ruleID,
		"regex":       `/\bBADGE-\d{6,8}\b/`,
		"description": "changed outside terraform",
		"lib":         "custom",
		"rulesets":    []string{ruleset, ruleset + "-extra"},
	}

	if _, err := restclient.Patch[map[string]any, map[string]any](context.Background(), client, path, body); err != nil {
		t.Fatalf("out-of-band patch of %s failed: %v", path, err)
	}
}

// guardRuleAPIClient builds a REST client from the same credentials the provider
// uses, applying the provider's defaults for workspace, org, and cloud domain.
func guardRuleAPIClient(t *testing.T) *restclient.Client {
	t.Helper()

	credentials, err := auth.GetCredentials()
	if err != nil || credentials == nil {
		t.Fatalf("load Cribl credentials: %v", err)
	}
	if credentials.Workspace == "" {
		credentials.Workspace = "main"
	}
	if credentials.CloudDomain == "" {
		credentials.CloudDomain = "cribl.cloud"
	}
	if credentials.OrganizationID == "" {
		t.Fatal("CRIBL_ORGANIZATION_ID must be set to patch a Guard Rule out of band")
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
	})
}
