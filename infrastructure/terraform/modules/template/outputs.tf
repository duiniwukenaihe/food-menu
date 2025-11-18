# =============================================================================
# Template Module Outputs
# =============================================================================

output "template_id" {
  description = "ID of the created VM template"
  value       = proxmox_virtual_environment_vm.template.vm_id
}

output "template_name" {
  description = "Name of the created VM template"
  value       = proxmox_virtual_environment_vm.template.name
}

output "template_node" {
  description = "Proxmox node where the template is located"
  value       = proxmox_virtual_environment_vm.template.node_name
}

output "template_storage" {
  description = "Storage pool where the template is stored"
  value       = var.storage
}

output "image_url" {
  description = "URL of the downloaded image"
  value       = var.image_url
}

output "cloud_init_config" {
  description = "Cloud-init configuration used for the template"
  value       = local_file.template_cloud_init.content
  sensitive   = true
}

output "template_ready" {
  description = "Whether the template is ready for use"
  value       = true
}