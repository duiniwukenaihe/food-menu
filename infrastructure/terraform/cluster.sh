#!/bin/bash
# Cluster Management Script for Kubernetes Infrastructure

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TERRAFORM_DIR="$SCRIPT_DIR"
CLUSTER_NAME=${CLUSTER_NAME:-"k8s-cluster"}
SSH_KEY="$HOME/.ssh/${CLUSTER_NAME}-key"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if terraform is installed
check_terraform() {
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed. Please install Terraform first."
        exit 1
    fi
}

# Function to check if AWS CLI is configured
check_aws() {
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS CLI is not configured. Please run 'aws configure' first."
        exit 1
    fi
}

# Function to initialize terraform
init_terraform() {
    log_info "Initializing Terraform..."
    cd "$TERRAFORM_DIR"
    terraform init
    log_success "Terraform initialized successfully"
}

# Function to deploy the cluster
deploy_cluster() {
    log_info "Deploying Kubernetes cluster..."
    cd "$TERRAFORM_DIR"
    
    if [ ! -f "terraform.tfvars" ]; then
        log_warning "terraform.tfvars not found. Copying from example..."
        cp terraform.tfvars.example terraform.tfvars
        log_warning "Please edit terraform.tfvars with your configuration before deploying."
        exit 1
    fi
    
    terraform plan -var-file="terraform.tfvars"
    read -p "Do you want to proceed with the deployment? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        terraform apply -var-file="terraform.tfvars" -auto-approve
        log_success "Cluster deployed successfully"
    else
        log_info "Deployment cancelled"
    fi
}

# Function to destroy the cluster
destroy_cluster() {
    log_info "Destroying Kubernetes cluster..."
    cd "$TERRAFORM_DIR"
    
    read -p "Are you sure you want to destroy the cluster? This action cannot be undone. (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        terraform destroy -var-file="terraform.tfvars" -auto-approve
        log_success "Cluster destroyed successfully"
    else
        log_info "Destruction cancelled"
    fi
}

# Function to get cluster status
cluster_status() {
    log_info "Getting cluster status..."
    cd "$TERRAFORM_DIR"
    
    if ! terraform state list | grep -q "aws_instance"; then
        log_warning "No cluster found. Please deploy the cluster first."
        return 1
    fi
    
    echo "=== Cluster Information ==="
    terraform output cluster_name
    terraform output kubernetes_version
    terraform output network_plugin
    terraform output cluster_endpoint
    
    echo -e "\n=== Master Nodes ==="
    terraform output -json master_public_ips | jq -r '.value[]' | nl
    
    echo -e "\n=== Worker Nodes ==="
    terraform output -json worker_public_ips | jq -r '.value[]' | nl
    
    echo -e "\n=== VPC Information ==="
    terraform output vpc_id
    terraform output -json public_subnet_ids | jq -r '.value[]' | nl
    terraform output -json private_subnet_ids | jq -r '.value[]' | nl
}

# Function to SSH to master node
ssh_master() {
    log_info "Connecting to master node..."
    cd "$TERRAFORM_DIR"
    
    MASTER_IP=$(terraform output -json master_public_ips | jq -r '.value[0]')
    
    if [ "$MASTER_IP" == "null" ] || [ -z "$MASTER_IP" ]; then
        log_error "Master node IP not found. Please deploy the cluster first."
        exit 1
    fi
    
    log_info "Connecting to master node at $MASTER_IP"
    ssh -i "$SSH_KEY" ubuntu@"$MASTER_IP"
}

# Function to get kubeconfig
get_kubeconfig() {
    log_info "Getting kubeconfig..."
    cd "$TERRAFORM_DIR"
    
    MASTER_IP=$(terraform output -json master_public_ips | jq -r '.value[0]')
    
    if [ "$MASTER_IP" == "null" ] || [ -z "$MASTER_IP" ]; then
        log_error "Master node IP not found. Please deploy the cluster first."
        exit 1
    fi
    
    log_info "Copying kubeconfig from master node..."
    scp -i "$SSH_KEY" ubuntu@"$MASTER_IP":/home/ubuntu/.kube/config "$HOME/.kube/config-${CLUSTER_NAME}"
    
    # Update the server endpoint
    sed -i.bak "s/127.0.0.1:6443/$MASTER_IP:6443/g" "$HOME/.kube/config-${CLUSTER_NAME}"
    
    # Set KUBECONFIG environment variable
    export KUBECONFIG="$HOME/.kube/config-${CLUSTER_NAME}"
    
    log_success "Kubeconfig copied to $HOME/.kube/config-${CLUSTER_NAME}"
    log_info "To use this kubeconfig, run: export KUBECONFIG=\"$HOME/.kube/config-${CLUSTER_NAME}\""
}

# Function to scale cluster
scale_cluster() {
    local master_count=$1
    local worker_count=$2
    
    if [ -z "$master_count" ] && [ -z "$worker_count" ]; then
        log_error "Please specify master and/or worker count: scale <master_count> <worker_count>"
        exit 1
    fi
    
    log_info "Scaling cluster..."
    cd "$TERRAFORM_DIR"
    
    # Update terraform.tfvars
    if [ -n "$master_count" ]; then
        sed -i "s/master_count = .*/master_count = $master_count/" terraform.tfvars
        log_info "Setting master count to $master_count"
    fi
    
    if [ -n "$worker_count" ]; then
        sed -i "s/worker_count = .*/worker_count = $worker_count/" terraform.tfvars
        log_info "Setting worker count to $worker_count"
    fi
    
    terraform plan -var-file="terraform.tfvars"
    read -p "Do you want to proceed with scaling? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        terraform apply -var-file="terraform.tfvars" -auto-approve
        log_success "Cluster scaled successfully"
    else
        log_info "Scaling cancelled"
    fi
}

# Function to show logs
show_logs() {
    local node_type=${1:-"master"}
    local node_index=${2:-0}
    
    log_info "Getting logs from $node_type node $node_index..."
    cd "$TERRAFORM_DIR"
    
    if [ "$node_type" == "master" ]; then
        NODE_IP=$(terraform output -json master_public_ips | jq -r ".value[$node_index]")
    else
        NODE_IP=$(terraform output -json worker_public_ips | jq -r ".value[$node_index]")
    fi
    
    if [ "$NODE_IP" == "null" ] || [ -z "$NODE_IP" ]; then
        log_error "$node_type node $node_index not found."
        exit 1
    fi
    
    log_info "Showing kubelet logs from $node_type node $node_index ($NODE_IP)..."
    ssh -i "$SSH_KEY" ubuntu@"$NODE_IP" "sudo journalctl -u kubelet -f"
}

# Function to backup cluster
backup_cluster() {
    log_info "Creating cluster backup..."
    cd "$TERRAFORM_DIR"
    
    BACKUP_DIR="$HOME/k8s-backups/$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    # Backup Terraform state
    if [ -f "terraform.tfstate" ]; then
        cp terraform.tfstate "$BACKUP_DIR/"
        log_success "Terraform state backed up"
    fi
    
    # Backup kubeconfig
    get_kubeconfig
    cp "$HOME/.kube/config-${CLUSTER_NAME}" "$BACKUP_DIR/"
    
    # Backup cluster resources
    export KUBECONFIG="$HOME/.kube/config-${CLUSTER_NAME}"
    kubectl get all --all-namespaces -o yaml > "$BACKUP_DIR/cluster-resources.yaml"
    
    log_success "Cluster backup created at $BACKUP_DIR"
}

# Function to show help
show_help() {
    echo "Kubernetes Cluster Management Script"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  init          Initialize Terraform"
    echo "  deploy        Deploy the cluster"
    echo "  destroy       Destroy the cluster"
    echo "  status        Show cluster status"
    echo "  ssh           SSH to master node"
    echo "  kubeconfig    Get kubeconfig from master"
    echo "  scale <m> <w> Scale cluster (m=master count, w=worker count)"
    echo "  logs <t> <i>  Show logs (t=node type, i=node index)"
    echo "  backup        Create cluster backup"
    echo "  help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy"
    echo "  $0 scale 3 5"
    echo "  $0 logs worker 0"
    echo "  $0 ssh"
    echo ""
    echo "Environment Variables:"
    echo "  CLUSTER_NAME  Name of the cluster (default: k8s-cluster)"
}

# Main script logic
main() {
    case "${1:-help}" in
        init)
            check_terraform
            check_aws
            init_terraform
            ;;
        deploy)
            check_terraform
            check_aws
            deploy_cluster
            ;;
        destroy)
            check_terraform
            check_aws
            destroy_cluster
            ;;
        status)
            check_terraform
            cluster_status
            ;;
        ssh)
            check_terraform
            ssh_master
            ;;
        kubeconfig)
            check_terraform
            get_kubeconfig
            ;;
        scale)
            check_terraform
            check_aws
            scale_cluster "$2" "$3"
            ;;
        logs)
            check_terraform
            show_logs "$2" "$3"
            ;;
        backup)
            check_terraform
            backup_cluster
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"