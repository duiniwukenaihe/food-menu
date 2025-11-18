# Infrastructure as Code

This directory contains the Infrastructure as Code (IaC) for the project.

## Structure

```
infrastructure/
└── terraform/          # Terraform configuration for Proxmox/Kubernetes
    ├── main.tf         # Main configuration
    ├── variables.tf    # Variable definitions
    ├── outputs.tf      # Output definitions
    ├── providers.tf    # Provider configuration
    ├── versions.tf     # Version constraints
    ├── terraform.tfvars.example  # Example configuration
    ├── README.md       # Terraform documentation
    └── modules/        # Terraform modules
        ├── template/           # VM template module
        ├── control-plane/      # Kubernetes control plane module
        └── worker-pool/        # Kubernetes worker pool module
```

## Getting Started

See the [Terraform README](./terraform/README.md) for detailed instructions on deploying the infrastructure.

## Quick Start

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your configuration
terraform init
terraform plan
terraform apply
```

## Requirements

- Terraform >= 1.5
- Proxmox VE 9 cluster
- Network connectivity to Proxmox API
- SSH access to Proxmox host
- SSH keys for VM access

## Environment Variables

For secure credential management:

```bash
# Proxmox authentication
export PROXMOX_VE_USERNAME="root@pam"
export PROXMOX_VE_PASSWORD="your-password"
# OR
export PROXMOX_VE_API_TOKEN="root@pam!terraform=your-token"

# SSH private key
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)"
```

## Documentation

- [Terraform Configuration](./terraform/README.md)
- [Template Module](./terraform/modules/template/README.md)
- [Control Plane Module](./terraform/modules/control-plane/README.md)
- [Worker Pool Module](./terraform/modules/worker-pool/README.md)
