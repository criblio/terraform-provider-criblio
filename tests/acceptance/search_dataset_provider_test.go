package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDatasetProvider(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	apiHTTPID := "tf_api_http_" + suffix
	elasticID := "tf_elastic_" + suffix
	s3ID := "tf_s3_" + suffix

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: searchDatasetProviderConfig(apiHTTPID, elasticID, s3ID, "created"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_searchdatasetprovider", "id", apiHTTPID),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_searchdatasetprovider", "apihttp.id", apiHTTPID),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_elastic_provider", "id", elasticID),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_elastic_provider", "api_elasticsearch.id", elasticID),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_s3_provider", "id", s3ID),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_s3_provider", "s3.id", s3ID),
					),
				},
				{
					Config: searchDatasetProviderConfig(apiHTTPID, elasticID, s3ID, "updated"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_searchdatasetprovider", "apihttp.description", "Example API HTTP search dataset provider updated"),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_elastic_provider", "api_elasticsearch.description", "Example Elasticsearch provider updated"),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_s3_provider", "s3.description", "Example S3 search dataset provider updated"),
					),
				},
				{
					Config:   searchDatasetProviderConfig(apiHTTPID, elasticID, s3ID, "updated"),
					PlanOnly: true,
				},
				{
					ResourceName:      "criblio_search_dataset_provider.my_searchdatasetprovider",
					ImportState:       true,
					ImportStateId:     apiHTTPID,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func searchDatasetProviderConfig(apiHTTPID, elasticID, s3ID, descriptionSuffix string) string {
	return fmt.Sprintf(`
resource "criblio_search_dataset_provider" "my_searchdatasetprovider" {
  apihttp = {
    authentication_method = "none"
    available_endpoints = [
      {
        data_field = "results"
        headers = [
          {
            name  = "Content-Type"
            value = "application/json"
          }
        ]
        method = "POST"
        name   = "search"
        url    = "http://localhost:8080/api/search"
      }
    ]
    description = "Example API HTTP search dataset provider %[4]s"
    id          = %[1]q
    type        = "api_http"
  }
}

resource "criblio_search_dataset_provider" "my_elastic_provider" {
  api_elasticsearch = {
    description = "Example Elasticsearch provider %[4]s"
    endpoint    = "https://localhost:9200"
    id          = %[2]q
    password    = "changeme"
    type        = "api_elasticsearch"
    username    = "elastic"
  }
}

resource "criblio_search_dataset_provider" "my_s3_provider" {
  s3 = {
    assume_role_arn           = "arn:aws:iam::123456789012:role/example-role"
    assume_role_external_id   = "example-external-id"
    aws_api_key               = "AKIAEXAMPLE"
    aws_authentication_method = "auto"
    aws_secret_key            = "example-secret"
    bucket                    = "example-bucket"
    bucket_path_suggestion    = "logs/"
    description               = "Example S3 search dataset provider %[4]s"
    enable_abac_tagging       = true
    enable_assume_role        = false
    endpoint                  = "https://s3.us-east-1.amazonaws.com"
    id                        = %[3]q
    region                    = "us-east-1"
    reject_unauthorized       = true
    reuse_connections         = false
    session_token             = "example-session-token"
    signature_version         = "v4"
    type                      = "s3"
  }
}
`, apiHTTPID, elasticID, s3ID, descriptionSuffix)
}
