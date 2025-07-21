![workflow status](https://github.com/criblio/terraform-provider-criblio/actions/workflows/e2e_checks.yaml/badge.svg) ![workflow status](https://github.com/criblio/terraform-provider-criblio/actions/workflows/tf_provider_release.yaml/badge.svg)

<!-- Start Summary [summary] -->
## Summary

Cribl API Reference: This API Reference lists available REST endpoints, along with their supported operations for accessing, creating, updating, or deleting resources. See our complementary product documentation at [docs.cribl.io](http://docs.cribl.io).
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
  * [Requirements](#requirements)
  * [Installation](#installation)
  * [Available Resources and Data Sources](#available-resources-and-data-sources)
  * [Testing the provider locally](#testing-the-provider-locally)
  * [Contributing](#contributing)
  * [License](#license)

<!-- End Table of Contents [toc] -->

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19
<!-- Start Installation [installation] -->
## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    criblio = {
      source  = "criblio/criblio"
      version = "1.0.55"
    }
  }
}

provider "criblio" {
  # Configuration options
}
```
<!-- End Installation [installation] -->

<!-- Start Available Resources and Data Sources [operations] -->
## Available Resources and Data Sources

### Resources

* [criblio_appscope_config](docs/resources/appscope_config.md)
* [criblio_commit](docs/resources/commit.md)
* [criblio_database_connection](docs/resources/database_connection.md)
* [criblio_deploy](docs/resources/deploy.md)
* [criblio_destination](docs/resources/destination.md)
* [criblio_event_breaker_ruleset](docs/resources/event_breaker_ruleset.md)
* [criblio_global_var](docs/resources/global_var.md)
* [criblio_grok](docs/resources/grok.md)
* [criblio_group](docs/resources/group.md)
* [criblio_hmac_function](docs/resources/hmac_function.md)
* [criblio_lookup_file](docs/resources/lookup_file.md)
* [criblio_pack](docs/resources/pack.md)
* [criblio_pack_pipeline](docs/resources/pack_pipeline.md)
* [criblio_parquet_schema](docs/resources/parquet_schema.md)
* [criblio_parser_lib_entry](docs/resources/parser_lib_entry.md)
* [criblio_pipeline](docs/resources/pipeline.md)
* [criblio_project](docs/resources/project.md)
* [criblio_regex](docs/resources/regex.md)
* [criblio_schema](docs/resources/schema.md)
* [criblio_source](docs/resources/source.md)
* [criblio_subscription](docs/resources/subscription.md)
### Data Sources

* [criblio_appscope_config](docs/data-sources/appscope_config.md)
* [criblio_config_version](docs/data-sources/config_version.md)
* [criblio_database_connection](docs/data-sources/database_connection.md)
* [criblio_destination](docs/data-sources/destination.md)
* [criblio_event_breaker_ruleset](docs/data-sources/event_breaker_ruleset.md)
* [criblio_global_var](docs/data-sources/global_var.md)
* [criblio_grok](docs/data-sources/grok.md)
* [criblio_group](docs/data-sources/group.md)
* [criblio_hmac_function](docs/data-sources/hmac_function.md)
* [criblio_lookup_file](docs/data-sources/lookup_file.md)
* [criblio_pack](docs/data-sources/pack.md)
* [criblio_pack_pipeline](docs/data-sources/pack_pipeline.md)
* [criblio_parquet_schema](docs/data-sources/parquet_schema.md)
* [criblio_parser_lib_entry](docs/data-sources/parser_lib_entry.md)
* [criblio_pipeline](docs/data-sources/pipeline.md)
* [criblio_project](docs/data-sources/project.md)
* [criblio_regex](docs/data-sources/regex.md)
* [criblio_schema](docs/data-sources/schema.md)
* [criblio_source](docs/data-sources/source.md)
* [criblio_subscription](docs/data-sources/subscription.md)
<!-- End Available Resources and Data Sources [operations] -->

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
