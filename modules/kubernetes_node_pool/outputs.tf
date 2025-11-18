output "nodes" {
  description = "Map of all created nodes with their details"
  value = {
    for key, node in proxmox_virtual_environment_vm.node : key => {
      vmid     = node.vm_id
      name     = node.name
      hostname = var.nodes[key].hostname
      ip       = var.nodes[key].ip
      cores    = node.cpu[0].cores
      memory   = node.memory[0].dedicated
      status   = "created"
    }
  }
}

output "node_ips" {
  description = "Map of node names to IP addresses"
  value = {
    for key, node in var.nodes : key => split("/", node.ip)[0]
  }
}

output "node_hostnames" {
  description = "Map of node names to hostnames"
  value = {
    for key, node in var.nodes : key => node.hostname
  }
}

output "node_vmids" {
  description = "Map of node names to VM IDs"
  value = {
    for key, node in var.nodes : key => node.vmid
  }
}

output "ssh_connection_strings" {
  description = "SSH connection strings for all nodes"
  value = {
    for key, node in var.nodes : key => "ssh ubuntu@${split("/", node.ip)[0]}"
  }
}

output "kubeadm_join_data" {
  description = "Data needed for kubeadm join commands (to be computed after cluster init)"
  value = {
    control_plane_endpoint = var.node_type == "control-plane" ? split("/", var.nodes[keys(var.nodes)[0]].ip)[0] : null
    node_ips              = [for node in var.nodes : split("/", node.ip)[0]]
    node_hostnames        = [for node in var.nodes : node.hostname]
  }
}

output "ansible_inventory" {
  description = "Ansible-style inventory data"
  value = {
    hosts = {
      for key, node in var.nodes : node.hostname => {
        ansible_host = split("/", node.ip)[0]
        ansible_user = "ubuntu"
        node_type    = var.node_type
      }
    }
  }
}
