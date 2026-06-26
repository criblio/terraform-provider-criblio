package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestStreamSyslogToLake(t *testing.T) {
	onPrem := os.Getenv("DEPLOYMENT") == "onprem"
	configVariables := config.Variables{
		"onprem": config.BoolVariable(onPrem),
	}

	t.Run("plan-diff", func(t *testing.T) {
		steps := []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		}
		if !onPrem {
			steps[0].Check = streamSyslogToLakeCheck()
			steps = append(steps, resource.TestStep{
				ImportState:       true,
				ImportStateId:     "syslog-workers",
				ResourceName:      "criblio_group.syslog_worker_group[0]",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cloud",
					"estimated_ingest_rate",
					"provisioned",
					"streamtags",
					"worker_remote_access",
				},
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
			})
		}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps:                     steps,
		})
	})
}

func streamSyslogToLakeCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("criblio_group.syslog_worker_group[0]", "id", "syslog-workers"),
		resource.TestCheckResourceAttr("criblio_group.syslog_worker_group[0]", "name", "syslog-workers"),
		resource.TestCheckResourceAttr("criblio_group.syslog_worker_group[0]", "product", "stream"),
		resource.TestCheckResourceAttr("criblio_source.syslog_source[0]", "id", "syslog-input"),
		resource.TestCheckResourceAttr("criblio_source.syslog_source[0]", "group_id", "syslog-workers"),
		resource.TestCheckResourceAttr("criblio_destination.cribl_lake[0]", "id", "cribl-lake-2"),
		resource.TestCheckResourceAttr("criblio_destination.cribl_lake[0]", "group_id", "syslog-workers"),
	)
}
