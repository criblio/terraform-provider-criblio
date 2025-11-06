package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/criblio/terraform-provider-criblio/tools/criblio-import/internal/terraform"
)

func main() {
	ctx := context.Background()

	// Check if groups.tf exists
	if _, err := os.Stat("groups.tf"); os.IsNotExist(err) {
		log.Fatal("❌ groups.tf not found. Please run generate-groups.go first")
	}

	fmt.Println("🚀 Importing groups using Terraform-Exec library")
	fmt.Println("================================================")
	fmt.Println("")

	// Create provider.tf if it doesn't exist
	if _, err := os.Stat("provider.tf"); os.IsNotExist(err) {
		fmt.Println("📝 Creating provider.tf...")
		providerConfig := `terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Authentication via environment variables or ~/.cribl/credentials
}
`
		if err := os.WriteFile("provider.tf", []byte(providerConfig), 0644); err != nil {
			log.Fatalf("❌ Failed to create provider.tf: %v", err)
		}
		fmt.Println("  ✓ Created provider.tf")
		fmt.Println("")
	}

	// Initialize Terraform importer
	importer, err := terraform.NewImporter(".")
	if err != nil {
		log.Fatalf("❌ Failed to create Terraform importer: %v", err)
	}

	// Initialize Terraform
	if err := importer.Init(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize Terraform: %v", err)
	}

	// Generate import specs from groups.tf
	fmt.Println("📖 Parsing groups.tf to generate import commands...")
	importSpecs, err := terraform.GenerateImportSpecsFromGroups("groups.tf")
	if err != nil {
		log.Fatalf("❌ Failed to parse groups.tf: %v", err)
	}

	if len(importSpecs) == 0 {
		log.Fatal("❌ No groups found in groups.tf")
	}

	fmt.Printf("  Found %d groups to import\n", len(importSpecs))
	fmt.Println("")

	// Import all resources
	if err := importer.ImportMultipleResources(ctx, importSpecs); err != nil {
		log.Printf("⚠️  Some imports may have failed: %v", err)
	}

	// Refresh state to populate all fields from API
	// This is critical because ImportState only sets the ID,
	// and we need all fields populated to avoid replacement warnings
	fmt.Println("")
	fmt.Println("🔄 Refreshing state to populate all fields from API...")
	if err := importer.Refresh(ctx); err != nil {
		log.Printf("⚠️  Refresh warning: %v", err)
		log.Println("  Note: You may need to run 'terraform refresh' manually")
	} else {
		fmt.Println("  ✓ State refreshed successfully")
	}
	fmt.Println("")

	// Validate configuration
	if err := importer.Validate(ctx); err != nil {
		log.Printf("⚠️  Validation warning: %v", err)
	} else {
		fmt.Println("")
	}

	// Show plan
	if err := importer.Plan(ctx); err != nil {
		log.Printf("⚠️  Plan warning: %v", err)
	}

	fmt.Println("")
	fmt.Println("✨ Import complete!")
	fmt.Println("")
	fmt.Println("💡 Next steps:")
	fmt.Println("   1. Review the plan output above")
	fmt.Println("   2. Run: terraform show (to see current state)")
	fmt.Println("   3. Run: terraform apply (to apply any changes)")
	fmt.Println("")
}

