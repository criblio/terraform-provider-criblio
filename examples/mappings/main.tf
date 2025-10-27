# Example: Stream mapping ruleset - routes workers to different fleets based on environment and role
resource "criblio_mapping_ruleset" "stream_mappings" {
  id      = "stream_mappings"
  product = "stream"

  conf = {
    functions = [
      {
        # Production leaders - high priority for mission-critical data
        conf = {
          add = [
            {
              name  = "groupId"
              value = "prod-leaders"
            }
          ]
        }
        description = "Production leaders"
        disabled    = false
        filter      = "env == \"prod\" && role == \"leader\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Production workers
        conf = {
          add = [
            {
              name  = "groupId"
              value = "prod-workers"
            }
          ]
        }
        description = "Production workers"
        disabled    = false
        filter      = "env == \"prod\" && role == \"worker\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Staging environment
        conf = {
          add = [
            {
              name  = "groupId"
              value = "staging-fleet"
            }
          ]
        }
        description = "Staging environment"
        disabled    = false
        filter      = "env == \"staging\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Development environment
        conf = {
          add = [
            {
              name  = "groupId"
              value = "dev-fleet"
            }
          ]
        }
        description = "Development environment"
        disabled    = false
        filter      = "env == \"dev\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Operations specific routing
        conf = {
          add = [
            {
              name  = "groupId"
              value = "ops-fleet"
            }
          ]
        }
        description = "Operations team workers"
        disabled    = false
        filter      = "team == \"operations\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Default fallback - catch all
        conf = {
          add = [
            {
              name  = "groupId"
              value = "default-fleet"
            }
          ]
        }
        description = "Default fallback"
        disabled    = false
        filter      = "true"
        final       = true
        group_id    = "default"
        id          = "eval"
      }
    ]
  }
}

# Example: Edge mapping ruleset - routes edge devices by region and device type
resource "criblio_mapping_ruleset" "edge_mappings" {
  id      = "edge_mappings"
  product = "edge"

  conf = {
    functions = [
      {
        # North America region - routers and switches
        conf = {
          add = [
            {
              name  = "groupId"
              value = "na-network-devices"
            }
          ]
        }
        description = "North America network devices"
        disabled    = false
        filter      = "region == \"na\" && device_type == \"network\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # North America region - IoT devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "na-iot-devices"
            }
          ]
        }
        description = "North America IoT devices"
        disabled    = false
        filter      = "region == \"na\" && device_type == \"iot\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Europe region - all devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "eu-edge-devices"
            }
          ]
        }
        description = "Europe edge devices"
        disabled    = false
        filter      = "region == \"eu\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Asia Pacific region
        conf = {
          add = [
            {
              name  = "groupId"
              value = "apac-edge-devices"
            }
          ]
        }
        description = "Asia Pacific edge devices"
        disabled    = false
        filter      = "region == \"apac\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # High priority devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "high-priority-fleet"
            }
          ]
        }
        description = "High priority devices"
        disabled    = false
        filter      = "priority == \"high\""
        final       = false
        group_id    = "default"
        id          = "eval"
      },
      {
        # Default fallback for unknown devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "default-edge-fleet"
            }
          ]
        }
        description = "Default edge fleet"
        disabled    = false
        filter      = "true"
        final       = true
        group_id    = "default"
        id          = "eval"
      }
    ]
  }
}
