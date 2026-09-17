---
# hand-maintained -- no generator in this repo produces data-source docs, see docs/resources/monitor.md's "How this was built"
page_title: "criblio_monitors Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Monitors DataSource
---

# criblio_monitors (Data Source)

Monitors DataSource

## Example Usage

```terraform
data "criblio_monitors" "all" {}
```

<!-- schema hand-maintained; kept in sync with docs/resources/monitor.md by hand, not by tfplugindocs -->
## Schema

### Read-Only

- `items` (Attributes List) (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only: same field list and nested schemas as the `criblio_monitor` data source -- see `docs/data-sources/monitor.md` and `docs/resources/monitor.md` for the complete, codegen-accurate field list.
