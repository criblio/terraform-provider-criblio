package tests

import (
	"os"
	"time"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackFromFile(t *testing.T) {
        if os.Getenv("DEPLOYMENT") == "onprem" {
                time.Sleep(2 * time.Second)
        }

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "id", "pack-from-file"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "description", "Pack from file"),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "display_name", "Pack from file"),
						resource.TestCheckResourceAttrWith("criblio_pack.my_pack", "filename", func(value string) error {
							baseName := filepath.Base(value)
							expectedName := "cribl-palo-alto-networks-source-1.0.0.crbl"
							if baseName != expectedName {
								return fmt.Errorf("expected filename base name %q, got %q (full path: %q)", expectedName, baseName, value)
							}
							return nil
						}),
						resource.TestCheckResourceAttr("criblio_pack.my_pack", "version", "1.0.0"),
					),
				},
				{
					ConfigDirectory: config.TestNameDirectory(),
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
