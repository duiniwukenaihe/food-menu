# Kubernetes Cluster Deployment Guide

## 🚀 Quick Deployment Steps

### 1. Prerequisites

Ensure you have the following installed and configured:
- [Terraform >= 1.0](https://learn.hashicorp.com/tutorials/terraform/install-cli)
- [AWS CLI](https://aws.amazon.com/cli/) with configured credentials
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [jq](https://stedolan.github.io/jq/download/)

### 2. Initial Setup

```bash
# Navigate to the infrastructure directory
cd infrastructure/terraform

# Copy the example configuration
cp terraform.tfvars.example terraform.tfvars

# Edit the configuration with your specific values
nano terraform.tfvars
```

### 3. Deploy the Cluster

```bash
# Initialize Terraform
terraform init

# Review the execution plan
terraform plan -var-file="terraform.tfvars"

# Apply the configuration
terraform apply -var-file="terraform.tfvars"
```

### 4. Access the Cluster

```bash
# Get the master node IP
MASTER_IP=$(terraform output -json master_public_ips | jq -r '.value[0]')

# SSH to the master node
ssh -i ~/.ssh/k8s-cluster-key ubuntu@$MASTER_IP

# Verify cluster status
kubectl get nodes
kubectl get pods --all-namespaces
```

### 5. Configure Local Access

```bash
# Copy kubeconfig to your local machine
scp -i ~/.ssh/k8s-cluster-key ubuntu@$MASTER_IP:/home/ubuntu/.kube/config ~/.kube/config-k8s-cluster

# Update the server endpoint
sed -i.bak "s/127.0.0.1:6443/$MASTER_IP:6443/g" ~/.kube/config-k8s-cluster

# Set the KUBECONFIG environment variable
export KUBECONFIG=$HOME/.kube/config-k8s-cluster

# Verify access from local machine
kubectl get nodes
```

## 📋 Configuration Options

### Network Plugin Selection

Choose your CNI plugin by setting the `network_plugin` variable:

```hcl
# Flannel (Default) - Simple and lightweight
network_plugin = "flannel"

# Cilium - High performance with eBPF
network_plugin = "cilium"

# Calico - Advanced networking and security
network_plugin = "calico"
```

### Instance Sizing

```hcl
# Development environment
master_instance_type = "t3.small"
worker_instance_type = "t3.medium"
master_count = 1
worker_count = 2

# Production environment
master_instance_type = "t3.medium"
worker_instance_type = "t3.large"
master_count = 3
worker_count = 5
```

### Storage Configuration

```hcl
# Root volumes for OS
root_volume_size = 30
root_volume_type = "gp3"

# Data volumes for worker nodes
data_volume_size = 100
data_volume_type = "gp3"
```

## 🔧 Advanced Configuration

### High Availability Setup

For a production HA setup:

```hcl
master_count = 3
worker_count = 5

# Use larger instances for HA
master_instance_type = "t3.medium"
worker_instance_type = "t3.large"

# Enable monitoring
enable_monitoring = true
```

### Custom Networking

```hcl
# Custom VPC CIDR
vpc_cidr = "10.100.0.0/16"

# Custom pod network
pod_network_cidr = "10.244.0.0/16"

# Multiple AZs for HA
availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
```

## 🛠️ Management Commands

### Using the Helper Script

```bash
# Deploy cluster
./cluster.sh deploy

# Check cluster status
./cluster.sh status

# SSH to master
./cluster.sh ssh

# Get kubeconfig
./cluster.sh kubeconfig

# Scale workers to 10 nodes
./cluster.sh scale 3 10

# Destroy cluster
./cluster.sh destroy
```

### Using Make Commands

```bash
# Quick deployment
make quick-deploy

# Check status
make status

# Get kubeconfig
make get-kubeconfig

# Scale workers
make scale-workers N=10

# Security scan
make security-scan

# Destroy everything
make quick-destroy
```

## 🔍 Verification Steps

### 1. Verify Infrastructure

```bash
# Check Terraform outputs
terraform output

# Verify AWS resources
aws ec2 describe-instances --filters Name=tag:Name,Values=k8s-prod-cluster-*
aws ec2 describe-vpcs --filters Name=tag:Name,Values=k8s-prod-cluster-vpc
```

### 2. Verify Kubernetes Cluster

```bash
# Check node status
kubectl get nodes -o wide

# Check system pods
kubectl get pods --all-namespaces

# Check cluster info
kubectl cluster-info

# Check component status
kubectl get componentstatuses
```

### 3. Verify Network Plugin

```bash
# For Flannel
kubectl get pods -n kube-system -l app=flannel

# For Cilium
kubectl get pods -n kube-system -l app.kubernetes.io/name=cilium

# For Calico
kubectl get pods -n kube-system -l k8s-app=calico-node
```

### 4. Test Workload Deployment

```bash
# Deploy a test application
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port=80 --type=NodePort

# Check the deployment
kubectl get pods
kubectl get services

# Test connectivity
NODE_PORT=$(kubectl get service nginx -o jsonpath='{.spec.ports[0].nodePort}')
curl http://$MASTER_IP:$NODE_PORT
```

## 🚨 Troubleshooting

### Common Issues

1. **Nodes Not Joining**
   ```bash
   # Check join token validity
   kubeadm token list
   
   # Generate new token
   kubeadm token create --print-join-command
   ```

2. **Network Plugin Issues**
   ```bash
   # Check CNI pods
   kubectl get pods -n kube-system
   
   # Check network policies
   kubectl get networkpolicies --all-namespaces
   ```

3. **SSH Access Issues**
   ```bash
   # Check security group rules
   aws ec2 describe-security-groups --group-ids $(terraform output -raw security_group_id)
   
   # Check key pair
   aws ec2 describe-key-pairs --key-names k8s-cluster-key
   ```

### Log Collection

```bash
# Master node logs
ssh -i ~/.ssh/k8s-cluster-key ubuntu@$MASTER_IP
sudo journalctl -u kubelet -f
sudo journalctl -u kube-apiserver -f

# Worker node logs
WORKER_IP=$(terraform output -json worker_public_ips | jq -r '.value[0]')
ssh -i ~/.ssh/k8s-cluster-key ubuntu@$WORKER_IP
sudo journalctl -u kubelet -f
```

## 📊 Monitoring and Observability

### CloudWatch Metrics

```bash
# Enable detailed monitoring
enable_monitoring = true

# View metrics in AWS Console
# EC2 → Instances → Select Instance → Monitoring
```

### Kubernetes Monitoring

```bash
# Install metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Check resource usage
kubectl top nodes
kubectl top pods --all-namespaces
```

## 💰 Cost Optimization

### Right-Sizing Instances

```bash
# Monitor resource usage
kubectl top nodes
kubectl describe nodes

# Adjust instance types accordingly
master_instance_type = "t3.medium"  # Based on actual usage
worker_instance_type = "t3.large"  # Based on actual usage
```

### Storage Optimization

```bash
# Use gp3 for better price/performance
root_volume_type = "gp3"
data_volume_type = "gp3"

# Monitor disk usage
ssh ubuntu@$MASTER_IP
df -h
```

## 🔄 Scaling Operations

### Scaling Workers

```bash
# Update terraform.tfvars
worker_count = 10

# Apply changes
terraform apply -var-file="terraform.tfvars"

# Verify new nodes
kubectl get nodes
```

### Scaling Masters (HA)

```bash
# Update terraform.tfvars
master_count = 3

# Apply changes
terraform apply -var-file="terraform.tfvars"

# Configure external load balancer for API server
# (Manual step required for HA setup)
```

## 🧹 Cleanup

### Destroy Cluster

```bash
# Using helper script
./cluster.sh destroy

# Using Terraform directly
terraform destroy -var-file="terraform.tfvars"

# Using Make
make quick-destroy
```

### Verify Cleanup

```bash
# Check for remaining resources
aws ec2 describe-instances --filters Name=tag:Project,Values=kubernetes-cluster
aws ec2 describe-vpcs --filters Name=tag:Project,Values=kubernetes-cluster

# Clean up local files
rm ~/.kube/config-k8s-cluster
```

## 📚 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [kubeadm Reference](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)
- [AWS EC2 Documentation](https://docs.aws.amazon.com/ec2/)

---

For support and questions, please refer to the main README.md file or create an issue in the repository.