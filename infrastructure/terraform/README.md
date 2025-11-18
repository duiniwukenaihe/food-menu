# Kubernetes Cluster Infrastructure with Terraform

This repository contains Terraform configurations for deploying a production-ready Kubernetes cluster on AWS using kubeadm, containerd, and cloud-init.

## 🏗️ Architecture Overview

The infrastructure provisions:

- **VPC** with public and private subnets across multiple Availability Zones
- **NAT Gateways** for outbound internet access from private subnets
- **Security Groups** with proper Kubernetes networking rules
- **EC2 Instances** for master and worker nodes
- **EBS Volumes** for persistent storage
- **Cloud-init Templates** for automated node provisioning
- **Kubeadm Scripts** for cluster initialization and worker joining

## 🚀 Quick Start

### Prerequisites

1. **Terraform >= 1.0**
2. **AWS CLI** configured with appropriate credentials
3. **SSH Key Pair** (or let Terraform generate one)

### Installation Steps

1. **Clone and Navigate**
   ```bash
   cd infrastructure/terraform
   ```

2. **Initialize Terraform**
   ```bash
   terraform init
   ```

3. **Review Configuration**
   ```bash
   # Review terraform.tfvars for your environment
   cat terraform.tfvars
   ```

4. **Plan the Deployment**
   ```bash
   terraform plan -var-file="terraform.tfvars"
   ```

5. **Apply the Configuration**
   ```bash
   terraform apply -var-file="terraform.tfvars"
   ```

6. **Get Cluster Access**
   ```bash
   # Get the master node IP
   terraform output -json | jq -r '.master_public_ips.value[0]'
   
   # SSH to master node
   ssh -i ~/.ssh/k8s-cluster-key ubuntu@<MASTER_IP>
   
   # Verify cluster status
   kubectl get nodes
   ```

## 📋 Configuration Reference

### Core Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `cluster_name` | Name of the Kubernetes cluster | `k8s-cluster` |
| `region` | AWS region for deployment | `us-west-2` |
| `master_count` | Number of master nodes | `1` |
| `worker_count` | Number of worker nodes | `2` |
| `network_plugin` | CNI plugin (flannel/cilium/calico) | `flannel` |
| `kubernetes_version` | Kubernetes version | `1.30.0` |

### Networking Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `vpc_cidr` | VPC CIDR block | `10.0.0.0/16` |
| `pod_network_cidr` | Pod network CIDR | `10.244.0.0/16` |
| `public_subnet_cidrs` | Public subnet CIDRs | `["10.0.1.0/24", "10.0.2.0/24"]` |
| `private_subnet_cidrs` | Private subnet CIDRs | `["10.0.11.0/24", "10.0.12.0/24"]` |

### Compute Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `master_instance_type` | Master node instance type | `t3.medium` |
| `worker_instance_type` | Worker node instance type | `t3.large` |
| `root_volume_size` | Root volume size (GB) | `30` |
| `data_volume_size` | Worker data volume size (GB) | `50` |

## 🔧 Network Plugin Configuration

### Flannel (Default)
- Simple overlay network
- Good for development and testing
- Uses VXLAN encapsulation
- Pod CIDR: `10.244.0.0/16`

### Cilium
- High-performance CNI
- Advanced networking and security
- eBPF-based networking
- Better for production workloads

### Calico
- Network policy enforcement
- BGP routing support
- Enterprise-grade features
- IPAM capabilities

## 📊 Scaling Operations

### Scaling Up Worker Nodes

1. **Update Configuration**
   ```bash
   # Edit terraform.tfvars
   worker_count = 10  # Increase from current value
   ```

2. **Apply Changes**
   ```bash
   terraform plan -var-file="terraform.tfvars"
   terraform apply -var-file="terraform.tfvars"
   ```

3. **Verify New Nodes**
   ```bash
   kubectl get nodes
   ```

### Scaling Up Master Nodes (High Availability)

1. **Update Configuration**
   ```bash
   master_count = 3  # For HA setup
   ```

2. **Apply Changes**
   ```bash
   terraform apply -var-file="terraform.tfvars"
   ```

3. **Configure Load Balancer**
   - Set up external load balancer for API server
   - Update kubeconfig to use load balancer endpoint

## 🔄 Image Update Workflow

### Update Kubernetes Version

1. **Update Version Variable**
   ```bash
   # In terraform.tfvars
   kubernetes_version = "1.30.1"
   ```

2. **Update Cloud-init Templates**
   - Edit `cloud-init/master.yaml.tpl` and `cloud-init/worker.yaml.tpl`
   - Update package versions in scripts

3. **Apply Changes**
   ```bash
   terraform apply -var-file="terraform.tfvars"
   ```

### Update Container Images

1. **Check Current Images**
   ```bash
   kubectl get pods -A -o jsonpath='{.items[*].spec.containers[*].image}' | tr ' ' '\n' | sort -u
   ```

2. **Update Images**
   ```bash
   # For system images, update cloud-init templates
   # For application images, update your deployments
   kubectl set image deployment/app app=new-image:tag
   ```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Nodes Not Joining Cluster
```bash
# Check join token validity
kubectl get tokens

# Generate new join token
kubeadm token create --print-join-command

# Check network connectivity
ping <MASTER_IP>
telnet <MASTER_IP> 6443
```

#### 2. Pod Network Issues
```bash
# Check CNI pods
kubectl get pods -n kube-system -l k8s-app=<network-plugin>

# Check network policies
kubectl get networkpolicies --all-namespaces

# Test pod connectivity
kubectl run test-pod --image=busybox --rm -it -- /bin/sh
```

#### 3. Storage Issues
```bash
# Check PV/PV status
kubectl get pv,pvc --all-namespaces

# Check EBS volumes
aws ec2 describe-volumes --filters Name=tag:Name,Values=<cluster-name>-*
```

#### 4. SSH Access Issues
```bash
# Check security group rules
aws ec2 describe-security-groups --group-ids <sg-id>

# Check key pair
aws ec2 describe-key-pairs --key-names <key-name>

# Debug SSH connection
ssh -v -i ~/.ssh/<key> ubuntu@<node-ip>
```

### Log Collection

#### Master Node Logs
```bash
# Kubernetes components
sudo journalctl -u kubelet
sudo journalctl -u kube-apiserver
sudo journalctl -u kube-controller-manager
sudo journalctl -u kube-scheduler

# Container runtime
sudo journalctl -u containerd
```

#### Worker Node Logs
```bash
# Kubelet
sudo journalctl -u kubelet

# CNI logs
sudo journalctl -u flannel  # or cilium/calico
```

### Recovery Procedures

#### Recover Failed Master
1. **Identify Failed Node**
   ```bash
   kubectl get nodes -o wide
   ```

2. **Remove from Cluster**
   ```bash
   kubectl delete node <failed-node-name>
   ```

3. **Terminate Instance**
   ```bash
   aws ec2 terminate-instances --instance-ids <instance-id>
   ```

4. **Recreate with Terraform**
   ```bash
   terraform apply -var-file="terraform.tfvars"
   ```

#### Recover Failed Worker
1. **Drain Node**
   ```bash
   kubectl drain <node-name> --ignore-daemonsets --delete-local-data
   ```

2. **Remove from Cluster**
   ```bash
   kubectl delete node <node-name>
   ```

3. **Recreate with Terraform**
   ```bash
   terraform apply -var-file="terraform.tfvars"
   ```

## 📈 Monitoring and Observability

### CloudWatch Integration
- Instance metrics enabled by default
- Custom metrics can be added via CloudWatch agent
- Logs can be shipped to CloudWatch Logs

### Kubernetes Monitoring
```bash
# Install metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Check node metrics
kubectl top nodes
kubectl top pods --all-namespaces
```

### Health Checks
```bash
# Cluster health
kubectl get componentstatuses
kubectl get cs

# Node health
kubectl get nodes -o wide

# Pod health
kubectl get pods --all-namespaces
```

## 🔒 Security Considerations

### Network Security
- Security groups restrict access to necessary ports only
- Private subnets for worker nodes when possible
- VPC flow logs for network monitoring

### Instance Security
- AMIs with latest security updates
- Encrypted EBS volumes
- SSH key-based authentication only

### Kubernetes Security
- RBAC enabled by default
- Network policies (when using Calico/Cilium)
- Pod Security Policies (if required)

## 💰 Cost Optimization

### Instance Types
- Use appropriate instance sizes for workloads
- Consider burstable instances (t3/t4) for dev/test
- Use Reserved Instances for production workloads

### Storage
- Use gp3 volumes for better performance/cost ratio
- Implement lifecycle policies for EBS snapshots
- Monitor and clean up unused volumes

### Monitoring Costs
```bash
# Check AWS costs
aws ce get-cost-and-usage --time-period Start=<start-date>,End=<end-date> --granularity MONTHLY

# Monitor resource utilization
kubectl top nodes
kubectl top pods --all-namespaces
```

## 🧪 Development and Testing

### Local Development
```bash
# Use local backend for state
terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}

# Test with smaller instance counts
master_count = 1
worker_count = 1
```

### Testing Workflow
1. **Plan and Review**
   ```bash
   terraform plan -detailed-exitcode -var-file="terraform.tfvars"
   ```

2. **Validate Configuration**
   ```bash
   terraform validate
   terraform fmt -check
   ```

3. **Security Scanning**
   ```bash
   # Install tfsec
   curl -sfL https://raw.githubusercontent.com/aquasecurity/tfsec/master/scripts/install_linux.sh | sh -s -- -b /usr/local/bin
   
   # Run security scan
   tfsec .
   ```

## 📚 Additional Resources

### Documentation
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubeadm Reference](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)

### Community
- [Kubernetes Slack](https://kubernetes.slack.com/)
- [Terraform Community](https://discuss.hashicorp.com/c/terraform-core)
- [AWS Containers](https://aws.amazon.com/containers/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For issues and questions:

1. Check the troubleshooting section above
2. Search existing GitHub issues
3. Create a new issue with detailed information
4. Include logs and configuration details

---

**Note**: This infrastructure is designed for educational and development purposes. For production use, consider additional security hardening, monitoring, and backup strategies.
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
