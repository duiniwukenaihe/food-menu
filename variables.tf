variable "proxmox_node" {
  description = "Proxmox node where resources will be created"
  type        = string
}

variable "proxmox_api_url" {
  description = "Proxmox API URL"
  type        = string
}

variable "proxmox_api_token_id" {
  description = "Proxmox API token ID"
  type        = string
  sensitive   = true
}

variable "proxmox_api_token_secret" {
  description = "Proxmox API token secret"
  type        = string
  sensitive   = true
}

variable "proxmox_tls_insecure" {
  description = "Whether to skip TLS verification for Proxmox API"
  type        = bool
  default     = false
}

variable "ubuntu_template_config" {
  description = "Ubuntu template configuration"
  type = object({
    vm_id               = optional(number, 9001)
    storage_pool        = string
    network_bridge      = optional(string, "vmbr0")
    cores               = optional(number, 2)
    memory              = optional(number, 2048)
    disk_size           = optional(string, "20G")
    cloud_init_storage  = optional(string, null)
  })
  default = {
    storage_pool = "local"
  }
}

variable "ubuntu_image_config" {
  description = "Ubuntu image download configuration"
  type = object({
    version             = optional(string, "22.04")
    architecture        = optional(string, "amd64")
    interactive_replace = optional(bool, true)
  })
  default = {}
}

variable "control_plane_nodes" {
  description = "Map of control plane nodes to create"
  type = map(object({
    vmid         = number
    hostname     = string
    ip           = string
    gateway      = string
    cores        = optional(number)
    memory       = optional(number)
    disk_size    = optional(string)
    nameserver   = optional(string, "8.8.8.8")
    searchdomain = optional(string, "local")
  }))
  default = {}
}

variable "worker_nodes" {
  description = "Map of worker nodes to create"
  type = map(object({
    vmid         = number
    hostname     = string
    ip           = string
    gateway      = string
    cores        = optional(number)
    memory       = optional(number)
    disk_size    = optional(string)
    nameserver   = optional(string, "8.8.8.8")
    searchdomain = optional(string, "local")
  }))
  default = {}
}

variable "control_plane_defaults" {
  description = "Default resource allocation for control plane nodes"
  type = object({
    cores     = optional(number, 4)
    memory    = optional(number, 8192)
    disk_size = optional(string, "50G")
  })
  default = {
    cores     = 4
    memory    = 8192
    disk_size = "50G"
  }
}

variable "worker_defaults" {
  description = "Default resource allocation for worker nodes"
  type = object({
    cores     = optional(number, 8)
    memory    = optional(number, 16384)
    disk_size = optional(string, "100G")
  })
  default = {
    cores     = 8
    memory    = 16384
    disk_size = "100G"
  }
}

variable "ssh_public_keys" {
  description = "List of SSH public keys to add to all Kubernetes nodes"
  type        = list(string)
  default     = []
}

variable "ssh_password" {
  description = "SSH password for ubuntu user on all Kubernetes nodes"
  type        = string
  sensitive   = true
  default     = null
}