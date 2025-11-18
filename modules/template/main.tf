terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.50.0"
      source  = "telmate/proxmox"
      version = ">= 3.0.0"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 3.0.0"
    }
  }
}

# Download and verify Ubuntu cloud image
resource "null_resource" "download_ubuntu_image" {
  triggers = {
    version      = var.ubuntu_version
    architecture = var.ubuntu_architecture
  }

  provisioner "local-exec" {
    command = "${path.module}/../../scripts/get-ubuntu-cloudimg.sh ${var.ubuntu_version} ${var.ubuntu_architecture} ${var.interactive_replace}"
    
    working_dir = path.module
  }
}

# Upload the QCOW2 image to Proxmox storage
resource "proxmox_virtual_environment_file" "ubuntu_image" {
  content_type = "iso"
  datastore_id = var.storage_pool
  node_name    = var.proxmox_node
  
  source_file = "${path.module}/../../downloads/ubuntu-${var.ubuntu_version}-server-cloudimg-${var.ubuntu_architecture}.img"
  
  file_name = "ubuntu-${var.ubuntu_version}-server-cloudimg-${var.ubuntu_architecture}.img"
  
  depends_on = [null_resource.download_ubuntu_image]
}

# Create the VM from the uploaded image
resource "proxmox_virtual_environment_vm" "ubuntu_template" {
  name        = "ubuntu-${var.ubuntu_version}-template"
  node_name   = var.proxmox_node
  vm_id       = var.vm_id
  
  # Basic VM configuration
  description = "Ubuntu ${var.ubuntu_version} Cloud Image Template"
  tags        = ["template", "ubuntu", "cloud-init"]
  
  # CPU and memory
  cpu {
    cores = var.cores
    type  = "host"
  }
  
  memory {
    dedicated = var.memory
  }
  
  # Network configuration
  network_device {
    bridge      = var.network_bridge
    model       = "virtio"
    firewall    = false
  }
  
  # Main disk from uploaded cloud image
  disk {
    datastore_id = var.storage_pool
    file_id      = proxmox_virtual_environment_file.ubuntu_image.id
    interface    = "virtio0"
    size         = var.disk_size
    iothread     = true
  }
  
  # Cloud-init drive
  disk {
    datastore_id = var.cloud_init_storage != null ? var.cloud_init_storage : var.storage_pool
    interface    = "ide0"
    size         = "4M"
    file_format  = "raw"
  }
  
  # QEMU Guest Agent
  agent {
    enabled = true
  }
  
  # Cloud-init configuration
  operating_system {
    type = "ubuntu"
  }
  
  # Boot configuration
  boot_order = ["virtio0", "ide0"]
  
  # Template settings
  template = true
  
  depends_on = [
    null_resource.download_ubuntu_image,
    proxmox_virtual_environment_file.ubuntu_image
  ]
}

# Configure cloud-init settings for the template
resource "proxmox_virtual_environment_vm" "ubuntu_template_config" {
  vm_id = proxmox_virtual_environment_vm.ubuntu_template.vm_id
  node_name = var.proxmox_node
  
  # Cloud-init user configuration
  user_data = base64encode(templatefile("${path.module}/cloud-init-user-data.tpl", {
    ssh_public_key = ""
  }))
  
  # Cloud-init network configuration
  network_data = base64encode(templatefile("${path.module}/cloud-init-network-data.tpl", {}))
  
  depends_on = [proxmox_virtual_environment_vm.ubuntu_template]
}