---
page_title: "criblio_search_dataset_intelligence Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  Reads saved intelligence for a Search dataset.
---

# criblio_search_dataset_intelligence (Data Source)

Reads the current generation status and saved Markdown for an existing Search dataset. AI and agentic Search must be enabled, AI consent accepted, and the caller must have read access to the dataset. Returns an error when no intelligence exists. Reading does not start generation or wait for completion.

## Example Usage

```terraform
data "criblio_search_dataset_intelligence" "example" {
  dataset_id = "application_logs"
}
```

## Schema

### Required

- `dataset_id` (String) ID of the Search dataset.

### Read-Only

- `content` (String) Generated Markdown, null when no output exists.
- `status` (String) Generation status: `processing`, `complete`, `partial`, or `failed`.
- `failure_reason` (String) Failure details, when present.
- `generated_at` (Number) Generation start time in Unix epoch milliseconds.
- `updated_at` (Number) Last manual edit time in Unix epoch milliseconds, when present.
