---
# hand-maintained -- no generator in this repo produces data-source docs, see docs/resources/monitor.md's "How this was built"
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
- `aetos_metadata` (Attributes Map) Arbitrary key-value metadata.
- `managed_by` (String) Stamped 'terraform' by the backend when provisioned via the Terraform provider (cribl/cribl#43133).
- `name` (String)
- `priority_scalar_value` (Attributes) (see resource docs for the full nested schema)
- `query` (Attributes Map) Monitor queries keyed by query label. (see resource docs for the full nested schema)
- `search_mode` (String) Logs monitors only.
- `silence` (List of String) IDs of silence windows that suppress this monitor's alerts.
- `team_scalar_value` (Attributes) (see resource docs for the full nested schema)
- `template_params` (Map of String) Per-query template parameters, keyed by query label -- still opaque JSON due to a codegen limitation with map-of-map schemas (INFRA-12649), not a design choice.
- `type` (String)
- `unit` (String)

`detection_config` does not exist as a field -- it's four typed sub-blocks (`anomaly_config`/`change_config`/`outlier_config`/`forecast_config`), exactly one populated matching `type`, correctly nested under the API's `detectionConfig` on the wire via an explicit discriminator on the sibling `type` field (not shape-matching, which is unsafe -- see `docs/resources/monitor.md`'s "How this was built").

`priority`/`team`/`firingCondition`/`firingRule`/`metadata`/`notification` are each modeled as a pair of named sub-blocks (a concrete one and an `_inheritable` one for profile-linked values) rather than narrowed to a single branch -- see `docs/resources/monitor.md` for the full nested schemas.

Full nested attribute schemas match the `criblio_monitor` resource -- see `docs/resources/monitor.md` for the complete, codegen-accurate field list.
