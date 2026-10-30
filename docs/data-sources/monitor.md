---
# hand-maintained
page_title: "criblio_monitor Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Monitor DataSource
---

# criblio_monitor (Data Source)

Monitor DataSource

## Example Usage

```terraform
data "criblio_monitor" "example" {
  id = "cpu-high"
}
```

<!-- schema hand-maintained; kept in sync with docs/resources/monitor.md by hand, not by tfplugindocs -->
## Schema

### Required

- `id` (String)

### Read-Only

- `anomaly_config` (Attributes) Present when `type = "anomaly"`. (see resource docs for the full nested schema)
- `change_config` (Attributes) Present when `type = "change"`. (see resource docs for the full nested schema)
- `outlier_config` (Attributes) Present when `type = "outlier"`. (see resource docs for the full nested schema)
- `forecast_config` (Attributes) Present when `type = "forecast"`. (see resource docs for the full nested schema)
- `dataset_id` (String)
- `description` (String)
- `enabled` (Boolean)
- `expr` (Attributes List) Query expressions evaluated to produce the monitor's series. (see resource docs for the full nested schema)
- `firing_condition` (Attributes) (see resource docs for the full nested schema)
- `firing_rule` (Attributes) (see resource docs for the full nested schema)
- `aetos_metadata` (Map of String) Arbitrary metadata key/value pairs.
- `managed_by` (String) External provisioner that owns this monitor.
- `profile_ref` (Attributes) Linked profile ID and fingerprint.
- `name` (String)
- `priority_scalar_value` (Attributes) (see resource docs for the full nested schema)
- `query` (Attributes Map) Monitor queries keyed by query label. (see resource docs for the full nested schema)
- `search_mode` (String) Logs monitors only.
- `silence` (List of String) IDs of silence windows that suppress this monitor's alerts.
- `team_scalar_value` (Attributes) (see resource docs for the full nested schema)
- `template_params` (String) JSON-encoded template parameter bindings keyed by query label.
- `type` (String)
- `unit` (String)

Detection configuration is returned in the block matching the monitor's `type`: `anomaly_config`, `change_config`, `outlier_config`, or `forecast_config`.

Profile-linked values are returned in `firing_condition_inheritable`, `firing_rule_inheritable`, `metadata_inheritable`, `priority_inheritable`, or `team_inheritable`, each containing `config_refs` with IDs and fingerprints. Notifications are returned in `monitor_policy_notification_config`, `monitor_target_notification_config`, or `notification_inheritable`.

Inherited queries contain `config_refs` in their query-label map entry. Inline queries contain `mode` and `promql`.

Full nested attribute schemas match the `criblio_monitor` resource -- see `docs/resources/monitor.md` for the complete, codegen-accurate field list.
