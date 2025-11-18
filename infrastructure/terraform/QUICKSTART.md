# Quick Start Guide

Get your Kubernetes cluster on Proxmox up and running in minutes!

## Prerequisites Checklist

- [ ] Proxmox VE 9 installed and accessible at `192.168.0.200` (or your configured IP)
- [ ] Proxmox API credentials (username/password or API token)
- [ ] SSH access to Proxmox host
- [ ] SSH key pair generated (`ssh-keygen -t rsa -b 4096`)
- [ ] Available IP addresses for VMs
- [ ] Sufficient resources in Proxmox (CPU, RAM, Storage)

## Step 1: Configure Credentials

Choose one of the following authentication methods:

### Option A: Environment Variables (Recommended)

Create a file `env.sh` (do not commit this):

```bash
#!/bin/bash
# Proxmox authentication
export PROXMOX_VE_USERNAME="root@pam"
export PROXMOX_VE_PASSWORD="your-password-here"

# SSH for Proxmox host
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)"
```

Then source it:
```bash
source env.sh
```

### Option B: Configuration File

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
nano terraform.tfvars  # or vim, code, etc.
```

## Step 2: Review and Customize Configuration

Edit `terraform.tfvars` to match your environment:

```hcl
# Proxmox connection
proxmox_endpoint    = "https://192.168.0.200:8006"
proxmox_node_name   = "pve"  # Your Proxmox node name

# Control plane nodes (masters)
control_plane_count = 3  # Odd number recommended (1, 3, 5)

# Worker nodes
worker_count = 3  # Adjust based on your needs

# Network settings
network_gateway = "192.168.0.1"
network_dns_servers = ["8.8.8.8", "8.8.4.4"]

# SSH key for VM access
ssh_public_key_file = "~/.ssh/id_rsa.pub"
```

## Step 3: Initialize Terraform

```bash
terraform init
```

Expected output:
```
Initializing the backend...
Initializing provider plugins...
- Installing bpg/proxmox v0.86.0...
Terraform has been successfully initialized!
```

## Step 4: Validate Configuration

```bash
terraform validate
```

Expected output:
```
Success! The configuration is valid.
```

## Step 5: Preview Changes

```bash
terraform plan
```

Review the output to see what will be created.

## Step 6: Apply Configuration (When Modules Are Ready)

**Note**: Currently, the module implementations are placeholders. Once the modules are implemented:

```bash
terraform apply
```

Type `yes` when prompted to create the infrastructure.

## Current Status

The Terraform scaffolding is complete with:
- ✅ Version constraints (Terraform >= 1.5, Proxmox provider v0.86.0)
- ✅ Provider configuration (API + SSH)
- ✅ Comprehensive variables
- ✅ Module structure (placeholders)
- ✅ Output definitions
- ✅ Example configuration

### What's Next?

To complete the infrastructure:

1. **Implement the template module** (`modules/template/`)
   - Creates VM template from Ubuntu cloud image
   - Configures cloud-init settings

2. **Implement the control-plane module** (`modules/control-plane/`)
   - Deploys Kubernetes master nodes
   - Initializes the cluster
   - Sets up HA (if multiple masters)

3. **Implement the worker-pool module** (`modules/worker-pool/`)
   - Deploys Kubernetes worker nodes
   - Joins them to the cluster

4. **Uncomment module blocks** in `main.tf`

5. **Uncomment outputs** in `outputs.tf`

## Useful Commands

```bash
# Initialize
terraform init

# Validate
terraform validate

# Format code
terraform fmt -recursive

# Plan changes
terraform plan

# Apply changes
terraform apply

# Show current state
terraform show

# List resources
terraform state list

# Show outputs
terraform output

# Destroy everything
terraform destroy
```

## Using the Makefile

```bash
# Show available commands
make help

# Initialize
make init

# Validate
make validate

# Format
make fmt

# Plan
make plan

# Apply
make apply

# Destroy
make destroy

# Clean terraform files
make clean
```

## Troubleshooting

### "command not found: terraform"

Install Terraform:
```bash
# Ubuntu/Debian
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt-get update
sudo apt-get install terraform
```

### "Error: error creating Proxmox client"

- Check Proxmox endpoint is correct and accessible
- Verify credentials are set correctly
- Test connection: `curl -k https://192.168.0.200:8006`

### "SSH authentication failed"

- Verify SSH key is correct
- Check public key is in Proxmox host's `~/.ssh/authorized_keys`
- Test manually: `ssh root@192.168.0.200`

## Resource Requirements

### Minimum (Development)
- Control Plane: 3 VMs × (2 CPU, 4GB RAM, 50GB disk)
- Workers: 3 VMs × (2 CPU, 4GB RAM, 50GB disk)
- **Total**: 12 CPUs, 24GB RAM, 300GB storage

### Recommended (Production)
- Control Plane: 3 VMs × (4 CPU, 8GB RAM, 100GB disk)
- Workers: 5 VMs × (8 CPU, 16GB RAM, 200GB disk)
- **Total**: 52 CPUs, 104GB RAM, 1.3TB storage

## Next Steps After Deployment

Once the infrastructure is deployed (modules implemented):

1. **Access the cluster**
   ```bash
   # Get connection info
   terraform output
   ```

2. **Configure kubectl**
   ```bash
   # Copy kubeconfig from control plane
   scp ubuntu@<control-plane-ip>:~/.kube/config ~/.kube/config
   ```

3. **Verify cluster**
   ```bash
   kubectl get nodes
   kubectl get pods -A
   ```

4. **Deploy applications**
   - Install ingress controller
   - Deploy your applications
   - Configure storage classes

## Support

- **Documentation**: See [README.md](README.md)
- **Environment Variables**: See [ENV_VARS.md](ENV_VARS.md)
- **Terraform Docs**: https://www.terraform.io/docs
- **Proxmox Provider**: https://registry.terraform.io/providers/bpg/proxmox/latest/docs

## Security Reminders

- 🔒 Never commit `terraform.tfvars` or `env.sh` with real credentials
- 🔒 Use API tokens instead of passwords when possible
- 🔒 Restrict file permissions: `chmod 600 env.sh terraform.tfvars`
- 🔒 Use proper TLS certificates in production (`proxmox_insecure = false`)
- 🔒 Regularly rotate credentials
