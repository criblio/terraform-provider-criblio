---
page_title: "criblio_pack_collector Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Reads a Collector within a Pack.
---

# criblio_pack_collector (Data Source)

Reads a Collector within a Pack, including its type-specific collector configuration.

## Example Usage

```terraform
data "criblio_pack_collector" "rest" {
  group_id = "default"
  pack     = "my-pack"
  id       = "rest-api"
}
```

## Schema

Required:

- `group_id` (String) Worker group ID.
- `pack` (String) Pack ID.
- `id` (String) Collector ID within the pack.

Read-only attributes include `environment`, `ignore_group_jobs_limit`, `resume_on_boot`,
`ttl`, `worker_affinity`, and the `input_collector_*` objects documented in the
[pack collector resource](../resources/pack_collector.md). The object matching the
collector type is populated; the other collector objects are null.
