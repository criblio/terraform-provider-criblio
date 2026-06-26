package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDestination(t *testing.T) {
	testCases := map[string]struct {
		resourceName      string
		id                string
		createDescription string
		updateDescription string
		createConfig      func(string, string) string
		updateConfig      func(string, string) string
		importIgnore      []string
		check             resource.TestCheckFunc
	}{
		"s3": {
			resourceName:      "criblio_destination.s3",
			id:                "phase3-destination-s3",
			createDescription: "S3 destination",
			updateDescription: "S3 destination updated",
			createConfig:      destinationS3Config,
			updateConfig:      destinationS3Config,
			importIgnore:      []string{"output_s3.aws_secret_key"},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("criblio_destination.s3", "group_id", "default"),
				resource.TestCheckResourceAttr("criblio_destination.s3", "output_s3.type", "s3"),
				resource.TestCheckResourceAttr("criblio_destination.s3", "output_s3.bucket", "`cribl-destination-test`"),
			),
		},
		"cribl_http": {
			resourceName:      "criblio_destination.cribl_http",
			id:                "phase3-destination-cribl-http",
			createDescription: "Cribl HTTP destination",
			updateDescription: "Cribl HTTP destination updated",
			createConfig:      destinationCriblHTTPConfig,
			updateConfig:      destinationCriblHTTPConfig,
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("criblio_destination.cribl_http", "group_id", "default"),
				resource.TestCheckResourceAttr("criblio_destination.cribl_http", "output_cribl_http.type", "cribl_http"),
				resource.TestCheckResourceAttr("criblio_destination.cribl_http", "output_cribl_http.url", "https://edge.example.com:10200"),
			),
		},
		"splunk": {
			resourceName:      "criblio_destination.splunk",
			id:                "phase3-destination-splunk",
			createDescription: "Splunk destination",
			updateDescription: "Splunk destination updated",
			createConfig:      destinationSplunkConfig,
			updateConfig:      destinationSplunkConfig,
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("criblio_destination.splunk", "group_id", "default"),
				resource.TestCheckResourceAttr("criblio_destination.splunk", "output_splunk.type", "splunk"),
				resource.TestCheckResourceAttr("criblio_destination.splunk", "output_splunk.host", "localhost"),
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: providerFactory,
				Steps: []resource.TestStep{
					{
						Config: tc.createConfig(tc.id, tc.createDescription),
						Check:  tc.check,
					},
					{
						Config: tc.updateConfig(tc.id, tc.updateDescription),
						Check: resource.ComposeAggregateTestCheckFunc(
							tc.check,
							resource.TestCheckResourceAttr(tc.resourceName, destinationDescriptionAttr(name), tc.updateDescription),
						),
					},
					{
						Config:   tc.updateConfig(tc.id, tc.updateDescription),
						PlanOnly: true,
					},
					{
						ResourceName:            tc.resourceName,
						ImportState:             true,
						ImportStateId:           fmt.Sprintf(`{"group_id":"default","id":%q}`, tc.id),
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: tc.importIgnore,
					},
				},
			})
		})
	}
}

func TestDestinationExamples(t *testing.T) {
	config := destinationExampleConfig(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  destinationExampleChecks(),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

func destinationExampleConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join("..", "..", "examples", "destination", "main.tf")
	config, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read destination example: %v", err)
	}
	return string(config)
}

func destinationExampleChecks() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "id", "cribl_http_prod"),
		resource.TestCheckResourceAttr("criblio_destination.my_cribl_http_destination", "output_cribl_http.type", "cribl_http"),

		resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "id", "chronicle_prod"),
		resource.TestCheckResourceAttr("criblio_destination.my_chronicle_destination", "output_chronicle.type", "chronicle"),

		resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "id", "cloudflare_r2_prod"),
		resource.TestCheckResourceAttr("criblio_destination.my_cloudflare_r2_destination", "output_cloudflare_r2.type", "cloudflare_r2"),

		resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "id", "databricks_prod"),
		resource.TestCheckResourceAttr("criblio_destination.my_databricks_destination", "output_databricks.type", "databricks"),

		resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "id", "microsoft_fabric_prod"),
		resource.TestCheckResourceAttr("criblio_destination.my_microsoft_fabric_destination", "output_microsoft_fabric.type", "microsoft_fabric"),

		resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "id", "sentinel_one_ai_siem_prod"),
		resource.TestCheckResourceAttr("criblio_destination.my_sentinel_one_ai_siem_destination", "output_sentinel_one_ai_siem.type", "sentinel_one_ai_siem"),

		resource.TestCheckResourceAttr("criblio_destination.grafana_cloud", "group_id", "default"),
		resource.TestCheckResourceAttr("criblio_destination.grafana_cloud", "id", "CriblCloud"),
		resource.TestCheckResourceAttr("criblio_destination.grafana_cloud", "output_grafana_cloud.type", "grafana_cloud"),
	)
}

func destinationDescriptionAttr(name string) string {
	switch name {
	case "s3":
		return "output_s3.description"
	case "cribl_http":
		return "output_cribl_http.description"
	case "splunk":
		return "output_splunk.description"
	default:
		return "description"
	}
}

func destinationS3Config(id, description string) string {
	return fmt.Sprintf(`
resource "criblio_destination" "s3" {
  group_id = "default"
  id       = %[1]q

  output_s3 = {
    id              = %[1]q
    type            = "s3"
    description     = %[2]q
    bucket          = "`+"`cribl-destination-test`"+`"
    region          = "us-east-1"
    stage_path      = "$CRIBL_HOME/state/outputs/%[1]s"
    aws_api_key     = "AKIAIOSFODNN7EXAMPLE"
    aws_secret_key  = "test-secret-key"
    format          = "json"
    compress        = "gzip"
    on_backpressure = "block"
    pipeline        = "passthru"
  }
}
`, id, description)
}

func destinationCriblHTTPConfig(id, description string) string {
	return fmt.Sprintf(`
resource "criblio_destination" "cribl_http" {
  group_id = "default"
  id       = %[1]q

  output_cribl_http = {
    id              = %[1]q
    type            = "cribl_http"
    description     = %[2]q
    url             = "https://edge.example.com:10200"
    compression     = "gzip"
    on_backpressure = "block"
    pipeline        = "passthru"
  }
}
`, id, description)
}

func destinationSplunkConfig(id, description string) string {
	return fmt.Sprintf(`
resource "criblio_destination" "splunk" {
  group_id = "default"
  id       = %[1]q

  output_splunk = {
    id              = %[1]q
    type            = "splunk"
    description     = %[2]q
    host            = "localhost"
    port            = 9997
    on_backpressure = "block"
    pipeline        = "passthru"
  }
}
`, id, description)
}
