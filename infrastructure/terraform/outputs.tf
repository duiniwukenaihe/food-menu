# =============================================================================
# Kubernetes Cluster Outputs
# =============================================================================

output "cluster_name" {
  description = "Name of the Kubernetes cluster"
  value       = var.cluster_name
}

output "kubernetes_version" {
  description = "Kubernetes version installed"
  value       = var.k8s_version
}

output "cluster_endpoint" {
  description = "Kubernetes API server endpoint"
  value       = module.kubernetes_nodes.master_ips[0]
}

output "master_nodes" {
  description = "Information about master nodes"
  value = {
    count = var.master_count
    vms   = module.kubernetes_nodes.master_vms
    ips   = module.kubernetes_nodes.master_ips
    names = module.kubernetes_nodes.master_names
  }
}

output "worker_nodes" {
  description = "Information about worker nodes"
  value = {
    count = var.worker_count
    vms   = module.kubernetes_nodes.worker_vms
    ips   = module.kubernetes_nodes.worker_ips
    names = module.kubernetes_nodes.worker_names
  }
}

output "network_plugin" {
  description = "Network plugin used"
  value       = var.network_plugin
}

output "pod_network_cidr" {
  description = "Pod network CIDR"
  value       = var.pod_network_cidr
}

output "ssh_access" {
  description = "SSH access information"
  value = {
    username    = var.username
    private_key = var.ssh_private_key_path
    master_ips  = module.kubernetes_nodes.master_ips
    worker_ips  = module.kubernetes_nodes.worker_ips
  }
}

output "template_info" {
  description = "VM template information"
  value = {
    template_id = module.template.template_id
    template_name = module.template.template_name
    image_url = var.image_url
  }
}

output "proxmox_info" {
  description = "Proxmox configuration information"
  value = {
    api_url = var.proxmox_api_url
    node    = var.proxmox_node
    storage = var.storage
    bridge  = var.bridge
  }
}

# =============================================================================
# Connection Commands
# =============================================================================

output "connect_commands" {
  description = "SSH connection commands for the cluster nodes"
  value = {
    master_ssh = [
      for i, ip in module.kubernetes_nodes.master_ips :
      "ssh -i ${var.ssh_private_key_path} ${var.username}@${ip}"
    ]
    worker_ssh = [
      for i, ip in module.kubernetes_nodes.worker_ips :
      "ssh -i ${var.ssh_private_key_path} ${var.username}@${ip}"
    ]
  }
}

output "kubectl_commands" {
  description = "Useful kubectl commands"
  value = {
    get_nodes = "kubectl get nodes -o wide"
    get_pods = "kubectl get pods --all-namespaces"
    cluster_info = "kubectl cluster-info"
    join_command = "sudo kubeadm token create --print-join-command"
  }
}