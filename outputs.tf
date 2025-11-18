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

output "control_plane_nodes" {
  description = "Control plane node details"
  value       = length(module.control_plane) > 0 ? module.control_plane[0].nodes : {}
}

output "control_plane_ips" {
  description = "Control plane node IP addresses"
  value       = length(module.control_plane) > 0 ? module.control_plane[0].node_ips : {}
}

output "control_plane_ssh_connections" {
  description = "SSH connection strings for control plane nodes"
  value       = length(module.control_plane) > 0 ? module.control_plane[0].ssh_connection_strings : {}
}

output "worker_nodes" {
  description = "Worker node details"
  value       = length(module.workers) > 0 ? module.workers[0].nodes : {}
}

output "worker_ips" {
  description = "Worker node IP addresses"
  value       = length(module.workers) > 0 ? module.workers[0].node_ips : {}
}

output "worker_ssh_connections" {
  description = "SSH connection strings for worker nodes"
  value       = length(module.workers) > 0 ? module.workers[0].ssh_connection_strings : {}
}

output "kubeadm_join_data" {
  description = "Data needed for kubeadm join commands"
  value = {
    control_plane = length(module.control_plane) > 0 ? module.control_plane[0].kubeadm_join_data : null
    workers       = length(module.workers) > 0 ? module.workers[0].kubeadm_join_data : null
  }
}

output "all_nodes" {
  description = "All Kubernetes nodes with connection details"
  value = merge(
    length(module.control_plane) > 0 ? { for k, v in module.control_plane[0].nodes : k => merge(v, { type = "control-plane" }) } : {},
    length(module.workers) > 0 ? { for k, v in module.workers[0].nodes : k => merge(v, { type = "worker" }) } : {}
  )
}

output "ansible_inventory" {
  description = "Complete Ansible inventory for all nodes"
  value = {
    control_plane = length(module.control_plane) > 0 ? module.control_plane[0].ansible_inventory : { hosts = {} }
    workers       = length(module.workers) > 0 ? module.workers[0].ansible_inventory : { hosts = {} }
  }
}