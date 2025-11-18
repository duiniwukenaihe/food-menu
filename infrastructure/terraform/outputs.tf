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