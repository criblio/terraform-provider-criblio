resource "criblio_pack" "full_config_pack" {
  id           = "pack-with-full-config"
  group_id     = "default"
  description  = "Pack with full pipeline configuration"
  disabled     = false
  display_name = "Pack with Full Config"
  version      = "1.0.0"
}

# Pack Pipeline with comprehensive configuration
# Demonstrates various pipeline functions including eval, serde, and serialize
resource "criblio_pack_pipeline" "AuditdLogs_main" {
  id         = "main"
  pack       = criblio_pack.full_config_pack.id
  group_id   = "default"
  depends_on = [criblio_pack.full_config_pack]

  conf = {
    streamtags = []

    functions = [
      {
        id          = "eval"
        filter      = "true"
        disabled    = false
        description = "Expose metadata environments"
        conf = jsonencode({
          add = [
            {
              name     = "domain"
              value    = "__metadata.env.CRIBL_CLOUD_DOMAIN"
              disabled = false
            },
            {
              name     = "organizationId"
              value    = "__metadata.env.CRIBL_CLOUD_TENANT_ID"
              disabled = false
            },
            {
              name     = "accountId"
              value    = "__metadata.env.CRIBL_CLOUD_ACCOUNT_ID"
              disabled = false
            },
            {
              name     = "service"
              value    = "'auditd'"
              disabled = false
            },
            {
              name     = "index"
              value    = "__metadata.env.CRIBL_CLOUD_DOMAIN"
              disabled = false
            },
            {
              name     = "sourcetype"
              value    = "'file'"
              disabled = false
            },
            {
              name     = "group"
              value    = "__metadata.cribl.group"
              disabled = false
            },
            {
              name     = "instanceId"
              value    = "__metadata.aws.identity.instanceId"
              disabled = false
            },
            {
              name     = "timestamp"
              value    = "_time"
              disabled = false
            }
          ]
        })
      },
      {
        id          = "serde"
        filter      = "true"
        disabled    = false
        description = "Extract _raw data to key value pairs"
        conf = jsonencode({
          mode         = "extract"
          type         = "kvp"
          clean_fields = true
          delim_char   = ","
          quote_char   = "\""
          null_value   = "-"
          src_field    = "_raw"
          escape_char  = "\\"
        })
      },
      {
        id       = "serialize"
        filter   = "true"
        disabled = true
        conf = jsonencode({
          dst_field    = "_raw"
          clean_fields = false
          type         = "kvp"
          fields = [
            "!_*",
            "!cribl_breaker",
            "!source",
            "!cribl_route",
            "!sourcetype",
            "!service",
            "*"
          ]
        })
      },
      {
        id          = "eval"
        filter      = "true"
        disabled    = false
        description = "Remove _raw"
        conf = jsonencode({
          remove = [
            "_raw"
          ]
        })
      },
    ]
  }
}

