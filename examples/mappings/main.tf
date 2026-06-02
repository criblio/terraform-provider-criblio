# Example: Stream mapping ruleset - routes workers to different fleets based on environment and role
resource "criblio_mapping_ruleset" "stream_mappings" {
  product = "stream"

  conf = {
    functions = [
      {
        # Production leaders - high priority for mission-critical data
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'prod-leaders'"
            }
          ]
        }
        description = "Production leaders"
        disabled    = false
        filter      = "env == \"prod\" && role == \"leader\""
        group_id    = "default"
      },
      {
        # Production workers
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'prod-workers'"
            }
          ]
        }
        description = "Production workers"
        disabled    = false
        filter      = "env == \"prod\" && role == \"worker\""
        group_id    = "default"
      },
      {
        # Staging environment
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'staging-fleet'"
            }
          ]
        }
        description = "Staging environment"
        disabled    = false
        filter      = "env == \"staging\""
        group_id    = "default"
      },
      {
        # Development environment
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'dev-fleet'"
            }
          ]
        }
        description = "Development environment"
        disabled    = false
        filter      = "env == \"dev\""
        group_id    = "default"
      },
      {
        # Operations specific routing
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'ops-fleet'"
            }
          ]
        }
        description = "Operations team workers"
        disabled    = false
        filter      = "team == \"operations\""
        group_id    = "default"
      },
      {
        # Default fallback - catch all
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'default-fleet'"
            }
          ]
        }
        description = "Default fallback"
        disabled    = false
        filter      = "true"
        group_id    = "default"
      }
    ]
  }
}

# Example: Edge mapping ruleset - routes edge devices by region and device type
resource "criblio_mapping_ruleset" "edge_mappings" {
  product = "edge"

  conf = {
    functions = [
      {
        # North America region - routers and switches
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'na-network-devices'"
            }
          ]
        }
        description = "North America network devices"
        disabled    = false
        filter      = "region == \"na\" && device_type == \"network\""
        group_id    = "default"
      },
      {
        # North America region - IoT devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'na-iot-devices'"
            }
          ]
        }
        description = "North America IoT devices"
        disabled    = false
        filter      = "region == \"na\" && device_type == \"iot\""
        group_id    = "default"
      },
      {
        # Europe region - all devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'eu-edge-devices'"
            }
          ]
        }
        description = "Europe edge devices"
        disabled    = false
        filter      = "region == \"eu\""
        group_id    = "default"
      },
      {
        # Asia Pacific region
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'apac-edge-devices'"
            }
          ]
        }
        description = "Asia Pacific edge devices"
        disabled    = false
        filter      = "region == \"apac\""
        group_id    = "default"
      },
      {
        # High priority devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'high-priority-fleet'"
            }
          ]
        }
        description = "High priority devices"
        disabled    = false
        filter      = "priority == \"high\""
        group_id    = "default"
      },
      {
        # Default fallback for unknown devices
        conf = {
          add = [
            {
              name  = "groupId"
              value = "'default-edge-fleet'"
            }
          ]
        }
        description = "Default edge fleet"
        disabled    = false
        filter      = "true"
        group_id    = "default"
      }
    ]
  }
}
