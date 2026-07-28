---
name: criblio-terraform-contributor
description: Use when helping internal contributors work on terraform-provider-criblio, especially first-time internal contributors, issue reporters, maintainers preparing PRs, reviewers evaluating generated-code changes, repo-access onboarding, or deciding whether a change belongs in upstream OpenAPI, Terraform overlays, tools/codegen, docs, examples, tests, or import-cli. Covers required repo write access, asking the SRE team to add internal contributors to the repo, repo orientation, generated-code boundaries, contribution workflow, local validation, internal issue guidance, and handoff to the resource-development skill for endpoint/resource implementation.
---

# Cribl Terraform Contributor

Orient contributors and keep changes aligned with the provider's generated-code workflow. For deep implementation of a new or changed Terraform resource/data source, switch to `criblio-terraform-resource-development`.

## Contribution Boundary

- Contributors who make code changes need write access to `criblio/terraform-provider-criblio`. If they do not have access, instruct them to ask the SRE team to add them to this repository before starting implementation work.
- Contributions are internal-only for now. If someone is not an internal contributor, do not guide them through implementation work in this repo.
- For maintainers or internal contributors, prefer changes to OpenAPI specs, `terraform-overlay.yml`, `terraform-mgmt-overlay.yml`, `tools/codegen`, custom companion files, docs/examples inputs, tests, or import-cli metadata over hand-editing generated provider files.
- Do not encourage edits to generated files as the primary fix unless the file is protected by `.codegen-ignore`, is an intentional companion/manual file, or the user explicitly asks for a temporary diagnostic patch.
- Preserve backward compatibility for Terraform attributes, import IDs, examples, state behavior, and existing acceptance-test semantics.

## First Response Workflow

1. Identify the contributor type:
   - Internal issue reporter: help them produce a reproducible issue, minimal Terraform config, environment details, logs, and expected/actual behavior.
   - Maintainer or internal contributor with repo write access: help them find the right source of truth and validation path.
   - Internal user without repo write access: tell them to ask the SRE team to add them to `criblio/terraform-provider-criblio`, then proceed with planning or issue preparation until access is granted.
   - Reviewer: inspect behavior, generated-code provenance, docs/examples/tests, and import/export impact.
2. Check repo state before making changes:
   - `git status --short --branch`
   - `git diff --name-status`
   - Avoid overwriting unrelated user changes.
3. Classify the change:
   - Provider auth/client/routing behavior: inspect `internal/auth`, `internal/restclient`, and provider configuration.
   - Resource/data source surface: inspect overlays/spec/codegen before generated files.
   - Generated output drift: run or inspect `make generate`.
   - Docs/examples/tests only: verify examples match actual schemas and tests cover the behavior.
   - Import/export: inspect `tools/import-cli/internal/registry` and `tools/import-cli/internal/export`.
4. Choose a minimal validation plan and tell the user what will be checked.

## Contributor Workflow

Use this workflow when someone wants to contribute a code change:

1. Confirm access and intent:
   - Ask whether the contributor has write access to `criblio/terraform-provider-criblio`.
   - If not, tell them to ask the SRE team to add them to the repo.
   - If they cannot get write access, help them prepare an internal issue or implementation plan instead of changing code.
2. Prepare the branch:
   - Check `git status --short --branch`.
   - Sync from the expected base branch if the user approves network operations.
   - Create a focused branch named after the ticket and topic, such as `INFRA-12345-fix-routes-import`.
3. Locate the source of truth:
   - For generated resource/data source behavior, inspect specs, overlays, and `tools/codegen` before generated files.
   - For auth/client behavior, inspect `internal/auth` and `internal/restclient`.
   - For docs/examples/test-only changes, still verify the generated schema supports the documented Terraform.
4. Implement through durable inputs:
   - Change upstream/spec, overlays, generator, companion files, tests, examples, docs inputs, or import-cli metadata as appropriate.
   - Run `make generate` after spec, overlay, or generator edits.
5. Inspect the result:
   - Review generated provider code, docs, examples, acceptance tests, provider registration, and import-cli coverage.
   - Confirm no generated file changed without a reproducible source change unless it is intentionally protected/manual.
6. Validate:
   - Run focused unit tests first.
   - Run broader build/test commands when the change touches shared generator, provider, auth, or import-cli behavior.
   - Run live acceptance tests only when credentials and environment are available.
7. Prepare PR:
   - Summarize the user-facing Terraform surface, source/spec changes, generated outputs, tests run, and any environment limitations.
   - Ask before pushing or opening the PR if network access or credentials are required.

## Repo Orientation

- `README.md`: user-facing install, auth, cloud/on-prem configuration, and supported-resource notes.
- `CONTRIBUTING.md`: issue-reporting expectations and generated-code contribution context.
- `terraform-overlay.yml`: provider API overlays for generated resources/data sources.
- `terraform-mgmt-overlay.yml`: management API overlays.
- `tools/codegen`: generator, templates, parser tests, generated docs/examples/tests behavior.
- `internal/provider`: generated and companion provider implementation.
- `internal/auth` and `internal/restclient`: auth, routing, request/response behavior.
- `tests/acceptance`: generated and custom acceptance tests.
- `tools/import-cli`: import/export discovery and HCL rendering.
- `examples` and `docs`: user-facing Terraform examples and generated registry docs.
- `.codegen-ignore`: files intentionally protected from regeneration.

## Safe Change Patterns

- Prefer source changes that survive regeneration:
  - OpenAPI/spec sync for true upstream API shape changes.
  - Overlay annotations for Terraform-specific naming, sensitivity, computed/default behavior, fixed singleton identities, and diff suppression.
  - Generator changes for repeated generated-code patterns.
  - Companion/manual files only for behavior the generator cannot express cleanly.
- After any spec, overlay, or generator edit, run `make generate` and inspect generated output.
- Verify generated docs and examples are valid against actual provider schemas. A copied example with a non-existent attribute is a user-facing bug.
- When changing resource behavior, check import-cli export output so generated HCL does not include read-only, computed-only, sensitive, or write-only fields.
- For on-prem support, remember that Search, Lake, Lakehouse, workspace management, and some notification/search resources can be cloud-only.

## Validation Commands

Use the narrowest useful checks first:

```bash
go test -count=1 ./tools/codegen/...
go test -count=1 ./internal/provider -run '<FocusedTestName>'
go test -count=1 ./tools/import-cli/...
go build ./...
go build ./tools/import-cli/...
```

Useful Make targets:

```bash
make generate
make unit-test
make unit-test-import-cli
make build-import-cli
make acceptance-test
```

Live acceptance tests require credentials/environment and should be scoped:

```bash
TF_ACC=1 go test -v -parallel 1 -run 'Test<ResourceName>' ./tests/acceptance/
```

If a validation command cannot run because credentials, network, Terraform, or a live Cribl environment are missing, say exactly what was not run and why.

## Review Checklist

When reviewing a contribution, prioritize findings over summary:

- Does the change edit generated files without changing the spec, overlay, generator, or protected companion source?
- Does `make generate` reproduce the generated diff?
- Do docs and examples match the actual schema?
- Are import IDs, data sources, and import-cli export paths still correct?
- Are sensitive/write-only/API-owned fields preserved or filtered correctly?
- Are plan-no-diff, import, and update/delete semantics covered by tests where applicable?
- Are cloud/on-prem and license restrictions handled with explicit skips or routing behavior?
- Did the change introduce broad generator behavior that needs parser/render tests?

## Internal Issue Guidance

For internal issue reports, help contributors include:

- Provider version, Terraform version, Go version if building locally, OS/architecture.
- Deployment type: Cribl.Cloud or on-prem, plus relevant workspace/search/management context.
- Minimal Terraform configuration with secrets removed.
- Exact command run: `terraform plan`, `terraform apply`, `terraform import`, import-cli command, etc.
- Expected behavior and actual behavior.
- Relevant logs/errors with tokens and secrets redacted.
- Whether the issue reproduces after `terraform refresh` or a second `terraform plan`.

## Handoff To Implementation Skill

Switch to `criblio-terraform-resource-development` when the user asks to add a resource/data source, migrate a resource, wire a new endpoint, alter overlay annotations, change codegen behavior, update import-cli support for a generated resource, or prepare a PR for endpoint support.
