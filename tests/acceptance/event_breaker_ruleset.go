package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestEventBreakerRuleset(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: ebConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "id", "test_eventbreakerruleset"),
						resource.TestCheckResourceAttr("criblio_event_breaker_ruleset.my_eventbreakerruleset", "description", "test"),
					),
				},
				{
					Config: ebConfig,
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

var ebConfig = `
resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  description    = "test"
  group_id       = "default"
  id             = "test_eventbreakerruleset"
  lib            = "custom"
  min_raw_length = 256
  rules = [
    {
      condition = "PASS_THROUGH_SOURCE_TYPE"
      disabled  = false
      fields = []
      max_event_bytes     = 51200
      name                = "test"
      parser_enabled      = false
      should_use_data_raw = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
      timestamp = {
        length = 150
        type   = "auto"
      }
      timestamp_anchor_regex = "/^/"
      timestamp_earliest     = "-420weeks"
      timestamp_latest       = "+1week"
      timestamp_timezone     = "local"
      type                   = "regex"
    }
  ]
  tags = "test"
}

output "event_breaker_ruleset" {
  value = criblio_event_breaker_ruleset.my_eventbreakerruleset
}

data "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  group_id = "...my_group_id..."
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider2"
  version = "999.99.9"
}

`
