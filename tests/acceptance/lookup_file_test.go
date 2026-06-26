package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLookupFile(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	lookupID := "test_lookup_file_" + suffix + ".csv"
	resourceName := "criblio_lookup_file.my_lookup_file"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: lookupFileConfig(lookupID, "created", "region,name\nUS,United States", "geo,created"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", lookupID),
						resource.TestCheckResourceAttr(resourceName, "group_id", "default_search"),
						resource.TestCheckResourceAttr(resourceName, "content", "region,name\nUS,United States"),
						resource.TestCheckResourceAttr(resourceName, "description", "Lookup file created"),
						resource.TestCheckResourceAttr(resourceName, "mode", "memory"),
						resource.TestCheckResourceAttr(resourceName, "tags", "geo,created"),
					),
				},
				{
					Config: lookupFileConfig(lookupID, "updated", "region,name\nCA,Canada", "geo,updated"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "content", "region,name\nCA,Canada"),
						resource.TestCheckResourceAttr(resourceName, "description", "Lookup file updated"),
						resource.TestCheckResourceAttr(resourceName, "tags", "geo,updated"),
					),
				},
				{
					Config:   lookupFileConfig(lookupID, "updated", "region,name\nCA,Canada", "geo,updated"),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     fmt.Sprintf(`{"group_id":"default_search","id":%q}`, lookupID),
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"content",
						"description",
						"pending_task",
						"tags",
						"version",
					},
				},
			},
		})
	})
}

func lookupFileConfig(id, descriptionSuffix, content, tags string) string {
	return fmt.Sprintf(`resource "criblio_lookup_file" "my_lookup_file" {
  content     = %q
  description = "Lookup file %s"
  group_id    = "default_search"
  id          = %q
  mode        = "memory"
  tags        = %q
}
`, content, descriptionSuffix, id, tags)
}
