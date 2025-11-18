# =============================================================================
# Kubernetes Node Pool Module Outputs
# =============================================================================

output "master_vms" {
  description = "Master VM resources"
  value       = proxmox_virtual_environment_vm.master
}

output "worker_vms" {
  description = "Worker VM resources"
  value       = proxmox_virtual_environment_vm.worker
}

output "master_names" {
  description = "Names of master nodes"
  value       = local.master_names
}

output "worker_names" {
  description = "Names of worker nodes"
  value       = local.worker_names
}

output "master_ips" {
  description = "IP addresses of master nodes"
  value       = local.master_ips
}

output "worker_ips" {
  description = "IP addresses of worker nodes"
  value       = local.worker_ips
}

output "master_count" {
  description = "Number of master nodes created"
  value       = var.master_count
}

output "worker_count" {
  description = "Number of worker nodes created"
  value       = var.worker_count
}

output "cluster_info" {
  description = "Cluster information"
  value = {
    name = var.cluster_name
    k8s_version = var.k8s_version
    network_plugin = var.network_plugin
    pod_network_cidr = var.pod_network_cidr
    master_count = var.master_count
    worker_count = var.worker_count
  }
}

output "ssh_access" {
  description = "SSH access information"
  value = {
    username = var.username
    master_ips = local.master_ips
    worker_ips = local.worker_ips
  }
  sensitive = true
}

output "kubernetes_config" {
  description = "Kubernetes configuration details"
  value = {
    image_repository = var.image_repository
    pod_network_cidr = var.pod_network_cidr
    network_plugin = var.network_plugin
    k8s_version = var.k8s_version
  }
}