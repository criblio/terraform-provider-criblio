package export

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

type credentialTransport struct{ body string }

func (tr credentialTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"items":[` + tr.body + `]}`)), Request: req}, nil
}

func TestExporterOmitsCredentialLiterals(t *testing.T) {
	reg := buildTestRegistry(t)
	tests := []struct{ name, typeName, body, field string }{
		{"HTTP Raw tokens", "criblio_source", `{"id":"example","type":"http_raw","host":"0.0.0.0","port":10080,"authTokens":["ISSUE236_DUMMY"]}`, "auth_tokens"},
		{"HTTP tokens", "criblio_source", `{"id":"example","type":"http","authTokens":["ISSUE236_DUMMY"]}`, "auth_tokens"},
		{"Pack HTTP Raw tokens", "criblio_pack_source", `{"id":"example","type":"http_raw","host":"0.0.0.0","port":10080,"authTokens":["ISSUE236_DUMMY"]}`, "auth_tokens"},
		{"Splunk HEC token", "criblio_source", `{"id":"example","type":"splunk_hec","authTokens":[{"token":"ISSUE236_DUMMY","authType":"manual"}]}`, "token"},
		{"TCP auth token", "criblio_source", `{"id":"example","type":"tcp","authToken":"ISSUE236_DUMMY"}`, "auth_token"},
		{"Kafka SASL password", "criblio_source", `{"id":"example","type":"kafka","sasl":{"password":"ISSUE236_DUMMY"}}`, "password"},
		{"Syslog TLS passphrase", "criblio_source", `{"id":"example","type":"syslog","tls":{"passphrase":"ISSUE236_DUMMY"}}`, "passphrase"},
		{"Elastic API key", "criblio_destination", `{"id":"example","type":"elastic","auth":{"manualAPIKey":"ISSUE236_DUMMY","authType":"manualAPIKey"}}`, "manual_apikey"},
		{"Chronicle credentials", "criblio_destination", `{"id":"example","type":"google_chronicle","serviceAccountCredentials":"ISSUE236_DUMMY"}`, "service_account_credentials"},
		{"S3 destination secret key", "criblio_destination", `{"id":"example","type":"s3","awsSecretKey":"ISSUE236_DUMMY"}`, "aws_secret_key"},
		{"REST collector password", "criblio_collector", `{"id":"example","type":"collection","collector":{"type":"rest","conf":{"password":"ISSUE236_DUMMY"}}}`, "password"},
		{"REST collector token", "criblio_collector", `{"id":"example","type":"collection","collector":{"type":"rest","conf":{"token":"ISSUE236_DUMMY"}}}`, "token"},
		{"S3 collector secret key", "criblio_collector", `{"id":"example","type":"collection","collector":{"type":"s3","conf":{"awsSecretKey":"ISSUE236_DUMMY"}}}`, "aws_secret_key"},
		{"Pack REST collector token", "criblio_pack_collector", `{"id":"example","type":"collection","collector":{"type":"rest","conf":{"token":"ISSUE236_DUMMY"}}}`, "token"},
		{"Pack Elastic API key", "criblio_pack_destination", `{"id":"example","type":"elastic","auth":{"manualAPIKey":"ISSUE236_DUMMY","authType":"manualAPIKey"}}`, "manual_apikey"},
		{"Webhook target password", "criblio_notification_target", `{"id":"example","type":"webhook","password":"ISSUE236_DUMMY"}`, "password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, ok := reg.ByTypeName(tt.typeName)
			if !ok {
				t.Fatal("missing registry entry")
			}
			client := &importclient.Client{REST: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "dummy", HTTPClient: &http.Client{Transport: credentialTransport{body: tt.body}}})}
			ids := map[string]string{"group_id": "default", "id": "example"}
			if strings.Contains(tt.typeName, "pack_") {
				ids["pack"] = "example-pack"
			}
			if tt.typeName == "criblio_pack_destination" {
				ids["group_id"] = "test-group"
			}
			item, skip := convertOneResource(context.Background(), client, discovery.Result{TypeName: tt.typeName}, e, ids, nil, []string{"default"}, false, IncludeOverride{}, &ExportResult{}, &conversionRequest{})
			if item == nil {
				t.Fatalf("conversion skipped: %s", skip)
			}
			dir := t.TempDir()
			if err := generator.WriteModuleDirectory(dir, []generator.ResourceItem{*item}, hcl.DefaultResourceBlockOptions()); err != nil {
				t.Fatal(err)
			}
			main, err := os.ReadFile(filepath.Join(generator.ModuleDir(dir, tt.typeName), "main.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(main), tt.field) {
				t.Fatalf("credential field missing in output: %s", main)
			}
			if strings.Contains(string(main), "ISSUE236_DUMMY") {
				t.Fatalf("credential leaked in generated main.tf: %s", main)
			}
			variables, err := os.ReadFile(filepath.Join(generator.ModuleDir(dir, tt.typeName), "variables.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(variables), "sensitive = true") || strings.Contains(string(variables), "ISSUE236_DUMMY") {
				t.Fatalf("expected sensitive variables without credential defaults: %s", variables)
			}
			if len(hcl.CollectSecretVariableNames(item.Attrs)) == 0 {
				t.Fatal("expected a credential variable reference")
			}
			if strings.Contains(string(main), "jsonencode(var.") {
				t.Fatalf("string credentials must not be JSON-encoded: %s", main)
			}
			infos := generator.RootModuleInfosFromItems([]generator.ResourceItem{*item})
			if err := generator.WriteRootFiles(dir, infos, []generator.ResourceItem{*item}, false); err != nil {
				t.Fatal(err)
			}
			if err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if entry.IsDir() {
					return nil
				}
				contents, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if strings.Contains(string(contents), "ISSUE236_DUMMY") {
					t.Errorf("credential leaked in %s", path)
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestExporterPreservesSecretReferences(t *testing.T) {
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_source")
	if !ok {
		t.Fatal("missing source registry entry")
	}
	client := &importclient.Client{REST: restclient.New(restclient.Config{
		BaseURL: "https://example.invalid", BearerToken: "dummy",
		HTTPClient: &http.Client{Transport: credentialTransport{body: `{"id":"example","type":"splunk_hec","authTokens":[{"tokenSecret":"stored-token-name","authType":"secret"}]}`}},
	})}
	item, skip := convertOneResource(context.Background(), client, discovery.Result{TypeName: entry.TypeName}, entry,
		map[string]string{"group_id": "default", "id": "example"}, nil, []string{"default"}, false, IncludeOverride{}, &ExportResult{}, &conversionRequest{})
	if item == nil {
		t.Fatalf("conversion skipped: %s", skip)
	}
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(file.Bytes()), "stored-token-name") || len(hcl.CollectSecretVariableNames(item.Attrs)) != 0 {
		t.Fatalf("secret names must remain literal references: %s", file.Bytes())
	}
}

func TestExporterHonorsSchemaSensitivityInSharedConversion(t *testing.T) {
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_source")
	if !ok {
		t.Fatal("missing source registry entry")
	}
	// This name is deliberately absent from the exporter's credential heuristics.
	entry.SensitivePaths["input_http_raw.host"] = true
	model, err := converter.ConvertRawItemWithIdentifiers(entry,
		json.RawMessage(`{"id":"example","type":"http_raw","host":"ISSUE236_DUMMY","port":10080}`),
		toRequestParams(map[string]string{"group_id": "default", "id": "example"}))
	if err != nil {
		t.Fatal(err)
	}
	var out ExportResult
	if err := appendResourceItemFromModel(&out, entry.TypeName, entry, map[string]string{"group_id": "default", "id": "example"}, model, nil, []string{"default"}, false, IncludeOverride{}); err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 || out.Items[0].Attrs["input_http_raw"].Map["host"].Kind != hcl.KindVariableRef {
		t.Fatal("schema-sensitive field must become a variable in the shared conversion path")
	}
}
