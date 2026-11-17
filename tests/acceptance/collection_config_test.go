package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// withoutCollectionDataSources retains resource and by-ID checks after the
// lifecycle's first collection assertion. Generated collection blocks use "all".
func withoutCollectionDataSources(t *testing.T, config string) string {
	t.Helper()
	file, diagnostics := hclsyntax.ParseConfig([]byte(config), "acceptance.tf", hcl.InitialPos)
	if diagnostics.HasErrors() {
		t.Fatalf("parse acceptance config: %s", diagnostics)
	}
	blocks := file.Body.(*hclsyntax.Body).Blocks
	// Remove from the end so earlier source ranges remain valid. HCL ranges
	// handle nested blocks, comments, and braces inside strings and heredocs.
	for i := len(blocks) - 1; i >= 0; i-- {
		block := blocks[i]
		if block.Type == "data" && len(block.Labels) == 2 && block.Labels[1] == "all" {
			rangeToRemove := block.Range()
			config = config[:rangeToRemove.Start.Byte] + config[rangeToRemove.End.Byte:]
		}
	}
	return config
}

func TestCollectionRemovalKeepsPlanEmpty(t *testing.T) {
	terraformPath, err := exec.LookPath("terraform")
	if err != nil {
		t.Skip("Terraform CLI is required for the local collection removal regression")
	}
	statePath := filepath.Join(t.TempDir(), "remote.tfstate")
	if err := os.WriteFile(statePath, []byte(`{"version":4,"terraform_version":"1.5.0","serial":1,"lineage":"collection-test","outputs":{"id":{"value":"test","type":"string"}},"resources":[]}`), 0600); err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(`data "terraform_remote_state" "all" {
  backend = "local"
  config = { path = %q }
}
output "constant" { value = "unchanged" }
`, statePath)
	// The built-in data source exercises Terraform's state removal without
	// external providers, downloads, credentials, or Cribl API calls.
	workDir := t.TempDir()
	configPath := filepath.Join(workDir, "main.tf")
	writeConfig := func(value string) {
		t.Helper()
		if err := os.WriteFile(configPath, []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	run := func(arguments ...string) {
		t.Helper()
		command := exec.CommandContext(t.Context(), terraformPath, arguments...)
		command.Dir = workDir
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("terraform %v: %v\n%s", arguments, err, output)
		}
	}
	writeConfig(config)
	run("init", "-input=false", "-no-color")
	run("apply", "-input=false", "-auto-approve", "-no-color")
	writeConfig(withoutCollectionDataSources(t, config))
	// Exit code 2 would indicate a plan change after removing the data block.
	run("plan", "-input=false", "-no-color", "-detailed-exitcode")
}

func TestWithoutCollectionDataSources(t *testing.T) {
	config := `resource "criblio_pipeline" "test" {
  id = "test"
  conf = { description = "braces } stay {", functions = [] }
}
data "criblio_pipeline" "by_id" { id = criblio_pipeline.test.id }
data "criblio_pipelines" "all" {
  group_id = "default"
  depends_on = [criblio_pipeline.test]
}
data "criblio_packs" "all" { group_id = "default" }
output "pipeline_id" { value = data.criblio_pipeline.by_id.id }
`
	result := withoutCollectionDataSources(t, config)
	file, diagnostics := hclsyntax.ParseConfig([]byte(result), "result.tf", hcl.InitialPos)
	if diagnostics.HasErrors() {
		t.Fatal(diagnostics)
	}
	blocks := file.Body.(*hclsyntax.Body).Blocks
	if len(blocks) != 3 || blocks[0].Type != "resource" || blocks[1].Labels[1] != "by_id" || blocks[2].Type != "output" {
		t.Fatalf("resource, by-ID data source, and output must remain: %s", result)
	}
	if again := withoutCollectionDataSources(t, result); again != result {
		t.Fatal("config without collection blocks must remain unchanged")
	}
}
