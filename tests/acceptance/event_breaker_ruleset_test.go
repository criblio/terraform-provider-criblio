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
				Config: eventBreakerRulesetConfig("test", "items.events"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "test_event_breaker"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.json_array_field", "items.events"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.parent_fields_to_copy.0", "source"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.json_extract_all", "false"),
				),
			},
			{
				Config: eventBreakerRulesetConfig("test updated", "records"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test updated"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.json_array_field", "records"),
				),
			},
			{Config: eventBreakerRulesetConfig("test updated", "records"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default","id":"test_event_breaker"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func eventBreakerRulesetConfig(description, jsonArrayField string) string {
	return `resource "criblio_event_breaker_ruleset" "my_eventbreakerruleset" {
  description    = "` + description + `"
  group_id       = "default"
  id             = "test_event_breaker"
  lib            = "custom"
  min_raw_length = 256
  rules = [{
    condition              = "true"
    disabled               = false
    event_breaker_regex    = "/[\\n\\r]+(?!\\s)/"
    fields                 = []
    max_event_bytes        = 51200
    name                   = "test"
    parser_enabled         = false
    should_use_data_raw    = false
    timestamp              = { length = 150, type = "auto" }
    timestamp_anchor_regex = "/^/"
    type                   = "regex"
  }, {
    name                   = "api-events"
    condition              = "true"
    type                   = "json_array"
    json_array_field       = "` + jsonArrayField + `"
    parent_fields_to_copy  = ["source"]
    json_extract_all       = false
    timestamp              = { length = 150, type = "auto" }
    timestamp_anchor_regex = "/^/"
  }]
  tags = "test"
}
`
}
