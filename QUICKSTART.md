# Quick Start Guide - Kubernetes on Proxmox

This guide will help you quickly deploy a Kubernetes cluster on Proxmox VE using Terraform.

## Prerequisites

- Proxmox VE 7.0 or higher
- Terraform 1.3 or higher
- Proxmox API token with appropriate permissions
- Network configuration (IP range, gateway, etc.)

## Step 1: Configure Proxmox API Token

Create an API token in Proxmox:

```bash
# In Proxmox web UI:
# Datacenter → Permissions → API Tokens → Add
# Or via CLI:
pveum user token add root@pam terraform -privsep 0
```

## Step 2: Clone and Configure

```bash
# Clone the repository
git clone <repository-url>
cd <repository>

# Copy the example configuration
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars with your settings
nano terraform.tfvars
```

## Step 3: Edit terraform.tfvars

Minimal configuration:

```hcl
# Proxmox Connection
proxmox_node     = "pve"
proxmox_api_url  = "https://your-proxmox-server:8006/api2/json"
proxmox_api_token_id     = "root@pam!terraform"
proxmox_api_token_secret = "your-token-secret"
proxmox_tls_insecure = true

# Storage
ubuntu_template_config = {
  storage_pool = "local-lvm"
}

# Control Plane Nodes
control_plane_nodes = {
  "master-1" = {
    vmid     = 100
    hostname = "k8s-master-1"
    ip       = "192.168.1.10/24"
    gateway  = "192.168.1.1"
  }
}

# Worker Nodes
worker_nodes = {
  "worker-1" = {
    vmid     = 110
    hostname = "k8s-worker-1"
    ip       = "192.168.1.20/24"
    gateway  = "192.168.1.1"
  }
}

# SSH Keys (recommended)
ssh_public_keys = [
  "ssh-rsa AAAAB3NzaC... your-key-here"
]
```

## Step 4: Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply
```

## Step 5: Access Your Nodes

After Terraform completes:

```bash
# View SSH connection strings
terraform output control_plane_ssh_connections
terraform output worker_ssh_connections

# Connect to a node
ssh ubuntu@192.168.1.10
```

## Step 6: Install Kubernetes

The VMs are provisioned but Kubernetes is not installed yet. You can:

### Option A: Manual Installation

```bash
# SSH to the first control plane node
ssh ubuntu@192.168.1.10

# Install kubeadm, kubelet, kubectl
# Follow official Kubernetes documentation
```

### Option B: Using Ansible (Recommended)

```bash
# Export Ansible inventory
terraform output ansible_inventory

# Use the inventory with your Ansible playbooks
ansible-playbook -i inventory.yml install-k8s.yml
```

## Scaling

### Add More Nodes

Edit `terraform.tfvars`:

```hcl
worker_nodes = {
  "worker-1" = { ... }
  "worker-2" = {  # New node
    vmid     = 111
    hostname = "k8s-worker-2"
    ip       = "192.168.1.21/24"
    gateway  = "192.168.1.1"
  }
}
```

Then apply:

```bash
terraform apply
```

### Remove Nodes

Remove the node from the map in `terraform.tfvars` and apply:

```bash
terraform apply
```

## Customization

### Override Default Resources

```hcl
control_plane_nodes = {
  "master-1" = {
    vmid      = 100
    hostname  = "k8s-master-1"
    ip        = "192.168.1.10/24"
    gateway   = "192.168.1.1"
    cores     = 8       # Override default 4
    memory    = 16384   # Override default 8192
    disk_size = "100G"  # Override default 50G
  }
}
```

### Change Global Defaults

```hcl
control_plane_defaults = {
  cores     = 6
  memory    = 12288
  disk_size = "80G"
}

worker_defaults = {
  cores     = 12
  memory    = 32768
  disk_size = "200G"
}
```

## Troubleshooting

### Check Terraform State

```bash
terraform show
terraform state list
```

### VM Not Starting

1. Check Proxmox web UI for VM status
2. Verify network configuration
3. Check cloud-init logs in VM:
   ```bash
   ssh ubuntu@<vm-ip>
   sudo cat /var/log/cloud-init-output.log
   ```

### Network Issues

1. Verify IP range doesn't conflict
2. Check gateway is reachable
3. Ensure bridge (vmbr0) exists in Proxmox

### SSH Access Issues

1. Verify SSH keys are correctly formatted
2. Check cloud-init applied keys:
   ```bash
   ssh ubuntu@<vm-ip>
   cat ~/.ssh/authorized_keys
   ```

## Clean Up

To destroy all resources:

```bash
terraform destroy
```

**Warning**: This will delete all VMs and the template!

## Next Steps

1. Install Kubernetes on the provisioned nodes
2. Configure kubectl on your local machine
3. Deploy your applications
4. Set up monitoring and logging

## Additional Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Proxmox VE Documentation](https://pve.proxmox.com/pve-docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Module README](modules/kubernetes_node_pool/README.md)

## Support

For issues related to:
- **Terraform Configuration**: Check module READMEs
- **Proxmox Issues**: Check Proxmox logs and documentation
- **Networking**: Verify your network configuration in Proxmox

---

**Happy Kubernetes Clustering! 🚀**
