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
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "id", "cribl_http_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.id", "cribl_http_prod"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.type", "cribl_http"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.description", "Send events to Cribl Worker HTTP endpoint with retries"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.compression", "gzip"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.url", "https://edge.example.com:10200"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.load_balanced", "true"),
						resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.on_backpressure", "block"),
					),
				},
			},
		})
	})
}
