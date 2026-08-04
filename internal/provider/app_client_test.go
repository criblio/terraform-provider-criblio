package provider

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCreateAppArchive(t *testing.T) {
	content, err := createAppArchive(AppModel{
		ID:          types.StringValue("created-app"),
		DisplayName: types.StringValue("Created App"),
		Version:     types.StringValue("1.2.3"),
	})
	if err != nil {
		t.Fatalf("create App archive: %v", err)
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		t.Fatalf("open gzip archive: %v", err)
	}
	defer gzipReader.Close()

	files := map[string][]byte{}
	staticDirectory := false
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read tar archive: %v", err)
		}
		if header.Name == "static/" && header.Typeflag == tar.TypeDir {
			staticDirectory = true
			continue
		}
		body, err := io.ReadAll(tarReader)
		if err != nil {
			t.Fatalf("read %s: %v", header.Name, err)
		}
		files[header.Name] = body
	}
	if len(files["static/index.html"]) == 0 {
		t.Fatal("generated App is missing static/index.html")
	}
	if !staticDirectory {
		t.Fatal("generated App is missing the static directory entry")
	}
	var manifest map[string]any
	if err := json.Unmarshal(files["package.json"], &manifest); err != nil {
		t.Fatalf("decode package.json: %v", err)
	}
	if manifest["name"] != "created-app" || manifest["version"] != "1.2.3" {
		t.Fatalf("package.json = %#v", manifest)
	}
}

func TestAppUpload(t *testing.T) {
	const content = "app archive"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/apps" {
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("filename"); got != "example.tgz" {
			t.Fatalf("filename = %q, want example.tgz", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/octet-stream" {
			t.Fatalf("Content-Type = %q, want application/octet-stream", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request: %v", err)
		}
		if string(body) != content {
			t.Fatalf("body = %q, want %q", body, content)
		}
		if err := json.NewEncoder(w).Encode(map[string]string{"source": "example.random.tgz"}); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	filename := filepath.Join(t.TempDir(), "example.tgz")
	if err := os.WriteFile(filename, []byte(content), 0o600); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	api := newAppAPI(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"}))
	source, err := api.upload(context.Background(), filename)
	if err != nil {
		t.Fatalf("upload App: %v", err)
	}
	if source != "example.random.tgz" {
		t.Fatalf("source = %q, want example.random.tgz", source)
	}
}
