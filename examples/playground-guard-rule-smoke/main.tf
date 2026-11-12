terraform {
  required_providers {
    criblio = {
      source = "criblio/criblio"
    }
  }
}

# organization_id and workspace_id come from ~/.cribl/credentials.
provider "criblio" {
  cloud_domain = "cribl-playground.cloud"
}

resource "criblio_guard_rule" "badge" {
  group_id    = "default"
  id          = "tf_smoke_badge_id"
  description = "Internal employee badge identifiers"
  regex       = "/\\bBADGE-\\d{8}\\b/"
  lib         = "custom"
  rulesets    = ["tf-smoke-internal-identifiers"]

  context_keywords = [
    {
      keyword   = "badge"
      placement = "before"
    },
  ]
}

data "criblio_guard_rule" "by_id" {
  group_id = "default"
  id       = criblio_guard_rule.badge.id
}

data "criblio_guard_rules" "all" {
  group_id   = "default"
  depends_on = [criblio_guard_rule.badge]
}

output "guard_rule" {
  value = criblio_guard_rule.badge
}

output "guard_rule_by_id_regex" {
  value = data.criblio_guard_rule.by_id.regex
}

output "guard_rule_count" {
  value = length(data.criblio_guard_rules.all.items)
}
