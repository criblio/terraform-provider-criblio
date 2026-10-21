package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	hclv2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

// preserveImportedSensitiveValues changes only the integration test's copy of
// exported HCL. Placeholder inputs must not overwrite imported credentials.
// Ignore the containing attribute/block for nested secrets: Terraform cannot
// address individual elements of every collection type (notably sets).
func preserveImportedSensitiveValues(outputDir string) error {
	filesByDir := make(map[string][]string)
	if err := filepath.WalkDir(outputDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == ".terraform" {
			return filepath.SkipDir
		}
		if !entry.IsDir() && filepath.Ext(path) == ".tf" {
			filesByDir[filepath.Dir(path)] = append(filesByDir[filepath.Dir(path)], path)
		}
		return nil
	}); err != nil {
		return err
	}
	for _, paths := range filesByDir {
		files := make(map[string]*hclwrite.File)
		sensitiveVariables := make(map[string]bool)
		for _, path := range paths {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			file, diags := hclwrite.ParseConfig(content, path, hclv2.InitialPos)
			if diags.HasErrors() {
				return fmt.Errorf("parse exported HCL: %s", diags.Error())
			}
			files[path] = file
			for _, block := range file.Body().Blocks() {
				if block.Type() != "variable" || len(block.Labels()) != 1 {
					continue
				}
				attribute := block.Body().GetAttribute("sensitive")
				if attribute == nil {
					continue
				}
				expr, diags := hclsyntax.ParseExpression(attribute.Expr().BuildTokens(nil).Bytes(), path, hclv2.InitialPos)
				if diags.HasErrors() {
					return fmt.Errorf("parse variable sensitivity: %s", diags.Error())
				}
				value, diags := expr.Value(nil)
				if !diags.HasErrors() && value.RawEquals(cty.True) {
					sensitiveVariables[block.Labels()[0]] = true
				}
			}
		}
		for path, file := range files {
			changed := false
			for _, block := range file.Body().Blocks() {
				if block.Type() != "resource" {
					continue
				}
				ignored := make(map[string]bool)
				for name, attribute := range block.Body().Attributes() {
					if usesSensitiveVariable(attribute.Expr(), sensitiveVariables) {
						ignored[name] = true
					}
				}
				for _, nested := range block.Body().Blocks() {
					if nested.Type() != "lifecycle" && bodyUsesSensitiveVariable(nested.Body(), sensitiveVariables) {
						ignored[nested.Type()] = true
					}
				}
				if len(ignored) == 0 {
					continue
				}
				if err := appendIgnoredAttributes(block.Body(), ignored); err != nil {
					return fmt.Errorf("preserve sensitive values in %s: %w", path, err)
				}
				changed = true
			}
			if changed {
				if err := os.WriteFile(path, file.Bytes(), 0600); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func usesSensitiveVariable(expr *hclwrite.Expression, sensitiveVariables map[string]bool) bool {
	for _, variable := range expr.Variables() {
		traversal := variable.BuildTokens(nil).Bytes()
		parsed, diags := hclsyntax.ParseTraversalAbs(traversal, "", hclv2.InitialPos)
		if diags.HasErrors() || len(parsed) < 2 || parsed.RootName() != "var" {
			continue
		}
		if name, ok := parsed[1].(hclv2.TraverseAttr); ok && sensitiveVariables[name.Name] {
			return true
		}
	}
	return false
}

func bodyUsesSensitiveVariable(body *hclwrite.Body, sensitiveVariables map[string]bool) bool {
	for _, attribute := range body.Attributes() {
		if usesSensitiveVariable(attribute.Expr(), sensitiveVariables) {
			return true
		}
	}
	for _, block := range body.Blocks() {
		if bodyUsesSensitiveVariable(block.Body(), sensitiveVariables) {
			return true
		}
	}
	return false
}

func appendIgnoredAttributes(body *hclwrite.Body, ignored map[string]bool) error {
	var lifecycle *hclwrite.Block
	for _, block := range body.Blocks() {
		if block.Type() == "lifecycle" {
			lifecycle = block
			break
		}
	}
	if lifecycle == nil {
		lifecycle = body.AppendNewBlock("lifecycle", nil)
	}
	var elements []hclwrite.Tokens
	if attribute := lifecycle.Body().GetAttribute("ignore_changes"); attribute != nil {
		source := attribute.Expr().BuildTokens(nil).Bytes()
		if strings.TrimSpace(string(source)) == "all" {
			return nil
		}
		expr, diags := hclsyntax.ParseExpression(source, "", hclv2.InitialPos)
		if diags.HasErrors() {
			return fmt.Errorf("parse ignore_changes: %s", diags.Error())
		}
		tuple, ok := expr.(*hclsyntax.TupleConsExpr)
		if !ok {
			return fmt.Errorf("ignore_changes must be a list or all")
		}
		for _, element := range tuple.Exprs {
			traversal, diags := hclv2.AbsTraversalForExpr(element)
			if diags.HasErrors() {
				return fmt.Errorf("parse ignored attribute: %s", diags.Error())
			}
			elements = append(elements, hclwrite.TokensForTraversal(traversal))
			// A top-level ignore already covers its nested sensitive values.
			if len(traversal) == 1 {
				delete(ignored, traversal.RootName())
			}
		}
	}
	var names []string
	for name := range ignored {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		elements = append(elements, hclwrite.TokensForTraversal(hclv2.Traversal{hclv2.TraverseRoot{Name: name}}))
	}
	lifecycle.Body().SetAttributeRaw("ignore_changes", hclwrite.TokensForTuple(elements))
	return nil
}

func TestPreserveImportedSensitiveValues(t *testing.T) {
	dir := t.TempDir()
	moduleDir := filepath.Join(dir, "default", "database_connection")
	require.NoError(t, os.MkdirAll(moduleDir, 0700))
	require.NoError(t, os.WriteFile(filepath.Join(moduleDir, "variables.tf"), []byte(`
variable "connection_string" { sensitive = true }
variable "password" { sensitive = true }
variable "description" { sensitive = false }
`), 0600))
	path := filepath.Join(moduleDir, "main.tf")
	require.NoError(t, os.WriteFile(path, []byte(`
resource "criblio_database_connection" "test" {
  connection_string = var.connection_string
  description = var.description
  tls = { passphrase = var.password, disabled = false }
  lifecycle {
    prevent_destroy = true
    ignore_changes = [host]
  }
}
resource "criblio_source" "nested" {
  input_http {
    auth_tokens = [var.password]
  }
}
resource "criblio_database_connection" "all" {
  connection_string = var.connection_string
  lifecycle { ignore_changes = all }
}
resource "criblio_database_connection" "public" {
  description = var.description
}
`), 0600))
	require.NoError(t, preserveImportedSensitiveValues(dir))
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	file, diags := hclwrite.ParseConfig(content, path, hclv2.InitialPos)
	require.False(t, diags.HasErrors(), diags.Error())
	blocks := file.Body().Blocks()
	ignore := blocks[0].Body().Blocks()[0].Body().GetAttribute("ignore_changes").Expr().BuildTokens(nil).Bytes()
	require.Equal(t, "[host, connection_string, tls]", strings.TrimSpace(string(ignore)))
	require.NotNil(t, blocks[0].Body().Blocks()[0].Body().GetAttribute("prevent_destroy"))
	ignore = blocks[1].Body().Blocks()[1].Body().GetAttribute("ignore_changes").Expr().BuildTokens(nil).Bytes()
	require.Equal(t, "[input_http]", strings.TrimSpace(string(ignore)))
	ignore = blocks[2].Body().Blocks()[0].Body().GetAttribute("ignore_changes").Expr().BuildTokens(nil).Bytes()
	require.Equal(t, "all", strings.TrimSpace(string(ignore)))
	require.Empty(t, blocks[3].Body().Blocks())
	// Running preparation twice must not duplicate lifecycle entries.
	require.NoError(t, preserveImportedSensitiveValues(dir))
	second, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, string(content), string(second))
}

// Use Terraform's built-in resource to verify state retention without a live
// Cribl organization or downloading a provider.
func TestPreserveImportedSensitiveValuesTerraform(t *testing.T) {
	terraform, err := exec.LookPath("terraform")
	if err != nil {
		t.Skip("Terraform is not installed")
	}
	dir := t.TempDir()
	config := []byte(`
variable "connection_string" {
  type = string
  sensitive = true
}
resource "terraform_data" "connection" {
  input = var.connection_string
}
`)
	path := filepath.Join(dir, "main.tf")
	require.NoError(t, os.WriteFile(path, config, 0600))
	tfvars := filepath.Join(dir, "terraform.tfvars")
	require.NoError(t, os.WriteFile(tfvars, []byte(`connection_string = "postgresql://user:dummy@db.example.com:5432/appdb"`), 0600))
	run := func(args ...string) ([]byte, error) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, terraform, args...)
		cmd.Dir = dir
		return cmd.CombinedOutput()
	}
	output, err := run("init", "-backend=false", "-input=false", "-no-color")
	require.NoError(t, err, "%s", output)
	output, err = run("apply", "-auto-approve", "-input=false", "-no-color")
	require.NoError(t, err, "%s", output)
	require.NoError(t, os.WriteFile(tfvars, []byte(`connection_string = "placeholder"`), 0600))
	require.NoError(t, preserveImportedSensitiveValues(dir))
	output, err = run("plan", "-detailed-exitcode", "-input=false", "-no-color")
	require.NoError(t, err, "placeholder must not change existing state: %s", output)
	require.Contains(t, string(output), "No changes.")
	// Without test preparation, the same placeholder would overwrite the value.
	require.NoError(t, os.WriteFile(path, config, 0600))
	output, err = run("plan", "-detailed-exitcode", "-input=false", "-no-color")
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "%s", output)
	require.Equal(t, 2, exitErr.ExitCode())
}
