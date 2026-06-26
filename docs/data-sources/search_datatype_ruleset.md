---
page_title: "criblio_search_datatype_ruleset Data Source - terraform-provider-criblio"
subcategory: ""
description: |-
  SearchDatatypeRuleset Data Source
---

# criblio_search_datatype_ruleset (Data Source)

SearchDatatypeRuleset Data Source

## Example Usage

```terraform
data "criblio_search_datatype_ruleset" "my_searchdatatyperuleset" {
  id = "default"
}
```

## Schema

### Optional

- `id` (String) Ruleset identifier.

### Read-Only

- `rules` (Attributes List) Rules evaluated in order for datatype routing. (see [below for nested schema](#nestedatt--rules))

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Read-Only:

- `datatype` (String) Target datatype id for events that match `kusto_expression`.
- `description` (String) Brief description of the rule.
- `disabled` (Boolean) If `true`, the rule is ignored for routing.
- `id` (String) Unique identifier for the rule.
- `kusto_expression` (String) Kusto predicate that selects events for this rule.
- `name` (String) Display name of the rule.
