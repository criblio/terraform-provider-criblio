# HCL Syntax Notes for Terraform Resource Generation

## Important: SingleNestedAttribute vs Blocks

When generating Terraform HCL, there's a critical distinction between:

### SingleNestedAttribute (uses `=` syntax)
```hcl
resource "criblio_group" "example" {
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
}
```

### Blocks (use block syntax)
```hcl
resource "criblio_source" "example" {
  input_http {
    port = 8080
    host = "0.0.0.0"
  }
}
```

## How to Determine Which to Use

Check the Terraform schema:
- `schema.SingleNestedAttribute` → Use `attr = { ... }`
- `schema.Block` or nested blocks → Use `block { ... }`

## Common Patterns in Cribl Provider

### SingleNestedAttribute (use `=`)
- `cloud` in `criblio_group`
- `git` in `criblio_group`
- `tls` in sources/destinations (sometimes)
- `pq` in sources/destinations (sometimes)

### Blocks (no `=`)
- `input_*` in `criblio_source`
- `output_*` in `criblio_destination`
- Most nested configurations in sources/destinations

## Example: Correct Syntax

```hcl
# CORRECT - SingleNestedAttribute
resource "criblio_group" "example" {
  id      = "my-group"
  product = "stream"
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
}

# CORRECT - Block
resource "criblio_source" "example" {
  id       = "my-source"
  group_id = "default"
  
  input_http {
    port = 8080
    host = "0.0.0.0"
  }
}
```

## Fixing the Generator

Always check the schema definition to determine the correct syntax:
- Look in `internal/provider/*_resource.go` for schema definitions
- SingleNestedAttribute → `attr = { ... }`
- Block → `block { ... }`

