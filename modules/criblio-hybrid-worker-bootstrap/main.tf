
data "criblio_instance_settings" "instancesettings" {
}


locals {
  safe_install_dir  = urlencode(var.install_dir)
  safe_group_tags   = (join("&", [for tag in var.group_tags : format("tag=%s", urlencode(tag))]))
  leader_auth_token = data.criblio_instance_settings.instancesettings.items[0].instance_settings_schema1.master.auth_token

  user_data_script = templatefile("${path.module}/templates/install-cribl-worker.tftpl", {
    safe_install_dir  = local.safe_install_dir
    safe_group_tags   = local.safe_group_tags
    os_cribl_user     = var.os_cribl_user
    os_cribl_group    = var.os_cribl_group
    group_id          = var.group_id
    organization_id   = var.organization_id
    workspace_id      = var.workspace_id
    cloud_domain      = var.cloud_domain
    leader_auth_token = local.leader_auth_token
  })
}

