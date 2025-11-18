output "template_name" {
  description = "Name of the created Ubuntu template"
  value       = "ubuntu-${var.ubuntu_version}-template"
}

output "template_vmid" {
  description = "VMID of the created Ubuntu template"
  value       = var.vm_id
}

output "template_storage_path" {
  description = "Storage path of the uploaded Ubuntu image"
  value       = "${var.storage_pool}:iso/ubuntu-${var.ubuntu_version}-server-cloudimg-${var.ubuntu_architecture}.img"
}

output "template_node" {
  description = "Proxmox node where the template was created"
  value       = var.proxmox_node
}