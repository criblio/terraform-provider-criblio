resource "criblio_guard_rule" "my_guard_rule" {
  group_id    = "default"
  id          = "employee_badge_id"
  description = "Internal employee badge identifiers"
  # The API validates this as a JavaScript regex literal, so keep the slashes.
  regex    = "/\\bBADGE-\\d{8}\\b/"
  rulesets = ["internal-identifiers"]

  context_keywords = [
    {
      keyword   = "badge"
      placement = "before"
    },
  ]
}
