package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestCommitAndDeploy(t *testing.T) {

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_commit.my_commit", "message", "test"),
						resource.TestCheckResourceAttr("criblio_commit.my_commit", "group", "default"),
						resource.TestCheckResourceAttrSet("criblio_commit.my_commit", "items.0.commit"),
						resource.TestCheckResourceAttr("criblio_deploy.my_deploy", "id", "default"),
						resource.TestCheckResourceAttrPair("criblio_deploy.my_deploy", "version", "criblio_commit.my_commit", "items.0.commit"),
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
