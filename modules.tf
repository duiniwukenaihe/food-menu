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

module "control_plane" {
  source = "./modules/kubernetes_node_pool"
  count  = length(var.control_plane_nodes) > 0 ? 1 : 0
  
  # Proxmox connection
  proxmox_node  = var.proxmox_node
  template_vmid = module.ubuntu_template.template_vmid
  storage_pool  = var.ubuntu_template_config.storage_pool
  
  # Network configuration
  network_bridge     = var.ubuntu_template_config.network_bridge
  cloud_init_storage = var.ubuntu_template_config.cloud_init_storage
  
  # Node configuration
  node_type         = "control-plane"
  nodes             = var.control_plane_nodes
  default_cores     = var.control_plane_defaults.cores
  default_memory    = var.control_plane_defaults.memory
  default_disk_size = var.control_plane_defaults.disk_size
  
  # SSH configuration
  ssh_public_keys = var.ssh_public_keys
  ssh_password    = var.ssh_password
  
  # Tags
  tags = ["kubernetes", "control-plane"]
  
  depends_on = [module.ubuntu_template]
}

module "workers" {
  source = "./modules/kubernetes_node_pool"
  count  = length(var.worker_nodes) > 0 ? 1 : 0
  
  # Proxmox connection
  proxmox_node  = var.proxmox_node
  template_vmid = module.ubuntu_template.template_vmid
  storage_pool  = var.ubuntu_template_config.storage_pool
  
  # Network configuration
  network_bridge     = var.ubuntu_template_config.network_bridge
  cloud_init_storage = var.ubuntu_template_config.cloud_init_storage
  
  # Node configuration
  node_type         = "worker"
  nodes             = var.worker_nodes
  default_cores     = var.worker_defaults.cores
  default_memory    = var.worker_defaults.memory
  default_disk_size = var.worker_defaults.disk_size
  
  # SSH configuration
  ssh_public_keys = var.ssh_public_keys
  ssh_password    = var.ssh_password
  
  # Tags
  tags = ["kubernetes", "worker"]
  
  depends_on = [module.ubuntu_template]
}