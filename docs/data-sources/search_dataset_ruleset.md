---
page_title: "criblio_search_dataset_ruleset Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  SearchDatasetRuleset Data Source
---

# criblio_search_dataset_ruleset (Data Source)

SearchDatasetRuleset Data Source

## Example Usage

```terraform
data "criblio_search_dataset_ruleset" "my_searchdatasetruleset" {
  id = "default"
}
```

## Schema

### Optional

- `id` (String) Unique identifier for the ruleset.

### Read-Only

- `rules` (Attributes List) Ordered rules for routing events to Datasets. The first matching rule wins. Each rule can include an extend expression. (see [below for nested schema](#nestedatt--rules))

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Read-Only:

- `dataset` (String) The target Dataset id.
- `description` (String) Brief description of the rule.
- `disabled` (Boolean) If `true`, the rule is ignored for routing.
- `extend_expression` (String) KQL-style extend expression.
- `extend_expression_enabled` (Boolean) Whether `extend_expression` is applied.
- `id` (String) Unique identifier for the rule.
- `kusto_expression` (String) Kusto predicate that selects events for this rule.
- `name` (String) Display name of the rule.
- `send_data_to` (String) Destination behavior, such as `destinationDataset` or `drop`.
