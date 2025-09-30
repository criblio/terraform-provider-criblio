resource "criblio_subscription" "my_subscription" {
  description = "Words describe things"
  disabled    = false
  filter      = "truthy"
  group_id    = "myExistingGroupId"
  id          = "myUniqueSubscriptionId"
  pipeline    = "myExistingPipelineId"
}