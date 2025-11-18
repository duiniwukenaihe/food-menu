# =============================================================================
# Kubernetes Node Pool Module - Creates Master and Worker Nodes
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
    null = {
      source  = "hashicorp/null"
      version = ">= 3.2.4"
    }
  }
}

# =============================================================================
# Master Nodes
# =============================================================================

# Create master nodes
resource "proxmox_virtual_environment_vm" "master" {
  count = var.master_count
  
  name        = "${var.cluster_name}-master-${count.index + 1}"
  description = "Kubernetes master node ${count.index + 1}"
  node_name   = var.proxmox_node
  vm_id       = var.master_vmid_start + count.index
  
  # CPU configuration
  cpu {
    cores = var.master_cores
    type  = "host"
  }
  
  # Memory configuration
  memory {
    dedicated = var.master_memory
  }
  
  # Network configuration
  network_device {
    bridge = var.bridge
    model  = "virtio"
  }
  
  # Disk configuration
  disk {
    datastore_id = var.storage
    interface    = "virtio0"
    size         = var.master_disk_size
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
    
    user_data_file_id = proxmox_virtual_environment_file.master_cloud_init[count.index].id
    vendor_data_file_id = proxmox_virtual_environment_file.master_cloud_init[count.index].id
  }
  
  # Tags
  tags = merge(var.tags, {
    "role" = "master"
    "node_index" = tostring(count.index + 1)
  })
  
  # Start after creation
  started = true
  
  depends_on = [
    proxmox_virtual_environment_file.master_cloud_init
  ]
}

# Create cloud-init files for master nodes
resource "proxmox_virtual_environment_file" "master_cloud_init" {
  count = var.master_count
  
  content_type = "snippets"
  datastore_id = var.storage
  node_name    = var.proxmox_node
  source_raw {
    data      = templatefile("${path.module}/templates/master.yaml.tpl", {
      hostname = "${var.cluster_name}-master-${count.index + 1}"
      username = var.username
      password = var.password
      ssh_public_key = var.ssh_public_key
      
      # Kubernetes Configuration
      k8s_version = var.k8s_version
      pod_network_cidr = var.pod_network_cidr
      network_plugin = var.network_plugin
      image_repository = var.image_repository
      is_master = true
      cluster_name = var.cluster_name
      master_count = var.master_count
      node_index = count.index + 1
    })
    file_name = "master-${count.index + 1}-cloud-init.yaml"
  }
}

# =============================================================================
# Worker Nodes
# =============================================================================

# Create worker nodes
resource "proxmox_virtual_environment_vm" "worker" {
  count = var.worker_count
  
  name        = "${var.cluster_name}-worker-${count.index + 1}"
  description = "Kubernetes worker node ${count.index + 1}"
  node_name   = var.proxmox_node
  vm_id       = var.worker_vmid_start + count.index
  
  # CPU configuration
  cpu {
    cores = var.worker_cores
    type  = "host"
  }
  
  # Memory configuration
  memory {
    dedicated = var.worker_memory
  }
  
  # Network configuration
  network_device {
    bridge = var.bridge
    model  = "virtio"
  }
  
  # Disk configuration
  disk {
    datastore_id = var.storage
    interface    = "virtio0"
    size         = var.worker_disk_size
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
    
    user_data_file_id = proxmox_virtual_environment_file.worker_cloud_init[count.index].id
    vendor_data_file_id = proxmox_virtual_environment_file.worker_cloud_init[count.index].id
  }
  
  # Tags
  tags = merge(var.tags, {
    "role" = "worker"
    "node_index" = tostring(count.index + 1)
  })
  
  # Start after creation
  started = true
  
  depends_on = [
    proxmox_virtual_environment_file.worker_cloud_init,
    proxmox_virtual_environment_vm.master
  ]
}

# Create cloud-init files for worker nodes
resource "proxmox_virtual_environment_file" "worker_cloud_init" {
  count = var.worker_count
  
  content_type = "snippets"
  datastore_id = var.storage
  node_name    = var.proxmox_node
  source_raw {
    data      = templatefile("${path.module}/templates/worker.yaml.tpl", {
      hostname = "${var.cluster_name}-worker-${count.index + 1}"
      username = var.username
      password = var.password
      ssh_public_key = var.ssh_public_key
      
      # Kubernetes Configuration
      k8s_version = var.k8s_version
      pod_network_cidr = var.pod_network_cidr
      network_plugin = var.network_plugin
      image_repository = var.image_repository
      is_master = false
      cluster_name = var.cluster_name
      node_index = count.index + 1
    })
    file_name = "worker-${count.index + 1}-cloud-init.yaml"
  }
}

# =============================================================================
# Cluster Initialization
# =============================================================================

# Initialize Kubernetes cluster on first master node
resource "null_resource" "cluster_init" {
  count = var.master_count > 0 ? 1 : 0
  
  provisioner "local-exec" {
    command = <<EOT
echo "Waiting for Kubernetes cluster initialization..."
echo "Master nodes: ${var.master_count}"
echo "Worker nodes: ${var.worker_count}"
echo "Kubernetes version: ${var.k8s_version}"
echo "Network plugin: ${var.network_plugin}"
EOT
  }
  
  depends_on = [
    proxmox_virtual_environment_vm.master,
    proxmox_virtual_environment_vm.worker
  ]
}

# =============================================================================
# Local Values
# =============================================================================

locals {
  master_names = [for vm in proxmox_virtual_environment_vm.master : vm.name]
  worker_names = [for vm in proxmox_virtual_environment_vm.worker : vm.name]
  
  # Get IP addresses from VM network interfaces
  master_ips = [
    for vm in proxmox_virtual_environment_vm.master :
    try(vm.network_device[0].ipv4_addresses[0], "dhcp")
  ]
  
  worker_ips = [
    for vm in proxmox_virtual_environment_vm.worker :
    try(vm.network_device[0].ipv4_addresses[0], "dhcp")
  ]
}