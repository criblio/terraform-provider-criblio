package restclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type testItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestMethods(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Authorization = %q, expected Bearer test-token", r.Header.Get("Authorization"))
		}

		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("GET path = %q", r.URL.Path)
			}
			writeJSON(t, w, testItem{ID: "cert-1", Name: "from-get"})
		case http.MethodPost:
			switch r.URL.Path {
			case "/api/v1/system/certificates":
				assertJSONBody(t, r, "from-post")
				writeJSON(t, w, testItem{ID: "cert-2", Name: "from-post"})
			case "/api/v1/system/no-response":
				assertJSONBody(t, r, "from-post-no-response")
				writeJSON(t, w, map[string]any{"items": []testItem{}})
			case "/api/v1/system/files":
				mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if err != nil {
					t.Errorf("failed to parse content type: %v", err)
				}
				if mediaType != "multipart/form-data" {
					t.Errorf("upload content type = %q, expected multipart/form-data", mediaType)
				}
				if err := r.ParseMultipartForm(1024); err != nil {
					t.Errorf("failed to parse multipart form: %v", err)
				}
				file, header, err := r.FormFile("file")
				if err != nil {
					t.Errorf("missing upload file: %v", err)
					return
				}
				defer file.Close()
				content, err := io.ReadAll(file)
				if err != nil {
					t.Errorf("failed to read upload file: %v", err)
				}
				if header.Filename != "lookup.csv" {
					t.Errorf("upload filename = %q, expected lookup.csv", header.Filename)
				}
				if string(content) != "a,b\n" {
					t.Errorf("upload content = %q, expected a,b\\n", string(content))
				}
				w.WriteHeader(http.StatusNoContent)
			default:
				t.Errorf("unexpected POST path %q", r.URL.Path)
			}
		case http.MethodPatch:
			if r.URL.Path == "/api/v1/system/no-response" {
				assertJSONBody(t, r, "from-patch-no-response")
				writeJSON(t, w, map[string]any{"items": []testItem{}})
				return
			}
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("PATCH path = %q", r.URL.Path)
			}
			assertJSONBody(t, r, "from-patch")
			writeJSON(t, w, testItem{ID: "cert-1", Name: "from-patch"})
		case http.MethodPut:
			if r.URL.Path == "/api/v1/system/files" {
				if r.Header.Get("Content-Type") != "text/csv" {
					t.Errorf("raw upload content type = %q, expected text/csv", r.Header.Get("Content-Type"))
				}
				content, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("failed to read raw upload body: %v", err)
				}
				if string(content) != "raw-a,raw-b\n" {
					t.Errorf("raw upload content = %q, expected raw-a,raw-b\\n", string(content))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("PUT path = %q", r.URL.Path)
			}
			assertJSONBody(t, r, "from-put")
			writeJSON(t, w, testItem{ID: "cert-1", Name: "from-put"})
		case http.MethodDelete:
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("DELETE path = %q", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	got, err := Get[testItem](context.Background(), client, "/system/certificates/cert-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "from-get" {
		t.Fatalf("Get name = %q, expected from-get", got.Name)
	}

	created, err := Post[testItem, testItem](context.Background(), client, "/system/certificates", testItem{Name: "from-post"})
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if created.ID != "cert-2" {
		t.Fatalf("Post ID = %q, expected cert-2", created.ID)
	}

	if err := PostNoResponse(context.Background(), client, "/system/no-response", testItem{Name: "from-post-no-response"}); err != nil {
		t.Fatalf("PostNoResponse returned error: %v", err)
	}

	patched, err := Patch[testItem, testItem](context.Background(), client, "/system/certificates/cert-1", testItem{Name: "from-patch"})
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	if patched.Name != "from-patch" {
		t.Fatalf("Patch name = %q, expected from-patch", patched.Name)
	}

	if err := PatchNoResponse(context.Background(), client, "/system/no-response", testItem{Name: "from-patch-no-response"}); err != nil {
		t.Fatalf("PatchNoResponse returned error: %v", err)
	}

	put, err := Put[testItem, testItem](context.Background(), client, "/system/certificates/cert-1", testItem{Name: "from-put"})
	if err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	if put.Name != "from-put" {
		t.Fatalf("Put name = %q, expected from-put", put.Name)
	}

	if err := Delete(context.Background(), client, "/system/certificates/cert-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if err := Upload(context.Background(), client, "/system/files", "lookup.csv", []byte("a,b\n")); err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}

	if err := PutRawNoResponse(context.Background(), client, "/system/files", "text/csv", []byte("raw-a,raw-b\n")); err != nil {
		t.Fatalf("PutRawNoResponse returned error: %v", err)
	}
}

func TestDecodeEnvelopeSingleAndSlice(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"count": 2,
			"items": []testItem{
				{ID: "one", Name: "first"},
				{ID: "two", Name: "second"},
			},
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	single, err := Get[testItem](context.Background(), client, "/system/certificates/one")
	if err != nil {
		t.Fatalf("Get single returned error: %v", err)
	}
	if single.ID != "one" {
		t.Fatalf("single ID = %q, expected one", single.ID)
	}

	list, err := Get[[]testItem](context.Background(), client, "/system/certificates")
	if err != nil {
		t.Fatalf("Get slice returned error: %v", err)
	}
	if len(*list) != 2 {
		t.Fatalf("slice length = %d, expected 2", len(*list))
	}
	if (*list)[1].ID != "two" {
		t.Fatalf("second ID = %q, expected two", (*list)[1].ID)
	}
}

func TestDecodeEnvelopeEmptySingleIsNotFound(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"count": 0,
			"items": []testItem{},
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	_, err := Get[testItem](context.Background(), client, "/system/certificates/missing")
	if err == nil {
		t.Fatal("Get empty envelope returned nil error")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound = false, expected true for empty single-resource envelope: %v", err)
	}

	list, err := Get[[]testItem](context.Background(), client, "/system/certificates")
	if err != nil {
		t.Fatalf("Get empty list returned error: %v", err)
	}
	if len(*list) != 0 {
		t.Fatalf("empty list length = %d, expected 0", len(*list))
	}
}

func TestDecodePlainJSONAndNoContent(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/empty") {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeJSON(t, w, testItem{ID: "plain", Name: "plain-json"})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	plain, err := Get[testItem](context.Background(), client, "/system/plain")
	if err != nil {
		t.Fatalf("Get plain returned error: %v", err)
	}
	if plain.ID != "plain" {
		t.Fatalf("plain ID = %q, expected plain", plain.ID)
	}

	empty, err := Patch[testItem, testItem](context.Background(), client, "/system/plain/empty", testItem{Name: "empty"})
	if err != nil {
		t.Fatalf("Patch empty returned error: %v", err)
	}
	if empty != nil {
		t.Fatalf("Patch empty = %#v, expected nil", empty)
	}
}

func TestErrors(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/missing":
			http.Error(w, "missing resource", http.StatusNotFound)
		case "/api/v1/bad":
			http.Error(w, "bad request", http.StatusBadRequest)
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	_, err := Get[testItem](context.Background(), client, "/missing")
	if err == nil {
		t.Fatal("Get missing returned nil error")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound = false, expected true for 404")
	}

	_, err = Get[testItem](context.Background(), client, "/bad")
	if err == nil {
		t.Fatal("Get bad returned nil error")
	}
	if IsNotFound(err) {
		t.Fatalf("IsNotFound = true, expected false for non-404")
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error type = %T, expected HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, expected 400", httpErr.StatusCode)
	}
	if !strings.Contains(httpErr.Body, "bad request") {
		t.Fatalf("body = %q, expected bad request", httpErr.Body)
	}

	if IsNotFound(fmt.Errorf("ordinary error")) {
		t.Fatalf("IsNotFound = true, expected false for ordinary error")
	}
}

func TestGatewayRouting(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	tests := []struct {
		name string
		path string
	}{
		{
			name: "v1 path",
			path: "/v1/organizations/org-id/workspaces",
		},
		{
			name: "api v1 path",
			path: "/api/v1/organizations/org-id/workspaces",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "gateway.cribl.cloud" {
			t.Errorf("Host = %q, expected gateway.cribl.cloud", r.Host)
		}
		if r.URL.Path != "/v1/organizations/org-id/workspaces" {
			t.Errorf("path = %q, expected /v1/organizations/org-id/workspaces", r.URL.Path)
		}
		writeJSON(t, w, testItem{ID: "workspace", Name: "gateway"})
	}))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
		HTTPClient: &http.Client{
			Transport: rewriteTransport{target: targetURL},
		},
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Get[testItem](context.Background(), client, test.path)
			if err != nil {
				t.Fatalf("Get gateway returned error: %v", err)
			}
			if got.Name != "gateway" {
				t.Fatalf("gateway name = %q, expected gateway", got.Name)
			}
		})
	}
}

func TestProviderWorkspaceRouting(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	t.Setenv("CRIBL_ORGANIZATION_ID", "")
	t.Setenv("CRIBL_WORKSPACE_ID", "")
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "provider-workspace-provider-org.cribl-playground.cloud" {
			t.Errorf("Host = %q, expected provider-workspace-provider-org.cribl-playground.cloud", r.Host)
		}
		if r.URL.Path != "/api/v1/system/certificates/cert-1" {
			t.Errorf("path = %q, expected /api/v1/system/certificates/cert-1", r.URL.Path)
		}
		writeJSON(t, w, testItem{ID: "cert-1", Name: "provider-config"})
	}))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}

	client := New(Config{
		ProviderOrgID:       "provider-org",
		ProviderWorkspaceID: "provider-workspace",
		ProviderCloudDomain: "cribl-playground.cloud",
		BearerToken:         "test-token",
		HTTPClient: &http.Client{
			Transport: rewriteTransport{target: targetURL},
		},
	})

	got, err := Get[testItem](context.Background(), client, "/system/certificates/cert-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "provider-config" {
		t.Fatalf("name = %q, expected provider-config", got.Name)
	}
}

func TestAuthenticationRequired(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	client := New(Config{BaseURL: "http://127.0.0.1:1"})
	_, err := Get[testItem](context.Background(), client, "/system/certificates")
	if err == nil {
		t.Fatal("Get returned nil error, expected authentication error")
	}
	if !strings.Contains(err.Error(), "authentication requires bearer token or credentials") {
		t.Fatalf("error = %q, expected authentication message", err.Error())
	}
}

type rewriteTransport struct {
	target *url.URL
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	out := req.Clone(req.Context())
	out.URL.Scheme = t.target.Scheme
	out.URL.Host = t.target.Host
	out.Host = req.URL.Host
	return http.DefaultTransport.RoundTrip(out)
}

func assertJSONBody(t *testing.T, r *http.Request, expectedName string) {
	t.Helper()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, expected application/json", r.Header.Get("Content-Type"))
	}
	var body testItem
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Errorf("failed to decode JSON body: %v", err)
		return
	}
	if body.Name != expectedName {
		t.Errorf("body name = %q, expected %q", body.Name, expectedName)
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("failed to write JSON response: %v", err)
	}
}
