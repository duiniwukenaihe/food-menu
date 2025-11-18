# =============================================================================
# Proxmox Connection Variables
# =============================================================================

variable "proxmox_endpoint" {
  description = "Proxmox VE API endpoint URL"
  type        = string
  default     = "https://192.168.0.200:8006"
}

variable "proxmox_username" {
  description = "Proxmox VE username (e.g., root@pam)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "proxmox_password" {
  description = "Proxmox VE password (can also be set via PROXMOX_VE_PASSWORD environment variable)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "proxmox_api_token" {
  description = "Proxmox VE API token (alternative to username/password, format: user@realm!tokenid=secret)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "proxmox_insecure" {
  description = "Skip TLS verification for self-signed certificates"
  type        = bool
  default     = true
}

variable "proxmox_ssh_username" {
  description = "SSH username for Proxmox host operations"
  type        = string
  default     = "root"
}

variable "proxmox_ssh_private_key" {
  description = "Path to SSH private key file for Proxmox host access"
  type        = string
  default     = ""
}

# =============================================================================
# Proxmox Resource Variables
# =============================================================================

variable "proxmox_node_name" {
  description = "Default Proxmox node name for VM placement"
  type        = string
  default     = "pve"
}

variable "proxmox_storage" {
  description = "Proxmox storage pool name for VM disks"
  type        = string
  default     = "local-lvm"
}

variable "proxmox_iso_storage" {
  description = "Proxmox storage for ISO images"
  type        = string
  default     = "local"
}

variable "proxmox_network_bridge" {
  description = "Network bridge for VM network interfaces"
  type        = string
  default     = "vmbr0"
}

variable "proxmox_vlan_tag" {
  description = "VLAN tag for VM network (leave empty for no VLAN)"
  type        = number
  default     = null
}

# =============================================================================
# VM Template Variables
# =============================================================================

variable "template_name" {
  description = "Name for the VM template"
  type        = string
  default     = "k8s-template"
}

variable "template_vm_id" {
  description = "VM ID for the template (must be unique)"
  type        = number
  default     = 9000
}

variable "template_image_url" {
  description = "URL to download the cloud-init image (e.g., Ubuntu cloud image)"
  type        = string
  default     = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}

variable "template_image_checksum" {
  description = "Checksum for the cloud image (optional)"
  type        = string
  default     = ""
}

# =============================================================================
# SSH Key Variables
# =============================================================================

variable "ssh_public_key" {
  description = "SSH public key for VM access"
  type        = string
  default     = ""
}

variable "ssh_public_key_file" {
  description = "Path to SSH public key file (alternative to ssh_public_key variable)"
  type        = string
  default     = ""
}

variable "ssh_authorized_keys" {
  description = "List of SSH public keys for VM access"
  type        = list(string)
  default     = []
}

# =============================================================================
# Default VM Sizing
# =============================================================================

variable "default_cpu_cores" {
  description = "Default number of CPU cores for VMs"
  type        = number
  default     = 2
}

variable "default_cpu_sockets" {
  description = "Default number of CPU sockets for VMs"
  type        = number
  default     = 1
}

variable "default_memory_mb" {
  description = "Default memory size in MB for VMs"
  type        = number
  default     = 4096
}

variable "default_disk_size_gb" {
  description = "Default disk size in GB for VMs"
  type        = number
  default     = 32
}

# =============================================================================
# Control Plane Configuration
# =============================================================================

variable "control_plane_count" {
  description = "Number of control plane nodes"
  type        = number
  default     = 3
}

variable "control_plane_cpu_cores" {
  description = "CPU cores for control plane nodes"
  type        = number
  default     = 2
}

variable "control_plane_memory_mb" {
  description = "Memory in MB for control plane nodes"
  type        = number
  default     = 4096
}

variable "control_plane_disk_size_gb" {
  description = "Disk size in GB for control plane nodes"
  type        = number
  default     = 50
}

variable "control_plane_nodes" {
  description = "Map of control plane node configurations"
  type = map(object({
    vm_id        = number
    name         = string
    node_name    = optional(string)
    cpu_cores    = optional(number)
    memory_mb    = optional(number)
    disk_size_gb = optional(number)
    ip_address   = optional(string)
  }))
  default = {}
}

# =============================================================================
# Worker Node Configuration
# =============================================================================

variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 3
}

variable "worker_cpu_cores" {
  description = "CPU cores for worker nodes"
  type        = number
  default     = 4
}

variable "worker_memory_mb" {
  description = "Memory in MB for worker nodes"
  type        = number
  default     = 8192
}

variable "worker_disk_size_gb" {
  description = "Disk size in GB for worker nodes"
  type        = number
  default     = 100
}

variable "worker_nodes" {
  description = "Map of worker node configurations"
  type = map(object({
    vm_id        = number
    name         = string
    node_name    = optional(string)
    cpu_cores    = optional(number)
    memory_mb    = optional(number)
    disk_size_gb = optional(number)
    ip_address   = optional(string)
  }))
  default = {}
}

# =============================================================================
# Network Configuration
# =============================================================================

variable "network_gateway" {
  description = "Network gateway for VMs"
  type        = string
  default     = "192.168.0.1"
}

variable "network_dns_servers" {
  description = "DNS servers for VMs"
  type        = list(string)
  default     = ["8.8.8.8", "8.8.4.4"]
}

variable "network_domain" {
  description = "Network domain name"
  type        = string
  default     = "local"
}

# =============================================================================
# Kubernetes Configuration
# =============================================================================

variable "kubernetes_version" {
  description = "Kubernetes version to install"
  type        = string
  default     = "1.28.0"
}

variable "kubernetes_pod_network_cidr" {
  description = "Pod network CIDR for Kubernetes"
  type        = string
  default     = "10.244.0.0/16"
}

variable "kubernetes_service_cidr" {
  description = "Service network CIDR for Kubernetes"
  type        = string
  default     = "10.96.0.0/12"
}

variable "kubernetes_cni" {
  description = "CNI plugin to use (calico, flannel, cilium, etc.)"
  type        = string
  default     = "calico"
}

variable "kubernetes_enable_monitoring" {
  description = "Enable monitoring stack (Prometheus, Grafana)"
  type        = bool
  default     = false
}

variable "kubernetes_enable_logging" {
  description = "Enable logging stack (EFK/ELK)"
  type        = bool
  default     = false
}

# =============================================================================
# Node Maps
# =============================================================================

variable "node_map" {
  description = "Mapping of logical node names to Proxmox cluster nodes for VM placement"
  type        = map(string)
  default = {
    "default" = "pve"
  }
}

# =============================================================================
# Tags and Metadata
# =============================================================================

variable "tags" {
  description = "Tags to apply to all resources"
  type        = list(string)
  default     = ["kubernetes", "terraform"]
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "k8s-cluster"
}
