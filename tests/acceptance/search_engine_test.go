package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchEngine(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	resourceName := "criblio_search_engine.my_searchengine"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: searchEngineConfig("test search engine"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_search_engine"),
					resource.TestCheckResourceAttr(resourceName, "description", "test search engine"),
					resource.TestCheckResourceAttr(resourceName, "tier_size", "small"),
				),
			},
			{
				Config: searchEngineConfig("test search engine updated"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "test search engine updated"),
			},
			{
				Config:   searchEngineConfig("test search engine updated"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "test_search_engine",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"active_workflow",
					"datasets",
					"deletion_started_at",
					"effective_status",
					"engine_type",
					"has_main",
					"is_compute_deprovisioned",
					"is_storage_deprovisioned",
					"last_provisioned_ms",
					"metrics_last_published_at",
					"status",
				},
			},
		},
	})
}

func searchEngineConfig(description string) string {
	return `resource "criblio_search_engine" "my_searchengine" {
  description = "` + description + `"
  id          = "test_search_engine"
  tier_size   = "small"
}
`
}
