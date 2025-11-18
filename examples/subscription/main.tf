resource "criblio_subscription" "my_subscription" {
  //count required for cribl internal testing
  //count is not required for most customer implementations
  count = var.onprem == false ? 1 : 0

  description = "test subscription"
  disabled    = true
  filter      = "test"
  group_id    = "default"
  id          = "my_subscription"
  pipeline    = "passthru"
}

resource "criblio_subscription" "my_subscription_with_enabled" {
  //count required for cribl internal testing
  //count is not required for most customer implementations
  count = var.onprem == false ? 1 : 0

  description = "test subscription with enabled"
  disabled    = false
  filter      = "test"
  group_id    = "default"
  id          = "my_subscription_with_enabled"
  pipeline    = "passthru"
}

output "subscription" {
  value = criblio_subscription.my_subscription
}
