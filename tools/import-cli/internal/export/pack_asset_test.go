package export

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttachPackAssetsDownloadsArchiveAndUsesFilename(t *testing.T) {
	archive := []byte("portable-crbl-content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/m/default/packs/test_pack/export", r.URL.Path)
		assert.Equal(t, "merge", r.URL.Query().Get("mode"))
		assert.Equal(t, "test_pack.crbl", r.URL.Query().Get("filename"))
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(archive)
	}))
	defer server.Close()

	client := &importclient.Client{REST: restclient.New(restclient.Config{
		BaseURL:     server.URL,
		BearerToken: "test",
		HTTPClient:  server.Client(),
	})}
	items := []generator.ResourceItem{{
		TypeName: "criblio_pack",
		Name:     "pack_default_test_pack",
		GroupID:  "default",
		Attrs: map[string]hcl.Value{
			"id":           {Kind: hcl.KindString, String: "test_pack"},
			"source":       {Kind: hcl.KindString, String: "file:/source-only/test_pack.crbl"},
			"description":  {Kind: hcl.KindString, String: "from source"},
			"display_name": {Kind: hcl.KindString, String: "Test Pack"},
		},
	}}

	count, err := AttachPackAssets(context.Background(), client, items)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.NotContains(t, items[0].Attrs, "source")
	assert.NotContains(t, items[0].Attrs, "description")
	assert.NotContains(t, items[0].Attrs, "display_name")
	require.Equal(t, hcl.KindExpression, items[0].Attrs["filename"].Kind)
	assert.Equal(t, `"${path.module}/files/pack_default_test_pack/test_pack.crbl"`, items[0].Attrs["filename"].Expr)
	require.Len(t, items[0].Files, 1)
	assert.Equal(t, "files/pack_default_test_pack/test_pack.crbl", items[0].Files[0].Path)
	assert.Equal(t, archive, items[0].Files[0].Content)
}

func TestPreparePackResourcesRemovesSourceArtifactFields(t *testing.T) {
	items := []generator.ResourceItem{{
		TypeName: "criblio_pack",
		Attrs: map[string]hcl.Value{
			"id":       {Kind: hcl.KindString, String: "test"},
			"source":   {Kind: hcl.KindString, String: "file:/source/test.crbl"},
			"filename": {Kind: hcl.KindString, String: "/tmp/test.crbl"},
		},
		Files: []generator.ResourceFile{{Path: "test.crbl", Content: []byte("test")}},
	}}

	assert.Equal(t, 1, PreparePackResources(items))
	assert.NotContains(t, items[0].Attrs, "source")
	assert.NotContains(t, items[0].Attrs, "filename")
	assert.Empty(t, items[0].Files)
}

func TestRemovePackChildResources(t *testing.T) {
	items := []generator.ResourceItem{
		{TypeName: "criblio_pack"},
		{TypeName: "criblio_pack_pipeline"},
		{TypeName: "criblio_pack_routes"},
		{TypeName: "criblio_source"},
	}
	filtered, removed := RemovePackChildResources(items)
	assert.Equal(t, 2, removed)
	assert.Equal(t, []generator.ResourceItem{{TypeName: "criblio_pack"}, {TypeName: "criblio_source"}}, filtered)
}

func TestAttachPackAssetsRejectsEmptyArchive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client := &importclient.Client{REST: restclient.New(restclient.Config{
		BaseURL:     server.URL,
		BearerToken: "test",
		HTTPClient:  server.Client(),
	})}
	items := []generator.ResourceItem{{
		TypeName: "criblio_pack",
		Name:     "pack_default_empty",
		GroupID:  "default",
		Attrs:    map[string]hcl.Value{"id": {Kind: hcl.KindString, String: "empty"}},
	}}

	_, err := AttachPackAssets(context.Background(), client, items)
	require.ErrorContains(t, err, "empty archive")
}
