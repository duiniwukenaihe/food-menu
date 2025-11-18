variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "k8s-cluster"
}

variable "region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24"]
}

variable "master_instance_type" {
  description = "EC2 instance type for master nodes"
  type        = string
  default     = "t3.medium"
}

variable "worker_instance_type" {
  description = "EC2 instance type for worker nodes"
  type        = string
  default     = "t3.large"
}

variable "master_count" {
  description = "Number of master nodes"
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

variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
}

variable "ssh_key_name" {
  description = "SSH key pair name for EC2 instances"
  type        = string
  default     = "k8s-keypair"
}

variable "ssh_public_key" {
  description = "SSH public key content"
  type        = string
  default     = ""
}

variable "network_plugin" {
  description = "Network plugin for Kubernetes (flannel, cilium, calico)"
  type        = string
  default     = "flannel"
  
  validation {
    condition = contains(["flannel", "cilium", "calico"], var.network_plugin)
    error_message = "The network_plugin must be one of: flannel, cilium, calico."
  }
}

variable "pod_network_cidr" {
  description = "Pod network CIDR"
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

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.30.0"
}

variable "tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    "Project"     = "kubernetes-cluster"
    "Environment" = "development"
    "ManagedBy"   = "terraform"
  }
}

variable "enable_monitoring" {
  description = "Enable CloudWatch monitoring for instances"
  type        = bool
  default     = true
}

variable "enable_ebs_optimization" {
  description = "Enable EBS optimization for instances"
  type        = bool
  default     = true
}

variable "root_volume_size" {
  description = "Root volume size in GB"
  type        = number
  default     = 30
}

variable "root_volume_type" {
  description = "Root volume type"
  type        = string
  default     = "gp3"
}

variable "data_volume_size" {
  description = "Data volume size in GB for worker nodes"
  type        = number
  default     = 50
}

variable "data_volume_type" {
  description = "Data volume type for worker nodes"
  type        = string
  default     = "gp3"
}

variable "availability_zones" {
  description = "Availability zones for subnets"
  type        = list(string)
  default     = []
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
