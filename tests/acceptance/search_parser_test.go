package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchParser(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_parser_lib_entry.test_parser_lib_entry_unique", "id", "test_parser_lib_entry_unique_123"),
						resource.TestCheckResourceAttr("criblio_parser_lib_entry.test_parser_lib_entry_unique", "group_id", "default_search"),
						resource.TestCheckResourceAttr("criblio_parser_lib_entry.test_parser_lib_entry_unique", "lib", "custom"),
						resource.TestCheckResourceAttr("criblio_parser_lib_entry.test_parser_lib_entry_unique", "description", "test_description"),
						resource.TestCheckResourceAttr("criblio_parser_lib_entry.test_parser_lib_entry_unique", "tags", "test_tags"),
					),
				},
			},
		})
	})
}
