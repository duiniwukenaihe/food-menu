# =============================================================================
# Kubernetes on Proxmox - Main Configuration
# =============================================================================

# -----------------------------------------------------------------------------
# VM Template Module
# -----------------------------------------------------------------------------
# This module creates a VM template from a cloud-init image that will be used
# as the base for all Kubernetes nodes (control plane and workers)
# 
# Uncomment and configure once the template module is created

# module "template" {
#   source = "./modules/template"
#
#   vm_id                  = var.template_vm_id
#   name                   = var.template_name
#   node_name              = var.proxmox_node_name
#   storage                = var.proxmox_storage
#   iso_storage            = var.proxmox_iso_storage
#   image_url              = var.template_image_url
#   image_checksum         = var.template_image_checksum
#   ssh_public_key         = var.ssh_public_key != "" ? var.ssh_public_key : (var.ssh_public_key_file != "" ? file(var.ssh_public_key_file) : "")
#   ssh_authorized_keys    = var.ssh_authorized_keys
#   tags                   = concat(var.tags, ["template"])
# }

# -----------------------------------------------------------------------------
# Control Plane Module
# -----------------------------------------------------------------------------
# This module creates the Kubernetes control plane nodes (masters)
# Control plane nodes run the Kubernetes API server, scheduler, and controller
# manager components
#
# Uncomment and configure once the control-plane module is created

# module "control_plane" {
#   source = "./modules/control-plane"
#
#   depends_on = [module.template]
#
#   template_id            = module.template.template_id
#   template_name          = module.template.template_name
#   
#   node_count             = var.control_plane_count
#   nodes                  = var.control_plane_nodes
#   
#   default_node_name      = var.proxmox_node_name
#   default_cpu_cores      = var.control_plane_cpu_cores
#   default_memory_mb      = var.control_plane_memory_mb
#   default_disk_size_gb   = var.control_plane_disk_size_gb
#   
#   storage                = var.proxmox_storage
#   network_bridge         = var.proxmox_network_bridge
#   vlan_tag               = var.proxmox_vlan_tag
#   
#   network_gateway        = var.network_gateway
#   network_dns_servers    = var.network_dns_servers
#   network_domain         = var.network_domain
#   
#   kubernetes_version     = var.kubernetes_version
#   pod_network_cidr       = var.kubernetes_pod_network_cidr
#   service_cidr           = var.kubernetes_service_cidr
#   cni_plugin             = var.kubernetes_cni
#   
#   tags                   = concat(var.tags, ["control-plane", "master"])
#   environment            = var.environment
#   project_name           = var.project_name
# }

# -----------------------------------------------------------------------------
# Worker Pool Module
# -----------------------------------------------------------------------------
# This module creates the Kubernetes worker nodes
# Worker nodes run the actual application workloads
#
# Uncomment and configure once the worker-pool module is created

# module "worker_pool" {
#   source = "./modules/worker-pool"
#
#   depends_on = [module.template, module.control_plane]
#
#   template_id            = module.template.template_id
#   template_name          = module.template.template_name
#   
#   node_count             = var.worker_count
#   nodes                  = var.worker_nodes
#   
#   default_node_name      = var.proxmox_node_name
#   default_cpu_cores      = var.worker_cpu_cores
#   default_memory_mb      = var.worker_memory_mb
#   default_disk_size_gb   = var.worker_disk_size_gb
#   
#   storage                = var.proxmox_storage
#   network_bridge         = var.proxmox_network_bridge
#   vlan_tag               = var.proxmox_vlan_tag
#   
#   network_gateway        = var.network_gateway
#   network_dns_servers    = var.network_dns_servers
#   network_domain         = var.network_domain
#   
#   control_plane_endpoint = module.control_plane.cluster_endpoint
#   join_token             = module.control_plane.join_token
#   
#   tags                   = concat(var.tags, ["worker"])
#   environment            = var.environment
#   project_name           = var.project_name
# }

# -----------------------------------------------------------------------------
# Local Variables
# -----------------------------------------------------------------------------

locals {
  # Common tags for all resources
  common_tags = concat(
    var.tags,
    [
      "environment:${var.environment}",
      "project:${var.project_name}",
      "managed-by:terraform"
    ]
  )

  # Determine SSH public key to use
  # Use try() to gracefully handle missing files
  ssh_key = var.ssh_public_key != "" ? var.ssh_public_key : (
    var.ssh_public_key_file != "" ? try(file(var.ssh_public_key_file), "") : ""
  )

  # All SSH keys combined
  all_ssh_keys = concat(
    local.ssh_key != "" ? [local.ssh_key] : [],
    var.ssh_authorized_keys
  )
}
