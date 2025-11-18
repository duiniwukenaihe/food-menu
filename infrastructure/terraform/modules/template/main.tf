# =============================================================================
# Template Module - Creates VM Template from Ubuntu Cloud Image
# =============================================================================

terraform {
  required_version = ">= 1.0"
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.86.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

# =============================================================================
# Data Sources
# =============================================================================

# Get the current template info if it exists
data "proxmox_virtual_environment_vm" "existing_template" {
  vm_id = var.template_id
  count  = var.template_id != null ? 1 : 0
}

# =============================================================================
# Local Files
# =============================================================================

# Generate cloud-init configuration for template
resource "local_file" "template_cloud_init" {
  content = templatefile("${path.module}/templates/cloud-init.yaml.tpl", {
    hostname = var.template_name
    username = var.username
    password = var.password
    ssh_public_key = var.ssh_public_key
  })
  
  filename = "${path.module}/templates/generated/template-cloud-init.yaml"
}

# =============================================================================
# Template Creation
# =============================================================================

# Download Ubuntu cloud image
resource "proxmox_virtual_environment_download_file" "ubuntu_image" {
  content_type = "iso"
  datastore_id = var.storage
  node_name    = var.proxmox_node
  url          = var.image_url
  
  # Verify checksum if provided
  checksum     = var.image_checksum != "" ? var.image_checksum : null
  checksum_algorithm = var.image_checksum != "" ? "sha256" : null
  
  # Overwrite if exists
  overwrite    = true
  
  timeouts {
    create = "30m"
  }
}

# Create VM template from downloaded image
resource "proxmox_virtual_environment_vm" "template" {
  name        = var.template_name
  description = "Kubernetes node template for ${var.cluster_name}"
  node_name   = var.proxmox_node
  vm_id       = var.template_id
  
  # Mark as template
  template    = true
  
  # Basic VM configuration
  cpu {
    cores = 2
    type  = "host"
  }
  
  memory {
    dedicated = 2048
  }
  
  # Network configuration
  network_device {
    bridge = "vmbr0"
    model  = "virtio"
  }
  
  # Disk configuration
  disk {
    datastore_id = var.storage
    file_id      = proxmox_virtual_environment_download_file.ubuntu_image.id
    interface    = "virtio0"
    size         = 20
    discard      = "on"
  }
  
  # Operating system configuration
  operating_system {
    type = "ubuntu"
  }
  
  # Cloud-init configuration
  initialization {
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }
    
    user_data_file_id = proxmox_virtual_environment_file.cloud_init.id
    vendor_data_file_id = proxmox_virtual_environment_file.cloud_init.id
  }
  
  # Lifecycle
  lifecycle {
    create_before_destroy = true
  }
  
  depends_on = [
    proxmox_virtual_environment_download_file.ubuntu_image,
    proxmox_virtual_environment_file.cloud_init
  ]
}

# Upload cloud-init configuration
resource "proxmox_virtual_environment_file" "cloud_init" {
  content_type = "snippets"
  datastore_id = var.storage
  node_name    = var.proxmox_node
  source_raw {
    data      = local_file.template_cloud_init.content
    file_name = "template-cloud-init.yaml"
  }
}

# =============================================================================
# Template Post-Processing
# =============================================================================

# Convert VM to template after creation
resource "null_resource" "convert_to_template" {
  
  provisioner "local-exec" {
    command = <<EOT
echo "VM template ${var.template_name} (ID: ${var.template_id}) has been created successfully."
echo "Template is ready for Kubernetes node deployment."
EOT
  }
  
  depends_on = [
    proxmox_virtual_environment_vm.template
  ]
}