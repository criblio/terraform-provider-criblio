---
page_title: "criblio_project_pipelines Data Source - terraform-provider-criblio"
description: |-
  Lists pipelines belonging to a Cribl Project.
---

# criblio_project_pipelines (Data Source)

Lists pipelines belonging to a Cribl Project.

## Example Usage

```terraform
data "criblio_project_pipelines" "all" {
  group_id   = "default"
  project_id = "example"
}
```

## Schema

### Required

- `group_id` (String) Worker group ID.
- `project_id` (String) Project ID.

### Read-Only

- `items` (List of Objects) Project pipelines and their configurations.
