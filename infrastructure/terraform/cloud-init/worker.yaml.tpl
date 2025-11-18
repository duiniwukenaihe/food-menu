#cloud-config
package_upgrade: true
packages:
  - apt-transport-https
  - ca-certificates
  - curl
  - gnupg
  - lsb-release

write_files:
  - path: /etc/sysctl.d/k8s.conf
    content: |
      net.bridge.bridge-nf-call-iptables  = 1
      net.bridge.bridge-nf-call-ip6tables = 1
      net.ipv4.ip_forward                 = 1
  - path: /etc/modules-load.d/containerd.conf
    content: |
      overlay
      br_netfilter
  - path: /tmp/kubeadm-join.sh
    permissions: '0755'
    content: |
      #!/bin/bash
      set -e
      
      echo "=== Kubernetes Worker Bootstrap Script ==="
      
      # Disable swap
      echo "Disabling swap..."
      sudo swapoff -a
      sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
      
      # Load kernel modules
      echo "Loading kernel modules..."
      sudo modprobe overlay
      sudo modprobe br_netfilter
      
      # Apply sysctl settings
      echo "Applying sysctl settings..."
      sudo sysctl --system
      
      # Install containerd
      echo "Installing containerd..."
      sudo apt-get update
      sudo apt-get install -y containerd.io
      
      # Configure containerd
      echo "Configuring containerd with registry mirror..."
      sudo mkdir -p /etc/containerd
      sudo containerd config default | sudo tee /etc/containerd/config.toml
      
      # Configure registry mirror for Alibaba Cloud
      sudo sed -i 's|k8s.gcr.io|registry.aliyuncs.com/google_containers|g' /etc/containerd/config.toml
      
      # Restart containerd
      sudo systemctl restart containerd
      sudo systemctl enable containerd
      
      # Add Kubernetes apt repository
      echo "Adding Kubernetes repository..."
      curl -fsSL https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg
      echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main' | sudo tee /etc/apt/sources.list.d/kubernetes.list
      
      # Install Kubernetes packages
      echo "Installing Kubernetes packages..."
      sudo apt-get update
      sudo apt-get install -y kubelet=1.30.0-1.1 kubeadm=1.30.0-1.1 kubectl=1.30.0-1.1
      sudo apt-mark hold kubelet kubeadm kubectl
      
      # Wait a bit for master to be ready
      echo "Waiting for master to be ready..."
      sleep 30
      
      # Join the cluster
      echo "Joining Kubernetes cluster..."
      ${join_command}
      
      echo "=== Worker node bootstrap completed ==="

runcmd:
  - [bash, /tmp/kubeadm-join.sh]

final_message: "Kubernetes worker node bootstrap completed successfully!"