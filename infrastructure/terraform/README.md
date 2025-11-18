# Terraform Infrastructure for Kubernetes on Proxmox

This directory contains the Terraform configuration for deploying a Kubernetes cluster on Proxmox VE 9.

## Overview

This infrastructure as code (IaC) setup provisions:
- A VM template based on Ubuntu cloud image
- Multiple Kubernetes control plane (master) nodes
- Multiple Kubernetes worker nodes
- Network and storage configuration

## Prerequisites

1. **Terraform**: Version >= 1.5
   ```bash
   terraform version
   ```

2. **Proxmox VE**: Version 9 running at `192.168.0.200` (or your configured endpoint)

3. **Network Requirements**:
   - Network connectivity to Proxmox API endpoint
   - SSH access to Proxmox host
   - Available IP addresses for VMs

4. **Credentials**:
   - Proxmox API credentials (username/password or API token)
   - SSH private key for Proxmox host access
   - SSH public key for VM access

## Quick Start

### 1. Initialize Configuration

Copy the example variables file:
```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your specific configuration:
```bash
# Use your preferred editor
vim terraform.tfvars
# or
nano terraform.tfvars
```

### 2. Set Environment Variables

For security, it's recommended to use environment variables for sensitive data:

```bash
# For username/password authentication
export PROXMOX_VE_USERNAME="root@pam"
export PROXMOX_VE_PASSWORD="your-password"

# OR for API token authentication (recommended)
export PROXMOX_VE_API_TOKEN="root@pam!terraform=12345678-1234-1234-1234-123456789abc"

# SSH private key (if not using file path)
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)"
```

### 3. Initialize Terraform

Download the required provider plugins:
```bash
terraform init
```

### 4. Review the Plan

See what resources will be created:
```bash
terraform plan
```

### 5. Apply Configuration

Once modules are implemented, create the infrastructure:
```bash
terraform apply
```

## Configuration

### Required Variables

The following variables must be set in `terraform.tfvars` or via environment variables:

- `proxmox_endpoint`: Proxmox API endpoint (default: https://192.168.0.200:8006)
- `proxmox_username` or `PROXMOX_VE_USERNAME`: Proxmox username
- `proxmox_password` or `PROXMOX_VE_PASSWORD`: Proxmox password
- OR `proxmox_api_token` or `PROXMOX_VE_API_TOKEN`: API token
- `ssh_public_key` or `ssh_public_key_file`: SSH public key for VM access

### Optional Variables

See `variables.tf` for all available options, including:
- VM sizing (CPU, memory, disk)
- Network configuration
- Kubernetes version and settings
- Node placement and distribution

## Authentication Methods

### Option 1: Username and Password

```bash
export PROXMOX_VE_USERNAME="root@pam"
export PROXMOX_VE_PASSWORD="your-password"
```

### Option 2: API Token (Recommended)

1. Create an API token in Proxmox:
   - Navigate to Datacenter → Permissions → API Tokens
   - Create a new token with appropriate privileges
   - Note the token ID and secret

2. Set the environment variable:
```bash
export PROXMOX_VE_API_TOKEN="root@pam!terraform=12345678-1234-1234-1234-123456789abc"
```

### SSH Key Setup

For Proxmox host operations:
```bash
# Option 1: Use existing SSH key file
proxmox_ssh_private_key = "~/.ssh/id_rsa"

# Option 2: Use environment variable
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)"
```

For VM access:
```bash
# Option 1: Direct key in variables
ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2E..."

# Option 2: Path to key file (recommended)
ssh_public_key_file = "~/.ssh/id_rsa.pub"
```

## Module Structure

The infrastructure is organized into modules:

```
infrastructure/terraform/
├── main.tf                 # Main configuration and module calls
├── variables.tf            # Variable definitions
├── outputs.tf              # Output definitions
├── providers.tf            # Provider configuration
├── versions.tf             # Version constraints
├── terraform.tfvars.example # Example configuration
├── README.md              # This file
└── modules/               # Module implementations (to be created)
    ├── template/          # VM template creation
    ├── control-plane/     # Control plane node deployment
    └── worker-pool/       # Worker node deployment
```

## Next Steps

The Terraform configuration is initialized and ready. To complete the infrastructure:

1. **Implement the template module** (`./modules/template`)
   - Create VM template from cloud image
   - Configure cloud-init settings

2. **Implement the control-plane module** (`./modules/control-plane`)
   - Deploy control plane VMs from template
   - Initialize Kubernetes control plane
   - Configure high availability

3. **Implement the worker-pool module** (`./modules/worker-pool`)
   - Deploy worker VMs from template
   - Join workers to the cluster

4. **Uncomment module blocks** in `main.tf`

5. **Uncomment outputs** in `outputs.tf`

## Troubleshooting

### terraform init fails

**Issue**: Provider download fails
```
Error: Failed to query available provider packages
```

**Solution**: Check your internet connection and firewall settings. The provider is downloaded from the Terraform Registry.

### Connection to Proxmox fails

**Issue**: Cannot connect to Proxmox API
```
Error: error creating Proxmox client: ...
```

**Solutions**:
- Verify the `proxmox_endpoint` URL is correct
- Check that Proxmox API is accessible from your machine
- Verify credentials are correct
- If using self-signed certificates, ensure `proxmox_insecure = true`

### SSH connection fails

**Issue**: Cannot SSH to Proxmox host
```
Error: SSH authentication failed
```

**Solutions**:
- Verify SSH key path is correct
- Ensure the public key is in Proxmox host's `~/.ssh/authorized_keys`
- Check SSH username is correct (usually `root`)
- Test SSH connection manually: `ssh -i ~/.ssh/id_rsa root@192.168.0.200`

## Security Best Practices

1. **Use API Tokens**: Prefer API tokens over username/password
2. **Environment Variables**: Store sensitive data in environment variables, not in `.tfvars` files
3. **TLS Certificates**: Use valid TLS certificates in production (set `proxmox_insecure = false`)
4. **SSH Keys**: Use strong SSH keys (RSA 4096-bit or Ed25519)
5. **Network Security**: Restrict access to Proxmox API and VMs using firewalls
6. **State Files**: Store Terraform state in a secure remote backend (S3, Terraform Cloud, etc.)

## Useful Commands

```bash
# Initialize and download providers
terraform init

# Validate configuration
terraform validate

# Format configuration files
terraform fmt -recursive

# Plan changes
terraform plan

# Apply changes
terraform apply

# Show current state
terraform show

# List resources
terraform state list

# Destroy infrastructure
terraform destroy

# Show outputs
terraform output
```

## Documentation

- [Terraform Documentation](https://www.terraform.io/docs)
- [Proxmox Provider Documentation](https://registry.terraform.io/providers/bpg/proxmox/latest/docs)
- [Proxmox VE API Documentation](https://pve.proxmox.com/wiki/Proxmox_VE_API)

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review Terraform and provider documentation
3. Consult Proxmox VE documentation
4. Open an issue in the project repository
