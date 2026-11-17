package tests

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchSourceHTTPRawAuthTokens(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests require TF_ACC")
	}
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Search sources require Cribl Cloud")
	}
	portValue := os.Getenv("CRIBL_TEST_SEARCH_SOURCE_HTTP_RAW_PORT")
	if portValue == "" {
		t.Skip("Set CRIBL_TEST_SEARCH_SOURCE_HTTP_RAW_PORT to an unused provisioned Search ingest port")
	}
	port, err := strconv.Atoi(portValue)
	if err != nil || port < 1 || port > 65535 {
		t.Fatal("CRIBL_TEST_SEARCH_SOURCE_HTTP_RAW_PORT must be a valid unused provisioned port")
	}
	id := acctest.RandomWithPrefix("test_raw_auth")
	const name = "criblio_search_source.test"
	config := func(tokens string) string {
		assignment := ""
		if tokens != "" {
			assignment = "auth_tokens = " + tokens
		}
		return fmt.Sprintf(`resource "criblio_search_source" "test" {
  id = %q
  type = "http_raw"
  host = "0.0.0.0"
  port = %d
  description = "Raw HTTP authentication acceptance test"
  tls = { disabled = false }
  %s
}`, id, port, assignment)
	}
	firstToken := acctest.RandStringFromCharSet(32, acctest.CharSetAlphaNum)
	rotatedToken := acctest.RandStringFromCharSet(32, acctest.CharSetAlphaNum)
	first := config(fmt.Sprintf(`[{enabled = true, description = "test token", token = %q}]`, firstToken))
	rotated := config(fmt.Sprintf(`[{token = %q}]`, rotatedToken))
	empty := config("[]")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{Config: config("null"), Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "0")},
			{Config: first, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "auth_tokens.#", "1"),
				resource.TestCheckResourceAttr(name, "auth_tokens.0.token", firstToken),
				resource.TestCheckResourceAttr(name, "auth_tokens.0.enabled", "true"),
			)},
			{Config: first, PlanOnly: true},
			{Config: rotated, Check: resource.TestCheckResourceAttr(name, "auth_tokens.0.token", rotatedToken)},
			{Config: rotated, PlanOnly: true},
			// Import reads the server's encrypted/masked token, not its plaintext.
			{ResourceName: name, ImportState: true, ImportStateId: id, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"auth_tokens"}},
			{Config: empty, Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "0")},
			{Config: empty, PlanOnly: true},
			{Config: first, Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "1")},
			{Config: config("null"), Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "0")},
			{Config: config("null"), PlanOnly: true},
			{Config: first, Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "1")},
			{Config: config(""), Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "0")},
			{Config: config(""), PlanOnly: true},
		},
	})
}
