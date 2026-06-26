package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackLookups(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack test for on-prem: uses source pack installation")
	}

	const (
		packResource   = "criblio_pack.lookups_pack"
		lookupResource = "criblio_pack_lookups.my_packlookups"
	)

	t.Run("plan-diff", func(t *testing.T) {
		config := exampleConfig(t, "pack-with-lookups")

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(packResource, "id", "pack-with-lookups"),
						resource.TestCheckResourceAttr(packResource, "group_id", "default"),
						resource.TestCheckResourceAttr(packResource, "description", "Pack with lookups"),
						resource.TestCheckResourceAttr(lookupResource, "id", "my_id"),
						resource.TestCheckResourceAttr(lookupResource, "group_id", "default"),
						resource.TestCheckResourceAttr(lookupResource, "content", "column1, column2, column3, column4"),
						resource.TestCheckResourceAttr(lookupResource, "description", "my_description"),
						resource.TestCheckResourceAttr(lookupResource, "mode", "memory"),
						resource.TestCheckResourceAttr(lookupResource, "tags", "my_tags"),
					),
				},
				{
					Config: config,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	})
}

func exampleConfig(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "..", "examples", name, "main.tf"))
	if err != nil {
		t.Fatalf("read example %q: %v", name, err)
	}

	return string(content)
}
