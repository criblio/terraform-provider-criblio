package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Defaults to default_logs and the user ID used in examples/search-dataset-acl.
// Always uses test-team in the target workspace, creating it if missing.
// The shared team fixture is retained after the test; only its dataset shares are removed.
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
	credentials, err := auth.GetCredentials()
	if err != nil {
		t.Fatal(err)
	}
	baseURL := auth.ConstructBaseURL(auth.ConstructBaseURLInput{}, credentials)
	client := restclient.New(restclient.Config{BaseURL: baseURL, Credentials: credentials})
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	if err := ensureDatasetACLTestTeam(ctx, client); err != nil {
		t.Fatalf("prepare test-team in acceptance workspace %s before changing ACLs: %v", baseURL, err)
	}
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

type datasetACLTestTeam struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Roles       []string `json:"roles"`
}

func ensureDatasetACLTestTeam(ctx context.Context, client *restclient.Client) error {
	const path = "/system/teams/test-team"
	exists := func() (bool, error) {
		teams, err := restclient.Get[[]datasetACLTestTeam](ctx, client, path)
		if restclient.IsNotFound(err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if teams != nil {
			for _, team := range *teams {
				if team.ID == "test-team" {
					return true, nil
				}
			}
		}
		return false, nil
	}
	if found, err := exists(); err != nil {
		return fmt.Errorf("look up workspace team: %w", err)
	} else if found {
		return nil
	}
	createErr := restclient.PostNoResponse(ctx, client, "/system/teams", datasetACLTestTeam{
		ID: "test-team", Name: "test-team", Description: "Shared Terraform acceptance fixture", Roles: []string{"search_user"},
	})
	// Another acceptance run may have created the shared fixture concurrently.
	if found, err := exists(); err == nil && found {
		return nil
	} else if createErr != nil {
		return fmt.Errorf("create test-team with Search access in the target workspace (requires team-management permission): %w", createErr)
	} else if err != nil {
		return fmt.Errorf("verify created workspace team: %w", err)
	}
	return fmt.Errorf("test-team is not visible in the target workspace after creation")
}
