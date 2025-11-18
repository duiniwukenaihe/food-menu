#!/bin/bash
# =============================================================================
# Kubernetes Worker Node Join Script
# =============================================================================

set -e

# Configuration
K8S_VERSION="${1:-1.30.0}"
MASTER_IP="${2:-k8s-cluster-master-1}"
JOIN_COMMAND_FILE="${3:-/tmp/join-command.sh}"

echo "=== Kubernetes Worker Node Join ==="
echo "Kubernetes Version: $K8S_VERSION"
echo "Master IP: $MASTER_IP"
echo "Join Command File: $JOIN_COMMAND_FILE"
echo

# Function to wait for service
wait_for_service() {
    local service=$1
    local timeout=$2
    echo "Waiting for $service to be ready..."
    
    for i in $(seq 1 $timeout); do
        if systemctl is-active --quiet $service; then
            echo "$service is ready"
            return 0
        fi
        echo "Attempt $i/$timeout: $service not ready yet..."
        sleep 10
    done
    
    echo "Timeout waiting for $service"
    return 1
}

# Function to install Kubernetes components
install_kubernetes() {
    echo "Installing Kubernetes components..."
    
    # Add Kubernetes repository
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v$K8S_VERSION/deb/Release.key | \
        sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    
    echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v$K8S_VERSION/deb/ /" | \
        sudo tee /etc/apt/sources.list.d/kubernetes.list
    
    # Update package list
    sudo apt update
    
    # Install Kubernetes packages
    sudo apt install -y kubelet=$K8S_VERSION-1.1.1 kubeadm=$K8S_VERSION-1.1.1 kubectl=$K8S_VERSION-1.1.1
    
    # Hold packages to prevent upgrades
    sudo apt-mark hold kubelet kubeadm kubectl
    
    echo "Kubernetes components installed successfully"
}

# Function to get join command from master
get_join_command() {
    echo "Attempting to get join command from master..."
    
    # Try to get join command from master via SSH
    if command -v ssh >/dev/null 2>&1; then
        for i in {1..10}; do
            if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 -o BatchMode=yes \
                ubuntu@$MASTER_IP "test -f /tmp/join-command.sh" 2>/dev/null; then
                JOIN_COMMAND=$(ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 -o BatchMode=yes \
                    ubuntu@$MASTER_IP "cat /tmp/join-command.sh" 2>/dev/null)
                if [ -n "$JOIN_COMMAND" ]; then
                    echo "✓ Retrieved join command from master"
                    echo "$JOIN_COMMAND"
                    return 0
                fi
            fi
            echo "Attempt $i/10: Master not ready yet, waiting 30 seconds..."
            sleep 30
        done
    fi
    
    echo "Could not retrieve join command from master"
    return 1
}

# Function to join cluster manually
join_cluster_manual() {
    echo "Manual join required. Please follow these steps:"
    echo
    echo "1. SSH to the master node:"
    echo "   ssh ubuntu@$MASTER_IP"
    echo
    echo "2. Generate join command:"
    echo "   sudo kubeadm token create --print-join-command"
    echo
    echo "3. Copy the output and run it on this worker node with sudo"
    echo
    echo "4. After joining, verify the node:"
    echo "   kubectl get nodes"
    echo
    echo "Press Enter when you have completed the manual join..."
    read -r
}

# Function to join cluster
join_cluster() {
    echo "Joining Kubernetes cluster..."
    
    if [ -n "$JOIN_COMMAND" ]; then
        echo "Using provided join command:"
        echo "$JOIN_COMMAND"
        sudo $JOIN_COMMAND
        
        echo "✓ Worker node joined successfully"
    else
        echo "No join command provided, attempting to retrieve from master..."
        
        if JOIN_COMMAND=$(get_join_command); then
            echo "Retrieved join command from master:"
            echo "$JOIN_COMMAND"
            sudo $JOIN_COMMAND
            
            echo "✓ Worker node joined successfully"
        else
            echo "Failed to retrieve join command from master"
            join_cluster_manual
        fi
    fi
}

# Function to verify node status
verify_join() {
    echo "Verifying node join status..."
    
    # Wait a bit for the node to register
    sleep 30
    
    # Check if kubelet is running
    if systemctl is-active --quiet kubelet; then
        echo "✓ Kubelet is running"
    else
        echo "✗ Kubelet is not running"
        return 1
    fi
    
    # Check node status (this requires kubectl access)
    if command -v kubectl >/dev/null 2>&1 && [ -f "$HOME/.kube/config" ]; then
        echo "Checking node status with kubectl..."
        kubectl get nodes
    else
        echo "kubectl not available or not configured"
        echo "To verify node status, run from master:"
        echo "kubectl get nodes"
    fi
    
    echo "Worker node setup completed"
}

# Main execution
main() {
    echo "Starting Kubernetes worker node setup..."
    
    # Wait for containerd
    wait_for_service "containerd" 30
    
    # Install Kubernetes components
    install_kubernetes
    
    # Enable and start kubelet
    sudo systemctl enable kubelet
    sudo systemctl start kubelet
    
    # Join cluster
    join_cluster
    
    # Verify join
    verify_join
    
    echo
    echo "=== Kubernetes Worker Setup Complete ==="
    echo "Worker node has joined the cluster"
    echo "Verify with: kubectl get nodes (from master)"
}

# Parse command line arguments
JOIN_COMMAND=""
while [[ $# -gt 0 ]]; do
    case $1 in
        --join-command)
            JOIN_COMMAND="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--join-command 'sudo kubeadm join ...']"
            exit 1
            ;;
    esac
done

# Run main function
main