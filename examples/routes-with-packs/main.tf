resource "criblio_pack" "billing_pipeline" {
  id                     = "billing_pipeline"
  group_id               = "default"
  description            = "Billing Pipeline"
  disabled               = false
  display_name           = "Billing Pipeline"
  allow_custom_functions = true
  filename               = "${abspath(path.module)}/cribl-palo-alto-networks-source-1.0.0.crbl"
}

resource "criblio_routes" "routes" {
  id       = "default"
  group_id = "default"

  routes = [
    {
      name     = "Billing Pipeline"
      filter   = "__inputId.startsWith('tcp:test') || __inputId.startsWith('http:test:')"
      final    = false
      pipeline = "pack:${criblio_pack.billing_pipeline.id}"
      output   = jsonencode("devnull")
    }
  ]

  depends_on = [criblio_pack.billing_pipeline]
}
