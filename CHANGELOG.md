# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Release dates and contents through v1.25.55 were reconstructed from Git tags,
commit history, and tag-to-tag diffs.

## [Unreleased]

### Fixed
- Retried HTTP 429 responses using the server-provided `Retry-After` delay, including Config Helper admission throttling during bulk group and fleet operations.

## [1.25.55] - 2026-08-19

### Fixed
- Corrected a type in the destination resource schema.

## [1.25.54] - 2026-08-14

### Changed
- Expanded provider-owned schema handling to cover all supported property patterns.

## [1.25.51] - 2026-08-13

### Fixed
- Fixed route comments and group handling.
- Made route ordering deterministic.

## [1.25.50] - 2026-08-12

### Added
- Added import CLI support for on-premises-to-Cloud migration.

## [1.25.48] - 2026-08-10

### Added
- Expanded bulk exporter coverage, including project pipelines.

## [1.25.47] - 2026-08-04

### Added
- Added the `criblio_app` resource.

## [1.25.45] - 2026-07-30

### Fixed
- Fixed output router handling when descriptions are omitted.

## [1.25.44] - 2026-07-30

### Fixed
- Corrected OAuth schema handling for the Cribl Lake collector.

## [1.25.20] - 2026-07-29

### Added
- Added the v2 `criblio_search_datatype` resource.

## [1.25.14] - 2026-07-27

### Added
- Added filesystem collector support.

## [1.25.9] - 2026-07-21

### Added
- Added lookup file support to the bulk exporter.

## [1.25.5] - 2026-07-20

### Fixed
- Fixed regressions introduced during the in-house provider migration.

## [1.25.4] - 2026-07-15

### Fixed
- Corrected the source model used by the in-house provider.

## [1.25.3] - 2026-07-14

### Fixed
- Fixed handling of file-based packs.

## [1.25.2] - 2026-07-09

### Fixed
- Hardened concurrent authentication token requests.

## [1.25.1] - 2026-06-29

### Changed
- Completed the architectural migration from Speakeasy-generated provider code to
  the provider-owned REST client, OpenAPI overlays, and code-generation pipeline.
- Removed the remaining Speakeasy SDK, hooks, configuration, and generated artifacts.

## [1.24.22] - 2026-06-29

### Changed
- Expanded pack acceptance coverage while completing migration of pack behavior to
  the provider-owned implementation.

## [1.23.210] - 2026-06-26

### Changed
- Migrated the remaining legacy generated resources to the provider-owned client and
  code-generation pipeline.

## [1.23.119] - 2026-06-25

### Changed
- Added management API support to the provider-owned code generator.

## [1.23.85] - 2026-06-24

### Changed
- Migrated credential, HMAC function, key, and secret resources away from the
  Speakeasy SDK.

## [1.23.48] - 2026-06-19

### Changed
- Migrated a second set of resources to the provider-owned REST client and fixed
  associated code-generation issues.

## [1.23.47] - 2026-06-19

### Changed
- Introduced provider-owned authentication and REST client packages.
- Introduced OpenAPI synchronization, overlay merging, and provider code-generation
  tooling.
- Began migrating resources away from Speakeasy-generated implementations.

### Fixed
- Added validation for pipeline function IDs.

## [1.23.46] - 2026-06-03

### Fixed
- Fixed custom banner handling when `link` and `display` are omitted.

## [1.23.45] - 2026-06-02

### Added
- Added an import CLI option to include default IDs.

## [1.23.36] - 2026-05-21

### Fixed
- Added Cribl Lake dataset descriptions to prevent duplicate dataset creation.

## [1.23.34] - 2026-05-18

### Changed
- Expanded recognized secret suffixes used when exporting configuration.

## [1.23.32] - 2026-05-12

### Fixed
- Added collector `discover_request_params` schema coverage and corrected configuration
  update behavior.

## [1.23.28] - 2026-05-01

### Added
- Added custom banner resource and data source support.

### Fixed
- Fixed import CLI handling for search dashboards and legacy packs.

## [1.23.24] - 2026-04-29

### Fixed
- Corrected a Cribl Lake collector schema typo.

## [1.23.22] - 2026-04-24

### Changed
- Expanded collector schema coverage with pagination and previously missing fields.

## [1.23.20] - 2026-04-22

### Added
- Added advanced settings support to notification resources.

## [1.23.18] - 2026-04-22

### Added
- Added resources and data sources for search dataset rulesets, datatype rulesets,
  search engines, and search sources.

## [1.22.3] - 2026-04-13

### Fixed
- Fixed import CLI behavior and added a user agent to its API requests.

## [1.22.1] - 2026-04-08

### Fixed
- Removed invalid templated fields from pack source and destination schemas.

## [1.22.0] - 2026-04-07

### Fixed
- Fixed syslog input sources within packs.

## [1.21.9] - 2026-04-02

### Fixed
- Corrected global variable generation and regenerated affected provider code.

## [1.21.4] - 2026-03-27

### Changed
- Made changes to configured `id` and `group_id` values require resource replacement.

## [1.21.3] - 2026-03-26

### Fixed
- Fixed source generation issues following the Speakeasy and OpenAPI upgrade.

## [1.21.1] - 2026-03-25

### Changed
- Upgraded Speakeasy-generated code and the OpenAPI specification for source and
  destination resources.

## [1.20.145] - 2026-03-23

### Fixed
- Fixed pack updates and pack-scoped resources when using legacy IDs or packs with
  missing source definitions.

## [1.20.144] - 2026-03-19

### Security
- Pinned CIRCL and gRPC dependencies to versions containing security fixes.

## [1.20.140] - 2026-03-18

### Fixed
- Normalized pack IDs to lowercase as required by the API.

## [1.20.138] - 2026-03-16

### Changed
- Renamed the import CLI binary and exposed bulk export through an `export` subcommand.

## [1.20.135] - 2026-03-03

### Added
- Added integration tests and an independent release workflow for the import CLI.

## [1.20.133] - 2026-03-02

### Changed
- Updated release tooling; no provider behavior changed.

## [1.20.107] - 2026-02-16

### Fixed
- Fixed collector and pipeline resource behavior.

## [1.20.103] - 2026-01-30

### Fixed
- Updated file-based packs in place when referenced by routes instead of replacing
  them.

## [1.20.101] - 2026-01-29

### Fixed
- Fixed multiple pack lifecycle and file-handling issues.

## [1.20.91] - 2026-01-19

### Fixed
- Fixed provider configuration behavior when values are supplied through provider
  variables.

## [1.20.82] - 2026-01-08

### Added
- Added the `criblio_system_info` data source for retrieving the Cribl product version.

## [1.20.75] - 2026-01-02

### Fixed
- Updated group schema handling and ignored API-provided default group values that
  caused configuration drift.

## [1.20.72] - 2025-12-19

### Changed
- Improved API retry handling with exponential backoff.

## [1.20.69] - 2025-12-17

### Added
- Added group-scoped certificate, key, and secret resources and data sources.

## [1.20.64] - 2025-12-16

### Fixed
- Fixed pipeline function configuration handling.

## [1.20.62] - 2025-12-16

### Fixed
- Fixed `client_id` and `client_secret` provider configuration value handling.

## [1.20.60] - 2025-12-12

### Fixed
- Restored search dashboard category support.

## [1.20.59] - 2025-12-12

### Fixed
- Corrected the search dashboard schema.

## [1.20.48] - 2025-12-08

### Fixed
- Fixed dashboard resources containing time inputs.

## [1.20.33] - 2025-11-27

### Changed
- Decoupled end-to-end tests from provider releases; no provider behavior changed.

## [1.20.32] - 2025-11-20

### Fixed
- Fixed file-based pack update behavior.

## [1.20.27] - 2025-11-19

### Fixed
- Fixed pipeline function configuration handling.

## [1.20.23] - 2025-11-18

### Fixed
- Added retries when deleting sources and destinations.
- Corrected pipeline schema handling.

## [1.18.27] - 2025-11-12

### Added
- Added group system settings resource and data source support.

### Fixed
- Fixed authentication in end-to-end release validation.

## [1.18.21] - 2025-10-31

### Fixed
- Added the additional update operation required for file-based packs.

## [1.18.18] - 2025-10-30

### Fixed
- Restored the default route only when deleting route configuration.

## [1.18.10] - 2025-10-29

### Added
- Added on-premises deployment support to the provider.

## [1.18.8] - 2025-10-27

### Added
- Added mapping ruleset resource and data source support.
- Added a plural mappings data source.

## [1.17.19] - 2025-10-23

### Added
- Added plural data sources for collectors, destinations, notification targets, and
  sources.

### Fixed
- Fixed SBOM generation in the release workflow.

## [1.17.6] - 2025-10-10

### Added
- Added certificate resource and data source support.
- Expanded management API schema coverage.

## [1.14.5] - 2025-09-25

### Added
- Added workspace management resource and data source support.

## [1.14.4] - 2025-09-19

### Added
- Added the `criblio_instance_settings` data source.

## [1.12.2] - 2025-09-18

### Added
- Added routes resource and data source support.
- Added authentication support for Cribl government environments.

## [1.11.23] - 2025-09-16

### Fixed
- Corrected collector schema and documentation.

## [1.11.20] - 2025-09-11

### Changed
- Improved API URL construction and authentication for the current Cloud URL format.

## [1.11.17] - 2025-09-10

### Added
- Added collector resource and data source support.

## [1.11.0] - 2025-09-05

### Fixed
- Fixed persistent plan differences in pipeline resources.

## [1.10.2] - 2025-08-25

### Fixed
- Added retries for Lakehouse connections.
- Fixed Cribl Lake destination behavior.

## [1.9.7] - 2025-08-20

### Added
- Added a global variable resource example.

## [1.9.6] - 2025-08-19

### Added
- Added a regex resource example.

## [1.9.4] - 2025-08-18

### Fixed
- Stopped sending the unsupported `fileNameSuffix` field to the Cribl Lake API.

## [1.9.3] - 2025-08-18

### Fixed
- Removed `fileNameSuffix` from Cribl Lake output handling.

## [1.9.2] - 2025-08-13

### Added
- Added the Lakehouse dataset connection resource.

## [1.5.4] - 2025-08-12

### Added
- Added Cribl Lakehouse resource and data source support.

## [1.4.10] - 2025-08-08

### Fixed
- Fixed persistent plan differences for saved searches and cleaned up their schema.

## [1.4.5] - 2025-08-06

### Added
- Added notification resource and data source support for saved searches.

## [1.3.31] - 2025-08-04

### Fixed
- Removed conditional search URL construction.

## [1.3.28] - 2025-08-01

### Added
- Added the pack routes data source.

### Fixed
- Fixed evaluation function behavior.

## [1.3.27] - 2025-07-29

### Changed
- Retagged the v1.3.25 release; both tags reference the same commit.

## [1.3.25] - 2025-07-29

### Added
- Added Cribl Lake dataset and notification target resources and data sources.
- Added pack breaker, destination, lookup, route, source, and variable resources and
  corresponding data sources where supported.
- Added search dashboard, dashboard category, dataset, dataset provider, macro, saved
  query, and usage group resources and data sources.

## [1.0.45] - 2025-07-03

### Changed
- Updated release metadata and documentation; no provider behavior changed.

## [1.0.41] - 2025-06-27

### Added
- Added resources and data sources for AppScope configurations, database connections,
  event breaker rulesets, global variables, Grok patterns, HMAC functions, lookup
  files, Parquet schemas, parser library entries, projects, regexes, schemas, and
  subscriptions.

## [1.0.12] - 2025-06-02

### Fixed
- Fixed API URL templating.
- Added worker group deletion support.

## [1.0.8] - 2025-05-23

### Added
- Added commit and deploy resources.
- Added data sources for configuration versions, destinations, groups, and sources.

## [1.0.2] - 2025-05-22

### Fixed
- Corrected the generated provider name.

## [1.0.1] - 2025-05-21

### Added
- Added an end-to-end Cribl Stream example.

## [1.0.0] - 2025-05-21

### Added
- Initial provider release.
- Added destination, group, pack, pack pipeline, pipeline, and source resources.
- Added pack, pack pipeline, and pipeline data sources.

[Unreleased]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.55...HEAD
[1.25.55]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.54...v1.25.55
[1.25.54]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.51...v1.25.54
[1.25.51]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.50...v1.25.51
[1.25.50]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.48...v1.25.50
[1.25.48]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.47...v1.25.48
[1.25.47]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.45...v1.25.47
[1.25.45]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.44...v1.25.45
[1.25.44]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.20...v1.25.44
[1.25.20]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.14...v1.25.20
[1.25.14]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.9...v1.25.14
[1.25.9]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.5...v1.25.9
[1.25.5]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.4...v1.25.5
[1.25.4]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.3...v1.25.4
[1.25.3]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.2...v1.25.3
[1.25.2]: https://github.com/criblio/terraform-provider-criblio/compare/v1.25.1...v1.25.2
[1.25.1]: https://github.com/criblio/terraform-provider-criblio/compare/v1.24.22...v1.25.1
[1.24.22]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.210...v1.24.22
[1.23.210]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.119...v1.23.210
[1.23.119]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.85...v1.23.119
[1.23.85]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.48...v1.23.85
[1.23.48]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.47...v1.23.48
[1.23.47]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.46...v1.23.47
[1.23.46]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.45...v1.23.46
[1.23.45]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.36...v1.23.45
[1.23.36]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.34...v1.23.36
[1.23.34]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.32...v1.23.34
[1.23.32]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.28...v1.23.32
[1.23.28]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.24...v1.23.28
[1.23.24]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.22...v1.23.24
[1.23.22]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.20...v1.23.22
[1.23.20]: https://github.com/criblio/terraform-provider-criblio/compare/v1.23.18...v1.23.20
[1.23.18]: https://github.com/criblio/terraform-provider-criblio/compare/v1.22.3...v1.23.18
[1.22.3]: https://github.com/criblio/terraform-provider-criblio/compare/v1.22.1...v1.22.3
[1.22.1]: https://github.com/criblio/terraform-provider-criblio/compare/v1.22.0...v1.22.1
[1.22.0]: https://github.com/criblio/terraform-provider-criblio/compare/v1.21.9...v1.22.0
[1.21.9]: https://github.com/criblio/terraform-provider-criblio/compare/v1.21.4...v1.21.9
[1.21.4]: https://github.com/criblio/terraform-provider-criblio/compare/v1.21.3...v1.21.4
[1.21.3]: https://github.com/criblio/terraform-provider-criblio/compare/v1.21.1...v1.21.3
[1.21.1]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.145...v1.21.1
[1.20.145]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.144...v1.20.145
[1.20.144]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.140...v1.20.144
[1.20.140]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.138...v1.20.140
[1.20.138]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.135...v1.20.138
[1.20.135]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.133...v1.20.135
[1.20.133]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.107...v1.20.133
[1.20.107]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.103...v1.20.107
[1.20.103]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.101...v1.20.103
[1.20.101]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.91...v1.20.101
[1.20.91]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.82...v1.20.91
[1.20.82]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.75...v1.20.82
[1.20.75]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.72...v1.20.75
[1.20.72]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.69...v1.20.72
[1.20.69]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.64...v1.20.69
[1.20.64]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.62...v1.20.64
[1.20.62]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.60...v1.20.62
[1.20.60]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.59...v1.20.60
[1.20.59]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.48...v1.20.59
[1.20.48]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.33...v1.20.48
[1.20.33]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.32...v1.20.33
[1.20.32]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.27...v1.20.32
[1.20.27]: https://github.com/criblio/terraform-provider-criblio/compare/v1.20.23...v1.20.27
[1.20.23]: https://github.com/criblio/terraform-provider-criblio/compare/v1.18.27...v1.20.23
[1.18.27]: https://github.com/criblio/terraform-provider-criblio/compare/v1.18.21...v1.18.27
[1.18.21]: https://github.com/criblio/terraform-provider-criblio/compare/v1.18.18...v1.18.21
[1.18.18]: https://github.com/criblio/terraform-provider-criblio/compare/v1.18.10...v1.18.18
[1.18.10]: https://github.com/criblio/terraform-provider-criblio/compare/v1.18.8...v1.18.10
[1.18.8]: https://github.com/criblio/terraform-provider-criblio/compare/v1.17.19...v1.18.8
[1.17.19]: https://github.com/criblio/terraform-provider-criblio/compare/v1.17.6...v1.17.19
[1.17.6]: https://github.com/criblio/terraform-provider-criblio/compare/v1.14.5...v1.17.6
[1.14.5]: https://github.com/criblio/terraform-provider-criblio/compare/v1.14.4...v1.14.5
[1.14.4]: https://github.com/criblio/terraform-provider-criblio/compare/v1.12.2...v1.14.4
[1.12.2]: https://github.com/criblio/terraform-provider-criblio/compare/v1.11.23...v1.12.2
[1.11.23]: https://github.com/criblio/terraform-provider-criblio/compare/v1.11.20...v1.11.23
[1.11.20]: https://github.com/criblio/terraform-provider-criblio/compare/v1.11.17...v1.11.20
[1.11.17]: https://github.com/criblio/terraform-provider-criblio/compare/v1.11.0...v1.11.17
[1.11.0]: https://github.com/criblio/terraform-provider-criblio/compare/v1.10.2...v1.11.0
[1.10.2]: https://github.com/criblio/terraform-provider-criblio/compare/v1.9.7...v1.10.2
[1.9.7]: https://github.com/criblio/terraform-provider-criblio/compare/v1.9.6...v1.9.7
[1.9.6]: https://github.com/criblio/terraform-provider-criblio/compare/v1.9.4...v1.9.6
[1.9.4]: https://github.com/criblio/terraform-provider-criblio/compare/v1.9.3...v1.9.4
[1.9.3]: https://github.com/criblio/terraform-provider-criblio/compare/v1.9.2...v1.9.3
[1.9.2]: https://github.com/criblio/terraform-provider-criblio/compare/v1.5.4...v1.9.2
[1.5.4]: https://github.com/criblio/terraform-provider-criblio/compare/v1.4.10...v1.5.4
[1.4.10]: https://github.com/criblio/terraform-provider-criblio/compare/v1.4.5...v1.4.10
[1.4.5]: https://github.com/criblio/terraform-provider-criblio/compare/v1.3.31...v1.4.5
[1.3.31]: https://github.com/criblio/terraform-provider-criblio/compare/v1.3.28...v1.3.31
[1.3.28]: https://github.com/criblio/terraform-provider-criblio/compare/v1.3.27...v1.3.28
[1.3.27]: https://github.com/criblio/terraform-provider-criblio/compare/v1.3.25...v1.3.27
[1.3.25]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.45...v1.3.25
[1.0.45]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.41...v1.0.45
[1.0.41]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.12...v1.0.41
[1.0.12]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.8...v1.0.12
[1.0.8]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.2...v1.0.8
[1.0.2]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/criblio/terraform-provider-criblio/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/criblio/terraform-provider-criblio/releases/tag/v1.0.0
