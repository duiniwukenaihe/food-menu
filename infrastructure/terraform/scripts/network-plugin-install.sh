#!/bin/bash
# =============================================================================
# Network Plugin Installation Script
# =============================================================================

set -e

# Configuration
NETWORK_PLUGIN="${1:-flannel}"
POD_NETWORK_CIDR="${2:-10.244.0.0/16}"

echo "=== Network Plugin Installation ==="
echo "Network Plugin: $NETWORK_PLUGIN"
echo "Pod Network CIDR: $POD_NETWORK_CIDR"
echo

# Function to install Flannel
install_flannel() {
    echo "Installing Flannel CNI..."
    
    # Download and apply Flannel manifest
    kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
    
    echo "Flannel CNI installed successfully"
}

# Function to install Calico
install_calico() {
    echo "Installing Calico CNI..."
    
    # Install Calico operator
    kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/tigera-operator.yaml
    
    # Create custom resources for Calico
    cat <<EOF | kubectl apply -f -
apiVersion: operator.tigera.io/v1
kind: Installation
metadata:
  name: default
spec:
  calicoNetwork:
    ipPools:
    - blockSize: 26
      cidr: $POD_NETWORK_CIDR
      encapsulation: VXLANCrossSubnet
      natOutgoing: Enabled
      nodeSelector: all()
EOF
    
    echo "Calico CNI installed successfully"
}

# Function to install Cilium
install_cilium() {
    echo "Installing Cilium CNI..."
    
    # Install Cilium using quick-install
    kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.14.0/install/kubernetes/quick-install.yaml
    
    echo "Cilium CNI installed successfully"
}

# Function to verify installation
verify_installation() {
    echo "Verifying network plugin installation..."
    
    # Wait for pods to be ready
    echo "Waiting for network plugin pods to be ready..."
    
    case $NETWORK_PLUGIN in
        "flannel")
            kubectl wait --for=condition=ready pod -l app=flannel -n kube-flannel --timeout=300s
            ;;
        "calico")
            kubectl wait --for=condition=ready pod -l k8s-app=calico-node --timeout=300s
            kubectl wait --for=condition=ready pod -l k8s-app=calico-kube-controllers --timeout=300s
            ;;
        "cilium")
            kubectl wait --for=condition=ready pod -l k8s-app=cilium --timeout=300s
            ;;
    esac
    
    # Show pod status
    echo "Network plugin pod status:"
    kubectl get pods --all-namespaces -l k8s-app=flannel,k8s-app=calico-node,k8s-app=calico-kube-controllers,k8s-app=cilium
    
    echo "Network plugin verification completed"
}

# Function to show network info
show_network_info() {
    echo "Network Information:"
    echo "Plugin: $NETWORK_PLUGIN"
    echo "Pod CIDR: $POD_NETWORK_CIDR"
    echo
    
    case $NETWORK_PLUGIN in
        "flannel")
            echo "Flannel Configuration:"
            echo "  - Backend: VXLAN"
            echo "  - Port: 8472/UDP"
            echo "  - Network: $POD_NETWORK_CIDR"
            ;;
        "calico")
            echo "Calico Configuration:"
            echo "  - Backend: VXLAN"
            echo "  - Network: $POD_NETWORK_CIDR"
            echo "  - IPAM: Calico IPAM"
            ;;
        "cilium")
            echo "Cilium Configuration:"
            echo "  - Backend: VXLAN"
            echo "  - Network: $POD_NETWORK_CIDR"
            echo "  - Features: BPF, Hubble"
            ;;
    esac
}

# Function to check prerequisites
check_prerequisites() {
    echo "Checking prerequisites..."
    
    # Check if kubectl is available
    if ! command -v kubectl >/dev/null 2>&1; then
        echo "Error: kubectl is not installed or not in PATH"
        exit 1
    fi
    
    # Check if cluster is accessible
    if ! kubectl cluster-info >/dev/null 2>&1; then
        echo "Error: Cannot access Kubernetes cluster"
        echo "Please ensure kubectl is configured correctly"
        exit 1
    fi
    
    # Check if cluster is already initialized
    if ! kubectl get nodes >/dev/null 2>&1; then
        echo "Error: Cannot get cluster nodes"
        echo "Please ensure the cluster is initialized"
        exit 1
    fi
    
    echo "✓ Prerequisites check passed"
}

# Main execution
main() {
    echo "Starting network plugin installation..."
    
    # Check prerequisites
    check_prerequisites
    
    # Install network plugin
    case $NETWORK_PLUGIN in
        "flannel")
            install_flannel
            ;;
        "calico")
            install_calico
            ;;
        "cilium")
            install_cilium
            ;;
        *)
            echo "Error: Unsupported network plugin: $NETWORK_PLUGIN"
            echo "Supported plugins: flannel, calico, cilium"
            exit 1
            ;;
    esac
    
    # Verify installation
    verify_installation
    
    # Show network info
    show_network_info
    
    echo
    echo "=== Network Plugin Installation Complete ==="
    echo "Network plugin $NETWORK_PLUGIN is installed and ready"
    echo "You can now deploy applications to the cluster"
}

# Show usage
usage() {
    echo "Usage: $0 [NETWORK_PLUGIN] [POD_NETWORK_CIDR]"
    echo
    echo "Arguments:"
    echo "  NETWORK_PLUGIN   Network plugin to install (flannel, calico, cilium)"
    echo "  POD_NETWORK_CIDR Pod network CIDR (default: 10.244.0.0/16)"
    echo
    echo "Examples:"
    echo "  $0 flannel"
    echo "  $0 calico 192.168.0.0/16"
    echo "  $0 cilium 10.244.0.0/16"
}

# Parse command line arguments
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
    exit 0
fi

# Run main function
main "$@"