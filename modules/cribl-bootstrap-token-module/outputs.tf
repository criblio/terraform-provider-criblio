# Outputs for Cribl Bootstrap Token Module

output "bootstrap_token" {
  description = "The retrieved bootstrap token (only available if auto_execute is true)"
  value       = var.auto_execute ? trimspace(data.local_file.auth_token[0].content) : null
  sensitive   = true
}

output "bootstrap_script" {
  description = "The generated bootstrap token script content"
  value       = local.bootstrap_script
  sensitive   = true
}

output "script_command" {
  description = "The command to manually execute the bootstrap script"
  value       = "bash -c '${local.bootstrap_script}'"
  sensitive   = true
}

output "script_file_path" {
  description = "Path to the created script file (if create_script_file is true)"
  value       = var.create_script_file ? "${path.root}/${var.script_filename}" : null
}

output "token_file_path" {
  description = "Path to the token output file"
  value       = var.token_output_file
} 