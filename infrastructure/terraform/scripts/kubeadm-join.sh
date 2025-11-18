#!/bin/bash
# Kubernetes Worker Node Join Script
# This script joins a worker node to an existing Kubernetes cluster

set -e

# Configuration variables
JOIN_COMMAND=${1:-""}
JOIN_TOKEN=${2:-""}
CA_CERT_HASH=${3:-""}
MASTER_IP=${4:-""}

echo "=== Kubernetes Worker Node Join ==="

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Validate inputs
if [ -z "$JOIN_COMMAND" ] && [ -z "$JOIN_TOKEN" ] && [ -z "$MASTER_IP" ]; then
    echo "Error: Either provide a complete join command, or provide JOIN_TOKEN and MASTER_IP"
    echo "Usage: $0 '<join_command>'"
    echo "Usage: $0 '' '<join_token>' '<ca_cert_hash>' '<master_ip>'"
    exit 1
fi

# Build join command if not provided
if [ -z "$JOIN_COMMAND" ]; then
    JOIN_COMMAND="kubeadm join $MASTER_IP:6443 --token $JOIN_TOKEN --discovery-token-ca-cert-hash sha256:$CA_CERT_HASH"
fi

echo "Join command: $JOIN_COMMAND"

# Disable swap
echo "Disabling swap..."
swapoff -a
sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# Load kernel modules
echo "Loading kernel modules..."
modprobe overlay
modprobe br_netfilter

# Apply sysctl settings
echo "Applying sysctl settings..."
sysctl --system

# Install containerd if not present
if ! command_exists containerd; then
    echo "Installing containerd..."
    apt-get update
    apt-get install -y containerd.io
fi

# Configure containerd
echo "Configuring containerd..."
mkdir -p /etc/containerd
containerd config default > /etc/containerd/config.toml

# Configure registry mirror for Alibaba Cloud
sed -i 's|k8s.gcr.io|registry.aliyuncs.com/google_containers|g' /etc/containerd/config.toml

# Restart containerd
systemctl restart containerd
systemctl enable containerd

# Install Kubernetes packages if not present
if ! command_exists kubeadm; then
    echo "Installing Kubernetes packages..."
    curl -fsSL https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg
    echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main' | tee /etc/apt/sources.list.d/kubernetes.list
    
    apt-get update
    apt-get install -y kubelet=1.30.0-1.1 kubeadm=1.30.0-1.1 kubectl=1.30.0-1.1
    apt-mark hold kubelet kubeadm kubectl
fi

# Pull required images
echo "Pulling required images..."
kubeadm config images pull --image-repository=registry.aliyuncs.com/google_containers

# Wait a bit for master to be ready
echo "Waiting for master to be ready..."
sleep 30

# Join the cluster
echo "Joining Kubernetes cluster..."
eval "$JOIN_COMMAND"

echo "=== Worker node joined successfully ==="
echo ""
echo "To verify the node joined the cluster, run on the master:"
echo "kubectl get nodes"

# Save join info
echo "$JOIN_COMMAND" > /tmp/kubeadm-join-executed.sh
chmod 644 /tmp/kubeadm-join-executed.sh