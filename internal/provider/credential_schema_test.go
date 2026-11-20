package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestCredentialResourceSchemas(t *testing.T) {
	tests := []struct {
		name        string
		newResource func() resource.Resource
		paths       []string
		publicPaths []string
	}{
		{"source", NewSourceResource, []string{"input_http_raw.auth_tokens", "input_http.auth_tokens", "input_splunk_hec.auth_tokens.token", "input_tcp.auth_token", "input_kafka.sasl.password", "input_syslog.tls.passphrase", "items.input_http_raw.auth_tokens", "input_sentinelone.api_token", "input_netskope.api_token", "items.input_sentinelone.api_token", "items.input_netskope.api_token"}, []string{"input_splunk_hec.auth_tokens.token_secret"}},
		{"pack_source", NewPackSourceResource, []string{"input_http_raw.auth_tokens", "input_kafka.sasl.password", "items.input_http_raw.auth_tokens", "input_sentinelone.api_token", "input_netskope.api_token", "items.input_sentinelone.api_token", "items.input_netskope.api_token"}, nil},
		{"destination", NewDestinationResource, []string{"output_elastic.auth.manual_apikey", "output_google_chronicle.service_account_credentials", "output_s3.aws_secret_key"}, []string{"output_cribl_http.auth_tokens.token_secret"}},
		{"pack_destination", NewPackDestinationResource, []string{"output_elastic.auth.manual_apikey", "output_google_chronicle.service_account_credentials"}, nil},
		{"collector", NewCollectorResource, []string{"input_collector_rest.collector.conf.password", "input_collector_rest.collector.conf.token", "input_collector_s3.collector.conf.aws_secret_key"}, []string{"input_collector_s3.collector.conf.aws_secret"}},
		{"pack_collector", NewPackCollectorResource, []string{"input_collector_rest.collector.conf.password", "input_collector_s3.collector.conf.aws_secret_key"}, nil},
		{"notification_target", NewNotificationTargetResource, []string{"webhook_target.password", "webhook_target.token", "smtp_target.password", "slack_target.url"}, []string{"webhook_target.url"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response resource.SchemaResponse
			test.newResource().Schema(context.Background(), resource.SchemaRequest{}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			for _, path := range test.paths {
				if !credentialSchemaAttribute(t, response.Schema.Attributes, path).IsSensitive() {
					t.Errorf("%s must be sensitive", path)
				}
			}
			for _, path := range test.publicPaths {
				if credentialSchemaAttribute(t, response.Schema.Attributes, path).IsSensitive() {
					t.Errorf("%s must remain non-sensitive", path)
				}
			}
		})
	}
}

func credentialSchemaAttribute(t *testing.T, attributes map[string]schema.Attribute, path string) schema.Attribute {
	t.Helper()
	parts := strings.Split(path, ".")
	for index, name := range parts {
		attribute, ok := attributes[name]
		if !ok {
			t.Fatalf("schema missing %s at %s", name, path)
		}
		if index == len(parts)-1 {
			return attribute
		}
		switch attribute := attribute.(type) {
		case schema.SingleNestedAttribute:
			attributes = attribute.Attributes
		case schema.ListNestedAttribute:
			attributes = attribute.NestedObject.Attributes
		default:
			t.Fatalf("cannot descend into %s in %s", name, path)
		}
	}
	t.Fatalf("empty attribute path")
	return nil
}

func TestSourceAPITokenDataSourceSchemas(t *testing.T) {
	for _, test := range []struct {
		name          string
		newDataSource func() datasource.DataSource
	}{
		{"source", NewSourceDataSource},
		{"sources", NewSourcesDataSource},
		{"pack_source", NewPackSourceDataSource},
		{"pack_sources", NewPackSourcesDataSource},
	} {
		t.Run(test.name, func(t *testing.T) {
			var response datasource.SchemaResponse
			test.newDataSource().Schema(context.Background(), datasource.SchemaRequest{}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			items := response.Schema.Attributes["items"].(datasourceschema.ListNestedAttribute)
			for _, name := range []string{"input_sentinelone", "input_netskope"} {
				input := items.NestedObject.Attributes[name].(datasourceschema.SingleNestedAttribute)
				if !input.Attributes["api_token"].IsSensitive() {
					t.Errorf("items.%s.api_token must be sensitive", name)
				}
				if input.Attributes["text_secret"].IsSensitive() {
					t.Errorf("items.%s.text_secret is a reference and must remain non-sensitive", name)
				}
			}
		})
	}
}
