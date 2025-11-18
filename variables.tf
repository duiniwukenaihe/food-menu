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