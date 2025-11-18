# =============================================================================
# Terraform Kubernetes Cluster on Proxmox
# =============================================================================
# This configuration creates a complete Kubernetes cluster on Proxmox VE
# using the user's specific configuration parameters
# =============================================================================

terraform {
  required_version = ">= 1.0"
  
  # Backend configuration can be added here for state management
  # backend "local" {
  #   path = "./terraform.tfstate"
  # }
}

# =============================================================================
# Data Sources
# =============================================================================

# Get the current user's SSH public key
data "local_file" "ssh_public_key" {
  filename = var.ssh_public_key_path
}

# =============================================================================
# Modules
# =============================================================================

# Create VM template from Ubuntu cloud image
module "template" {
  source = "./modules/template"
  
  # Proxmox Configuration
  proxmox_api_url = var.proxmox_api_url
  proxmox_user    = var.proxmox_user
  proxmox_password = var.proxmox_password
  proxmox_node    = var.proxmox_node
  proxmox_ssh_key = var.proxmox_ssh_key
  
  # Template Configuration
  template_id   = var.template_id
  template_name = "${var.cluster_name}-template"
  image_url     = var.image_url
  image_name    = var.image_name
  storage       = var.storage
  
  # SSH Configuration
  ssh_public_key = data.local_file.ssh_public_key.content
  
  # User Configuration
  username = var.username
  password = var.password
}

# Create Kubernetes node pool (master and worker nodes)
module "kubernetes_nodes" {
  source = "./modules/kubernetes_node_pool"
  
  depends_on = [module.template]
  
  # Proxmox Configuration
  proxmox_api_url = var.proxmox_api_url
  proxmox_user    = var.proxmox_user
  proxmox_password = var.proxmox_password
  proxmox_node    = var.proxmox_node
  proxmox_ssh_key = var.proxmox_ssh_key
  
  # Template Reference
  template_id = module.template.template_id
  
  # Network Configuration
  bridge = var.bridge
  
  # SSH Configuration
  ssh_public_key = data.local_file.ssh_public_key.content
  
  # User Configuration
  username = var.username
  password = var.password
  
  # Master Node Configuration
  master_count     = var.master_count
  master_vmid_start = var.master_vmid_start
  master_cores     = var.master_cores
  master_memory    = var.master_memory
  master_disk_size = var.master_disk_size
  
  # Worker Node Configuration
  worker_count     = var.worker_count
  worker_vmid_start = var.worker_vmid_start
  worker_cores     = var.worker_cores
  worker_memory    = var.worker_memory
  worker_disk_size = var.worker_disk_size
  
  # Kubernetes Configuration
  k8s_version      = var.k8s_version
  pod_network_cidr = var.pod_network_cidr
  network_plugin   = var.network_plugin
  image_repository = var.image_repository
  
  # Cluster Configuration
  cluster_name = var.cluster_name
  environment  = var.environment
  
  # Tags
  tags = var.tags
}

# =============================================================================
# Local Files for Cloud-Init
# =============================================================================

# Generate cloud-init configuration for master nodes
resource "local_file" "master_cloud_init" {
  for_each = toset([for i in range(var.master_count) : "master-${i}"])
  
  content = templatefile("${path.module}/cloud-init/master.yaml.tpl", {
    hostname = "${var.cluster_name}-${each.key}"
    username = var.username
    password = var.password
    ssh_public_key = data.local_file.ssh_public_key.content
    
    # Kubernetes Configuration
    k8s_version = var.k8s_version
    pod_network_cidr = var.pod_network_cidr
    network_plugin = var.network_plugin
    image_repository = var.image_repository
    is_master = true
    cluster_name = var.cluster_name
  })
  
  filename = "${path.module}/cloud-init/generated/${each.key}-cloud-init.yaml"
}

# Generate cloud-init configuration for worker nodes
resource "local_file" "worker_cloud_init" {
  for_each = toset([for i in range(var.worker_count) : "worker-${i}"])
  
  content = templatefile("${path.module}/cloud-init/worker.yaml.tpl", {
    hostname = "${var.cluster_name}-${each.key}"
    username = var.username
    password = var.password
    ssh_public_key = data.local_file.ssh_public_key.content
    
    # Kubernetes Configuration
    k8s_version = var.k8s_version
    pod_network_cidr = var.pod_network_cidr
    network_plugin = var.network_plugin
    image_repository = var.image_repository
    is_master = false
    cluster_name = var.cluster_name
  })
  
  filename = "${path.module}/cloud-init/generated/${each.key}-cloud-init.yaml"
}

# =============================================================================
# Null Resources for Additional Setup
# =============================================================================

# Download and verify Ubuntu cloud image
resource "null_resource" "download_image" {
  
  provisioner "local-exec" {
    command = "${path.module}/scripts/get-ubuntu-cloudimg.sh ${var.image_url} ${var.image_name}"
    
    working_dir = path.module
  }
}

# =============================================================================
# Lifecycle Management
# =============================================================================

# Prevent accidental deletion of the template
resource "terraform_data" "template_protection" {
  depends_on = [module.template]
  
  lifecycle {
    prevent_destroy = true
  }
}