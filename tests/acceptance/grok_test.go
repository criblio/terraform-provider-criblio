package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestGrok(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: grokConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_grok.my_grok", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_grok.my_grok", "id", "test_grok"),
					),
				},
				{
					Config: grokConfig,
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

var grokConfig = `
resource "criblio_grok" "my_grok" {
  group_id = "default"
  id       = "test_grok"
  content  = <<-EOT
SYSLOG5424PRINTASCII [!-~]+

SYSLOGBASE2 (?:%%{SYSLOGTIMESTAMP:timestamp}|%%{TIMESTAMP_ISO8601:timestamp8601}) (?:%%{SYSLOGFACILITY} )?%%{SYSLOGHOST:logsource}+(?: %%{SYSLOGPROG}:|)
SYSLOGPAMSESSION %%{SYSLOGBASE} (?=%%{GREEDYDATA:message})%%{WORD:pam_module}\(%%{DATA:pam_caller}\): session %%{WORD:pam_session_state} for user %%{USERNAME:username}(?: by %%{GREEDYDATA:pam_by})?

CRON_ACTION [A-Z ]+
CRONLOG %%{SYSLOGBASE} \(%%{USER:user}\) %%{CRON_ACTION:action} \(%%{DATA:message}\)

SYSLOGLINE %%{SYSLOGBASE2} %%{GREEDYDATA:message}

# IETF 5424 syslog(8) format (see http://www.rfc-editor.org/info/rfc5424)
SYSLOG5424PRI <%%{NONNEGINT:syslog5424_pri}>
SYSLOG5424SD \[%%{DATA}\]+
SYSLOG5424BASE %%{SYSLOG5424PRI}%%{NONNEGINT:syslog5424_ver} +(?:%%{TIMESTAMP_ISO8601:syslog5424_ts}|-) +(?:%%{IPORHOST:syslog5424_host}|-) +(-|%%{SYSLOG5424PRINTASCII:syslog5424_app}) +(-|%%{SYSLOG5424PRINTASCII:syslog5424_proc}) +(-|%%{SYSLOG5424PRINTASCII:syslog5424_msgid}) +(?:%%{SYSLOG5424SD:syslog5424_sd}|-|)

SYSLOG5424LINE %%{SYSLOG5424BASE} +%%{GREEDYDATA:syslog5424_msg}
EOT
}

/*
output "grok" {
  value = criblio_grok.my_grok
}

data "criblio_grok" "my_grok" {
  group_id = "...my_group_id..."
}
*/

`
