output "cluster_name" {
  description = "Name of the Kubernetes cluster"
  value       = var.cluster_name
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "master_instance_ids" {
  description = "IDs of the master instances"
  value       = module.k8s_nodes.master_instance_ids
}

output "worker_instance_ids" {
  description = "IDs of the worker instances"
  value       = module.k8s_nodes.worker_instance_ids
}

output "master_public_ips" {
  description = "Public IPs of the master instances"
  value       = module.k8s_nodes.master_public_ips
}

output "worker_public_ips" {
  description = "Public IPs of the worker instances"
  value       = module.k8s_nodes.worker_public_ips
}

output "master_private_ips" {
  description = "Private IPs of the master instances"
  value       = module.k8s_nodes.master_private_ips
}

output "worker_private_ips" {
  description = "Private IPs of the worker instances"
  value       = module.k8s_nodes.worker_private_ips
}

output "kubeadm_join_command" {
  description = "Kubeadm join command for worker nodes"
  value       = module.k8s_nodes.kubeadm_join_command
  sensitive   = true
}

output "kubeadm_join_token" {
  description = "Kubeadm join token"
  value       = module.k8s_nodes.kubeadm_join_token
  sensitive   = true
}

output "kubeadm_ca_cert_hash" {
  description = "Kubeadm CA cert hash"
  value       = module.k8s_nodes.kubeadm_ca_cert_hash
  sensitive   = true
}

output "security_group_id" {
  description = "ID of the Kubernetes security group"
  value       = module.k8s_nodes.security_group_id
}

output "ssh_key_name" {
  description = "SSH key pair name"
  value       = var.ssh_key_name
}

output "network_plugin" {
  description = "Network plugin used"
  value       = var.network_plugin
}

output "pod_network_cidr" {
  description = "Pod network CIDR"
  value       = var.pod_network_cidr
}

output "kubernetes_version" {
  description = "Kubernetes version"
  value       = var.kubernetes_version
}

output "cluster_endpoint" {
  description = "Kubernetes API server endpoint"
  value       = module.k8s_nodes.cluster_endpoint
}

output "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  value       = module.k8s_nodes.kubeconfig_path
}
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
