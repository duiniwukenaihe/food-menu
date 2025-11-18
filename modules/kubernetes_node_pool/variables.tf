variable "proxmox_node" {
  description = "Proxmox node where VMs will be created"
  type        = string
}

variable "template_vmid" {
  description = "VM ID of the template to clone from"
  type        = number
}

variable "nodes" {
  description = "Map of nodes to create with their configurations"
  type = map(object({
    vmid      = number
    hostname  = string
    ip        = string
    gateway   = string
    cores     = optional(number)
    memory    = optional(number)
    disk_size = optional(string)
    nameserver = optional(string, "8.8.8.8")
    searchdomain = optional(string, "local")
  }))
}

variable "default_cores" {
  description = "Default number of CPU cores if not specified per node"
  type        = number
  default     = 2
}

variable "default_memory" {
  description = "Default memory in MB if not specified per node"
  type        = number
  default     = 2048
}

variable "default_disk_size" {
  description = "Default disk size if not specified per node"
  type        = string
  default     = "32G"
}

variable "storage_pool" {
  description = "Storage pool for VM disks"
  type        = string
}

variable "network_bridge" {
  description = "Network bridge for VMs"
  type        = string
  default     = "vmbr0"
}

variable "ssh_public_keys" {
  description = "List of SSH public keys to add to VMs"
  type        = list(string)
  default     = []
}

variable "ssh_password" {
  description = "SSH password for default user (ubuntu)"
  type        = string
  sensitive   = true
  default     = null
}

variable "cloud_init_storage" {
  description = "Storage pool for cloud-init drive (defaults to main storage)"
  type        = string
  default     = null
}

variable "tags" {
  description = "Tags to apply to VMs"
  type        = list(string)
  default     = ["kubernetes"]
}

variable "node_type" {
  description = "Type of node (control-plane or worker) for tagging"
  type        = string
  default     = "node"
}

variable "start_on_boot" {
  description = "Whether to start VMs automatically on boot"
  type        = bool
  default     = true
}
