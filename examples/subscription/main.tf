
resource "criblio_subscription" "my_subscription" {
  description = "test subscription"
  disabled    = true
  filter      = "test"
  group_id    = "default"
  id          = "my_subscription"
  pipeline    = "passthru"
}

resource "criblio_subscription" "my_subscription_with_enabled" {
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
