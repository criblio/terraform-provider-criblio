package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSystemLimits(t *testing.T) {
	groupID := "default"
	resourceName := "criblio_system_limits.test"
	var schemaResponse frameworkresource.SchemaResponse
	provider.NewSystemLimitsResource().Schema(context.Background(), frameworkresource.SchemaRequest{}, &schemaResponse)
	var ignored []string
	for name := range schemaResponse.Schema.Attributes {
		if name != "group_id" && name != "lookup_allowed_paths" {
			ignored = append(ignored, name)
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{Config: systemLimitsConfig(groupID, `["/tmp"]`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "group_id", groupID),
				resource.TestCheckResourceAttr(resourceName, "lookup_allowed_paths.0", "/tmp"),
			)},
			{Config: systemLimitsConfig(groupID, `[]`), Check: resource.TestCheckResourceAttr(resourceName, "lookup_allowed_paths.#", "0")},
			{Config: systemLimitsConfig(groupID, `[]`), PlanOnly: true},
			{ResourceName: resourceName, ImportState: true, ImportStateId: groupID,
				ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "group_id", ImportStateVerifyIgnore: ignored},
		},
	})
}

func systemLimitsConfig(groupID, lookupPaths string) string {
	return fmt.Sprintf(`resource "criblio_system_limits" "test" {
  group_id             = %q
  lookup_allowed_paths = %s
}`, groupID, lookupPaths)
}
