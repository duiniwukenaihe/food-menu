#!/bin/bash
# =============================================================================
# Common Cloud-Init Functions for Kubernetes Nodes
# =============================================================================

set -e

# Function to log messages
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Function to check if service is running
wait_for_service() {
    local service=$1
    local timeout=${2:-30}
    
    log "Waiting for $service to be ready..."
    
    for i in $(seq 1 $timeout); do
        if systemctl is-active --quiet $service; then
            log "$service is ready"
            return 0
        fi
        sleep 10
    done
    
    log "Timeout waiting for $service"
    return 1
}

# Function to disable swap
disable_swap() {
    log "Disabling swap..."
    
    # Turn off swap
    swapoff -a
    
    # Remove swap entries from fstab
    sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
    
    log "Swap disabled successfully"
}

# Function to load kernel modules
load_kernel_modules() {
    log "Loading required kernel modules..."
    
    # Load modules
    modprobe br_netfilter
    modprobe overlay
    
    # Make modules persistent
    cat > /etc/modules-load.d/kubernetes.conf <<EOF
br_netfilter
overlay
EOF
    
    log "Kernel modules loaded successfully"
}

# Function to configure sysctl settings
configure_sysctl() {
    log "Configuring sysctl settings..."
    
    cat > /etc/sysctl.d/99-kubernetes.conf <<EOF
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
vm.swappiness                       = 0
EOF
    
    # Apply settings
    sysctl --system
    
    log "Sysctl settings configured successfully"
}

# Function to configure containerd
configure_containerd() {
    log "Configuring containerd..."
    
    # Create config directory
    mkdir -p /etc/containerd
    
    # Generate default config
    containerd config default | tee /etc/containerd/config.toml
    
    # Modify config for Kubernetes
    sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
    sed -i "s|k8s.gcr.io/pause:3.9|${IMAGE_REGISTRY:-registry.aliyuncs.com/google_containers}/pause:3.9|g" /etc/containerd/config.toml
    
    # Restart containerd
    systemctl restart containerd
    systemctl enable containerd
    
    log "Containerd configured successfully"
}

# Function to configure Docker
configure_docker() {
    log "Configuring Docker..."
    
    # Create Docker daemon config
    mkdir -p /etc/docker
    
    cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "registry-mirrors": [
    "https://registry.aliyuncs.com"
  ]
}
EOF
    
    # Restart Docker
    systemctl restart docker
    systemctl enable docker
    
    log "Docker configured successfully"
}

# Function to add user to docker group
add_user_to_docker() {
    local username=${1:-ubuntu}
    
    log "Adding $username to docker group..."
    usermod -aG docker $username
    
    log "User $username added to docker group"
}

# Function to install Kubernetes components
install_kubernetes() {
    local k8s_version=${1:-1.30.0}
    
    log "Installing Kubernetes components version $k8s_version..."
    
    # Add Kubernetes repository
    curl -fsSL https://pkgs.k8s.io/core:/stable/v$k8s_version/deb/Release.key | \
        gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    
    echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable/v$k8s_version/deb/ /" | \
        tee /etc/apt/sources.list.d/kubernetes.list
    
    # Update package list
    apt-get update
    
    # Install Kubernetes packages
    apt-get install -y kubelet=$k8s_version-1.1.1 kubeadm=$k8s_version-1.1.1 kubectl=$k8s_version-1.1.1
    
    # Hold packages to prevent upgrades
    apt-mark hold kubelet kubeadm kubectl
    
    # Enable kubelet
    systemctl enable kubelet
    
    log "Kubernetes components installed successfully"
}

# Function to create Kubernetes directories
create_kubernetes_directories() {
    log "Creating Kubernetes directories..."
    
    mkdir -p /etc/kubernetes
    mkdir -p /var/lib/kubelet
    mkdir -p /var/lib/kubernetes
    mkdir -p /opt/cni/bin
    
    log "Kubernetes directories created"
}

# Function to validate system requirements
validate_system() {
    log "Validating system requirements..."
    
    # Check if running as root
    if [[ $EUID -ne 0 ]]; then
        log "This script must be run as root"
        exit 1
    fi
    
    # Check OS version
    if ! grep -q "Ubuntu" /etc/os-release; then
        log "This script is designed for Ubuntu"
        exit 1
    fi
    
    # Check if swap is disabled
    if swapon --show | grep -q .; then
        log "Swap is still enabled"
        exit 1
    fi
    
    # Check kernel modules
    if ! lsmod | grep -q br_netfilter; then
        log "br_netfilter module not loaded"
        exit 1
    fi
    
    if ! lsmod | grep -q overlay; then
        log "overlay module not loaded"
        exit 1
    fi
    
    log "System validation passed"
}

# Function to setup common prerequisites
setup_prerequisites() {
    log "Setting up common prerequisites..."
    
    # Update system
    apt-get update
    apt-get upgrade -y
    
    # Install required packages
    apt-get install -y \
        curl \
        wget \
        apt-transport-https \
        ca-certificates \
        gnupg \
        lsb-release \
        software-properties-common \
        htop \
        vim \
        git \
        unzip \
        jq \
        containerd \
        docker.io
    
    # Disable swap
    disable_swap
    
    # Load kernel modules
    load_kernel_modules
    
    # Configure sysctl
    configure_sysctl
    
    # Configure container runtime
    configure_containerd
    configure_docker
    
    # Create directories
    create_kubernetes_directories
    
    log "Common prerequisites setup completed"
}

# Export functions for use in other scripts
export -f log
export -f wait_for_service
export -f disable_swap
export -f load_kernel_modules
export -f configure_sysctl
export -f configure_containerd
export -f configure_docker
export -f add_user_to_docker
export -f install_kubernetes
export -f create_kubernetes_directories
export -f validate_system
export -f setup_prerequisites

log "Common functions loaded successfully"