package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSecret(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_secret.my_secret"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: secretConfig("API key for ingestion service", "token-abc123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "sec-test-001"),
					resource.TestCheckResourceAttr(resourceName, "description", "API key for ingestion service"),
					resource.TestCheckResourceAttr(resourceName, "secret_type", "text"),
				),
			},
			{
				Config: secretConfig("Updated API key for ingestion service", "token-def456"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "Updated API key for ingestion service"),
			},
			{
				Config:   secretConfig("Updated API key for ingestion service", "token-def456"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           `{"group_id":"default","id":"sec-test-001"}`,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key", "password", "secret_key", "value"},
			},
		},
	})
}

func secretConfig(description, value string) string {
	return `resource "criblio_secret" "my_secret" {
  description = "` + description + `"
  group_id    = "default"
  id          = "sec-test-001"
  secret_type = "text"
  tags        = "env:prod,team:security"
  value       = "` + value + `"
}
`
}
