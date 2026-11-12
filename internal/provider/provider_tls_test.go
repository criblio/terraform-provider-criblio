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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProviderReconfigurationReusesAuthentication(t *testing.T) {
	for _, mode := range []string{"onprem", "cloud"} {
		t.Run(mode, func(t *testing.T) {
			for _, trust := range []string{"default", "skip", "custom-ca"} {
				t.Run(trust, func(t *testing.T) {
					t.Setenv("CRIBL_BEARER_TOKEN", "")
					t.Setenv("CRIBL_CLIENT_ID", "test-client")
					t.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
					t.Setenv("CRIBL_ONPREM_USERNAME", "test-user")
					t.Setenv("CRIBL_ONPREM_PASSWORD", "test-password")
					t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "")
					t.Setenv("CRIBL_CA_CERT_FILE", "")
					auth.ClearTokenCache()
					t.Cleanup(auth.ClearTokenCache)
					var logins, apiConnectionCount atomic.Int32
					var apiConnections sync.Map
					handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var body any
						switch r.URL.Path {
						case "/api/v1/auth/login", "/oauth/token":
							logins.Add(1)
							body = map[string]any{"token": "test-token", "access_token": "test-token", "expires_in": 3600}
						case "/api/v1/master/groups/default":
							if _, loaded := apiConnections.LoadOrStore(r.RemoteAddr, struct{}{}); !loaded {
								apiConnectionCount.Add(1)
							}
							body = map[string]any{"id": "default"}
						default:
							t.Errorf("unexpected path: %s", r.URL.Path)
							http.NotFound(w, r)
							return
						}
						if err := json.NewEncoder(w).Encode(body); err != nil {
							t.Error(err)
						}
					})
					var server *httptest.Server
					if trust == "default" || mode == "cloud" {
						server = httptest.NewServer(handler)
					} else {
						server = httptest.NewTLSServer(handler)
					}
					defer server.Close()
					onpremURL := ""
					if mode == "onprem" {
						onpremURL = server.URL
					}
					t.Setenv("CRIBL_ONPREM_SERVER_URL", onpremURL)
					values := map[string]any{"server_url": server.URL, "cloud_domain": server.URL}
					switch trust {
					case "skip":
						values["insecure_skip_verify"] = true
					case "custom-ca":
						if mode == "cloud" {
							values["ca_cert_file"] = "/missing/onprem-only-ca.pem"
						} else {
							values["ca_cert_file"] = writeTLSCAFile(t, server)
						}
					}
					for range 3 {
						response := configureTLSProvider(t, values)
						if response.Diagnostics.HasError() {
							t.Fatal(response.Diagnostics)
						}
						clients := response.ResourceData.(*ProviderClients)
						if _, err := restclient.GetRaw(context.Background(), clients.RC, "/master/groups/default"); err != nil {
							t.Fatal(err)
						}
					}
					if (mode == "cloud" || trust == "default") && apiConnectionCount.Load() != 1 {
						t.Fatalf("API opened %d connections across configurations, want 1 reused connection", apiConnectionCount.Load())
					}
					if got := logins.Load(); got != 1 {
						t.Fatalf("provider reconfiguration fetched %d tokens, want 1 cached token", got)
					}
				})
			}
		})
	}
}

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
					if mode == "cloud-oauth" || trust == "default" || trust == "explicit-false" {
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
					if mode == "onprem-bearer" || trust == "env-skip" || trust == "env-ca" {
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
	t.Setenv("CRIBL_CLIENT_ID", "test-client")
	t.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	t.Setenv("CRIBL_ONPREM_SERVER_URL", "https://onprem.example")
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

func TestProviderCloudIgnoresOnPremTLSConfiguration(t *testing.T) {
	t.Setenv("CRIBL_CLIENT_ID", "test-client")
	t.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	t.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "invalid-onprem-setting")
	t.Setenv("CRIBL_CA_CERT_FILE", "/missing/onprem-ca.pem")
	server := httptest.NewTLSServer(http.NotFoundHandler())
	defer server.Close()
	for _, values := range []map[string]any{
		{},
		{"insecure_skip_verify": true},
		{"ca_cert_file": "/missing/explicit-onprem-ca.pem"},
		{"insecure_skip_verify": tftypes.UnknownValue, "ca_cert_file": tftypes.UnknownValue},
	} {
		values["server_url"] = server.URL
		values["cloud_domain"] = server.URL
		values["bearer_auth"] = "test-token"
		response := configureTLSProvider(t, values)
		if response.Diagnostics.HasError() {
			t.Fatalf("cloud validated on-prem TLS settings: %v", response.Diagnostics)
		}
		clients := response.ResourceData.(*ProviderClients)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := restclient.GetRaw(ctx, clients.RC, "/master/groups/default")
		cancel()
		if err == nil || !strings.Contains(err.Error(), "certificate signed by unknown authority") {
			t.Fatalf("cloud request error = %v, want default TLS verification", err)
		}
	}
}

func TestProviderTLSTransportIsolation(t *testing.T) {
	t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "invalid")
	t.Setenv("CRIBL_CA_CERT_FILE", "/missing/env-ca.pem")
	// Explicit provider settings override environment values.
	server := httptest.NewTLSServer(http.NotFoundHandler())
	defer server.Close()
	model := CriblioProviderModel{InsecureSkipVerify: types.BoolValue(false), CACertFile: types.StringValue(writeTLSCAFile(t, server))}
	transport, _, err := model.httpTransport()
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
	other, _, err := model.httpTransport()
	if err != nil {
		t.Fatal(err)
	}
	defer other.CloseIdleConnections()
	if other.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("changing one provider affected another provider")
	}
}

func TestProviderTLSAuthenticationScopeUsesCAContents(t *testing.T) {
	t.Setenv("CRIBL_INSECURE_SKIP_VERIFY", "")
	t.Setenv("CRIBL_CA_CERT_FILE", "")
	server := httptest.NewTLSServer(http.NotFoundHandler())
	defer server.Close()
	firstFile := writeTLSCAFile(t, server)
	secondFile := writeTLSCAFile(t, server)
	scopeFor := func(filename string) string {
		t.Helper()
		model := CriblioProviderModel{CACertFile: types.StringValue(filename)}
		transport, scope, err := model.httpTransport()
		if err != nil {
			t.Fatal(err)
		}
		transport.CloseIdleConnections()
		return scope
	}
	firstScope := scopeFor(firstFile)
	if secondScope := scopeFor(secondFile); secondScope != firstScope {
		t.Fatal("identical CA contents at different paths have different scopes")
	}
	contents, err := os.ReadFile(secondFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(secondFile, append(contents, []byte("\n# updated CA bundle\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	if updatedScope := scopeFor(secondFile); updatedScope == firstScope {
		t.Fatal("replacing CA contents at the same path reused the old scope")
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
	transport, _, err := model.httpTransport()
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
