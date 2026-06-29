---
name: criblio-terraform-resource-development
description: Use when adding or changing a Cribl API endpoint and defining the matching Terraform provider resource or data source in terraform-provider-criblio. Covers upstream OpenAPI sync, terraform-overlay.yml and terraform-mgmt-overlay.yml annotations, generated resources/data sources/docs/examples/acceptance tests, import-cli export support, branch creation, validation, and PR preparation for both CRUD and singleton/read-only endpoints.
---

# Cribl Terraform Resource Development

Use this skill when an engineer asks to add a new endpoint to the provider, generate a Terraform resource or data source, migrate a resource into the custom codegen path, or prepare the branch and PR for endpoint support.

## Core Rules

- Prefer upstream OpenAPI and provider-owned codegen over hand-written provider files.
- Preserve backward compatibility with existing Terraform attributes, import IDs, examples, and state behavior.
- Use existing generated resources as examples before adding new generator behavior.
- Ask before fetching from the network, pushing a branch, or opening a PR.
- Do not hand-edit generated files unless the change is intentionally protected by `.codegen-ignore` or belongs in a companion file.
- Treat `make generate` as the source of truth for generated provider files, docs, examples, and generated acceptance tests.

## Start Here

Gather the missing facts before editing. Ask only for what cannot be inferred from the repo:

- Terraform name: resource/data source name, expected `criblio_<name>` shape, and legacy names if any.
- API shape: CRUD object, singleton settings object, read-only data source, list-only data source, or action endpoint.
- Scope: stream workspace `/m/{groupId}`, search workspace `/m/{groupId}/search`, system/global, pack-scoped, or management API.
- Spec source: whether the upstream spec is already updated, or whether the agent should run the local sync script.
- Identity: import ID format, path params, fixed IDs, case-insensitive IDs, and update/delete support.
- Fields: sensitive/write-only fields, API-owned computed fields, fields omitted or reformatted by GET, defaults, and null-vs-empty behavior.
- Tests: live acceptance environment constraints, cloud-only or on-prem-only behavior, license restrictions, and required example coverage.

## Branch Workflow

1. Check `git status --short` and avoid overwriting unrelated user changes.
2. If the user wants a new branch and the working tree allows it, create a branch named for the ticket and resource, for example `INFRA-12345-add-<resource>`.
3. Keep commits focused: spec/overlay/generator changes, generated output, tests/import-cli/docs, and cleanup should all be explainable in the PR.
4. Before pushing or opening a PR, ask for approval if credentials/network access are needed.

## Working Loop

After gathering facts, restate the development plan in a short checklist before making broad changes. Then implement end to end:

1. Update or sync the spec.
2. Add overlay annotations and generator support if needed.
3. Run generation.
4. Inspect generated provider code, docs, examples, and tests.
5. Add import-cli support and export filtering.
6. Run focused validation.
7. Commit, push, and create the PR when the user approves network operations.

## Spec And Overlay Workflow

1. Prefer updated upstream spec:
   - If the user has a new upstream spec, copy or merge it using the repo's existing process.
   - If the user wants the agent to fetch it, run `make sync-openapi` after approval when network access may be required.
2. If the upstream spec is missing or incomplete, patch the narrowest overlay:
   - Use `terraform-overlay.yml` for provider API resources.
   - Use `terraform-mgmt-overlay.yml` for management API resources.
3. Annotate operations:
   - `x-terraform-resource: true`
   - `x-terraform-resource-name: <Name>`
   - `x-terraform-read: <Name>`
   - `x-terraform-update: <Name>`
   - `x-terraform-delete: <Name>`
   - `x-terraform-list: <Name>`
   - `x-terraform-list-name: <Names>`
4. Annotate schemas:
   - `x-terraform-name` for legacy or awkward Terraform field names.
   - `x-terraform-sensitive: true` for secrets, tokens, passwords, and webhook URLs.
   - `x-terraform-prefer-state: true` for write-only or masked fields and values omitted by GET.
   - `x-terraform-computed: true` for API-owned fields.
   - `x-terraform-optional-computed: true` for configurable fields that the API may default.
   - Fixed singleton values with OpenAPI `const` or `default` plus `x-terraform-fixed-value`.
   - Diff suppression/prefer-state annotations for API-normalized strings, arrays, maps, or JSON.
5. For oneOf/union APIs, preserve existing block names and hoisted fields. If there is no discriminator, verify generated flat dispatch behavior against similar resources such as source, destination, notification target, or search dataset provider.

## Generate And Inspect

1. Run `make generate`.
2. Inspect the generated files that apply:
   - `internal/provider/<name>_types.go`
   - `internal/provider/<name>_client.go`
   - `internal/provider/<name>_resource.go`
   - `internal/provider/<name>_data_source.go`
   - `internal/provider/<plural>_data_source.go`
   - `docs/resources/<name>.md`
   - `docs/data-sources/<name>.md`
   - `examples/resources/<name>/resource.tf`
   - `examples/data-sources/<name>/data-source.tf`
   - `tests/acceptance/<name>_test.go`
   - `tests/acceptance/<name>_data_source_test.go`
3. Verify provider registration in `internal/provider/provider.go`:
   - Resource constructor appears in `Resources`.
   - Single and list data source constructors appear in `DataSources` when applicable.
   - If registration is missing, fix generator behavior when practical; only patch `provider.go` directly as a short-term compatibility fix.
4. Inspect `applyAPIToState` and schema generation for sensitive fields, prefer-state fields, path params, fixed IDs, optional/computed defaults, and null-vs-empty collections.
5. If the generated shape is wrong, prefer changing OpenAPI overlay or `tools/codegen` instead of editing generated provider files.

## Tests, Docs, Examples

- Generated acceptance tests must be real tests, not permanent scaffolds.
- Resource tests should cover create, read, update, plan-no-diff, import, and delete when the API supports them.
- Singleton/read-only data source tests should cover successful read and stable state.
- List data source tests are useful only when the environment can reliably contain at least one item; otherwise pair them with a resource-created item or skip with a clear environment guard.
- Examples should be valid, minimal, and compatible with both docs and acceptance tests.
- Search resources normally use `group_id = "default_search"`; stream resources normally use `group_id = "default"`.
- Cloud-only, on-prem-only, and license-gated endpoints need acceptance skips or separate test helpers, not broken default tests.
- Sensitive values must not appear in plan diffs, generated docs as real secrets, or import-cli HCL output.
- Upstream properties whose API names start with `__` are Cribl-internal helper fields. They must not be generated as Terraform attributes because framework `tfsdk` tags cannot start with underscores and these fields are not user-facing configuration.

## Import CLI And Export

Update import-cli in the same change for every importable/exportable resource:

- `tools/import-cli/internal/registry/import_metadata.go` for import ID format and discovery hints.
- `tools/import-cli/internal/registry/modeltypes.go` for generated model mappings.
- `tools/import-cli/internal/export/config.go` for computed-only, read-only, sensitive, write-only, or API-owned fields that must be excluded from generated HCL.
- Import/export tests when adding special filtering, oneOf blocks, dynamic objects, singleton identities, or custom discovery.

Verify import-cli discovery does not silently skip API items with missing IDs. If discovery needs a fallback ID extraction rule, add it deliberately and test it.

## Validation

Run the narrowest useful checks first, then broaden:

```bash
go test ./tools/codegen/...
go test ./tests/acceptance -run TestDoesNotExist -count=0
go build ./...
go build ./tools/import-cli/...
go test ./tools/import-cli/...
```

Run live acceptance tests when credentials and environment are available:

```bash
TF_ACC=1 go test -v -parallel 1 -run 'Test<ResourceName>' ./tests/acceptance/
```

For import-cli changes, also run the relevant unit tests and, when available, the integration export/import flow for the workspace in question.

## PR Checklist

Before asking the user to review or opening a PR:

- `make generate` has been run after every spec/overlay/generator edit.
- Generated docs, examples, resource tests, and data source tests are present or intentionally skipped with a reason.
- Provider registration includes the resource and applicable data sources.
- Import-cli can discover/export the resource without sensitive or computed-only drift.
- Acceptance tests either pass or have a documented environment/license skip.
- No new generated resource is added to `.codegen-ignore` unless a companion or hand-written implementation is intentional.
- The PR description includes endpoint/spec changes, generated Terraform surface, tests run, and known environment limitations.
