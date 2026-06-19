package tests

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestKey(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_key.my_key"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: keyConfig("My Key Metadata", 1759325416),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "key-001"),
					resource.TestCheckResourceAttr(resourceName, "description", "My Key Metadata"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "aes-256-cbc"),
					resource.TestCheckResourceAttr(resourceName, "kms", "local"),
					resource.TestCheckResourceAttr(resourceName, "use_iv", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
				),
			},
			{
				Config: keyConfig("Updated Key Metadata", 1800000000),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "Updated Key Metadata"),
			},
			{
				Config:   keyConfig("Updated Key Metadata", 1800000000),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: keyImportStateID(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func keyImportStateID(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}
		keyID := resourceState.Primary.Attributes["key_id"]
		if keyID == "" {
			return "", fmt.Errorf("resource %s has no key_id", resourceName)
		}
		id := resourceState.Primary.Attributes["id"]
		if id == "" {
			return "", fmt.Errorf("resource %s has no id", resourceName)
		}
		return fmt.Sprintf(`{"group_id":"default","id":%q,"key_id":%q}`, id, keyID), nil
	}
}

func keyConfig(description string, expires int) string {
	return `resource "criblio_key" "my_key" {
  algorithm   = "aes-256-cbc"
  description = "` + description + `"
  expires     = ` + strconv.Itoa(expires) + `
  group_id    = "default"
  id          = "key-001"
  keyclass    = 0
  kms         = "local"
  use_iv      = true
}
`
}
