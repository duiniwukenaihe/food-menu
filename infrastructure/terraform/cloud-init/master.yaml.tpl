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
  - path: /etc/systemd/system/containerd.service.d/registry-mirror.conf
    content: |
      [Service]
      Environment="CONTAINERD_CONFIG=/etc/containerd/config.toml"
  - path: /tmp/kubeadm-init.sh
    permissions: '0755'
    content: |
      #!/bin/bash
      set -e
      
      echo "=== Kubernetes Master Bootstrap Script ==="
      
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
      
      # Initialize Kubernetes cluster
      echo "Initializing Kubernetes cluster..."
      sudo kubeadm init --pod-network-cidr=${pod_network_cidr} --apiserver-advertise-address=${api_server_advertise_address} --ignore-preflight-errors=all
      
      # Configure kubectl for root user
      echo "Configuring kubectl..."
      mkdir -p $HOME/.kube
      sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
      sudo chown $(id -u):$(id -g) $HOME/.kube/config
      
      # Configure kubectl for non-root users
      sudo mkdir -p /home/ubuntu/.kube
      sudo cp -i /etc/kubernetes/admin.conf /home/ubuntu/.kube/config
      sudo chown ubuntu:ubuntu /home/ubuntu/.kube/config
      
      # Install network plugin
      echo "Installing network plugin: ${network_plugin}..."
      case ${network_plugin} in
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
          echo "Unknown network plugin: ${network_plugin}"
          exit 1
          ;;
      esac
      
      # Generate join command and save to shared location
      echo "Generating join command..."
      JOIN_TOKEN=$(sudo kubeadm token generate)
      sudo kubeadm token create $JOIN_TOKEN --ttl 24h --description "Auto-generated token for worker nodes"
      
      # Get CA cert hash
      CA_CERT_HASH=$(openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')
      
      # Create join command
      JOIN_COMMAND="sudo kubeadm join $(hostname -I | awk '{print $1}'):6443 --token $JOIN_TOKEN --discovery-token-ca-cert-hash sha256:$CA_CERT_HASH"
      
      # Save join info for workers
      echo "$JOIN_COMMAND" > /tmp/kubeadm-join-command.sh
      echo "$JOIN_TOKEN" > /tmp/kubeadm-join-token
      echo "$CA_CERT_HASH" > /tmp/kubeadm-ca-cert-hash
      chmod 644 /tmp/kubeadm-join-*
      
      echo "=== Master node bootstrap completed ==="
      echo "Join token: $JOIN_TOKEN"
      echo "CA cert hash: $CA_CERT_HASH"
      echo "Join command: $JOIN_COMMAND"

runcmd:
  - [bash, /tmp/kubeadm-init.sh]

final_message: "Kubernetes master node bootstrap completed successfully!"