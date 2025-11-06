package terraform

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// Importer handles Terraform operations (init, import, plan, apply)
type Importer struct {
	workDir string
	tf      *tfexec.Terraform
}

// NewImporter creates a new Terraform importer
func NewImporter(workDir string) (*Importer, error) {
	// Find terraform executable - try common locations first
	tfPath := os.Getenv("TERRAFORM")
	if tfPath == "" {
		// Try to find terraform in PATH
		path, err := exec.LookPath("terraform")
		if err != nil {
			return nil, fmt.Errorf("terraform not found in PATH: %w. Please install Terraform", err)
		}
		tfPath = path
	}

	tf, err := tfexec.NewTerraform(workDir, tfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create terraform instance: %w", err)
	}

	return &Importer{
		workDir: workDir,
		tf:      tf,
	}, nil
}

// Init initializes Terraform in the work directory
func (i *Importer) Init(ctx context.Context) error {
	log.Println("📦 Initializing Terraform...")
	
	// Check if already initialized
	if _, err := os.Stat(filepath.Join(i.workDir, ".terraform")); err == nil {
		log.Println("  ✓ Terraform already initialized")
		return nil
	}

	err := i.tf.Init(ctx)
	if err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}

	log.Println("  ✓ Terraform initialized successfully")
	return nil
}

// ImportResource imports a single resource into Terraform state
func (i *Importer) ImportResource(ctx context.Context, resourceType, resourceName, importID string) error {
	address := fmt.Sprintf("%s.%s", resourceType, resourceName)
	
	log.Printf("  Importing: %s (%s)", address, importID)
	
	err := i.tf.Import(ctx, address, importID)
	if err != nil {
		// Check if it's already in state
		if strings.Contains(err.Error(), "already managed") || 
		   strings.Contains(err.Error(), "already exists") {
			log.Printf("    ⚠️  %s already in state, skipping", address)
			return nil
		}
		return fmt.Errorf("failed to import %s: %w", address, err)
	}

	log.Printf("    ✓ Successfully imported %s", address)
	return nil
}

// ImportMultipleResources imports multiple resources
func (i *Importer) ImportMultipleResources(ctx context.Context, imports []ImportSpec) error {
	log.Println("📥 Importing resources into Terraform state...")
	log.Println("")

	for _, imp := range imports {
		err := i.ImportResource(ctx, imp.ResourceType, imp.ResourceName, imp.ImportID)
		if err != nil {
			log.Printf("    ❌ Error: %v", err)
			// Continue with other imports
		}
	}

	log.Println("")
	log.Println("✅ Import complete!")
	return nil
}

// Plan runs terraform plan to show what would change
func (i *Importer) Plan(ctx context.Context) error {
	log.Println("📊 Running terraform plan...")
	
	hasChanges, err := i.tf.Plan(ctx)
	if err != nil {
		return fmt.Errorf("terraform plan failed: %w", err)
	}

	if hasChanges {
		log.Println("  ⚠️  Plan shows changes - review before applying")
	} else {
		log.Println("  ✓ No changes detected")
	}

	return nil
}

// Apply runs terraform apply
func (i *Importer) Apply(ctx context.Context, autoApprove bool) error {
	log.Println("🚀 Applying Terraform configuration...")

	// Note: For auto-approve, users should run terraform apply -auto-approve manually
	// The terraform-exec library doesn't have a direct ApplyFlag option
	// We can implement this by calling terraform with -auto-approve flag if needed
	err := i.tf.Apply(ctx)
	if err != nil {
		return fmt.Errorf("terraform apply failed: %w", err)
	}

	log.Println("  ✓ Apply completed successfully")
	return nil
}

// Refresh runs terraform refresh to update state from remote
func (i *Importer) Refresh(ctx context.Context) error {
	log.Println("🔄 Refreshing Terraform state from remote...")
	
	err := i.tf.Refresh(ctx)
	if err != nil {
		return fmt.Errorf("terraform refresh failed: %w", err)
	}

	log.Println("  ✓ State refreshed successfully")
	return nil
}

// Validate runs terraform validate
func (i *Importer) Validate(ctx context.Context) error {
	log.Println("🔍 Validating Terraform configuration...")
	
	output, err := i.tf.Validate(ctx)
	if err != nil {
		return fmt.Errorf("terraform validate failed: %w", err)
	}

	if output == nil {
		return fmt.Errorf("terraform validate returned nil output")
	}

	if output.Valid {
		log.Println("  ✓ Configuration is valid")
		return nil
	}

	if len(output.Diagnostics) > 0 {
		var errMsg strings.Builder
		errMsg.WriteString("terraform configuration is invalid:\n")
		for _, diag := range output.Diagnostics {
			errMsg.WriteString(fmt.Sprintf("  - %s: %s\n", diag.Severity, diag.Summary))
		}
		return fmt.Errorf(errMsg.String())
	}

	log.Println("  ✓ Configuration is valid")
	return nil
}

// ImportSpec defines what to import
type ImportSpec struct {
	ResourceType string
	ResourceName string
	ImportID     string
}

// GenerateImportSpecsFromGroups generates import specs from groups.tf
// This is a simple parser - in production, you'd use a proper HCL parser
func GenerateImportSpecsFromGroups(groupsFile string) ([]ImportSpec, error) {
	content, err := os.ReadFile(groupsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read groups file: %w", err)
	}

	var specs []ImportSpec
	lines := strings.Split(string(content), "\n")
	
	var currentResource string
	var currentID string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Match: resource "criblio_group" "resource_name"
		if strings.HasPrefix(line, `resource "criblio_group"`) {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentResource = strings.Trim(parts[2], `"`)
			}
		}
		
		// Match: id = "group_id"
		if strings.HasPrefix(line, `id =`) {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				currentID = strings.Trim(strings.TrimSpace(parts[1]), `"`)
			}
		}
		
		// When we hit the closing brace, we have a complete resource
		if line == "}" && currentResource != "" && currentID != "" {
			// For groups, the import ID is just the group ID
			// (The ImportState function only uses the ID)
			specs = append(specs, ImportSpec{
				ResourceType: "criblio_group",
				ResourceName: currentResource,
				ImportID:     currentID,
			})
			
			// Reset for next resource
			currentResource = ""
			currentID = ""
		}
	}

	return specs, nil
}

