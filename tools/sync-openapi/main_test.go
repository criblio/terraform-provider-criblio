package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const validSpec = `openapi: 3.0.2
paths: {}
`

func TestParseVersionsMap(t *testing.T) {
	html := `<script>var tagMappings,apiTitle,versions={"v4.18.0":"https://example.test/4.18.yml","v4.17.1":"https://example.test/4.17.yml"};window.location.pathname</script>`

	versions, err := parseVersionsMap(html)
	if err != nil {
		t.Fatalf("parseVersionsMap() error = %v", err)
	}
	if got := versions["v4.18.0"]; got != "https://example.test/4.18.yml" {
		t.Fatalf("versions[v4.18.0] = %q", got)
	}
}

func TestSelectVersionLatest(t *testing.T) {
	versions := map[string]string{
		"v4.17.1": "https://example.test/4.17.yml",
		"v4.18.0": "https://example.test/4.18.yml",
		"v4.9.3":  "https://example.test/4.9.yml",
	}

	key, url, err := selectVersion(versions, "latest")
	if err != nil {
		t.Fatalf("selectVersion() error = %v", err)
	}
	if key != "v4.18.0" {
		t.Fatalf("key = %q, want v4.18.0", key)
	}
	if url != "https://example.test/4.18.yml" {
		t.Fatalf("url = %q", url)
	}
}

func TestSelectVersionMinorUsesHighestPatch(t *testing.T) {
	versions := map[string]string{
		"v4.18.0": "https://example.test/4.18.0.yml",
		"v4.18.2": "https://example.test/4.18.2.yml",
		"v4.18":   "https://example.test/4.18.yml",
		"v4.17.1": "https://example.test/4.17.1.yml",
	}

	key, url, err := selectVersion(versions, "4.18")
	if err != nil {
		t.Fatalf("selectVersion() error = %v", err)
	}
	if key != "v4.18.2" {
		t.Fatalf("key = %q, want v4.18.2", key)
	}
	if url != "https://example.test/4.18.2.yml" {
		t.Fatalf("url = %q", url)
	}
}

func TestSyncOpenAPIFromDocs(t *testing.T) {
	specServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, validSpec)
	}))
	defer specServer.Close()

	docsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<script>var tagMappings,apiTitle,versions={"v4.18.0":%q};window.location.pathname</script>`, specServer.URL+"/openapi.yml")
	}))
	defer docsServer.Close()

	out := filepath.Join(t.TempDir(), "upstream-openapi.yml")
	resolved, size, err := syncOpenAPI(context.Background(), testConfig(config{
		source:  "docs",
		docsURL: docsServer.URL,
		version: "latest",
		output:  out,
	}))
	if err != nil {
		t.Fatalf("syncOpenAPI() error = %v", err)
	}
	if resolved != specServer.URL+"/openapi.yml" {
		t.Fatalf("resolved = %q", resolved)
	}
	if size != len(validSpec) {
		t.Fatalf("size = %d, want %d", size, len(validSpec))
	}
	assertFile(t, out, validSpec)
}

func TestSyncOpenAPIFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, validSpec)
	}))
	defer server.Close()

	out := filepath.Join(t.TempDir(), "upstream-openapi.yml")
	_, _, err := syncOpenAPI(context.Background(), testConfig(config{
		source: "url",
		url:    server.URL + "/openapi.yml",
		output: out,
	}))
	if err != nil {
		t.Fatalf("syncOpenAPI() error = %v", err)
	}
	assertFile(t, out, validSpec)
}

func TestSyncOpenAPIFromFile(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.yml")
	if err := os.WriteFile(source, []byte(validSpec), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	out := filepath.Join(dir, "upstream-openapi.yml")
	resolved, _, err := syncOpenAPI(context.Background(), testConfig(config{
		source: "file",
		file:   source,
		output: out,
	}))
	if err != nil {
		t.Fatalf("syncOpenAPI() error = %v", err)
	}
	if resolved != source {
		t.Fatalf("resolved = %q, want %q", resolved, source)
	}
	assertFile(t, out, validSpec)
}

func TestSyncOpenAPIRejectsMissingOpenAPIField(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.yml")
	if err := os.WriteFile(source, []byte("paths: {}\n"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, _, err := syncOpenAPI(context.Background(), testConfig(config{
		source: "file",
		file:   source,
		output: filepath.Join(dir, "upstream-openapi.yml"),
	}))
	if err == nil {
		t.Fatal("syncOpenAPI() error = nil")
	}
}

func TestResolveDocsSpecURLReturnsMissingVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<script>var tagMappings,apiTitle,versions={"v4.18.0":"https://example.test/openapi.yml"};window.location.pathname</script>`)
	}))
	defer server.Close()

	_, err := resolveDocsSpecURL(context.Background(), server.Client(), server.URL, "4.17")
	if err == nil {
		t.Fatal("resolveDocsSpecURL() error = nil")
	}
}

func testConfig(cfg config) config {
	if cfg.client == nil {
		cfg.client = http.DefaultClient
	}
	return cfg
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("file = %q, want %q", got, want)
	}
}
