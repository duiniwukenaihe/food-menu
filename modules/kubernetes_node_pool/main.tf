terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.50.0"
    }
  }
}

locals {
  ssh_keys_list = var.ssh_public_keys
  ssh_password_hash = var.ssh_password != null ? sha512(var.ssh_password) : ""
}

resource "proxmox_virtual_environment_vm" "node" {
  for_each = var.nodes
  
  name      = each.value.hostname
  node_name = var.proxmox_node
  vm_id     = each.value.vmid
  
  description = "Kubernetes ${var.node_type} node: ${each.value.hostname}"
  tags        = concat(var.tags, [var.node_type])
  
  on_boot = var.start_on_boot
  
  clone {
    vm_id = var.template_vmid
    full  = true
  }
  
  cpu {
    cores = coalesce(each.value.cores, var.default_cores)
    type  = "host"
  }
  
  memory {
    dedicated = coalesce(each.value.memory, var.default_memory)
  }
  
  network_device {
    bridge   = var.network_bridge
    model    = "virtio"
    firewall = false
  }
  
  disk {
    datastore_id = var.storage_pool
    interface    = "virtio0"
    size         = coalesce(each.value.disk_size, var.default_disk_size)
    iothread     = true
    discard      = "on"
    ssd          = true
  }
  
  agent {
    enabled = true
    timeout = "15m"
  }
  
  operating_system {
    type = "l26"
  }
  
  initialization {
    datastore_id = var.cloud_init_storage != null ? var.cloud_init_storage : var.storage_pool
    
    user_account {
      username = "ubuntu"
      password = var.ssh_password
      keys     = var.ssh_public_keys
    }
    
    ip_config {
      ipv4 {
        address = each.value.ip
        gateway = each.value.gateway
      }
    }
    
    dns {
      servers = [coalesce(each.value.nameserver, "8.8.8.8")]
      domain  = coalesce(each.value.searchdomain, "local")
    }
    
    user_data_file_id = proxmox_virtual_environment_file.cloud_init_user_data[each.key].id
  }
  
  boot_order = ["virtio0"]
  
  lifecycle {
    ignore_changes = [
      initialization[0].user_data_file_id,
    ]
  }
}

resource "proxmox_virtual_environment_file" "cloud_init_user_data" {
  for_each = var.nodes
  
  content_type = "snippets"
  datastore_id = var.cloud_init_storage != null ? var.cloud_init_storage : var.storage_pool
  node_name    = var.proxmox_node
  
  source_raw {
    data = templatefile("${path.module}/cloud-init-user-data.tpl", {
      hostname        = each.value.hostname
      searchdomain    = coalesce(each.value.searchdomain, "local")
      ssh_public_keys = var.ssh_public_keys
      ssh_password    = var.ssh_password != null ? local.ssh_password_hash : ""
      lock_passwd     = var.ssh_password != null ? "false" : "true"
    })
    
    file_name = "cloud-init-user-data-${each.value.hostname}.yaml"
  }
}
