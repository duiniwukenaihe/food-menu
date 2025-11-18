#!/bin/bash
# =============================================================================
# Kubernetes Master Initialization Script
# =============================================================================

set -e

# Configuration
K8S_VERSION="${1:-1.30.0}"
POD_NETWORK_CIDR="${2:-10.244.0.0/16}"
NETWORK_PLUGIN="${3:-flannel}"
IMAGE_REPOSITORY="${4:-registry.aliyuncs.com/google_containers}"

echo "=== Kubernetes Master Initialization ==="
echo "Kubernetes Version: $K8S_VERSION"
echo "Pod Network CIDR: $POD_NETWORK_CIDR"
echo "Network Plugin: $NETWORK_PLUGIN"
echo "Image Repository: $IMAGE_REPOSITORY"
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

# Function to initialize cluster
initialize_cluster() {
    echo "Initializing Kubernetes cluster..."
    
    # Create kubeadm config
    cat <<EOF | sudo tee /tmp/kubeadm-config.yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 0.0.0.0
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///var/run/containerd/containerd.sock
  kubeletExtraArgs:
    cgroup-driver: systemd
---
apiVersion: kubeadm.k8s.io/v1beta3
kind: ClusterConfiguration
kubernetesVersion: v$K8S_VERSION
imageRepository: $IMAGE_REPOSITORY
networking:
  serviceSubnet: "10.96.0.0/12"
  podSubnet: "$POD_NETWORK_CIDR"
  dnsDomain: "cluster.local"
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
serverTLSBootstrap: true
rotateCertificates: true
EOF

    # Initialize cluster
    sudo kubeadm init --config=/tmp/kubeadm-config.yaml --upload-certs
    
    echo "Cluster initialized successfully"
}

# Function to configure kubectl
configure_kubectl() {
    echo "Configuring kubectl..."
    
    # Create .kube directory
    mkdir -p $HOME/.kube
    
    # Copy admin.conf
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config
    
    echo "kubectl configured for user $(whoami)"
}

# Function to install network plugin
install_network_plugin() {
    echo "Installing network plugin: $NETWORK_PLUGIN"
    
    case $NETWORK_PLUGIN in
        "flannel")
            kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
            ;;
        "calico")
            kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/tigera-operator.yaml
            kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/custom-resources.yaml
            ;;
        "cilium")
            kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.14.0/install/kubernetes/quick-install.yaml
            ;;
        *)
            echo "Unsupported network plugin: $NETWORK_PLUGIN"
            echo "Supported plugins: flannel, calico, cilium"
            exit 1
            ;;
    esac
    
    echo "Network plugin $NETWORK_PLUGIN installed"
}

# Function to generate join command
generate_join_command() {
    echo "Generating join command for worker nodes..."
    
    JOIN_COMMAND=$(sudo kubeadm token create --print-join-command)
    echo "$JOIN_COMMAND" > /tmp/join-command.sh
    chmod +x /tmp/join-command.sh
    
    echo "Join command saved to /tmp/join-command.sh"
    echo "Content:"
    cat /tmp/join-command.sh
}

# Main execution
main() {
    echo "Starting Kubernetes master initialization..."
    
    # Wait for containerd
    wait_for_service "containerd" 30
    
    # Install Kubernetes components
    install_kubernetes
    
    # Enable and start kubelet
    sudo systemctl enable kubelet
    sudo systemctl start kubelet
    
    # Initialize cluster
    initialize_cluster
    
    # Configure kubectl
    configure_kubectl
    
    # Install network plugin
    install_network_plugin
    
    # Generate join command
    generate_join_command
    
    # Wait for master to be ready
    echo "Waiting for master node to be ready..."
    sleep 30
    
    # Verify cluster status
    echo "Verifying cluster status..."
    kubectl get nodes
    kubectl get pods --all-namespaces
    
    echo
    echo "=== Kubernetes Master Setup Complete ==="
    echo "Master node is ready and operational"
    echo "Use 'kubectl get nodes' to check cluster status"
    echo "Use 'kubectl get pods --all-namespaces' to check pods"
    echo "Join command is available at /tmp/join-command.sh"
}

# Run main function
main "$@"