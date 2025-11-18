# Kubernetes on Proxmox with Terraform

A complete Infrastructure as Code (IaC) solution for deploying Kubernetes clusters on Proxmox VE using Terraform. This project provides a production-ready, scalable, and maintainable way to manage Kubernetes infrastructure.

## 🚀 Features

- **Complete Automation**: End-to-end Kubernetes cluster deployment
- **Proxmox Integration**: Native Proxmox VE API integration with SSH fallback
- **Modular Design**: Clean separation of concerns with reusable modules
- **Cloud-Init Ready**: Automated node provisioning and configuration
- **Multiple Network Plugins**: Support for Flannel, Calico, and Cilium
- **Flexible Configuration**: Comprehensive variable support with validation
- **Production Ready**: Includes security, monitoring, and best practices
- **Cluster Scaling**: Easy horizontal scaling of master and worker nodes

## 📋 Prerequisites

### System Requirements

- Proxmox VE 7.0+ with API access
- Ubuntu 22.04 LTS cloud image support
- Terraform 1.0+
- SSH key access to Proxmox host
- Sufficient resources for planned cluster size

### Required Software

```bash
# Install Terraform (Ubuntu/Debian)
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform

# Verify installation
terraform version
```

### Proxmox Configuration

1. **Enable API Access**:
   - Create Proxmox user with appropriate permissions
   - Generate API token or use username/password authentication
   - Configure network access to Proxmox API (port 8006)

2. **SSH Access**:
   - Enable SSH access to Proxmox host
   - Configure SSH key authentication
   - Verify user permissions for VM management

## 🏗️ Project Structure

```
infrastructure/terraform/
├── main.tf                          # Main orchestration file
├── variables.tf                      # Variable definitions
├── outputs.tf                       # Output definitions
├── versions.tf                      # Terraform and provider versions
├── providers.tf                     # Provider configurations
├── terraform.tfvars                 # User configuration (actual values)
├── terraform.tfvars.example         # Configuration template
├── modules/                         # Reusable modules
│   ├── template/                    # VM template creation
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── templates/
│   │       └── cloud-init.yaml.tpl
│   └── kubernetes_node_pool/        # Kubernetes node management
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       └── templates/
│           ├── master.yaml.tpl
│           └── worker.yaml.tpl
├── cloud-init/                      # Cloud-init templates
│   ├── master.yaml.tpl
│   ├── worker.yaml.tpl
│   └── common.sh
├── scripts/                         # Utility scripts
│   ├── get-ubuntu-cloudimg.sh        # Image download with verification
│   ├── kubeadm-init.sh              # Master initialization
│   ├── kubeadm-join.sh              # Worker node join
│   └── network-plugin-install.sh     # Network plugin setup
└── README.md                        # This documentation
```

## ⚙️ Configuration

### Quick Start

1. **Copy Configuration Template**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. **Edit Configuration**:
   Edit `terraform.tfvars` with your specific values:

   ```hcl
   # Proxmox Configuration
   proxmox_api_url = "https://192.168.0.200:8006/api2/json"
   proxmox_user = "root@pam"
   proxmox_password = "your_password_here"
   proxmox_node = "proxmox1"
   proxmox_host = "192.168.0.200"
   
   # Node Configuration
   master_count = 1
   worker_count = 2
   master_cores = 4
   worker_cores = 8
   master_memory = 8192
   worker_memory = 16384
   
   # Kubernetes Configuration
   k8s_version = "1.30.0"
   network_plugin = "flannel"
   pod_network_cidr = "10.244.0.0/16"
   ```

### Environment Variables

For enhanced security, you can use environment variables:

```bash
# Proxmox Authentication
export PROXMOX_VE_PASSWORD="your_password"
export PROXMOX_VE_API_TOKEN="user@realm!tokenid=secret"

# SSH Configuration
export TF_VAR_proxmox_ssh_key="~/.ssh/proxmox_key"
```

## 🚀 Deployment Guide

### 1. Initialize Terraform

```bash
cd infrastructure/terraform
terraform init
```

### 2. Plan Deployment

```bash
terraform plan
```

### 3. Deploy Cluster

```bash
terraform apply
```

### 4. Verify Deployment

```bash
# Check outputs
terraform output

# SSH to master node
ssh -i ~/.ssh/id_rsa ubuntu@$(terraform output -json | jq -r '.master_ips.value[0]')

# Check cluster status
kubectl get nodes
kubectl get pods --all-namespaces
```

## 📊 Cluster Management

### Scaling Operations

#### Add Worker Nodes

Update `terraform.tfvars`:
```hcl
worker_count = 5  # Increase from current count
```

Apply changes:
```bash
terraform apply
```

#### Add Master Nodes (HA)

Update `terraform.tfvars`:
```hcl
master_count = 3  # For high availability
```

Apply changes:
```bash
terraform apply
```

### Node Management

#### Drain Node for Maintenance
```bash
# Get node name
kubectl get nodes

# Drain node
kubectl drain <node-name> --ignore-daemonsets --delete-emptydir-data

# Maintenance operations...

# Uncordon node
kubectl uncordon <node-name>
```

#### Remove Node
```bash
# Drain node first
kubectl drain <node-name> --ignore-daemonsets --delete-emptydir-data

# Delete node
kubectl delete node <node-name>

# Update Terraform configuration
# Decrease master_count or worker_count
terraform apply
```

## 🔧 Network Configuration

### Supported Network Plugins

#### Flannel (Default)
- Simple, lightweight overlay network
- Good for development and testing
- VXLAN backend with UDP port 8472

#### Calico
- Production-ready network policy
- BGP routing support
- Advanced security features

#### Cilium
- eBPF-based networking
- High performance
- Advanced observability

### Network Plugin Selection

Edit `terraform.tfvars`:
```hcl
network_plugin = "calico"  # Options: flannel, calico, cilium
```

## 🔒 Security Considerations

### Authentication Methods

#### API Token Authentication (Recommended)
```hcl
proxmox_api_token_id = "root@pam!token_id"
proxmox_api_token_secret = "your_token_secret"
```

#### Username/Password Authentication
```hcl
proxmox_user = "root@pam"
proxmox_password = "your_password"
```

### SSH Key Management

Generate SSH keys:
```bash
ssh-keygen -t rsa -b 4096 -C "terraform-proxmox"
```

Configure in `terraform.tfvars`:
```hcl
ssh_private_key_path = "~/.ssh/terraform_proxmox"
ssh_public_key_path = "~/.ssh/terraform_proxmox.pub"
```

### Network Security

- Use VLANs for network isolation
- Configure firewall rules appropriately
- Enable TLS verification for production
- Use dedicated service accounts

## 📈 Monitoring and Logging

### Cluster Monitoring

Access monitoring dashboards:
```bash
# Port forward to access dashboards
kubectl port-forward svc/grafana 3000:3000
kubectl port-forward svc/prometheus 9090:9090
```

### Log Collection

View cluster logs:
```bash
# System logs
journalctl -u kubelet -f

# Kubernetes logs
kubectl logs -f deployment/coredns -n kube-system

# Node logs
kubectl logs -f -l k8s-app=calico-node -n kube-system
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Proxmox API Connection Failed

**Symptoms**: 
```
Error: failed to create proxmox provider: POST https://192.168.0.200:8006/api2/json/tokens...
```

**Solutions**:
- Verify Proxmox API URL and credentials
- Check network connectivity to Proxmox host
- Ensure user has sufficient permissions
- Verify TLS certificate settings

#### 2. VM Creation Failed

**Symptoms**:
```
Error: error creating VM: VM creation failed
```

**Solutions**:
- Check available resources on Proxmox node
- Verify storage pool availability
- Ensure VM ID range is not in use
- Check network bridge configuration

#### 3. Kubernetes Cluster Initialization Failed

**Symptoms**:
```
[kubelet-check] It seems like the kubelet isn't running or healthy
```

**Solutions**:
- Check containerd service status
- Verify swap is disabled
- Check kernel modules are loaded
- Review sysctl configuration

#### 4. Network Plugin Issues

**Symptoms**:
```
Network plugin is not ready: cni config uninitialized
```

**Solutions**:
- Verify pod network CIDR doesn't conflict
- Check network plugin pod status
- Review CNI configuration
- Ensure firewall allows required ports

### Debug Commands

#### Proxmox Issues
```bash
# Check Proxmox API
curl -k -X POST "https://192.168.0.200:8006/api2/json/access/ticket" \
  -d "username=root@pam&password=your_password"

# Check VM status
qm list
```

#### Kubernetes Issues
```bash
# Check cluster status
kubectl get nodes -o wide
kubectl get pods --all-namespaces

# Check system services
systemctl status kubelet
systemctl status containerd
systemctl status docker

# Check logs
journalctl -u kubelet -f
kubectl describe pod <pod-name> -n <namespace>
```

## 🔄 Maintenance Operations

### Image Updates

Update Ubuntu cloud image:
```bash
# Download new image
./scripts/get-ubuntu-cloudimg.sh

# Update template
terraform apply -replace=module.template
```

### Kubernetes Upgrades

Upgrade Kubernetes version:
```hcl
k8s_version = "1.30.1"  # Update to desired version
```

Apply upgrade:
```bash
terraform apply
```

### Backup and Recovery

#### Backup Cluster Configuration
```bash
# Backup etcd
kubectl get nodes
kubectl get pods --all-namespaces

# Save Terraform state
cp terraform.tfstate terraform.tfstate.backup
```

#### Restore from Backup
```bash
# Restore Terraform state
cp terraform.tfstate.backup terraform.tfstate

# Verify cluster state
terraform plan
```

## 📚 Advanced Configuration

### Custom Node Pools

Define specific node configurations:
```hcl
# High-performance workers
worker_cores = 16
worker_memory = 32768
worker_disk_size = "200G"

# GPU-enabled nodes (if supported)
# Additional configuration in modules
```

### Multi-Node Proxmox Clusters

Configure node mapping:
```hcl
# Distribute VMs across Proxmox nodes
proxmox_node_map = {
  "master-1" = "proxmox1"
  "worker-1" = "proxmox2"
  "worker-2" = "proxmox3"
}
```

### Storage Configuration

Configure multiple storage pools:
```hcl
# Different storage for different purposes
master_storage = "fast-ssd"
worker_storage = "large-hdd"
template_storage = "local"
```

## 🔗 Integration Examples

### CI/CD Pipeline Integration

```yaml
# GitLab CI example
deploy_k8s:
  stage: deploy
  script:
    - cd infrastructure/terraform
    - terraform init
    - terraform apply -auto-approve
  environment: production
```

### Ansible Integration

```yaml
# Ansible playbook example
- name: Deploy Kubernetes cluster
  hosts: localhost
  tasks:
    - name: Apply Terraform configuration
      community.general.terraform:
        project_path: infrastructure/terraform
        state: present
        variables_file: terraform.tfvars
```

## 📖 API Reference

### Terraform Outputs

Key outputs available after deployment:

```bash
# Cluster information
terraform output cluster_name
terraform output kubernetes_version
terraform output cluster_endpoint

# Node information
terraform output master_nodes
terraform output worker_nodes

# SSH access
terraform output ssh_access
terraform output connect_commands
```

### Module Interfaces

#### Template Module
- **Inputs**: Proxmox config, image URL, storage settings
- **Outputs**: Template ID, template name, storage location

#### Kubernetes Node Pool Module
- **Inputs**: Template ID, node counts, resource specifications
- **Outputs**: VM details, IP addresses, node names

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone repository
git clone <repository-url>
cd infrastructure/terraform

# Install development dependencies
terraform init
terraform fmt -check
terraform validate

# Run tests (if available)
terraform test
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- **Issues**: [GitHub Issues](https://github.com/your-repo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-repo/discussions)
- **Documentation**: [Wiki](https://github.com/your-repo/wiki)

## 🗺️ Roadmap

- [ ] Add Helm chart deployment support
- [ ] Implement backup and restore automation
- [ ] Add monitoring stack (Prometheus, Grafana)
- [ ] Support for additional network plugins
- [ ] Multi-cloud deployment support
- [ ] GitOps integration (ArgoCD, Flux)
- [ ] Automated security scanning
- [ ] Cost optimization features

---

**Happy Kubernetes clustering! 🎉**