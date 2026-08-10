---
page_title: "criblio_project_pipeline Data Source - terraform-provider-criblio"
description: |-
  Reads a pipeline belonging to a Cribl Project.
---

# criblio_project_pipeline (Data Source)

Reads a pipeline belonging to a Cribl Project.

## Example Usage

```terraform
data "criblio_project_pipeline" "main" {
  group_id   = "default"
  project_id = "example"
  id         = "main"
}
```

## Schema

### Required

- `group_id` (String) Worker group ID.
- `project_id` (String) Project ID.
- `id` (String) Pipeline ID.

### Read-Only

- `conf` (Attributes) Pipeline configuration, including its output, functions, groups, tags, description, and asynchronous function timeout.
