---
page_title: "criblio_notification_template Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Reads a Handlebars notification template in the current workspace.
---

# criblio_notification_template (Data Source)

Reads a custom or built-in Handlebars notification template. Requires V2 Notifications to be enabled. The provider selects `engine=handlebars` automatically.

## Example Usage

```terraform
data "criblio_notification_template" "email" {
  id = "default-email"
}
```

## Schema

### Required

- `id` (String) Template identifier.

### Read-Only

- `content` (String) Handlebars template content.
- `created_by` (String) User who created the template.
- `description` (String) Template description.
- `engine` (String) Template engine (`handlebars`).
- `last_updated` (Number) Last update timestamp in milliseconds since the Unix epoch.
- `type` (String) Notification target type.
