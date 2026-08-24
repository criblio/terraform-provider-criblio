package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchSource(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	resourceName := "criblio_search_source.my_searchsource"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: searchSourceConfig("test search source", 31170),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "test_search_source"),
						resource.TestCheckResourceAttr(resourceName, "type", "cribl_http"),
						resource.TestCheckResourceAttr(resourceName, "description", "test search source"),
						resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
						resource.TestCheckResourceAttr(resourceName, "auth_tokens.0.enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "auth_tokens.0.token", "test_search_source_auth_token"),
						resource.TestCheckResourceAttr(resourceName, "tls.cert_path", "$CRIBL_CLOUD_CRT"),
						resource.TestCheckResourceAttr(resourceName, "tls.min_version", "TLSv1.2"),
						resource.TestCheckResourceAttr(resourceName, "tls.priv_key_path", "$CRIBL_CLOUD_KEY"),
					),
				},
				{
					Config: searchSourceConfig("test search source updated", 31170),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "test search source updated"),
					),
				},
				{
					Config:   searchSourceConfig("test search source updated", 31170),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     "test_search_source",
					ImportStateVerify: true,
				},
			},
		})
	})
}

func searchSourceConfig(description string, port int) string {
	return `resource "criblio_secret" "search_source_token" {
  description = "Search source acceptance-test token"
  group_id    = "default_search"
  id          = "test_search_source_auth_token"
  secret_type = "text"
  value       = "test-search-source-token-value"
}

resource "criblio_search_source" "my_searchsource" {
  cribl_api   = "/cribl/_bulk"
  description = "` + description + `"
  disabled    = false
  host        = "0.0.0.0"
  id          = "test_search_source"
  port        = ` + fmt.Sprint(port) + `
  type        = "cribl_http"

  auth_tokens = [
    {
      description = "Search source acceptance-test token"
      enabled     = true
      token       = criblio_secret.search_source_token.id
    }
  ]

  tls = {
    disabled = false
  }
}
`
}
