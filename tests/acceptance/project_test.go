package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProject(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is 'prohibited by current license'")
	}

	resourceName := "criblio_project.my_project"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: projectConfig("test project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_lifecycle_project"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "description", "test project"),
					resource.TestCheckResourceAttr(resourceName, "destinations.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "subscriptions.#", "0"),
				),
			},
			{
				Config: projectConfig("updated project"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "updated project"),
			},
			{
				Config:   projectConfig("updated project"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default","id":"test_lifecycle_project"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func projectConfig(description string) string {
	return `resource "criblio_project" "my_project" {
  description = "` + description + `"
  destinations = []
  group_id = "default"
  id = "test_lifecycle_project"
  subscriptions = []
}
`
}
