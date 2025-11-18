# =============================================================================
# Proxmox Configuration Variables
# =============================================================================

variable "proxmox_api_url" {
  description = "Proxmox API URL"
  type        = string
  default     = "https://192.168.0.200:8006/api2/json"
  
  validation {
    condition     = can(regex("^https://", var.proxmox_api_url))
    error_message = "Proxmox API URL must start with https://."
  }
}

variable "proxmox_user" {
  description = "Proxmox username (e.g., root@pam)"
  type        = string
  default     = "root@pam"
  
  validation {
    condition     = can(regex("@", var.proxmox_user))
    error_message = "Proxmox username must be in format user@realm."
  }
}

variable "proxmox_password" {
  description = "Proxmox password (can also be set via PROXMOX_VE_PASSWORD environment variable)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "proxmox_node" {
  description = "Proxmox node name for VM placement"
  type        = string
  default     = "proxmox1"
}

variable "proxmox_host" {
  description = "Proxmox host IP address"
  type        = string
  default     = "192.168.0.200"
  
  validation {
    condition     = can(regex("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$", var.proxmox_host))
    error_message = "Proxmox host must be a valid IP address."
  }
}

variable "proxmox_ssh_user" {
  description = "SSH username for Proxmox host operations"
  type        = string
  default     = "root"
}

variable "proxmox_ssh_key" {
  description = "Path to SSH private key for Proxmox host access"
  type        = string
  default     = "~/.ssh/id_rsa"
}

variable "proxmox_api_token_id" {
  description = "Proxmox API token ID (alternative to username/password)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "proxmox_api_token_secret" {
  description = "Proxmox API token secret (alternative to username/password)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "proxmox_insecure" {
  description = "Skip TLS verification for self-signed certificates"
  type        = bool
  default     = true
}

# =============================================================================
# Image and Template Configuration Variables
# =============================================================================

variable "image_url" {
  description = "Ubuntu cloud image URL"
  type        = string
  default     = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
  
  validation {
    condition     = can(regex("^https://", var.image_url))
    error_message = "Image URL must start with https://."
  }
}

variable "image_base_url" {
  description = "Ubuntu cloud image base URL"
  type        = string
  default     = "https://cloud-images.ubuntu.com/jammy/current"
  
  validation {
    condition     = can(regex("^https://", var.image_base_url))
    error_message = "Image base URL must start with https://."
  }
}

variable "image_name" {
  description = "Ubuntu cloud image filename"
  type        = string
  default     = "jammy-server-cloudimg-amd64.img"
}

variable "template_id" {
  description = "VM template ID"
  type        = number
  default     = 9001
  
  validation {
    condition     = var.template_id >= 100 && var.template_id <= 9999
    error_message = "Template ID must be between 100 and 9999."
  }
}

variable "storage" {
  description = "Proxmox storage pool name for VM disks"
  type        = string
  default     = "local"
}

variable "bridge" {
  description = "Network bridge for VM network interfaces"
  type        = string
  default     = "vmbr0"
}

# =============================================================================
# User Configuration Variables
# =============================================================================

variable "username" {
  description = "Default username for VM instances"
  type        = string
  default     = "ubuntu"
}

variable "password" {
  description = "Default password for VM instances"
  type        = string
  default     = ""
  sensitive   = true
}

variable "ssh_private_key_path" {
  description = "Path to SSH private key for VM access"
  type        = string
  default     = "~/.ssh/id_rsa"
}

variable "ssh_public_key_path" {
  description = "Path to SSH public key for VM access"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

# =============================================================================
# Master Node Configuration Variables
# =============================================================================

variable "master_count" {
  description = "Number of master nodes"
  type        = number
  default     = 1
  
  validation {
    condition     = var.master_count >= 1 && var.master_count <= 3
    error_message = "Master count must be between 1 and 3."
  }
}

variable "master_vmid_start" {
  description = "Starting VM ID for master nodes"
  type        = number
  default     = 501
  
  validation {
    condition     = var.master_vmid_start >= 100 && var.master_vmid_start <= 499
    error_message = "Master VM ID start must be between 100 and 499."
  }
}

variable "master_cores" {
  description = "CPU cores for master nodes"
  type        = number
  default     = 4
  
  validation {
    condition     = var.master_cores >= 2 && var.master_cores <= 16
    error_message = "Master cores must be between 2 and 16."
  }
}

variable "master_memory" {
  description = "Memory in MB for master nodes"
  type        = number
  default     = 8192
  
  validation {
    condition     = var.master_memory >= 2048 && var.master_memory <= 32768
    error_message = "Master memory must be between 2048 and 32768 MB."
  }
}

variable "master_disk_size" {
  description = "Disk size for master nodes"
  type        = string
  default     = "50G"
  
  validation {
    condition     = can(regex("^[0-9]+[GT]?$", var.master_disk_size))
    error_message = "Master disk size must be in format like '50G' or '100'."
  }
}

# =============================================================================
# Worker Node Configuration Variables
# =============================================================================

variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
  
  validation {
    condition     = var.worker_count >= 0 && var.worker_count <= 20
    error_message = "Worker count must be between 0 and 20."
  }
}

variable "worker_vmid_start" {
  description = "Starting VM ID for worker nodes"
  type        = number
  default     = 601
  
  validation {
    condition     = var.worker_vmid_start >= 500 && var.worker_vmid_start <= 999
    error_message = "Worker VM ID start must be between 500 and 999."
  }
}

variable "worker_cores" {
  description = "CPU cores for worker nodes"
  type        = number
  default     = 8
  
  validation {
    condition     = var.worker_cores >= 2 && var.worker_cores <= 32
    error_message = "Worker cores must be between 2 and 32."
  }
}

variable "worker_memory" {
  description = "Memory in MB for worker nodes"
  type        = number
  default     = 16384
  
  validation {
    condition     = var.worker_memory >= 2048 && var.worker_memory <= 65536
    error_message = "Worker memory must be between 2048 and 65536 MB."
  }
}

variable "worker_disk_size" {
  description = "Disk size for worker nodes"
  type        = string
  default     = "100G"
  
  validation {
    condition     = can(regex("^[0-9]+[GT]?$", var.worker_disk_size))
    error_message = "Worker disk size must be in format like '100G' or '200'."
  }
}

# =============================================================================
# Kubernetes Configuration Variables
# =============================================================================

variable "k8s_version" {
  description = "Kubernetes version to install"
  type        = string
  default     = "1.30.0"
  
  validation {
    condition     = can(regex("^1\\.[0-9]+\\.[0-9]+$", var.k8s_version))
    error_message = "Kubernetes version must be in format like '1.30.0'."
  }
}

variable "pod_network_cidr" {
  description = "Pod network CIDR for Kubernetes"
  type        = string
  default     = "10.244.0.0/16"
  
  validation {
    condition     = can(regex("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+/[0-9]+$", var.pod_network_cidr))
    error_message = "Pod network CIDR must be in valid CIDR format like '10.244.0.0/16'."
  }
}

variable "network_plugin" {
  description = "Network plugin for Kubernetes (flannel, cilium, calico)"
  type        = string
  default     = "flannel"
  
  validation {
    condition     = contains(["flannel", "cilium", "calico"], var.network_plugin)
    error_message = "Network plugin must be one of: flannel, cilium, calico."
  }
}

variable "image_repository" {
  description = "Container image repository for Kubernetes"
  type        = string
  default     = "registry.aliyuncs.com/google_containers"
  
  validation {
    condition     = can(regex("^[a-zA-Z0-9.-]+/[a-zA-Z0-9._-]+$", var.image_repository))
    error_message = "Image repository must be in valid format like 'registry.example.com/repo'."
  }
}

# =============================================================================
# Additional Configuration Variables
# =============================================================================

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "k8s-cluster"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
  
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default = {
    "Project"     = "kubernetes-cluster"
    "Environment" = "development"
    "ManagedBy"   = "terraform"
  }
}