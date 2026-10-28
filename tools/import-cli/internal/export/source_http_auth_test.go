package export

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Plan exported configurations without loading local credentials or using an API.
type sourceAuthExportTestProvider struct{ provider.CriblioProvider }

func (p *sourceAuthExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func (p *sourceAuthExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewSourceResource, provider.NewPackSourceResource}
}

func TestHTTPRawAuthExportConfiguration(t *testing.T) {
	reg := buildTestRegistry(t)
	for _, typeName := range []string{"criblio_source", "criblio_pack_source"} {
		t.Run(typeName, func(t *testing.T) {
			entry, ok := reg.ByTypeName(typeName)
			if !ok {
				t.Fatal("missing source registry entry")
			}
			ids := map[string]string{"group_id": "default", "id": "http-auth-test"}
			if typeName == "criblio_pack_source" {
				ids["pack"] = "test-pack"
			}
			var steps []resource.TestStep
			for _, test := range []struct {
				name  string
				ext   string
				count int
			}{
				{"empty placeholders", `,"authTokensExt":[{}, {"authType":null,"token":null,"tokenSecret":null,"description":null,"metadata":null}]`, 0},
				{"mixed", `,"authTokensExt":[{}, {"authType":"secret","tokenSecret":"stored-token"}, {"token":"HTTP_EXT_DUMMY","description":"legacy extended token","metadata":[{"name":"env","value":"test"}]}]`, 2},
				{"without extended tokens", "", 0},
			} {
				body := `{"id":"http-auth-test","type":"http_raw","host":"0.0.0.0","port":10080,"authTokens":["HTTP_LEGACY_DUMMY"]` + test.ext + `}`
				client := &importclient.Client{REST: restclient.New(restclient.Config{
					BaseURL: "https://example.invalid", BearerToken: "test-token",
					HTTPClient: &http.Client{Transport: credentialTransport{body: body}},
				})}
				item, skip := convertOneResource(context.Background(), client, discovery.Result{TypeName: typeName}, entry, ids, nil, []string{"default"}, false, IncludeOverride{}, &ExportResult{}, &conversionRequest{})
				if item == nil {
					t.Fatalf("%s: conversion skipped: %s", test.name, skip)
				}
				// The shared conversion path must produce the same token attributes.
				model, err := converter.ConvertRawItemWithIdentifiers(entry, json.RawMessage(body), toRequestParams(ids))
				if err != nil {
					t.Fatal(err)
				}
				var out ExportResult
				if err := appendResourceItemFromModel(&out, typeName, entry, ids, model, nil, []string{"default"}, false, IncludeOverride{}); err != nil {
					t.Fatal(err)
				}
				if len(out.Items) != 1 {
					t.Fatalf("%s: shared conversion returned %d items", test.name, len(out.Items))
				}
				for _, attrs := range []map[string]hcl.Value{item.Attrs, out.Items[0].Attrs} {
					block := attrs["input_http_raw"].Map
					pairs := block["auth_tokens_ext"]
					if len(pairs.List) != test.count {
						t.Fatalf("%s: exported %d extended tokens, want %d", test.name, len(pairs.List), test.count)
					}
					if test.count > 0 {
						if pairs.List[0].Map["token_secret"].String != "stored-token" || pairs.List[1].Map["token"].Kind != hcl.KindVariableRef {
							t.Fatal("export lost extended credentials or token sensitivity")
						}
						if value := pairs.List[1].Map["auth_type"]; !value.IsNull() {
							t.Fatal("export invented an auth_type for a legacy extended token")
						}
						if len(pairs.List[1].Map["metadata"].List) != 1 {
							t.Fatal("export lost token metadata")
						}
					}
				}
				file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
				if err != nil {
					t.Fatal(err)
				}
				config := string(file.Bytes())
				if strings.Contains(config, "HTTP_LEGACY_DUMMY") || strings.Contains(config, "HTTP_EXT_DUMMY") {
					t.Fatal("export exposed a credential literal")
				}
				if test.count == 0 && strings.Contains(config, "auth_tokens_ext") {
					t.Fatalf("%s: export emitted empty placeholders", test.name)
				}
				for _, name := range hcl.CollectSecretVariableNames(item.Attrs) {
					config += fmt.Sprintf("\nvariable %q {\n type = string\n sensitive = true\n default = \"test-token\"\n}\n", name)
				}
				steps = append(steps, resource.TestStep{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"criblio": providerserver.NewProtocol6WithError(&sourceAuthExportTestProvider{}),
				},
				Steps: steps,
			})
		})
	}
}
