package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEventBreakerRuleset(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_event_breaker_ruleset.my_eventbreakerruleset"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: eventBreakerRulesetConfig("phase2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "phase2_event_breaker"),
					resource.TestCheckResourceAttr(resourceName, "description", "phase2"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
				),
			},
			{
				Config: eventBreakerRulesetConfig("phase2 updated"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "phase2 updated"),
			},
			{Config: eventBreakerRulesetConfig("phase2 updated"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default","id":"phase2_event_breaker"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func eventBreakerRulesetConfig(description string) string {
	return `resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  description    = "` + description + `"
  group_id       = "default"
  id             = "phase2_event_breaker"
  lib            = "custom"
  min_raw_length = 256
  rules = [{
    condition              = "true"
    disabled               = false
    event_breaker_regex    = "/[\\n\\r]+(?!\\s)/"
    fields                 = []
    max_event_bytes        = 51200
    name                   = "phase2"
    parser_enabled         = false
    should_use_data_raw    = false
    timestamp              = { type = "auto" }
    timestamp_anchor_regex = "/^/"
    type                   = "regex"
  }]
  tags = "phase2"
}
`
}
