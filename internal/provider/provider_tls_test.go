package provider

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProviderTLSRequests(t *testing.T) {
	for _, mode := range []string{"onprem-login", "onprem-bearer", "cloud-oauth"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("CRIBL_BEARER_TOKEN", "")
			t.Setenv("CRIBL_CLIENT_ID", "test-client")
			t.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
			t.Setenv("CRIBL_ONPREM_USERNAME", "test-user")
			t.Setenv("CRIBL_ONPREM_PASSWORD", "test-password")
			t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "")
			t.Setenv("CRIBL_CA_CERT_FILE", "")
			var loginCount, apiCount atomic.Int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var response any
				switch r.URL.Path {
				case "/api/v1/auth/login", "/oauth/token":
					loginCount.Add(1)
					if r.Method != http.MethodPost {
						t.Errorf("authentication method = %s, want POST", r.Method)
					}
					response = map[string]any{"token": "test-token", "access_token": "test-token", "expires_in": 3600}
				case "/api/v1/master/groups/default":
					apiCount.Add(1)
					if r.Header.Get("Authorization") != "Bearer test-token" {
						t.Errorf("unexpected authorization header: %q", r.Header.Get("Authorization"))
					}
					response = map[string]any{"id": "default"}
				default:
					t.Errorf("unexpected request path: %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			onpremURL := server.URL
			if mode == "cloud-oauth" {
				onpremURL = ""
			}
			t.Setenv("CRIBL_ONPREM_SERVER_URL", onpremURL)
			caFile := writeTLSCAFile(t, server)
			for _, trust := range []string{"default", "skip", "custom-ca", "env-skip", "env-ca", "explicit-false"} {
				t.Run(trust, func(t *testing.T) {
					values := map[string]any{"server_url": server.URL, "cloud_domain": server.URL}
					if mode == "onprem-bearer" {
						values["bearer_auth"] = "test-token"
					}
					switch trust {
					case "skip":
						values["insecure_skip_verify"] = true
					case "custom-ca":
						t.Setenv("CRIBL_CA_CERT_FILE", "/missing/env-ca.pem")
						values["ca_cert_file"] = caFile
					case "env-skip":
						t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "true")
					case "env-ca":
						t.Setenv("CRIBL_CA_CERT_FILE", caFile)
					case "explicit-false":
						t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "true")
						values["insecure_skip_verify"] = false
					}
					configured := configureTLSProvider(t, values)
					if configured.Diagnostics.HasError() {
						t.Fatal(configured.Diagnostics)
					}
					clients := configured.ResourceData.(*ProviderClients)
					beforeLogin, beforeAPI := loginCount.Load(), apiCount.Load()
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					_, err := restclient.GetRaw(ctx, clients.RC, "/master/groups/default")
					if trust == "default" || trust == "explicit-false" {
						if err == nil || !strings.Contains(err.Error(), "certificate signed by unknown authority") {
							t.Fatalf("request error = %v, want untrusted certificate", err)
						}
						if loginCount.Load() != beforeLogin || apiCount.Load() != beforeAPI {
							t.Fatal("untrusted TLS connection reached the HTTP handler")
						}
						return
					}
					if err != nil {
						t.Fatal(err)
					}
					// A second request should reuse the cached token with this client.
					if _, err := restclient.GetRaw(ctx, clients.RC, "/master/groups/default"); err != nil {
						t.Fatal(err)
					}
					wantLogins := int32(1)
					if mode == "onprem-bearer" {
						wantLogins = 0
					}
					if loginCount.Load()-beforeLogin != wantLogins || apiCount.Load()-beforeAPI != 2 {
						t.Fatalf("login/API requests = %d/%d, want %d/2", loginCount.Load()-beforeLogin, apiCount.Load()-beforeAPI, wantLogins)
					}
				})
			}
		})
	}
}

func TestProviderTLSConfigurationErrors(t *testing.T) {
	invalidCA := filepath.Join(t.TempDir(), "invalid.pem")
	if err := os.WriteFile(invalidCA, []byte("not a certificate"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, envName, envValue, want string
		values                        map[string]any
	}{
		{"invalid-boolean", "CRIBL_INSECURE_SKIP_VERIFY", "yes", "must be a boolean", nil},
		{"missing-ca", "", "", "read ca_cert_file", map[string]any{"ca_cert_file": filepath.Join(t.TempDir(), "missing.pem")}},
		{"invalid-ca", "", "", "contains no valid PEM certificates", map[string]any{"ca_cert_file": invalidCA}},
		{"unknown-skip", "", "", "must be known", map[string]any{"insecure_skip_verify": tftypes.UnknownValue}},
		{"unknown-ca", "", "", "must be known", map[string]any{"ca_cert_file": tftypes.UnknownValue}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "")
			t.Setenv("CRIBL_CA_CERT_FILE", "")
			if test.envName != "" {
				t.Setenv(test.envName, test.envValue)
			}
			response := configureTLSProvider(t, test.values)
			if !response.Diagnostics.HasError() || !strings.Contains(response.Diagnostics[0].Detail(), test.want) {
				t.Fatalf("diagnostics = %v, want %q", response.Diagnostics, test.want)
			}
			if response.ResourceData != nil {
				t.Fatal("client created despite invalid TLS configuration")
			}
		})
	}
}

func TestProviderTLSTransportIsolation(t *testing.T) {
	t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "invalid")
	t.Setenv("CRIBL_CA_CERT_FILE", "/missing/env-ca.pem")
	// Explicit false and an empty CA path override environment values.
	model := CriblioProviderModel{InsecureSkipVerify: types.BoolValue(false), CACertFile: types.StringValue("")}
	transport, err := model.httpTransport()
	if err != nil {
		t.Fatal(err)
	}
	defer transport.CloseIdleConnections()
	if transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("explicit false did not override the environment")
	}
	base := http.DefaultTransport.(*http.Transport)
	if transport == base || (base.TLSClientConfig != nil && base.TLSClientConfig == transport.TLSClientConfig) {
		t.Fatal("provider shares mutable transport or TLS configuration with the default transport")
	}
	transport.TLSClientConfig.InsecureSkipVerify = true
	other, err := model.httpTransport()
	if err != nil {
		t.Fatal(err)
	}
	defer other.CloseIdleConnections()
	if other.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("changing one provider affected another provider")
	}
}

func TestProviderTLSCustomCAPreservesHostnameVerification(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request with invalid hostname reached the HTTP handler")
	}))
	defer server.Close()
	model := CriblioProviderModel{
		InsecureSkipVerify: types.BoolValue(false),
		CACertFile:         types.StringValue(writeTLSCAFile(t, server)),
	}
	transport, err := model.httpTransport()
	if err != nil {
		t.Fatal(err)
	}
	defer transport.CloseIdleConnections()
	transport.TLSClientConfig.ServerName = "wrong.example.internal"
	client := &http.Client{Transport: transport}
	response, err := client.Get(server.URL)
	if response != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}
	if err == nil || !strings.Contains(err.Error(), "wrong.example.internal") {
		t.Fatalf("request error = %v, want hostname verification failure", err)
	}
}

func configureTLSProvider(t *testing.T, values map[string]any) provider.ConfigureResponse {
	t.Helper()
	ctx := context.Background()
	p := &CriblioProvider{}
	var schema provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schema)
	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	fields := make(map[string]tftypes.Value, len(typ.AttributeTypes))
	for name, fieldType := range typ.AttributeTypes {
		fields[name] = tftypes.NewValue(fieldType, values[name])
	}
	var response provider.ConfigureResponse
	p.Configure(ctx, provider.ConfigureRequest{Config: tfsdk.Config{
		Schema: schema.Schema,
		Raw:    tftypes.NewValue(typ, fields),
	}}, &response)
	return response
}

func writeTLSCAFile(t *testing.T, server *httptest.Server) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "ca.pem")
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	if err := os.WriteFile(filename, certificate, 0600); err != nil {
		t.Fatal(err)
	}
	return filename
}
