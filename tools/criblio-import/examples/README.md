# Example: Generate Groups Terraform Resources

This simple example script demonstrates:
1. **Automatic Authentication** - Using the SDK with automatic auth via `CriblTerraformHook`
2. **Exporting Groups** - Fetching all groups from Cribl API
3. **Generating Terraform** - Creating `criblio_group` resource blocks

## Prerequisites

1. Set up authentication (choose one method):

```bash
# Option 1: Bearer Token
export CRIBL_BEARER_TOKEN="your-token"
export CRIBL_WORKSPACE_ID="main"
export CRIBL_ORGANIZATION_ID="org"

# Option 2: OAuth
export CRIBL_CLIENT_ID="your-client-id"
export CRIBL_CLIENT_SECRET="your-client-secret"
export CRIBL_ORGANIZATION_ID="org"
export CRIBL_WORKSPACE_ID="main"

# Option 3: Credentials File (~/.cribl/credentials)
# See AUTHENTICATION_NOTES.md for format
```

## Usage

```bash
# From the repository root
cd tools/criblio-import/examples

# Run the script
go run generate-groups.go
```

## What It Does

1. **Authenticates** - Automatically via `CriblTerraformHook` (no code needed!)
2. **Fetches Groups** - Gets all groups for both "stream" and "edge" products
3. **Generates Terraform** - Creates `groups.tf` with `criblio_group` resources

## Output

The script generates a `groups.tf` file with Terraform resources like:

```hcl
resource "criblio_group" "my_group" {
  id      = "my-group"
  product = "stream"
  name    = "My Group"
  on_prem = true
  provisioned = true
  
  streamtags = [
    "prod",
    "edge",
  ]
  
  cloud {
    provider = "aws"
    region   = "us-east-1"
  }
}
```

## Next Steps

After generating the file:

```bash
# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply if satisfied
terraform apply
```

## Customization

You can modify the script to:
- Export only specific groups
- Add more fields from the API
- Generate different output formats
- Add resource dependencies

## Notes

- Authentication is handled automatically - no manual token management!
- The SDK handles URL construction, retries, and error handling
- This pattern can be extended to other resource types (sources, destinations, etc.)

