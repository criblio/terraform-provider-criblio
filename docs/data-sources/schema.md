---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "criblio_schema Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Schema DataSource
---

# criblio_schema (Data Source)

Schema DataSource

## Example Usage

```terraform
data "criblio_schema" "my_schema" {
  group_id = "...my_group_id..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The consumer group to which this instance belongs. Defaults to 'Cribl'.

### Read-Only

- `description` (String)
- `id` (String) The ID of this resource.
- `schema` (String) JSON schema matching standards of draft version 2019-09
