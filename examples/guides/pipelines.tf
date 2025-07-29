resource "criblio_pipeline" "my_pipeline" {
  id = "my-first-tf-pipeline"
  group_id = var.group-hybrid
  conf = {
    description = "Test pipeline for Terraform provider"
    functions = [
      {
        id = "eval"
        description = "Set organizationId test"
        disabled = false
        filter = "true"
        final = false
        group_id = var.group-hybrid
        conf = {
          add = jsonencode([
            {
              name = "organizationId"
              value = "tenantId || __forwardedAttrs.__metadata.env.TENANT_ID"
          }
          ])
        }
      },
        {
            id = "serde"
            description = "parse"
            disabled = false
            filter = "true"
            final = false
            group_id = var.group-hybrid
            conf = {
              mode = jsonencode("extract")
              type = jsonencode("csv")
              srcField = jsonencode("_raw")
              fields = jsonencode([
                "future_use_0",
                "receive_time",
                "serial_number",
                "type",
                "threat_content_type",
                "future_use_1",
                "generated_time",
                "source_ip",
                "destination_ip",
                "nat_source_ip",
                "nat_destination_ip",
                "rule_name",
                "source_user",
                "destination_user",
                "application",
                "virtual_system",
                "source_zone",
                "destination_zone",
                "inbound_interface",
                "outbound_interface",
                "log_action",
                "future_use_2",
                "session_id",
                "repeat_count",
                "source_port",
                "destination_port",
                "nat_source_port",
                "nat_destination_port",
                "flags",
                "protocol",
                "action",
                "bytes",
                "bytes_sent",
                "bytes_received",
                "packets",
                "start_time",
                "elapsed_time",
                "category",
                "future_use_3",
                "sequence_number",
                "action_flags",
                "source_location",
                "destination_location",
                "future_use_4",
                "packets_sent",
                "packets_received",
                "session_end_reason",
                "device_group_hierarchy_level_1",
                "device_group_hierarchy_level_2",
                "device_group_hierarchy_level_3",
                "device_group_hierarchy_level_4",
                "virtual_system_name",
                "device_name",
                "action_source",
                "source_vm_uuid",
                "destination_vm_uuid",
                "tunnel_id_imsi",
                "monitor_tag_imei",
                "parent_session_id",
                "parent_start_time",
                "tunnel_type",
                "sctp_association_id",
                "sctp_chunks",
                "sctp_chunks_sent",
                "chunkchunk"
              ])
            }
        }

    ]
    output = "default"
    streamtags = [""]
  }
  depends_on = [criblio_group.my_group_defaulthybrid]
}

locals {
  pipeline_versions = data.criblio_config_version.my_pipelineconfigversion.items
  latest_version    = length(local.pipeline_versions) > 0 ? local.pipeline_versions[length(local.pipeline_versions) - 1] : null
  timestamp = "${timestamp()}"
}

data "criblio_config_version" "my_pipelineconfigversion" {
  id = var.group-hybrid
  depends_on = [criblio_commit.my_pipecommit]
}

resource "criblio_commit" "my_pipecommit" {
  effective = true
  group     = var.group-hybrid
  message   = "terraform commit pipeline"
}
resource "criblio_deploy" "my_pipedeploy_safe" {
  id = var.group-hybrid
  version = length(data.criblio_config_version.my_pipelineconfigversion.items) > 0 ? data.criblio_config_version.my_pipelineconfigversion.items[length(data.criblio_config_version.my_pipelineconfigversion.items) - 1] : "default"
}