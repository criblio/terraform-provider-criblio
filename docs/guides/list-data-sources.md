---
page_title: "List Data Sources"
subcategory: ""
description: |-
  Discover existing Cribl resources and reference them in Terraform.
---

# List data sources

List data sources discover resources created through the UI, API, or another Terraform configuration. Each exposes `items`, a collection of objects with the same attribute names used by the corresponding resource. Required attributes select the Worker Group, Pack, Project, Lake, or Organization. API offset/limit pagination is handled automatically; users receive the entire collection. Collection data sources return an empty collection when there are no objects.

`criblio_keys` requests the full collection without pagination parameters. Cribl's generic pagination sorts by `id`, while encryption keys use `keyId`, causing paginated requests to fail on affected versions.

`criblio_project_pipelines` also requests the full collection without pagination parameters. On affected versions, Cribl's project route middleware fails to rewrite URLs containing query strings, causing paginated requests to return 404.

`criblio_search_sources` requests the full collection without pagination parameters for compatibility with Search versions that ignore offset/limit. Changing source status in those responses can otherwise prevent repeated-page detection and cause excessive requests.

Existing data source names remain supported. `criblio_routes`, `criblio_pack_routes`, `criblio_pack_breakers`, `criblio_pack_lookups`, and `criblio_pack_vars` read individual objects, so their collection counterparts use distinct names below. Collectors retain their existing dynamic `items` representation to accommodate heterogeneous saved state.

```terraform
data "criblio_packs" "existing" {
  group_id = "default"
}

output "pack_ids" {
  value = [for pack in data.criblio_packs.existing.items : pack.id]
}

data "criblio_pack_pipelines" "existing" {
  group_id = "default"
  pack     = "my-pack"
}
```

## Managed collections

| Resource | List data source |
| --- | --- |
| `criblio_app` | [`criblio_apps`](../data-sources/apps.md) |
| `criblio_appscope_config` | [`criblio_appscope_configs`](../data-sources/appscope_configs.md) |
| `criblio_certificate` | [`criblio_certificates`](../data-sources/certificates.md) |
| `criblio_collector` | [`criblio_collectors`](../data-sources/collectors.md) |
| `criblio_cribl_lake_dataset` | [`criblio_cribl_lake_datasets`](../data-sources/cribl_lake_datasets.md) |
| `criblio_cribl_lake_house` | [`criblio_cribl_lake_houses`](../data-sources/cribl_lake_houses.md) |
| `criblio_custom_banner` | [`criblio_custom_banners`](../data-sources/custom_banners.md) |
| `criblio_database_connection` | [`criblio_database_connections`](../data-sources/database_connections.md) |
| `criblio_destination` | [`criblio_destinations`](../data-sources/destinations.md) |
| `criblio_event_breaker_ruleset` | [`criblio_event_breaker_rulesets`](../data-sources/event_breaker_rulesets.md) |
| `criblio_global_var` | [`criblio_global_vars`](../data-sources/global_vars.md) |
| `criblio_grok` | [`criblio_groks`](../data-sources/groks.md) |
| `criblio_group` | [`criblio_groups`](../data-sources/groups.md) |
| `criblio_guard_rule` | [`criblio_guard_rules`](../data-sources/guard_rules.md) |
| `criblio_hmac_function` | [`criblio_hmac_functions`](../data-sources/hmac_functions.md) |
| `criblio_key` | [`criblio_keys`](../data-sources/keys.md) |
| `criblio_lookup_file` | [`criblio_lookup_files`](../data-sources/lookup_files.md) |
| `criblio_mapping_ruleset` | [`criblio_mappings`](../data-sources/mappings.md) |
| `criblio_monitor` | [`criblio_monitors`](../data-sources/monitors.md) |
| `criblio_monitor_profile` | [`criblio_monitor_profiles`](../data-sources/monitor_profiles.md) |
| `criblio_monitor_shared_config` | [`criblio_monitor_shared_configs`](../data-sources/monitor_shared_configs.md) |
| `criblio_monitor_silence` | [`criblio_monitor_silences`](../data-sources/monitor_silences.md) |
| `criblio_notification` | [`criblio_notifications`](../data-sources/notifications.md) |
| `criblio_notification_policy` | [`criblio_notification_policies`](../data-sources/notification_policies.md) |
| `criblio_notification_target` | [`criblio_notification_targets`](../data-sources/notification_targets.md) |
| `criblio_notification_template` | [`criblio_notification_templates`](../data-sources/notification_templates.md) |
| `criblio_pack` | [`criblio_packs`](../data-sources/packs.md) |
| `criblio_pack_breakers` | [`criblio_pack_breaker_rulesets`](../data-sources/pack_breaker_rulesets.md) |
| `criblio_pack_collector` | [`criblio_pack_collectors`](../data-sources/pack_collectors.md) |
| `criblio_pack_destination` | [`criblio_pack_destinations`](../data-sources/pack_destinations.md) |
| `criblio_pack_lookups` | [`criblio_pack_lookup_files`](../data-sources/pack_lookup_files.md) |
| `criblio_pack_pipeline` | [`criblio_pack_pipelines`](../data-sources/pack_pipelines.md) |
| `criblio_pack_routes` | [`criblio_pack_routing_tables`](../data-sources/pack_routing_tables.md) |
| `criblio_pack_source` | [`criblio_pack_sources`](../data-sources/pack_sources.md) |
| `criblio_pack_vars` | [`criblio_pack_global_vars`](../data-sources/pack_global_vars.md) |
| `criblio_parquet_schema` | [`criblio_parquet_schemas`](../data-sources/parquet_schemas.md) |
| `criblio_parser_lib_entry` | [`criblio_parser_lib_entries`](../data-sources/parser_lib_entries.md) |
| `criblio_pipeline` | [`criblio_pipelines`](../data-sources/pipelines.md) |
| `criblio_project` | [`criblio_projects`](../data-sources/projects.md) |
| `criblio_project_pipeline` | [`criblio_project_pipelines`](../data-sources/project_pipelines.md) |
| `criblio_regex` | [`criblio_regexes`](../data-sources/regexes.md) |
| `criblio_routes` | [`criblio_routing_tables`](../data-sources/routing_tables.md) |
| `criblio_schema` | [`criblio_schemas`](../data-sources/schemas.md) |
| `criblio_search_dashboard` | [`criblio_search_dashboards`](../data-sources/search_dashboards.md) |
| `criblio_search_dashboard_category` | [`criblio_search_dashboard_categories`](../data-sources/search_dashboard_categories.md) |
| `criblio_search_dataset` | [`criblio_search_datasets`](../data-sources/search_datasets.md) |
| `criblio_search_dataset_provider` | [`criblio_search_dataset_providers`](../data-sources/search_dataset_providers.md) |
| `criblio_search_dataset_ruleset` | [`criblio_search_dataset_rulesets`](../data-sources/search_dataset_rulesets.md) |
| `criblio_search_datatype` | [`criblio_search_datatypes`](../data-sources/search_datatypes.md) |
| `criblio_search_datatype_ruleset` | [`criblio_search_datatype_rulesets`](../data-sources/search_datatype_rulesets.md) |
| `criblio_search_engine` | [`criblio_search_engines`](../data-sources/search_engines.md) |
| `criblio_search_macro` | [`criblio_search_macros`](../data-sources/search_macros.md) |
| `criblio_search_saved_query` | [`criblio_search_saved_queries`](../data-sources/search_saved_queries.md) |
| `criblio_search_source` | [`criblio_search_sources`](../data-sources/search_sources.md) |
| `criblio_secret` | [`criblio_secrets`](../data-sources/secrets.md) |
| `criblio_source` | [`criblio_sources`](../data-sources/sources.md) |
| `criblio_subscription` | [`criblio_subscriptions`](../data-sources/subscriptions.md) |
| `criblio_workspace` | [`criblio_workspaces`](../data-sources/workspaces.md) |

## Exclusions

These resources do not expose a listing API. Their existing read data sources remain available where applicable.

| Resource | Reason |
| --- | --- |
| `criblio_commit` | Action; no durable collection to list. |
| `criblio_worker_group_upgrade` | Action; no durable collection to list. |
| `criblio_deploy` | Action; no durable collection to list. |
| `criblio_group_system_settings` | One settings object per Worker Group; use `criblio_group_system_settings`. |
| `criblio_lakehouse_dataset_connection` | The connect/disconnect API has no collection GET endpoint. |
| `criblio_search_dataset_acl` | Access control belongs to one dataset; use `criblio_search_dataset_acl`. No ACL collection API. |
| `criblio_search_dataset_intelligence` | Settings belong to one dataset; use `criblio_search_dataset_intelligence`. No settings collection API. |
| `criblio_system_limits` | One limits object per Worker Group; use `criblio_system_limits`. |
| `criblio_workspace_system_settings` | One settings object per workspace; use `criblio_workspace_system_settings`. |

## Acceptance tests

Collection reads stop with an error after 1,000 pages instead of returning partial
results. This also bounds APIs that ignore pagination while returning changing
responses. Pack-scoped reads share successful Pack ID listings within each REST
client and Worker Group for one minute, including concurrent requests. Failed
listings are not cached; long-lived provider processes refresh expired entries.

Every collection has a read-only acceptance test. Run with `TF_ACC=1` and the provider's normal credentials. `CRIBL_TEST_GROUP_ID` optionally selects an existing group (defaults to `default`, or `default_search` for Search). Pack-scoped tests require `CRIBL_TEST_PACK_ID`; Project-scoped tests require `CRIBL_TEST_PROJECT_ID`; workspace discovery requires `CRIBL_ORGANIZATION_ID`. Tests skip when a required scope is not configured, and cloud-only APIs skip on-prem deployments. Empty collections are valid. Provider unit tests exercise collection decoding and pagination against local HTTP servers, including server-imposed page sizes and pagination errors.

Notification collection tests also skip the on-prem acceptance fixture: its license prohibits notifications and targets, and it does not enable V2 Notifications for policies and templates. Licensed on-prem deployments with the required features can still use these data sources.

`make acceptance-test` runs both phases simultaneously by default, with at most one resource test and four read-only data source tests active. Both phases share the configured Cribl environment, so concurrent resource creation and cleanup can affect discovery results. Run only discovery with `make acceptance-test-data-sources`, or adjust its concurrency with `DATA_SOURCE_TEST_PARALLELISM=2 make acceptance-test`. Run only resource lifecycles with `make acceptance-test-resources`. To run the phases sequentially, run the resource target first and then the data source target.

Lifecycle tests retain their first collection assertion and by-ID checks. Later update and unchanged-plan steps omit redundant collection blocks; tests that assert updated collection contents retain those reads.
