# =============================================================================
# Template Module Variables
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
  description = "VM template ID"
  type        = number
}

variable "template_name" {
  description = "VM template name"
  type        = string
}

variable "image_url" {
  description = "Ubuntu cloud image URL"
  type        = string
}

variable "image_name" {
  description = "Ubuntu cloud image filename"
  type        = string
}

variable "image_checksum" {
  description = "SHA256 checksum for the image (optional)"
  type        = string
  default     = ""
}

variable "storage" {
  description = "Proxmox storage pool name"
  type        = string
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

# Cluster Configuration
variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "k8s-cluster"
}