# Kubernetes Node Pool Module

This Terraform module creates a pool of Kubernetes nodes (control-plane or worker) on Proxmox VE by cloning from an Ubuntu template.

## Features

- **Elastic Scaling**: Add or remove nodes by modifying the node map
- **Per-Node Configuration**: Individual VMIDs, hostnames, IPs, CPU, memory, and disk sizing
- **Cloud-Init Integration**: Automated provisioning with static networking
- **QEMU Guest Agent**: Enabled for better VM management
- **Sensible Defaults**: Masters (4vCPU/8GB), Workers (8vCPU/16GB)
- **SSH Access**: Support for SSH keys and password authentication
- **Connection Details**: Outputs IP addresses, hostnames, and SSH connection strings
- **Kubeadm Integration**: Outputs data needed for cluster joining

## Usage

```hcl
module "control_plane" {
  source = "./modules/kubernetes_node_pool"
  
  proxmox_node  = "pve"
  template_vmid = 9001
  storage_pool  = "local-lvm"
  
  node_type      = "control-plane"
  default_cores  = 4
  default_memory = 8192
  
  nodes = {
    "master-1" = {
      vmid     = 100
      hostname = "k8s-master-1"
      ip       = "192.168.1.10/24"
      gateway  = "192.168.1.1"
      disk_size = "50G"
    }
    "master-2" = {
      vmid     = 101
      hostname = "k8s-master-2"
      ip       = "192.168.1.11/24"
      gateway  = "192.168.1.1"
      disk_size = "50G"
    }
  }
  
  ssh_public_keys = [
    "ssh-rsa AAAAB3Nza... user@host"
  ]
}

module "workers" {
  source = "./modules/kubernetes_node_pool"
  
  proxmox_node  = "pve"
  template_vmid = 9001
  storage_pool  = "local-lvm"
  
  node_type      = "worker"
  default_cores  = 8
  default_memory = 16384
  
  nodes = {
    "worker-1" = {
      vmid     = 110
      hostname = "k8s-worker-1"
      ip       = "192.168.1.20/24"
      gateway  = "192.168.1.1"
      disk_size = "100G"
    }
    "worker-2" = {
      vmid     = 111
      hostname = "k8s-worker-2"
      ip       = "192.168.1.21/24"
      gateway  = "192.168.1.1"
      disk_size = "100G"
    }
  }
  
  ssh_public_keys = [
    "ssh-rsa AAAAB3Nza... user@host"
  ]
}
```

## Variables

### Required Variables

| Name | Description | Type |
|------|-------------|------|
| `proxmox_node` | Proxmox node where VMs will be created | `string` |
| `template_vmid` | VM ID of the template to clone from | `number` |
| `nodes` | Map of nodes to create | `map(object)` |
| `storage_pool` | Storage pool for VM disks | `string` |

### Optional Variables

| Name | Description | Default |
|------|-------------|---------|
| `default_cores` | Default CPU cores | `2` |
| `default_memory` | Default memory in MB | `2048` |
| `default_disk_size` | Default disk size | `"32G"` |
| `network_bridge` | Network bridge | `"vmbr0"` |
| `ssh_public_keys` | List of SSH public keys | `[]` |
| `ssh_password` | SSH password for ubuntu user | `null` |
| `cloud_init_storage` | Storage for cloud-init | `null` |
| `tags` | Tags to apply to VMs | `["kubernetes"]` |
| `node_type` | Node type for tagging | `"node"` |
| `start_on_boot` | Auto-start on boot | `true` |

### Node Object Structure

```hcl
nodes = {
  "node-key" = {
    vmid         = number          # Required: VM ID
    hostname     = string          # Required: Hostname
    ip           = string          # Required: IP address with CIDR (e.g., "192.168.1.10/24")
    gateway      = string          # Required: Gateway IP
    cores        = number          # Optional: CPU cores (uses default_cores if omitted)
    memory       = number          # Optional: Memory in MB (uses default_memory if omitted)
    disk_size    = string          # Optional: Disk size (uses default_disk_size if omitted)
    nameserver   = string          # Optional: DNS server (default: "8.8.8.8")
    searchdomain = string          # Optional: Search domain (default: "local")
  }
}
```

## Outputs

| Name | Description |
|------|-------------|
| `nodes` | Complete details of all created nodes |
| `node_ips` | Map of node names to IP addresses |
| `node_hostnames` | Map of node names to hostnames |
| `node_vmids` | Map of node names to VM IDs |
| `ssh_connection_strings` | SSH connection strings for all nodes |
| `kubeadm_join_data` | Data for kubeadm join commands |
| `ansible_inventory` | Ansible-style inventory data |

## Scaling

To add more nodes, simply add entries to the `nodes` map:

```hcl
nodes = {
  "master-1" = { ... }
  "master-2" = { ... }
  "master-3" = { ... }  # New node
}
```

To remove nodes, delete entries from the map. Terraform will destroy the corresponding VMs.

## Dependencies

This module depends on:
- A Proxmox template created by the `template` module
- Proxmox API access with appropriate permissions
- Network configuration (bridge, IP range, gateway)

## Notes

- VMs are fully cloned from the template (not linked clones)
- Cloud-init handles initial provisioning and network configuration
- QEMU guest agent is installed and enabled automatically
- SSH keys are recommended over password authentication
- All nodes in a pool share the same SSH credentials
