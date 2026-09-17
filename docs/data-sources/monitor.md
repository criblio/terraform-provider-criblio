---
page_title: "criblio_monitor Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Reads a Lakehouse engine metrics monitor by ID.
---

# criblio_monitor (Data Source)

Reads a monitor from `/products/lakehouse_engine_metrics/monitors/{id}`.

## Example Usage

```terraform
data "criblio_monitor" "example" {
  id = "cpu-high"
}
```

## Schema

### Required

- `id` (String) Unique identifier for the monitor.

All other monitor attributes are computed from the API response, including its dataset, query, firing configuration, notification configuration, name, and type.
