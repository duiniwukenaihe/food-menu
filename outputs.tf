output "template_name" {
  description = "Name of the created Ubuntu template"
  value       = module.ubuntu_template.template_name
}

output "template_vmid" {
  description = "VMID of the created Ubuntu template"
  value       = module.ubuntu_template.template_vmid
}

output "template_storage_path" {
  description = "Storage path of the uploaded Ubuntu image"
  value       = module.ubuntu_template.template_storage_path
}

output "template_node" {
  description = "Proxmox node where the template was created"
  value       = module.ubuntu_template.template_node
}