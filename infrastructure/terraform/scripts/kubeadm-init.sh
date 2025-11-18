#!/bin/bash
# Kubernetes Master Node Initialization Script
# This script initializes a Kubernetes master node using kubeadm

set -e

# Configuration variables
POD_NETWORK_CIDR=${1:-"10.244.0.0/16"}
API_SERVER_ADVERTISE_ADDRESS=${2:-"$(hostname -I | awk '{print $1}')"}
NETWORK_PLUGIN=${3:-"flannel"}

echo "=== Kubernetes Master Node Initialization ==="
echo "Pod Network CIDR: $POD_NETWORK_CIDR"
echo "API Server Address: $API_SERVER_ADVERTISE_ADDRESS"
echo "Network Plugin: $NETWORK_PLUGIN"

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to wait for command to succeed
wait_for_success() {
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if "$@" >/dev/null 2>&1; then
            echo "Success on attempt $attempt"
            return 0
        fi
        echo "Attempt $attempt failed, retrying in 10 seconds..."
        sleep 10
        attempt=$((attempt + 1))
    done
    
    echo "Failed after $max_attempts attempts"
    return 1
}

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

# Initialize Kubernetes cluster
echo "Initializing Kubernetes cluster..."
kubeadm init \
    --pod-network-cidr="$POD_NETWORK_CIDR" \
    --apiserver-advertise-address="$API_SERVER_ADVERTISE_ADDRESS" \
    --image-repository=registry.aliyuncs.com/google_containers \
    --ignore-preflight-errors=all

# Configure kubectl
echo "Configuring kubectl..."
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

# Configure kubectl for ubuntu user
mkdir -p /home/ubuntu/.kube
cp -i /etc/kubernetes/admin.conf /home/ubuntu/.kube/config
chown ubuntu:ubuntu /home/ubuntu/.kube/config

# Install network plugin
echo "Installing network plugin: $NETWORK_PLUGIN"
case $NETWORK_PLUGIN in
    "flannel")
        kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
        ;;
    "cilium")
        kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.14.5/install/kubernetes/cilium.yaml
        ;;
    "calico")
        kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.4/manifests/tigera-operator.yaml
        kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.4/manifests/custom-resources.yaml
        ;;
    *)
        echo "Unknown network plugin: $NETWORK_PLUGIN"
        echo "Supported plugins: flannel, cilium, calico"
        exit 1
        ;;
esac

# Wait for network plugin to be ready
echo "Waiting for network plugin to be ready..."
kubectl wait --for=condition=ready pods -l k8s-app=flannel --timeout=300s -n kube-system 2>/dev/null || \
kubectl wait --for=condition=ready pods -l app.kubernetes.io/name=cilium --timeout=300s -n kube-system 2>/dev/null || \
kubectl wait --for=condition=ready pods -l k8s-app=calico-node --timeout=300s -n kube-system 2>/dev/null

# Generate join command and save
echo "Generating join command..."
kubeadm token create --print-join-command > /tmp/kubeadm-join-command.sh
chmod 644 /tmp/kubeadm-join-command.sh

# Generate and save join token
JOIN_TOKEN=$(kubeadm token generate)
kubeadm token create "$JOIN_TOKEN" --ttl 24h --description "Auto-generated token for worker nodes"
echo "$JOIN_TOKEN" > /tmp/kubeadm-join-token
chmod 644 /tmp/kubeadm-join-token

# Save CA cert hash
CA_CERT_HASH=$(openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')
echo "$CA_CERT_HASH" > /tmp/kubeadm-ca-cert-hash
chmod 644 /tmp/kubeadm-ca-cert-hash

echo "=== Master node initialization completed successfully ==="
echo "Join command saved to: /tmp/kubeadm-join-command.sh"
echo "Join token: $JOIN_TOKEN"
echo "CA cert hash: $CA_CERT_HASH"
echo ""
echo "To add worker nodes, run:"
echo "sudo $(cat /tmp/kubeadm-join-command.sh)"

# Save cluster info for verification
kubectl cluster-info > /tmp/cluster-info.txt
kubectl get nodes > /tmp/nodes-info.txt