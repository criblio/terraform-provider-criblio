package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestCommitAndDeploy(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: cydConfig,
					Check:  resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("criblio_commit.my_commit", "message", "test"),
					//resource.TestCheckResourceAttr("criblio_commit.my_commit", "group", "syslog-workers"),
					//resource.TestCheckResourceAttr("criblio_commit.my_deploy", "id", "syslog-workers"),
					),
				},
				{
					Config: cydConfig,
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

var cydConfig = `
/*
resource "criblio_commit" "my_commit" {
  effective = true
  group     = "syslog-workers"
  message   = "test"
}

resource "criblio_deploy" "my_deploy" {
  id      = "default"
  version = data.criblio_config_version.my_configversion.items[0]
}

data "criblio_config_version" "my_configversion" {
  id         = "syslog-workers"
  depends_on = [criblio_commit.my_commit]
}

output "deploy" {
  value = criblio_deploy.my_deploy
}

output "config_version" {
  value = data.criblio_config_version.my_configversion
}

output "commit" {
  value = criblio_commit.my_commit
}
*/

provider "criblio" {
    server_url = "https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
}
`
