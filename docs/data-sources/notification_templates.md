---
page_title: "criblio_notification_templates Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Lists Handlebars notification templates in the current workspace.
---

# criblio_notification_templates (Data Source)

Lists custom and built-in Handlebars notification templates. Requires V2 Notifications to be enabled. Legacy condition templates are excluded.

## Example Usage

```terraform
data "criblio_notification_templates" "all" {}
```

## Schema

### Read-Only

- `items` (Attributes List) Notification templates. Each item contains the fields below.

### Nested Schema for `items`

Read-Only:

- `id` (String) Template identifier.
- `content` (String) Handlebars template content.
- `created_by` (String) User who created the template.
- `description` (String) Template description.
- `engine` (String) Template engine (`handlebars`).
- `last_updated` (Number) Last update timestamp in milliseconds since the Unix epoch.
- `type` (String) Notification target type.
