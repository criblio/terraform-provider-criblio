resource "criblio_search_usage_group" "my_search_usage_group" {
  id          = "power_users"
  description = "Members with elevated search limits."
  enabled     = true

  # Each rule is keyed by its limit type, and repeats that type in `type`.
  rules = {
    max_results_per_search = {
      type  = "maxResultsPerSearch"
      limit = 100000
    }
    max_running_time_per_search = {
      type  = "maxRunningTimePerSearch"
      limit = 300
    }
    # Scheduling limits accept a number or a percentage string.
    max_concurrent_adhoc_searches_per_user = {
      type  = "maxConcurrentAdhocSearchesPerUser"
      limit = "50%"
    }
    max_concurrent_searches = {
      type  = "maxConcurrentSearches"
      limit = 10
    }
  }

  # Members, keyed by the user ID from your identity provider. The API accepts
  # only IDs of users that already exist in the organization and silently drops
  # anything else, and a user belongs to one Usage Group at a time: assigning one
  # here removes them from the group they were in, and deleting this group
  # returns them to the catch-all `default` group.
  users = {
    "oidc|MyProvider|00u0000000000000000" = {
      id    = "oidc|MyProvider|00u0000000000000000"
      email = "search.admin@example.com"
    }
  }
}
