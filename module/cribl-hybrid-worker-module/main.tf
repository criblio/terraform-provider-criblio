
locals {
  # Construct the Cribl installation URL
  cribl_install_url = var.workspace_url != "" ? "${var.workspace_url}/init/install-worker.sh?group=${var.group}&token=${var.auth_token}&user=${var.cribl_user}&user_group=${var.cribl_group}&install_dir=${var.install_dir}" : "https://${var.workspace}-${var.organization}.${var.domain}/init/install-worker.sh?group=${var.group}&token=${var.auth_token}&user=${var.cribl_user}&user_group=${var.cribl_group}&install_dir=${var.install_dir}"
}

# Generate the installation script
locals {
  user_data_script = templatefile("${path.module}/templates/install-cribl-worker.tpl", {
    cribl_install_url = local.cribl_install_url
  })
} 