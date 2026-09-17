---
page_title: "criblio_monitors Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Lists Lakehouse engine metrics monitors.
---

# criblio_monitors (Data Source)

Lists monitors from `/products/lakehouse_engine_metrics/monitors`.

## Example Usage

```terraform
data "criblio_monitors" "all" {}
```

## Schema

### Read-Only

- `items` (List of Objects) Monitor configurations returned by the API.
