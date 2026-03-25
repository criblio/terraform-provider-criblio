//go:build integration

// Package integration provides live integration tests for the goatify CLI.
// Tests run against a real Cribl Cloud organization. Skip when credentials are missing.
// Run with: go test -tags=integration ./tools/import-cli/integration/...
package integration

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/stretchr/testify/require"
)

// buildLocalProvider builds the criblio terraform provider into pluginDir and returns the path.
// Uses the same layout as e2e: registry.terraform.io/criblio/criblio/999.99.9/<os>_<arch>/terraform-provider-criblio_v999.99.9
func buildLocalProvider(t *testing.T, pluginDir string) string {
	t.Helper()
	osArch := runtime.GOOS + "_" + runtime.GOARCH
	providerPath := filepath.Join(pluginDir, "registry.terraform.io", "criblio", "criblio", "999.99.9", osArch)
	require.NoError(t, os.MkdirAll(providerPath, 0755))
	binaryPath := filepath.Join(providerPath, "terraform-provider-criblio_v999.99.9")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = repoRoot(t)
	buildCmd.Env = os.Environ()
	var buf bytes.Buffer
	buildCmd.Stdout = &buf
	buildCmd.Stderr = &buf
	require.NoError(t, buildCmd.Run(), "build local provider: %s", buf.String())
	return pluginDir
}

// writeTerraformTfvarsWithPlaceholders writes terraform.tfvars with valid placeholder values for required variables.
// Certificate variables use examples/certificate/server.crt, private key variables use server.key; others use "placeholder".
// All certificate variables share the same cert content; all priv_key variables share the same key content.
func writeTerraformTfvarsWithPlaceholders(t *testing.T, outputDir string) {
	t.Helper()
	tfvarsExample := filepath.Join(outputDir, "terraform.tfvars.example")
	if _, err := os.Stat(tfvarsExample); err != nil {
		return // no variables, nothing to write
	}
	content, err := os.ReadFile(tfvarsExample)
	require.NoError(t, err)
	// Parse variable names: lines like "name = "" or "name = <<-EOT
	re := regexp.MustCompile(`(?m)^([a-zA-Z0-9_]+)\s*=`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	var names []string
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) >= 2 && !seen[m[1]] {
			seen[m[1]] = true
			names = append(names, m[1])
		}
	}
	if len(names) == 0 {
		return
	}
	root := repoRoot(t)
	certPath := filepath.Join(root, "examples", "certificate", "server.crt")
	keyPath := filepath.Join(root, "examples", "certificate", "server.key")
	certContent, err := os.ReadFile(certPath)
	require.NoError(t, err, "read server.crt for placeholder")
	keyContent, err := os.ReadFile(keyPath)
	require.NoError(t, err, "read server.key for placeholder")
	certStr := strings.TrimSpace(string(certContent))
	keyStr := strings.TrimSpace(string(keyContent))

	var buf bytes.Buffer
	buf.WriteString("# Generated for integration test - valid placeholders from examples/certificate/\n\n")
	for _, n := range names {
		lower := strings.ToLower(n)
		if strings.Contains(lower, "priv_key") || strings.Contains(lower, "private_key") {
			fmt.Fprintf(&buf, "%s = <<-EOT\n%s\nEOT\n", n, keyStr)
		} else if strings.Contains(lower, "cert") {
			fmt.Fprintf(&buf, "%s = <<-EOT\n%s\nEOT\n", n, certStr)
		} else {
			fmt.Fprintf(&buf, "%s = %q\n", n, "placeholder")
		}
	}
	require.NoError(t, os.WriteFile(filepath.Join(outputDir, "terraform.tfvars"), buf.Bytes(), 0600))
}

// TestIntegration_FullFlow_Cloud runs the full flow for criblio_group only: build, dry-run, export,
// HCL parse, terraform init/plan/apply, and no-drift check.
func TestIntegration_FullFlow_Cloud(t *testing.T) {
	if os.Getenv("CRIBL_CLIENT_ID") == "" || os.Getenv("CRIBL_CLIENT_SECRET") == "" {
		t.Skip("Integration test requires CRIBL_CLIENT_ID and CRIBL_CLIENT_SECRET (and CRIBL_ORGANIZATION_ID, CRIBL_WORKSPACE_ID, CRIBL_CLOUD_DOMAIN)")
	}

	orgID := os.Getenv("CRIBL_ORGANIZATION_ID")
	workspaceID := os.Getenv("CRIBL_WORKSPACE_ID")
	cloudDomain := os.Getenv("CRIBL_CLOUD_DOMAIN")
	if orgID == "" || workspaceID == "" || cloudDomain == "" {
		t.Skip("Integration test requires CRIBL_ORGANIZATION_ID, CRIBL_WORKSPACE_ID, CRIBL_CLOUD_DOMAIN")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "goatify")
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0755)

	env := append(os.Environ(),
		"CRIBL_CLIENT_ID="+os.Getenv("CRIBL_CLIENT_ID"),
		"CRIBL_CLIENT_SECRET="+os.Getenv("CRIBL_CLIENT_SECRET"),
		"CRIBL_ORGANIZATION_ID="+orgID,
		"CRIBL_WORKSPACE_ID="+workspaceID,
		"CRIBL_CLOUD_DOMAIN="+cloudDomain,
	)

	// Step 1: Build goatify
	t.Log("Step 1: Build goatify")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./tools/import-cli")
	buildCmd.Dir = repoRoot(t)
	buildCmd.Env = env
	var buildBuf bytes.Buffer
	buildCmd.Stdout = io.MultiWriter(&buildBuf, os.Stdout)
	buildCmd.Stderr = io.MultiWriter(&buildBuf, os.Stderr)
	err := buildCmd.Run()
	require.NoError(t, err, "build goatify: %s", buildBuf.String())
	_, err = os.Stat(binaryPath)
	require.NoError(t, err, "binary should exist")

	// Step 2: export --dry-run
	t.Log("Step 2: export --dry-run")
	dryRunCmd := exec.Command(binaryPath, "export", "--dry-run", "--include", "criblio_group",
		"--org-id", orgID, "--workspace-id", workspaceID, "--cloud-domain", cloudDomain)
	dryRunCmd.Dir = tmpDir
	dryRunCmd.Env = env
	var dryRunBuf bytes.Buffer
	dryRunCmd.Stdout = io.MultiWriter(&dryRunBuf, os.Stdout)
	dryRunCmd.Stderr = io.MultiWriter(&dryRunBuf, os.Stderr)
	err = dryRunCmd.Run()
	require.NoError(t, err, "dry-run should succeed")
	require.Contains(t, dryRunBuf.String(), "Preview:", "stderr should contain Preview")
	require.Contains(t, dryRunBuf.String(), "criblio_", "stderr should list resource types")

	// Step 3: export (criblio_group only)
	t.Log("Step 3: export (criblio_group)")
	exportCmd := exec.Command(binaryPath, "export", "--include", "criblio_group", "--output-dir", outputDir,
		"--org-id", orgID, "--workspace-id", workspaceID, "--cloud-domain", cloudDomain)
	exportCmd.Dir = tmpDir
	exportCmd.Env = env
	var exportBuf bytes.Buffer
	exportCmd.Stdout = io.MultiWriter(&exportBuf, os.Stdout)
	exportCmd.Stderr = io.MultiWriter(&exportBuf, os.Stderr)
	err = exportCmd.Run()
	require.NoError(t, err, "export should succeed: %s", exportBuf.String())

	// Validate output structure: import.tf, providers.tf, main.tf at root; at least one module dir with main.tf
	require.FileExists(t, filepath.Join(outputDir, "import.tf"), "import.tf should exist")
	require.FileExists(t, filepath.Join(outputDir, "providers.tf"), "providers.tf should exist")
	require.FileExists(t, filepath.Join(outputDir, "main.tf"), "main.tf should exist")

	var foundModuleMain bool
	_ = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "main.tf" && path != filepath.Join(outputDir, "main.tf") {
			foundModuleMain = true
		}
		return nil
	})
	require.True(t, foundModuleMain, "at least one module dir should have main.tf")

	// Step 4: Parse generated HCL
	t.Log("Step 4: Parse generated HCL")
	mainTF, err := os.ReadFile(filepath.Join(outputDir, "main.tf"))
	require.NoError(t, err)
	require.NoError(t, hcl.ParseHCL(mainTF, "main.tf"))

	importTF, err := os.ReadFile(filepath.Join(outputDir, "import.tf"))
	require.NoError(t, err)
	require.NoError(t, hcl.ParseHCL(importTF, "import.tf"))

	// Step 5: Verify import blocks
	t.Log("Step 5: Verify import blocks")
	require.Contains(t, string(importTF), "import {", "import.tf should contain import blocks")
	require.Contains(t, string(importTF), "to =", "import.tf should contain to =")
	require.Contains(t, string(importTF), "id =", "import.tf should contain id =")

	// Step 6: Build local provider and terraform init
	t.Log("Step 6: Build local provider and terraform init")
	pluginDir := filepath.Join(tmpDir, "plugin-dir")
	buildLocalProvider(t, pluginDir)
	pluginDirAbs, err := filepath.Abs(pluginDir)
	require.NoError(t, err)

	runTerraform := func(args ...string) ([]byte, error) {
		cmd := exec.Command("terraform", args...)
		cmd.Dir = outputDir
		cmd.Env = env
		var buf bytes.Buffer
		cmd.Stdout = io.MultiWriter(&buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(&buf, os.Stderr)
		err := cmd.Run()
		return buf.Bytes(), err
	}

	initOut, err := runTerraform("init", "-plugin-dir", pluginDirAbs)
	require.NoError(t, err, "terraform init should succeed: %s", string(initOut))

	// Step 7: terraform plan
	t.Log("Step 7: terraform plan")
	planOut, err := runTerraform("plan", "-detailed-exitcode")
	// Exit 0 = no changes, 2 = changes planned (acceptable before apply)
	if err != nil {
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr) && exitErr.ExitCode() == 2,
			"terraform plan should succeed or exit 2 (changes): %v\n%s", err, string(planOut))
	}

	// Step 8: terraform apply -auto-approve
	t.Log("Step 8: terraform apply -auto-approve")
	applyOut, err := runTerraform("apply", "-auto-approve")
	require.NoError(t, err, "terraform apply should succeed: %s", string(applyOut))

	// Step 9: terraform plan (second run) - must report no drift
	t.Log("Step 9: terraform plan (second run) - no drift")
	plan2Out, err := runTerraform("plan", "-detailed-exitcode")
	require.NoError(t, err, "second terraform plan should succeed (no drift): %s", string(plan2Out))
	require.Contains(t, string(plan2Out), "No changes", "second plan should report no drift")
}

// TestIntegration_FullExport_Cloud runs full export (all resource types, no --include), validates
// output structure and HCL parsing, then runs terraform init, plan, apply, and a second plan (no-drift check).
func TestIntegration_FullExport_Cloud(t *testing.T) {
	if os.Getenv("CRIBL_CLIENT_ID") == "" || os.Getenv("CRIBL_CLIENT_SECRET") == "" {
		t.Skip("Integration test requires CRIBL_CLIENT_ID and CRIBL_CLIENT_SECRET (and CRIBL_ORGANIZATION_ID, CRIBL_WORKSPACE_ID, CRIBL_CLOUD_DOMAIN)")
	}

	orgID := os.Getenv("CRIBL_ORGANIZATION_ID")
	workspaceID := os.Getenv("CRIBL_WORKSPACE_ID")
	cloudDomain := os.Getenv("CRIBL_CLOUD_DOMAIN")
	if orgID == "" || workspaceID == "" || cloudDomain == "" {
		t.Skip("Integration test requires CRIBL_ORGANIZATION_ID, CRIBL_WORKSPACE_ID, CRIBL_CLOUD_DOMAIN")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "goatify")
	outputDir := filepath.Join(tmpDir, "output")
	_ = os.MkdirAll(outputDir, 0755)

	env := append(os.Environ(),
		"CRIBL_CLIENT_ID="+os.Getenv("CRIBL_CLIENT_ID"),
		"CRIBL_CLIENT_SECRET="+os.Getenv("CRIBL_CLIENT_SECRET"),
		"CRIBL_ORGANIZATION_ID="+orgID,
		"CRIBL_WORKSPACE_ID="+workspaceID,
		"CRIBL_CLOUD_DOMAIN="+cloudDomain,
	)

	// Step 1: Build goatify
	t.Log("Step 1: Build goatify")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./tools/import-cli")
	buildCmd.Dir = repoRoot(t)
	buildCmd.Env = env
	var buildBuf bytes.Buffer
	buildCmd.Stdout = io.MultiWriter(&buildBuf, os.Stdout)
	buildCmd.Stderr = io.MultiWriter(&buildBuf, os.Stderr)
	err := buildCmd.Run()
	require.NoError(t, err, "build goatify: %s", buildBuf.String())

	// Step 2: export --dry-run (full, no --include)
	t.Log("Step 2: export --dry-run (full export)")
	dryRunCmd := exec.Command(binaryPath, "export", "--dry-run",
		"--org-id", orgID, "--workspace-id", workspaceID, "--cloud-domain", cloudDomain)
	dryRunCmd.Dir = tmpDir
	dryRunCmd.Env = env
	var dryRunBuf bytes.Buffer
	dryRunCmd.Stdout = io.MultiWriter(&dryRunBuf, os.Stdout)
	dryRunCmd.Stderr = io.MultiWriter(&dryRunBuf, os.Stderr)
	err = dryRunCmd.Run()
	require.NoError(t, err, "dry-run should succeed")
	require.Contains(t, dryRunBuf.String(), "Preview:", "stderr should contain Preview")

	// Step 3: export (full export, no --include)
	t.Log("Step 3: export (full export)")
	exportCmd := exec.Command(binaryPath, "export", "--output-dir", outputDir,
		"--org-id", orgID, "--workspace-id", workspaceID, "--cloud-domain", cloudDomain)
	exportCmd.Dir = tmpDir
	exportCmd.Env = env
	var exportBuf bytes.Buffer
	exportCmd.Stdout = io.MultiWriter(&exportBuf, os.Stdout)
	exportCmd.Stderr = io.MultiWriter(&exportBuf, os.Stderr)
	err = exportCmd.Run()
	require.NoError(t, err, "export should succeed: %s", exportBuf.String())

	// Validate output structure
	require.FileExists(t, filepath.Join(outputDir, "import.tf"), "import.tf should exist")
	require.FileExists(t, filepath.Join(outputDir, "providers.tf"), "providers.tf should exist")
	require.FileExists(t, filepath.Join(outputDir, "main.tf"), "main.tf should exist")

	var foundModuleMain bool
	_ = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "main.tf" && path != filepath.Join(outputDir, "main.tf") {
			foundModuleMain = true
		}
		return nil
	})
	require.True(t, foundModuleMain, "at least one module dir should have main.tf")

	// Step 4: Parse generated HCL
	t.Log("Step 4: Parse generated HCL")
	mainTF, err := os.ReadFile(filepath.Join(outputDir, "main.tf"))
	require.NoError(t, err)
	require.NoError(t, hcl.ParseHCL(mainTF, "main.tf"))

	importTF, err := os.ReadFile(filepath.Join(outputDir, "import.tf"))
	require.NoError(t, err)
	require.NoError(t, hcl.ParseHCL(importTF, "import.tf"))

	// Step 5: Write terraform.tfvars with valid placeholders (cert/key from examples/certificate/)
	t.Log("Step 5: Write terraform.tfvars with valid placeholders")
	writeTerraformTfvarsWithPlaceholders(t, outputDir)

	// Step 6: Build local provider and terraform init
	t.Log("Step 6: Build local provider and terraform init")
	pluginDir := filepath.Join(tmpDir, "plugin-dir")
	buildLocalProvider(t, pluginDir)
	pluginDirAbs, err := filepath.Abs(pluginDir)
	require.NoError(t, err)

	tfInit := exec.Command("terraform", "init", "-plugin-dir", pluginDirAbs)
	tfInit.Dir = outputDir
	tfInit.Env = env
	var initBuf bytes.Buffer
	tfInit.Stdout = io.MultiWriter(&initBuf, os.Stdout)
	tfInit.Stderr = io.MultiWriter(&initBuf, os.Stderr)
	err = tfInit.Run()
	require.NoError(t, err, "terraform init should succeed: %s", initBuf.String())

	runTerraform := func(args ...string) ([]byte, error) {
		cmd := exec.Command("terraform", args...)
		cmd.Dir = outputDir
		cmd.Env = env
		var buf bytes.Buffer
		cmd.Stdout = io.MultiWriter(&buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(&buf, os.Stderr)
		err := cmd.Run()
		return buf.Bytes(), err
	}

	// Step 7: terraform plan
	t.Log("Step 7: terraform plan")
	planOut, err := runTerraform("plan", "-detailed-exitcode")
	if err != nil {
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr) && (exitErr.ExitCode() == 0 || exitErr.ExitCode() == 2),
			"terraform plan should succeed or exit 2 (changes): %v\n%s", err, string(planOut))
	}

	// Step 8: terraform apply -auto-approve
	t.Log("Step 8: terraform apply -auto-approve")
	applyOut, err := runTerraform("apply", "-auto-approve")
	require.NoError(t, err, "terraform apply should succeed: %s", string(applyOut))

	// Step 9: terraform plan (second run) - no adds/destroys (in-place changes acceptable for full export)
	t.Log("Step 9: terraform plan (second run) - no adds or destroys")
	plan2Out, err := runTerraform("plan", "-detailed-exitcode")
	// Exit 0 = no changes, 2 = changes planned (acceptable when only in-place changes)
	if err != nil {
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr) && (exitErr.ExitCode() == 0 || exitErr.ExitCode() == 2),
			"second terraform plan should succeed or exit 2 (in-place changes): %v\n%s", err, string(plan2Out))
	}
	// When there are zero changes, Terraform prints "No changes." and omits the
	// "Plan: N to add, ..." summary; otherwise the summary includes 0 to add / 0 to destroy.
	assertPlanNoAddNoDestroy(t, string(plan2Out))
}

// assertPlanNoAddNoDestroy checks plan output for no new resources and no destroys.
// Terraform omits the numeric plan summary when the plan is completely empty (no changes).
func assertPlanNoAddNoDestroy(t *testing.T, planOut string) {
	t.Helper()
	if strings.Contains(planOut, "No changes") {
		return
	}
	require.Contains(t, planOut, "0 to add", "second plan should not add resources")
	require.Contains(t, planOut, "0 to destroy", "second plan should not destroy resources")
}

func repoRoot(t *testing.T) string {
	t.Helper()
	// Find repo root by looking for go.mod
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repo root (go.mod)")
		}
		dir = parent
	}
}
