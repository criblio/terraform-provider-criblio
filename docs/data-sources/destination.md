---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "criblio_destination Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Destination DataSource
---

# criblio_destination (Data Source)

Destination DataSource

## Example Usage

```terraform
data "criblio_destination" "my_destination" {
  group_id = "...my_group_id..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The consumer group to which this instance belongs. Defaults to 'Cribl'.
