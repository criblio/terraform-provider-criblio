---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "criblio_pack_pipeline Resource - terraform-provider-criblio"
subcategory: ""
description: |-
  PackPipeline Resource
---

# criblio_pack_pipeline (Resource)

PackPipeline Resource

## Example Usage

```terraform
resource "criblio_pack_pipeline" "my_packpipeline" {
  conf = {
    async_func_timeout = 9066
    description        = "...my_description..."
    functions = [
      {
        conf = {
          # ...
        }
        description = "...my_description..."
        disabled    = false
        filter      = "...my_filter..."
        final       = true
        group_id    = "...my_group_id..."
        id          = "...my_id..."
      }
    ]
    groups = {
      key = {
        description = "...my_description..."
        disabled    = true
        name        = "...my_name..."
      }
    }
    output = "...my_output..."
    streamtags = [
      "..."
    ]
  }
  group_id = "...my_group_id..."
  id       = "...my_id..."
  pack     = "...my_pack..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `conf` (Attributes) (see [below for nested schema](#nestedatt--conf))
- `group_id` (String) The consumer group to which this instance belongs. Defaults to 'Cribl'.
- `id` (String) Unique ID to PATCH for pack
- `pack` (String) pack ID to POST

<a id="nestedatt--conf"></a>
### Nested Schema for `conf`

Optional:

- `async_func_timeout` (Number) Time (in ms) to wait for an async function to complete processing of a data item
- `description` (String)
- `functions` (Attributes List) List of Functions to pass data through (see [below for nested schema](#nestedatt--conf--functions))
- `groups` (Attributes Map) (see [below for nested schema](#nestedatt--conf--groups))
- `output` (String) The output destination for events processed by this Pipeline. Default: "default"
- `streamtags` (List of String) Tags for filtering and grouping in @{product}

<a id="nestedatt--conf--functions"></a>
### Nested Schema for `conf.functions`

Optional:

- `conf` (Attributes) Not Null (see [below for nested schema](#nestedatt--conf--functions--conf))
- `description` (String) Simple description of this step
- `disabled` (Boolean) If true, data will not be pushed through this function
- `filter` (String) Filter that selects data to be fed through this Function. Default: "true"
- `final` (Boolean) If enabled, stops the results of this Function from being passed to the downstream Functions
- `group_id` (String) Group ID
- `id` (String) Function ID. Not Null

<a id="nestedatt--conf--functions--conf"></a>
### Nested Schema for `conf.functions.conf`



<a id="nestedatt--conf--groups"></a>
### Nested Schema for `conf.groups`

Optional:

- `description` (String) Short description of this group
- `disabled` (Boolean) Whether this group is disabled
- `name` (String) Not Null

## Import

Import is supported using the following syntax:

```shell
terraform import criblio_pack_pipeline.my_criblio_pack_pipeline '{"group_id": "", "pack": ""}'
```
