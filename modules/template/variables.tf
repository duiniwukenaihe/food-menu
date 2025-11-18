variable "proxmox_node" {
  description = "Proxmox node where the template will be created"
  type        = string
}

variable "vm_id" {
  description = "VM ID for the template"
  type        = number
  default     = 9001
}

variable "storage_pool" {
  description = "Storage pool for the template disk"
  type        = string
}

variable "network_bridge" {
  description = "Network bridge for the template"
  type        = string
  default     = "vmbr0"
}

variable "cores" {
  description = "Number of CPU cores"
  type        = number
  default     = 2
}

variable "memory" {
  description = "Memory in MB"
  type        = number
  default     = 2048
}

variable "disk_size" {
  description = "Disk size"
  type        = string
  default     = "20G"
}

variable "cloud_init_storage" {
  description = "Storage pool for cloud-init drive (defaults to main storage)"
  type        = string
  default     = null
}

variable "ubuntu_version" {
  description = "Ubuntu version"
  type        = string
  default     = "22.04"
}

variable "ubuntu_architecture" {
  description = "Ubuntu architecture"
  type        = string
  default     = "amd64"
}

variable "interactive_replace" {
  description = "Whether to prompt for replacement when image exists"
  type        = bool
  default     = true
}