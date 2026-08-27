package tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCollectorS3InputRemainsOptional(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{{
			Config: `
resource "criblio_collector" "s3" {
  id       = "collector-without-input"
  group_id = "default"

  input_collector_s3 = {
    collector = {
      type = "s3"
      conf = {
        aws_authentication_method = "auto"
        bucket                    = "test-bucket"
      }
    }
  }
}`,
			PlanOnly: true,
		}},
	})
}

func TestCollectorS3AcceptsValidInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{{
			Config: `
resource "criblio_collector" "s3" {
  id       = "collector-with-input"
  group_id = "default"

  input_collector_s3 = {
    input = {
      type = "collection"
    }
    collector = {
      type = "s3"
      conf = {
        aws_authentication_method = "auto"
        bucket                    = "test-bucket"
      }
    }
  }
}`,
			PlanOnly: true,
		}},
	})
}

func TestCollectorRejectsInvalidInputType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{{
			Config: `
resource "criblio_collector" "s3" {
  id       = "collector-invalid-input-type"
  group_id = "default"

  input_collector_s3 = {
    input = {
      type = "invalid"
    }
    collector = {
      type = "s3"
      conf = {
        aws_authentication_method = "auto"
        bucket                    = "test-bucket"
      }
    }
  }
}`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`input_collector_s3\.input\.type`),
		}},
	})
}

func TestCollectorRestInputRemainsOptional(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{{
			Config: `
resource "criblio_collector" "rest" {
  id       = "legacy-rest-discovery"
  group_id = "default"

  input_collector_rest = {
    collector = {
      type = "rest"
      conf = {
        authentication = "none"
        collect_method = "get"
        collect_url    = "'https://example.com/data'"
      }
    }
  }
}`,
			PlanOnly: true,
		}},
	})
}
