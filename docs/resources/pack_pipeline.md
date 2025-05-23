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
          key = jsonencode("value")
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

- `conf` (Attributes) Requires replacement if changed. (see [below for nested schema](#nestedatt--conf))
- `group_id` (String) Group Id. Requires replacement if changed.
- `id` (String) Requires replacement if changed.
- `pack` (String) pack ID to POST. Requires replacement if changed.

<a id="nestedatt--conf"></a>
### Nested Schema for `conf`

Optional:

- `async_func_timeout` (Number) Time (in ms) to wait for an async function to complete processing of a data item. Requires replacement if changed.
- `description` (String) Requires replacement if changed.
- `functions` (Attributes List) List of Functions to pass data through. Requires replacement if changed. (see [below for nested schema](#nestedatt--conf--functions))
- `groups` (Attributes Map) Requires replacement if changed. (see [below for nested schema](#nestedatt--conf--groups))
- `output` (String) The output destination for events processed by this Pipeline. Default: "default"; Requires replacement if changed.
- `streamtags` (List of String) Tags for filtering and grouping in @{product}. Requires replacement if changed.

<a id="nestedatt--conf--functions"></a>
### Nested Schema for `conf.functions`

Optional:

- `conf` (Map of String) Not Null; Requires replacement if changed.
- `description` (String) Simple description of this step. Requires replacement if changed.
- `disabled` (Boolean) If true, data will not be pushed through this function. Requires replacement if changed.
- `filter` (String) Filter that selects data to be fed through this Function. Default: "true"; Requires replacement if changed.
- `final` (Boolean) If enabled, stops the results of this Function from being passed to the downstream Functions. Requires replacement if changed.
- `group_id` (String) Group ID. Requires replacement if changed.
- `id` (String) Function ID. Not Null; Requires replacement if changed.


<a id="nestedatt--conf--groups"></a>
### Nested Schema for `conf.groups`

Optional:

- `description` (String) Short description of this group. Requires replacement if changed.
- `disabled` (Boolean) Whether this group is disabled. Requires replacement if changed.
- `name` (String) Not Null; Requires replacement if changed.

## Import

Import is supported using the following syntax:

```shell
terraform import criblio_pack_pipeline.my_criblio_pack_pipeline "{ \"group_id\": \"\",  \"id\": \"\",  \"pack\": \"\"}"
```
