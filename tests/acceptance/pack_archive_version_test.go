package tests

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
)

func TestPackArchiveVersion(t *testing.T) {
	const address = "criblio_pack.file"
	id := acctest.RandomWithPrefix("tf-archive")
	sourceID := id + "-input"
	dir := t.TempDir()
	config := func(version, metadata string) string {
		return fmt.Sprintf(`
resource "criblio_pack" "file" {
  id       = %q
  group_id = "default"
  filename = %q
  %s
}
resource "criblio_source" "attachment" {
  group_id = "default"
  id = %q
  input_datagen = {
    id = %q
    type = "datagen"
    disabled = true
    pipeline = "pack:${criblio_pack.file.id}"
    samples = [{ sample = "apache_common.log", events_per_sec = 1 }]
  }
}
`, id, filepath.Join(dir, id+"-"+version+".crbl"), metadata, sourceID, sourceID)
	}
	writePackArchive(t, dir, id, "1.0.1")
	writePackArchive(t, dir, id, "1.0.2")
	var client *restclient.Client
	checkState := func(version string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(address, "version", version),
			resource.TestCheckResourceAttr(address, "items.0.version", version),
		)
	}
	checkRemote := func(version string) {
		t.Helper()
		packPath := "/m/default/packs/" + id
		pack, err := restclient.Get[struct{ Version string }](t.Context(), client, packPath)
		require.NoError(t, err)
		require.Equal(t, version, pack.Version)
		archive, err := restclient.GetRaw(t.Context(), client, packPath+"/export?mode=merge")
		require.NoError(t, err)
		z, err := gzip.NewReader(bytes.NewReader(archive))
		require.NoError(t, err)
		defer z.Close()
		files := tar.NewReader(z)
		for {
			entry, err := files.Next()
			require.NoError(t, err, "exported pack must contain a root package.json")
			if path.Clean(entry.Name) == "package.json" {
				var manifest struct{ Version, Author string }
				require.NoError(t, json.NewDecoder(files).Decode(&manifest))
				require.Equal(t, version, manifest.Version)
				require.Equal(t, "Terraform", manifest.Author)
				break
			}
		}
		source, err := restclient.Get[struct{ Pipeline string }](t.Context(), client, "/m/default/system/inputs/"+sourceID)
		require.NoError(t, err)
		require.Equal(t, "pack:"+id, source.Pipeline)
	}
	v2 := config("1.0.2", "")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		PreCheck: func() {
			credentials, err := auth.GetCredentials()
			require.NoError(t, err)
			client = restclient.New(restclient.Config{BaseURL: auth.ConstructBaseURL(auth.ConstructBaseURLInput{}, credentials), Credentials: credentials})
		},
		Steps: []resource.TestStep{
			{Config: config("1.0.1", ""), Check: checkState("1.0.1")},
			{PreConfig: func() { checkRemote("1.0.1") }, Config: v2, Check: checkState("1.0.2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate)}}},
			{PreConfig: func() { checkRemote("1.0.2") }, RefreshState: true, Check: checkState("1.0.2")},
			{Config: v2, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}}},
			{Config: config("1.0.2", `tags = { streamtags = ["inventory"] }`), Check: checkState("1.0.2")},
			{PreConfig: func() { checkRemote("1.0.2") }, Config: config("1.0.2", `version = "1.0.3"`), Check: checkState("1.0.3")},
			{PreConfig: func() { checkRemote("1.0.3") }, Config: config("1.0.2", `version = "1.0.3"`), PlanOnly: true},
		},
	})
}

func writePackArchive(t *testing.T, dir, id, version string) {
	t.Helper()
	file, err := os.Create(filepath.Join(dir, id+"-"+version+".crbl"))
	require.NoError(t, err)
	z := gzip.NewWriter(file)
	files := tar.NewWriter(z)
	for _, dir := range []string{"default/", "default/pipelines/", "default/pipelines/main/"} {
		require.NoError(t, files.WriteHeader(&tar.Header{Name: dir, Mode: 0755, Typeflag: tar.TypeDir}))
	}
	for name, content := range map[string]string{
		"package.json":                    fmt.Sprintf(`{"name":%q,"version":%q,"displayName":"Version test","author":"Terraform"}`, id, version),
		"default/pipelines/main/conf.yml": "functions: []\n",
		"default/pipelines/route.yml":     "routes: []\n",
	} {
		require.NoError(t, files.WriteHeader(&tar.Header{Name: name, Mode: 0644, Size: int64(len(content))}))
		_, err := io.WriteString(files, content)
		require.NoError(t, err)
	}
	require.NoError(t, files.Close())
	require.NoError(t, z.Close())
	require.NoError(t, file.Close())
}
