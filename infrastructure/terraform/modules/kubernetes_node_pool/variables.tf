# =============================================================================
# Kubernetes Node Pool Module Variables
# =============================================================================

# Proxmox Configuration
variable "proxmox_api_url" {
  description = "Proxmox API URL"
  type        = string
}

variable "proxmox_user" {
  description = "Proxmox username"
  type        = string
}

variable "proxmox_password" {
  description = "Proxmox password"
  type        = string
  sensitive   = true
}

variable "proxmox_node" {
  description = "Proxmox node name"
  type        = string
}

variable "proxmox_ssh_key" {
  description = "Path to SSH private key for Proxmox host access"
  type        = string
}

# Template Configuration
variable "template_id" {
  description = "VM template ID to clone from"
  type        = number
}

# Network Configuration
variable "bridge" {
  description = "Network bridge for VM network interfaces"
  type        = string
}

variable "storage" {
  description = "Proxmox storage pool name for VM disks"
  type        = string
  default     = "local"
}

# User Configuration
variable "username" {
  description = "Default username for VM instances"
  type        = string
}

variable "password" {
  description = "Default password for VM instances"
  type        = string
  sensitive   = true
}

variable "ssh_public_key" {
  description = "SSH public key content"
  type        = string
}

# Master Node Configuration
variable "master_count" {
  description = "Number of master nodes"
  type        = number
}

variable "master_vmid_start" {
  description = "Starting VM ID for master nodes"
  type        = number
}

variable "master_cores" {
  description = "CPU cores for master nodes"
  type        = number
}

variable "master_memory" {
  description = "Memory in MB for master nodes"
  type        = number
}

variable "master_disk_size" {
  description = "Disk size for master nodes"
  type        = string
}

# Worker Node Configuration
variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
}

variable "worker_vmid_start" {
  description = "Starting VM ID for worker nodes"
  type        = number
}

variable "worker_cores" {
  description = "CPU cores for worker nodes"
  type        = number
}

variable "worker_memory" {
  description = "Memory in MB for worker nodes"
  type        = number
}

variable "worker_disk_size" {
  description = "Disk size for worker nodes"
  type        = string
}

# Kubernetes Configuration
variable "k8s_version" {
  description = "Kubernetes version to install"
  type        = string
}

variable "pod_network_cidr" {
  description = "Pod network CIDR for Kubernetes"
  type        = string
}

variable "network_plugin" {
  description = "Network plugin for Kubernetes"
  type        = string
}

variable "image_repository" {
  description = "Container image repository for Kubernetes"
  type        = string
}

# Cluster Configuration
variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

# Tags
variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}