# Ubuntu Template Module

This Terraform module automates the creation of Ubuntu 22.04 base templates for Proxmox VE.

## Features

- Downloads the latest Ubuntu 22.04 cloud image with SHA256 verification
- Uploads the QCOW2 image to Proxmox storage
- Creates a VM template with optimal settings
- Configures cloud-init for easy VM deployment
- Supports interactive replacement when images exist
- Configurable for different OS versions and storage pools

## Usage

### Basic Usage

```hcl
module "ubuntu_template" {
  source = "./modules/template"
  
  proxmox_node = "pve"
  storage_pool = "local-lvm"
  
  # Optional configurations
  vm_id      = 9001
  cores      = 2
  memory     = 2048
  disk_size  = "20G"
}
```

### Advanced Configuration

```hcl
module "ubuntu_template" {
  source = "./modules/template"
  
  # Proxmox settings
  proxmox_node        = "pve"
  storage_pool        = "local-lvm"
  network_bridge      = "vmbr0"
  cloud_init_storage  = "local-lvm"
  
  # VM specifications
  vm_id     = 9001
  cores     = 4
  memory    = 4096
  disk_size = "40G"
  
  # Ubuntu image settings
  ubuntu_version      = "22.04"
  ubuntu_architecture = "amd64"
  interactive_replace = false  # Non-automated mode
}
```

## Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| proxmox_node | Proxmox node where the template will be created | `string` | - |
| vm_id | VM ID for the template | `number` | `9001` |
| storage_pool | Storage pool for the template disk | `string` | - |
| network_bridge | Network bridge for the template | `string` | `"vmbr0"` |
| cores | Number of CPU cores | `number` | `2` |
| memory | Memory in MB | `number` | `2048` |
| disk_size | Disk size | `string` | `"20G"` |
| cloud_init_storage | Storage pool for cloud-init drive | `string` | `null` |
| ubuntu_version | Ubuntu version | `string` | `"22.04"` |
| ubuntu_architecture | Ubuntu architecture | `string` | `"amd64"` |
| interactive_replace | Whether to prompt for replacement when image exists | `bool` | `true` |

## Outputs

| Name | Description |
|------|-------------|
| template_name | Name of the created Ubuntu template |
| template_vmid | VMID of the created Ubuntu template |
| template_storage_path | Storage path of the uploaded Ubuntu image |
| template_node | Proxmox node where the template was created |

## Requirements

- Terraform >= 1.0
- Proxmox VE >= 7.0
- Proxmox API token with appropriate permissions
- curl and sha256sum utilities available

## Provider Configuration

Configure the Proxmox provider in your root module:

```hcl
terraform {
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
      version = ">= 3.0.0"
    }
  }
}

provider "proxmox" {
  pm_api_url          = var.proxmox_api_url
  pm_api_token_id     = var.proxmox_api_token_id
  pm_api_token_secret = var.proxmox_api_token_secret
  pm_tls_insecure     = var.proxmox_tls_insecure
}
```

## Template Features

The created template includes:

- **VirtIO drivers** for optimal performance
- **QEMU Guest Agent** for better integration
- **Cloud-init support** with DHCP networking
- **Automatic package updates** on first boot
- **SSH access** with key-based authentication
- **Optimized disk configuration** with growpart support

## Script Details

The `get-ubuntu-cloudimg.sh` script handles:

1. Downloading the official Ubuntu cloud image
2. Verifying SHA256 checksums
3. Interactive replacement prompts (when enabled)
4. Error handling and cleanup

## Security Considerations

- Store Proxmox API tokens securely (use environment variables or secure storage)
- Review SSH key configurations before deploying
- Consider network security settings for your environment
- Template images are downloaded from official Ubuntu repositories

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure Proxmox API token has sufficient permissions
2. **Storage Not Found**: Verify storage pool exists on the target node
3. **Network Bridge**: Confirm bridge name matches your Proxmox configuration
4. **Download Failures**: Check internet connectivity and Ubuntu repository status

### Debug Mode

Enable Terraform debug logging:
```bash
export TF_LOG=DEBUG
terraform apply
```

## License

This module is licensed under the MIT License.