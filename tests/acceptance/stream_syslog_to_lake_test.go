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
	if !onPrem {
		t.Skip("Skipping cloud Stream syslog-to-lake example: destroy can race source output reference cleanup")
	}

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
				Config: streamSyslogToLakeGroupImportConfig(),
			})
		}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps:                     steps,
		})
	})
}

func streamSyslogToLakeGroupImportConfig() string {
	return `resource "criblio_group" "syslog_worker_group" {
  count   = 1
  id      = "syslog-workers"
  product = "stream"
}
`
}
