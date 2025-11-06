# Auto-Import Feature Using HashiCorp Terraform-Exec

The CLI tool now supports automatic import of generated resources using HashiCorp's `terraform-exec` library.

## Features

✅ **Automatic Terraform Init** - Initializes Terraform if needed  
✅ **Automatic Import** - Imports all generated resources into state  
✅ **Validation** - Validates generated Terraform configuration  
✅ **Plan** - Shows what would change after import  
✅ **Optional Apply** - Can automatically apply changes  

## Usage

### Basic Auto-Import

```bash
# Generate and automatically import
go run generate-groups-with-import.go

# Or with environment variable
AUTO_IMPORT=true go run generate-groups.go
```

### With CLI Tool (when fully implemented)

```bash
# Generate Terraform files and auto-import
criblio-import --output ./configs --auto-import

# Generate, import, and automatically apply
criblio-import --output ./configs --auto-import --auto-apply
```

## How It Works

The `terraform-exec` library allows the Go program to:
1. Execute Terraform commands programmatically
2. Parse Terraform output
3. Handle errors gracefully
4. Chain operations (init → import → validate → plan → apply)

### Implementation Details

**Library Used:** `github.com/hashicorp/terraform-exec`

**Key Functions:**
- `tfexec.NewTerraform()` - Creates Terraform instance
- `tf.Init()` - Runs `terraform init`
- `tf.Import()` - Runs `terraform import`
- `tf.Validate()` - Runs `terraform validate`
- `tf.Plan()` - Runs `terraform plan`
- `tf.Apply()` - Runs `terraform apply`

### Example Flow

```go
// 1. Create importer
importer, err := terraform.NewImporter(".")
if err != nil {
    log.Fatal(err)
}

// 2. Initialize
importer.Init(ctx)

// 3. Import resources
importSpecs := []terraform.ImportSpec{
    {ResourceType: "criblio_group", ResourceName: "default", ImportID: "default"},
    // ... more resources
}
importer.ImportMultipleResources(ctx, importSpecs)

// 4. Validate
importer.Validate(ctx)

// 5. Plan
importer.Plan(ctx)

// 6. Apply (optional)
importer.Apply(ctx, autoApprove)
```

## Benefits

1. **No Manual Steps** - Everything happens automatically
2. **Error Handling** - Graceful handling of already-imported resources
3. **Progress Feedback** - Clear logging of each step
4. **CI/CD Ready** - Can be used in automation pipelines
5. **State Management** - Properly manages Terraform state

## Error Handling

The importer handles common scenarios:
- **Already Imported**: Skips resources already in state
- **Terraform Not Found**: Clear error message to install Terraform
- **Invalid Configuration**: Validation errors before import
- **Import Failures**: Continues with other resources

## Integration with Full CLI Tool

When the full CLI tool is implemented, the `--auto-import` flag will:

1. Generate Terraform files from Cribl config
2. Automatically initialize Terraform
3. Import all generated resources
4. Validate the configuration
5. Show a plan of changes
6. Optionally apply if `--auto-apply` is set

## Example Output

```
🔐 Initializing SDK client...
📥 Fetching groups from Cribl...
✅ Found 5 groups

📝 Generating Terraform resources...
✅ Generated Terraform file: groups.tf

🔄 Auto-importing resources into Terraform state...
📦 Initializing Terraform...
  ✓ Terraform initialized successfully

📥 Importing resources into Terraform state...
  Importing: criblio_group.default (default)
    ✓ Successfully imported criblio_group.default
  Importing: criblio_group.defaultHybrid (defaultHybrid)
    ✓ Successfully imported criblio_group.defaultHybrid
  ...

✅ Import complete!

🔍 Validating Terraform configuration...
  ✓ Configuration is valid

📊 Running terraform plan...
  ✓ No changes detected
```

## Next Steps

1. Test the example: `cd examples && go run generate-groups-with-import.go`
2. Review the generated state: `terraform show`
3. Apply if needed: `terraform apply`

