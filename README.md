![workflow status](https://github.com/criblio/terraform-provider-criblio/actions/workflows/e2e_checks.yaml/badge.svg) ![workflow status](https://github.com/criblio/terraform-provider-criblio/actions/workflows/tf_provider_release.yaml/badge.svg)


# terraform-provider-criblio

A Terraform provider for managing Cribl resources.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
    }
  }
}

provider "criblio" {
  # Configuration options
}
```

## Authentication

The Cribl provider supports multiple authentication methods. You can configure authentication using environment variables, credentials file, or directly in the provider configuration block.

### Environment Variables

You can set the following environment variables:

```bash
# Direct authentication
export CRIBL_BEARER_TOKEN="your-bearer-token"

# OAuth authentication
export CRIBL_CLIENT_ID="your-client-id"
export CRIBL_CLIENT_SECRET="your-client-secret"
export CRIBL_ORGANIZATION_ID="your-organization-id"
export CRIBL_WORKSPACE_ID="your-workspace-id"

```

### Credentials File

You can store your credentials in `~/.cribl/credentials` or `~/.cribl` (legacy) with the following format:

```ini
[default]
client_id = your-client-id
client_secret = your-client-secret
organization_id = your-organization-id
workspace = your-workspace-id

[profile2]
client_id = another-client-id
client_secret = another-client-secret
organization_id = another-organization-id
workspace = another-workspace-id
```

To use a specific profile, set the `CRIBL_PROFILE` environment variable:

```bash
export CRIBL_PROFILE="profile2"
```

### Provider Configuration

You can also configure authentication directly in your Terraform configuration:

```hcl
provider "criblio" {
  # Using bearer token
  bearer_token = "your-bearer-token"

  # Or using OAuth credentials
  client_id        = "your-client-id"
  client_secret    = "your-client-secret"
  organization_id  = "your-organization-id"
  workspace_id     = "your-workspace-id"

}
```

### Authentication Methods

#### 1. Bearer Token

The simplest way to authenticate is using a bearer token:

```hcl
provider "criblio" {
  bearer_token = "your-bearer-token"
}
```

#### 2. OAuth Authentication (Recommended)

For OAuth authentication, you can use client credentials:

```hcl
provider "criblio" {
  client_id       = "your-client-id"
  client_secret   = "your-client-secret"
  organization_id = "your-organization-id"
  workspace_id    = "your-workspace-id"
}
```

### Example with Environment Variables

```hcl
# main.tf
provider "criblio" {
  # Credentials will be read from environment variables
}

# Use the provider
resource "criblio_pipeline" "example" {
  name = "example-pipeline"
  # ... other configuration
}
```

```bash
# Set environment variables
export CRIBL_CLIENT_ID="your-client-id"
export CRIBL_CLIENT_SECRET="your-client-secret"
export CRIBL_ORGANIZATION_ID="your-organization-id"
export CRIBL_WORKSPACE_ID="your-workspace-id"

# Run Terraform
terraform init
terraform plan
```

<!-- No Installation [installation] -->


<!-- No Summary [summary] -->


<!-- No Table of Contents [toc] -->

<!-- Start Available Resources and Data Sources [operations] -->
## Available Resources and Data Sources

### Resources

* [criblio_appscope_config](docs/resources/appscope_config.md)
* [criblio_commit](docs/resources/commit.md)
* [criblio_cribl_lake_dataset](docs/resources/cribl_lake_dataset.md)
* [criblio_database_connection](docs/resources/database_connection.md)
* [criblio_deploy](docs/resources/deploy.md)
* [criblio_destination](docs/resources/destination.md)
* [criblio_event_breaker_ruleset](docs/resources/event_breaker_ruleset.md)
* [criblio_global_var](docs/resources/global_var.md)
* [criblio_grok](docs/resources/grok.md)
* [criblio_group](docs/resources/group.md)
* [criblio_hmac_function](docs/resources/hmac_function.md)
* [criblio_lookup_file](docs/resources/lookup_file.md)
* [criblio_notification_target](docs/resources/notification_target.md)
* [criblio_pack](docs/resources/pack.md)
* [criblio_pack_breakers](docs/resources/pack_breakers.md)
* [criblio_pack_destination](docs/resources/pack_destination.md)
* [criblio_pack_lookups](docs/resources/pack_lookups.md)
* [criblio_pack_pipeline](docs/resources/pack_pipeline.md)
* [criblio_pack_routes](docs/resources/pack_routes.md)
* [criblio_pack_source](docs/resources/pack_source.md)
* [criblio_pack_vars](docs/resources/pack_vars.md)
* [criblio_parquet_schema](docs/resources/parquet_schema.md)
* [criblio_parser_lib_entry](docs/resources/parser_lib_entry.md)
* [criblio_pipeline](docs/resources/pipeline.md)
* [criblio_project](docs/resources/project.md)
* [criblio_regex](docs/resources/regex.md)
* [criblio_schema](docs/resources/schema.md)
* [criblio_search_dashboard](docs/resources/search_dashboard.md)
* [criblio_search_dashboard_category](docs/resources/search_dashboard_category.md)
* [criblio_search_dataset](docs/resources/search_dataset.md)
* [criblio_search_dataset_provider](docs/resources/search_dataset_provider.md)
* [criblio_search_macro](docs/resources/search_macro.md)
* [criblio_search_saved_query](docs/resources/search_saved_query.md)
* [criblio_search_usage_group](docs/resources/search_usage_group.md)
* [criblio_source](docs/resources/source.md)
* [criblio_subscription](docs/resources/subscription.md)
### Data Sources

* [criblio_appscope_config](docs/data-sources/appscope_config.md)
* [criblio_config_version](docs/data-sources/config_version.md)
* [criblio_cribl_lake_dataset](docs/data-sources/cribl_lake_dataset.md)
* [criblio_database_connection](docs/data-sources/database_connection.md)
* [criblio_destination](docs/data-sources/destination.md)
* [criblio_event_breaker_ruleset](docs/data-sources/event_breaker_ruleset.md)
* [criblio_global_var](docs/data-sources/global_var.md)
* [criblio_grok](docs/data-sources/grok.md)
* [criblio_group](docs/data-sources/group.md)
* [criblio_hmac_function](docs/data-sources/hmac_function.md)
* [criblio_lookup_file](docs/data-sources/lookup_file.md)
* [criblio_notification_target](docs/data-sources/notification_target.md)
* [criblio_pack](docs/data-sources/pack.md)
* [criblio_pack_breakers](docs/data-sources/pack_breakers.md)
* [criblio_pack_destination](docs/data-sources/pack_destination.md)
* [criblio_pack_lookups](docs/data-sources/pack_lookups.md)
* [criblio_pack_pipeline](docs/data-sources/pack_pipeline.md)
* [criblio_pack_routes](docs/data-sources/pack_routes.md)
* [criblio_pack_source](docs/data-sources/pack_source.md)
* [criblio_pack_vars](docs/data-sources/pack_vars.md)
* [criblio_parquet_schema](docs/data-sources/parquet_schema.md)
* [criblio_parser_lib_entry](docs/data-sources/parser_lib_entry.md)
* [criblio_pipeline](docs/data-sources/pipeline.md)
* [criblio_project](docs/data-sources/project.md)
* [criblio_regex](docs/data-sources/regex.md)
* [criblio_schema](docs/data-sources/schema.md)
* [criblio_search_dashboard](docs/data-sources/search_dashboard.md)
* [criblio_search_dashboard_category](docs/data-sources/search_dashboard_category.md)
* [criblio_search_dataset](docs/data-sources/search_dataset.md)
* [criblio_search_dataset_provider](docs/data-sources/search_dataset_provider.md)
* [criblio_search_macro](docs/data-sources/search_macro.md)
* [criblio_search_saved_query](docs/data-sources/search_saved_query.md)
* [criblio_search_usage_group](docs/data-sources/search_usage_group.md)
* [criblio_source](docs/data-sources/source.md)
* [criblio_subscription](docs/data-sources/subscription.md)
<!-- End Available Resources and Data Sources [operations] -->

<!-- No SDK Installation -->
<!-- No SDK Example Usage -->
<!-- Start Testing the provider locally [usage] -->
## Testing the provider locally

#### Local Provider

Should you want to validate a change locally, the `--debug` flag allows you to execute the provider against a terraform instance locally.

This also allows for debuggers (e.g. delve) to be attached to the provider.

```sh
go run main.go --debug
# Copy the TF_REATTACH_PROVIDERS env var
# In a new terminal
cd examples/your-example
TF_REATTACH_PROVIDERS=... terraform init
TF_REATTACH_PROVIDERS=... terraform apply
```

#### Compiled Provider

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

1. Execute `go build` to construct a binary called `terraform-provider-criblio`
2. Ensure that the `.terraformrc` file is configured with a `dev_overrides` section such that your local copy of terraform can see the provider binary

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/criblio/criblio" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
<!-- End Testing the provider locally [usage] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->
## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the license included in the repository.
