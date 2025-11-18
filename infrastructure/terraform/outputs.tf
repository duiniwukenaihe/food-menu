# =============================================================================
# Template Outputs
# =============================================================================

# Uncomment once the template module is created and enabled

# output "template_id" {
#   description = "The VM ID of the created template"
#   value       = module.template.template_id
# }

# output "template_name" {
#   description = "The name of the created template"
#   value       = module.template.template_name
# }

# =============================================================================
# Control Plane Outputs
# =============================================================================

# Uncomment once the control-plane module is created and enabled

# output "control_plane_nodes" {
#   description = "Information about control plane nodes"
#   value = {
#     count      = length(module.control_plane.node_ids)
#     node_ids   = module.control_plane.node_ids
#     node_names = module.control_plane.node_names
#     ip_addresses = module.control_plane.ip_addresses
#   }
# }

# output "control_plane_endpoint" {
#   description = "Kubernetes API server endpoint"
#   value       = module.control_plane.cluster_endpoint
# }

# output "control_plane_ssh_connection_strings" {
#   description = "SSH connection strings for control plane nodes"
#   value       = module.control_plane.ssh_connection_strings
# }

# =============================================================================
# Worker Pool Outputs
# =============================================================================

# Uncomment once the worker-pool module is created and enabled

# output "worker_nodes" {
#   description = "Information about worker nodes"
#   value = {
#     count      = length(module.worker_pool.node_ids)
#     node_ids   = module.worker_pool.node_ids
#     node_names = module.worker_pool.node_names
#     ip_addresses = module.worker_pool.ip_addresses
#   }
# }

# output "worker_ssh_connection_strings" {
#   description = "SSH connection strings for worker nodes"
#   value       = module.worker_pool.ssh_connection_strings
# }

# =============================================================================
# Cluster Information
# =============================================================================

# Uncomment once modules are created and enabled

# output "cluster_info" {
#   description = "Complete cluster information"
#   value = {
#     control_plane = {
#       count      = length(module.control_plane.node_ids)
#       endpoint   = module.control_plane.cluster_endpoint
#       nodes      = module.control_plane.node_names
#       ips        = module.control_plane.ip_addresses
#     }
#     workers = {
#       count = length(module.worker_pool.node_ids)
#       nodes = module.worker_pool.node_names
#       ips   = module.worker_pool.ip_addresses
#     }
#     kubernetes_version = var.kubernetes_version
#     cni_plugin        = var.kubernetes_cni
#     environment       = var.environment
#   }
# }

# output "all_node_connection_info" {
#   description = "Connection information for all nodes in the cluster"
#   value = {
#     control_plane = {
#       for idx, name in module.control_plane.node_names : name => {
#         vm_id      = module.control_plane.node_ids[idx]
#         ip_address = module.control_plane.ip_addresses[idx]
#         ssh_command = "ssh ubuntu@${module.control_plane.ip_addresses[idx]}"
#       }
#     }
#     workers = {
#       for idx, name in module.worker_pool.node_names : name => {
#         vm_id      = module.worker_pool.node_ids[idx]
#         ip_address = module.worker_pool.ip_addresses[idx]
#         ssh_command = "ssh ubuntu@${module.worker_pool.ip_addresses[idx]}"
#       }
#     }
#   }
#   sensitive = false
# }

# =============================================================================
# Placeholder Outputs (Active until modules are implemented)
# =============================================================================

output "status" {
  description = "Current infrastructure status"
  value = {
    message = "Terraform configuration initialized. Module implementations pending."
    next_steps = [
      "1. Implement the template module in ./modules/template",
      "2. Implement the control-plane module in ./modules/control-plane",
      "3. Implement the worker-pool module in ./modules/worker-pool",
      "4. Uncomment module blocks in main.tf",
      "5. Uncomment outputs in outputs.tf"
    ]
  }
}

output "configuration_summary" {
  description = "Summary of the current configuration"
  value = {
    proxmox_endpoint    = var.proxmox_endpoint
    proxmox_node        = var.proxmox_node_name
    storage             = var.proxmox_storage
    network_bridge      = var.proxmox_network_bridge
    control_plane_count = var.control_plane_count
    worker_count        = var.worker_count
    kubernetes_version  = var.kubernetes_version
    cni_plugin          = var.kubernetes_cni
    environment         = var.environment
    project_name        = var.project_name
  }
}
