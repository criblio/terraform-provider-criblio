# Cribl Bootstrap Token Module
# This module provides a reusable way to get Cribl bootstrap tokens

terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

# Generate the bootstrap token script
locals {
  bootstrap_script = templatefile("${path.module}/templates/get_bootstrap_token.tpl", {
    client_id     = var.client_id
    client_secret = var.client_secret
    organization  = var.organization
    workspace     = var.workspace
    workspace_url = var.workspace_url
    environment   = var.environment
    group         = var.group
  })
}

# Create the script file
resource "local_file" "bootstrap_script" {
  count    = var.create_script_file ? 1 : 0
  filename = "${path.root}/${var.script_filename}"
  content  = local.bootstrap_script
  
  file_permission = "0755"
}

# Execute the script to get the token
resource "null_resource" "get_bootstrap_token" {
  count = var.auto_execute ? 1 : 0
  
  # Trigger re-execution when inputs change
  triggers = {
    client_id         = var.client_id
    client_secret     = var.client_secret
    organization      = var.organization
    workspace         = var.workspace
    group             = var.group
    environment       = var.environment
    token_output_file = var.token_output_file  # Store for destroy provisioner
  }

  provisioner "local-exec" {
    command = local.bootstrap_script_command
  }

  # Clean up the token file when destroyed
  provisioner "local-exec" {
    when    = destroy
    command = "rm -f ${self.triggers.token_output_file}"
  }

  # Note: depends_on must be static, so we handle dependencies via triggers instead
}

# Read the token from the output file
data "local_file" "auth_token" {
  count      = var.auto_execute ? 1 : 0
  filename   = var.token_output_file
  depends_on = [null_resource.get_bootstrap_token]
}

locals {
  # Command to execute the bootstrap script
  # When not creating a file, we need to write the script to a temp file first to avoid quoting issues
  bootstrap_script_command = var.create_script_file ? "./${var.script_filename} > ${var.token_output_file}" : "echo '${base64encode(local.bootstrap_script)}' | base64 -d | bash > ${var.token_output_file}"
} 