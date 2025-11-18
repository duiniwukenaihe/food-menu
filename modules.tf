module "ubuntu_template" {
  source = "./modules/template"
  
  # Proxmox connection settings
  proxmox_node = var.proxmox_node
  
  # Template configuration
  vm_id               = var.ubuntu_template_config.vm_id
  storage_pool        = var.ubuntu_template_config.storage_pool
  network_bridge      = var.ubuntu_template_config.network_bridge
  cores               = var.ubuntu_template_config.cores
  memory              = var.ubuntu_template_config.memory
  disk_size           = var.ubuntu_template_config.disk_size
  cloud_init_storage  = var.ubuntu_template_config.cloud_init_storage
  
  # Ubuntu image configuration
  ubuntu_version      = var.ubuntu_image_config.version
  ubuntu_architecture = var.ubuntu_image_config.architecture
  interactive_replace = var.ubuntu_image_config.interactive_replace
}