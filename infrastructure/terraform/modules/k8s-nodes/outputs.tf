# TLS Private Key for SSH (only if creating VPC and no key provided)
resource "tls_private_key" "k8s" {
  count = var.create_vpc && var.ssh_public_key == "" ? 1 : 0
  
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = var.create_vpc ? aws_vpc.main[0].id : var.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = var.create_vpc ? aws_subnet.public[*].id : var.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = var.create_vpc ? aws_subnet.private[*].id : var.private_subnet_ids
}

output "master_instance_ids" {
  description = "IDs of the master instances"
  value       = aws_instance.master[*].id
}

output "worker_instance_ids" {
  description = "IDs of the worker instances"
  value       = aws_instance.worker[*].id
}

output "master_public_ips" {
  description = "Public IPs of the master instances"
  value       = aws_instance.master[*].public_ip
}

output "worker_public_ips" {
  description = "Public IPs of the worker instances"
  value       = aws_instance.worker[*].public_ip
}

output "master_private_ips" {
  description = "Private IPs of the master instances"
  value       = aws_instance.master[*].private_ip
}

output "worker_private_ips" {
  description = "Private IPs of the worker instances"
  value       = aws_instance.worker[*].private_ip
}

output "security_group_id" {
  description = "ID of the Kubernetes security group"
  value       = aws_security_group.k8s.id
}

output "cluster_endpoint" {
  description = "Kubernetes API server endpoint"
  value       = var.master_count > 0 ? "${aws_instance.master[0].public_ip}:6443" : ""
}

output "kubeadm_join_command" {
  description = "Kubeadm join command for worker nodes"
  value       = var.master_count > 0 ? "sudo kubeadm join ${aws_instance.master[0].private_ip}:6443 --token ${local.kubeadm_join_token} --discovery-token-ca-cert-hash ${local.kubeadm_ca_cert_hash}" : ""
  sensitive   = true
}

output "kubeadm_join_token" {
  description = "Kubeadm join token"
  value       = local.kubeadm_join_token
  sensitive   = true
}

output "kubeadm_ca_cert_hash" {
  description = "Kubeadm CA cert hash"
  value       = local.kubeadm_ca_cert_hash
  sensitive   = true
}

output "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  value       = var.master_count > 0 ? "/home/ubuntu/.kube/config" : ""
}

output "ssh_private_key" {
  description = "SSH private key (only if generated)"
  value       = var.create_vpc && var.ssh_public_key == "" ? tls_private_key.k8s[0].private_key : ""
  sensitive   = true
}