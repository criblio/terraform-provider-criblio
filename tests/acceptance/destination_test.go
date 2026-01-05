package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDestination(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						// Cribl HTTP destination
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "id", "cribl_http_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.id", "cribl_http_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.type", "cribl_http"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.description", "Send events to Cribl Worker HTTP endpoint with retries"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.compression", "gzip"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.url", "https://edge.example.com:10200"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.load_balanced", "true"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.on_backpressure", "block"),
						// Chronicle destination
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "id", "chronicle_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.id", "chronicle_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.type", "chronicle"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.description", "Send events to Google Chronicle"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.authentication_method", "serviceAccountSecret"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.gcp_instance", "123e4567-e89b-12d3-a456-426614174000"),
						resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.on_backpressure", "block"),
						// Cloudflare R2 destination
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "id", "cloudflare_r2_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "output_cloudflare_r2.id", "cloudflare_r2_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "output_cloudflare_r2.type", "cloudflare_r2"),
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "output_cloudflare_r2.description", "Write objects to Cloudflare R2"),
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "output_cloudflare_r2.bucket", "my-r2-bucket"),
						resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "output_cloudflare_r2.on_backpressure", "block"),
						// Databricks destination
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "id", "databricks_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "output_databricks.id", "databricks_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "output_databricks.type", "databricks"),
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "output_databricks.description", "Write data to Databricks"),
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "output_databricks.catalog", "main"),
						resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "output_databricks.on_backpressure", "block"),
						// Microsoft Fabric destination
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "id", "microsoft_fabric_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "output_microsoft_fabric.id", "microsoft_fabric_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "output_microsoft_fabric.type", "microsoft_fabric"),
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "output_microsoft_fabric.description", "Produce events to Microsoft Fabric"),
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "output_microsoft_fabric.topic", "app-events"),
						resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "output_microsoft_fabric.on_backpressure", "block"),
						// SentinelOne AI SIEM destination
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "id", "sentinel_one_ai_siem_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "output_sentinel_one_ai_siem.id", "sentinel_one_ai_siem_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "output_sentinel_one_ai_siem.type", "sentinel_one_ai_siem"),
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "output_sentinel_one_ai_siem.description", "Send events to SentinelOne AI SIEM"),
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "output_sentinel_one_ai_siem.base_url", "https://api.sentinelone.com"),
						resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "output_sentinel_one_ai_siem.on_backpressure", "block"),
					),
				},
			},
		})
	})
}
