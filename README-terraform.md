# Ubuntu Template Module for Proxmox VE

This repository contains a Terraform module that automates the creation of Ubuntu 22.04 base templates for Proxmox Virtual Environment (VE).

## 🚀 Quick Start

### Prerequisites

- Proxmox VE 7.0+
- Terraform 1.0+
- Proxmox API token with appropriate permissions
- curl and sha256sum utilities

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ubuntu-template-proxmox
   ```

2. **Configure your environment**
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your Proxmox configuration
   ```

3. **Initialize Terraform**
   ```bash
   terraform init
   ```

4. **Apply the configuration**
   ```bash
   terraform apply
   ```

## 📁 Project Structure

```
.
├── main.tf                    # Root Terraform configuration
├── variables.tf               # Root variables
├── outputs.tf                 # Root outputs
├── modules.tf                 # Module instantiation
├── terraform.tfvars.example   # Example configuration
├── modules/
│   └── template/              # Ubuntu template module
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       ├── cloud-init-user-data.tpl
│       ├── cloud-init-network-data.tpl
│       └── README.md
├── scripts/
│   ├── get-ubuntu-cloudimg.sh  # Image download script
│   └── validate.sh             # Validation script
└── downloads/                  # Directory for downloaded images (gitignored)
```

## 🔧 Configuration

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `proxmox_node` | Proxmox node name | `"pve"` |
| `proxmox_api_url` | Proxmox API URL | `"https://pve.example.com:8006/api2/json"` |
| `proxmox_api_token_id` | API token ID | `"root@pam!terraform"` |
| `proxmox_api_token_secret` | API token secret | `"your-secret-token"` |
| `storage_pool` | Storage pool name | `"local-lvm"` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `vm_id` | Template VM ID | `9001` |
| `cores` | CPU cores | `2` |
| `memory` | Memory in MB | `2048` |
| `disk_size` | Disk size | `"20G"` |
| `network_bridge` | Network bridge | `"vmbr0"` |
| `ubuntu_version` | Ubuntu version | `"22.04"` |
| `interactive_replace` | Prompt to replace existing images | `true` |

## 🎯 Features

### Automated Image Management

- **Download**: Fetches the latest Ubuntu 22.04 cloud image from official repositories
- **Verification**: Validates SHA256 checksums to ensure image integrity
- **Interactive Mode**: Prompts for confirmation when replacing existing images
- **Non-Interactive Mode**: Supports automated deployments

### Template Configuration

- **Optimized Settings**: Pre-configured with virtio drivers, QEMU guest agent
- **Cloud-Init Ready**: Fully configured for automated VM deployment
- **Network Configuration**: DHCP networking with bridge support
- **Security**: SSH key-based authentication, passwordless sudo

### Storage Management

- **Flexible Storage**: Configurable storage pools for different environments
- **Cloud-Init Drive**: Separate ISO storage for cloud-init configuration
- **Disk Resizing**: Automatic disk expansion on first boot

## 📖 Usage Examples

### Basic Template Creation

```hcl
module "ubuntu_template" {
  source = "./modules/template"
  
  proxmox_node = "pve"
  storage_pool = "local-lvm"
}
```

### Advanced Configuration

```hcl
module "ubuntu_template" {
  source = "./modules/template"
  
  proxmox_node       = "pve"
  storage_pool       = "local-lvm"
  vm_id              = 9001
  cores              = 4
  memory             = 4096
  disk_size          = "40G"
  network_bridge     = "vmbr0"
  cloud_init_storage = "local"
  
  ubuntu_version      = "22.04"
  ubuntu_architecture = "amd64"
  interactive_replace = false
}
```

## 🔄 Workflow

1. **Image Download**: The `get-ubuntu-cloudimg.sh` script downloads the Ubuntu cloud image
2. **Verification**: SHA256 checksum is verified against official Ubuntu signatures
3. **Upload**: Image is uploaded to specified Proxmox storage
4. **VM Creation**: Template VM is created with optimal settings
5. **Cloud-Init**: Cloud-init configuration is applied
6. **Template Conversion**: VM is converted to a template for cloning

## 🛠️ Scripts

### get-ubuntu-cloudimg.sh

Handles Ubuntu cloud image download and verification:

```bash
# Usage
./scripts/get-ubuntu-cloudimg.sh [version] [architecture] [interactive_replace]

# Examples
./scripts/get-ubuntu-cloudimg.sh 22.04 amd64 true
./scripts/get-ubuntu-cloudimg.sh 22.04 amd64 false  # Non-interactive
```

### validate.sh

Validates module configuration and file structure:

```bash
./scripts/validate.sh
```

## 📊 Outputs

The module provides the following outputs:

- `template_name`: Name of the created template
- `template_vmid`: VM ID of the template
- `template_storage_path`: Storage location of the image
- `template_node`: Proxmox node where template was created

## 🔐 Security Considerations

- Store API tokens securely (environment variables, secret management)
- Review SSH key configurations before deployment
- Consider network security settings
- Images are downloaded from official Ubuntu repositories

## 🐛 Troubleshooting

### Common Issues

1. **API Authentication**: Verify token permissions and API URL
2. **Storage Issues**: Confirm storage pool exists and is accessible
3. **Network Configuration**: Check bridge name and network settings
4. **Download Failures**: Verify internet connectivity and repository status

### Debug Mode

```bash
export TF_LOG=DEBUG
terraform apply
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙋‍♂️ Support

For issues and questions:

1. Check the [troubleshooting guide](#-troubleshooting)
2. Review existing GitHub issues
3. Create a new issue with detailed information

---

**Note**: This module is designed for Proxmox VE environments and requires appropriate permissions and network access to Ubuntu's cloud image repositories.