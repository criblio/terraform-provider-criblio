package main

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/codegen/parser"
)

func TestWorkerGroupUpgradeGeneration(t *testing.T) {
	// Build both merge outputs in a temporary directory so a fresh checkout
	// needs no generated spec and the test never rewrites provider source files.
	dir := t.TempDir()
	specPath := filepath.Join(dir, "merged-spec.yml")
	command := exec.CommandContext(t.Context(), "go", "run", "./tools/merge-spec",
		"--output", specPath,
		"--cloud-only-paths-output", filepath.Join(dir, "cloud_only_paths.go"))
	command.Dir = filepath.Join("..", "..")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("merge test spec: %v\n%s", err, output)
	}
	definitions, err := parser.ParseFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	var definition *parser.ResourceDef
	for i := range definitions {
		if definitions[i].StructName == "WorkerGroupUpgrade" {
			definition = &definitions[i]
			break
		}
	}
	if definition == nil {
		t.Fatal("WorkerGroupUpgrade annotation missing from merged spec")
	}
	if !definition.Action || definition.Create.Path != "/system/distributed/upgrade/{group}" || definition.Create.CreateHook != "createWorkerGroupUpgrade" {
		t.Fatalf("incorrect action operation: %+v", definition.Create)
	}
	if len(definition.Fields) != 8 {
		t.Fatalf("unexpected upgrade fields: %+v", definition.Fields)
	}
	for _, field := range definition.Fields {
		if field.TerraformName == "version_to" {
			t.Fatal("operator-only versionTo must not be exposed")
		}
		if field.TerraformName != "job_ids" && (!field.ForceNew || !field.StrictForceNew) {
			t.Fatalf("input changes and removals must replace the action: %+v", field)
		}
		if (field.TerraformName == "trigger" || field.TerraformName == "job_ids" || field.TerraformName == "group_id") && field.RequestField {
			t.Fatalf("local/path field must not be serialized: %+v", field)
		}
	}
	client := renderTemplate(t, "client", *definition)
	assertContains(t, client, "return createWorkerGroupUpgrade(ctx, a.client, model)")
	assertNotContains(t, client, `"fmt"`)
	assertNotContains(t, client, "restclient.Post[")
	assertContains(t, renderTemplate(t, "types", *definition), `output["upgradeMode"] = value`)
	resource := renderTemplate(t, "resource", *definition)
	assertContains(t, resource, `stringdefault.StaticString("rolling")`)
	assertContains(t, resource, `stringvalidator.OneOf("rolling", "regular", "batch")`)
	assertContains(t, resource, "No import state operation is available for action resource")
	for _, file := range outputFiles(*definition) {
		if file.Kind == "data_source" || file.Kind == "test" || file.Kind == "data_source_test" {
			t.Fatalf("action must not generate discovery or live acceptance tests: %+v", file)
		}
	}
	doc := renderTemplate(t, "doc", *definition)
	assertContains(t, doc, "successful apply confirms submission")
	assertContains(t, doc, "Cribl-managed cloud groups and Outpost Groups are unsupported")
	assertContains(t, doc, `trigger            = "upgrade-1"`)
	assertContains(t, doc, `group_id           = var.worker_group_id`)
	assertNotContains(t, doc, `group_id           = "default"`)
}
