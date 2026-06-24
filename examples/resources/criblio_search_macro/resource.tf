resource "criblio_search_macro" "my_searchmacro" {
  description = "Filters to high-severity events."
  group_id    = "default_search"
  id          = "error_filter"
  replacement = "severity >= \"Error\""
  tags        = "errors,prod"
}
