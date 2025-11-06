# criblio-import

A CLI tool to import Cribl configuration from `/bulk/diag/download` and convert to Terraform modules.

## Installation

### From Source
```bash
cd tools/criblio-import
go build -o criblio-import ./cmd/criblio-import
sudo mv criblio-import /usr/local/bin/
```

### Using Go Install
```bash
go install github.com/criblio/terraform-provider-criblio/tools/criblio-import/cmd/criblio-import@latest
```

## Quick Start

1. **Set up authentication** (choose one method):

```bash
# Method 1: Environment Variables
export CRIBL_CLIENT_ID="your-client-id"
export CRIBL_CLIENT_SECRET="your-client-secret"
export CRIBL_ORGANIZATION_ID="your-org-id"
export CRIBL_WORKSPACE_ID="main"

# Method 2: Bearer Token
export CRIBL_BEARER_TOKEN="your-token"
export CRIBL_WORKSPACE_ID="main"
export CRIBL_ORGANIZATION_ID="your-org-id"

# Method 3: Credentials File (~/.cribl/credentials)
# See README.md in root for format
```

2. **Run the import**:

```bash
criblio-import --output ./terraform-configs
```

3. **Review and apply**:

```bash
cd terraform-configs
terraform init
terraform plan  # Review changes
terraform apply # Apply if satisfied
```

## Usage

```bash
criblio-import [flags]

Flags:
  -h, --help                   Show help message
  -o, --output string          Output directory for Terraform files (required)
  -f, --format string          Output format: resources or modules (default: modules)
  -i, --include strings        Only include these resource types (comma-separated)
  -e, --exclude strings        Exclude these resource types (comma-separated)
  -v, --var-file string        Generate variables.tf for sensitive fields
  -d, --dry-run                Preview what would be generated without writing files
      --verbose                Enable verbose logging
      --workspace-id string    Cribl workspace ID
      --organization-id string Cribl organization ID
      --bearer-token string    Bearer token for authentication
      --client-id string       OAuth client ID
      --client-secret string   OAuth client secret
      --onprem-server-url string On-prem server URL
      --onprem-username string On-prem username
      --onprem-password string On-prem password
```

## Examples

### Import All Resources
```bash
criblio-import --output ./configs
```

### Import Only Sources and Destinations
```bash
criblio-import --output ./configs --include sources,destinations
```

### Preview Before Importing
```bash
criblio-import --output ./preview --dry-run
```

### Generate with Variables for Sensitive Data
```bash
criblio-import --output ./configs --var-file variables.tf
```

### On-Prem Deployment
```bash
criblio-import \
  --output ./configs \
  --onprem-server-url http://localhost:9000 \
  --onprem-username admin \
  --onprem-password admin
```

## Supported Resource Types

- Sources (`criblio_source`)
- Destinations (`criblio_destination`)
- Pipelines (`criblio_pipeline`)
- Routes (`criblio_routes`)
- Packs (`criblio_pack`)
- Pack Components (breakers, lookups, vars, etc.)
- Global Variables (`criblio_global_var`)
- Certificates (`criblio_certificate`)
- Grok Patterns (`criblio_grok`)
- Regex Patterns (`criblio_regex`)
- And more...

## Output Structure

The tool generates a Terraform configuration structure:

```
output/
├── main.tf                 # Provider configuration
├── variables.tf            # Input variables (if --var-file specified)
├── modules/
│   ├── sources/
│   │   └── main.tf        # All source resources
│   ├── destinations/
│   │   └── main.tf        # All destination resources
│   ├── pipelines/
│   │   └── main.tf        # All pipeline resources
│   └── packs/
│       └── main.tf        # All pack resources
└── README.md              # Generated documentation
```

## Troubleshooting

### Authentication Issues
- Verify credentials using `curl` or the provider directly
- Check environment variables with `env | grep CRIBL`
- Review `~/.cribl/credentials` file format

### Import Errors
- Run with `--verbose` flag for detailed logs
- Check API endpoint accessibility
- Verify workspace/organization IDs

### Generated Terraform Issues
- Run `terraform fmt` to format files
- Use `terraform validate` to check syntax
- Review warnings in import logs

## Limitations

1. Some read-only or computed fields may not be imported
2. Resource dependencies may need manual ordering
3. Sensitive fields should be moved to variables
4. Complex nested configurations may require review

## Contributing

See the main repository CONTRIBUTING.md for guidelines.

## License

Same as the main repository.

