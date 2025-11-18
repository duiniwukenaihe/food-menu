#cloud-config
# =============================================================================
# Cloud-Init Configuration for Kubernetes Master Nodes
# =============================================================================

# Hostname configuration
hostname: ${hostname}
preserve_hostname: false
fqdn: ${hostname}.local

# User configuration
users:
  - name: ${username}
    groups: [sudo, adm, docker]
    shell: /bin/bash
    lock_passwd: false
    passwd: ${password}
    ssh_authorized_keys:
      - ${ssh_public_key}
    sudo: ["ALL=(ALL) NOPASSWD:ALL"]

# SSH configuration
ssh_pwauth: true
ssh_deletekeys: false
ssh_genkeytypes: [rsa, ecdsa, ed25519]

# Network configuration
network:
  version: 2
  ethernets:
    ens18:
      dhcp4: true
      dhcp6: false

# Package management
package_update: true
package_upgrade: true
packages:
  - curl
  - wget
  - apt-transport-https
  - ca-certificates
  - gnupg
  - lsb-release
  - software-properties-common
  - htop
  - vim
  - git
  - unzip
  - jq
  - containerd
  - docker.io

# Package configuration
package_reboot_if_required: false

# Time and timezone
timezone: UTC

# Locale
locale: en_US.UTF-8

# Write files
write_files:
  # Disable swap
  - path: /etc/systemd/system/disable-swap.service
    permissions: '0644'
    content: |
      [Unit]
      Description=Disable swap
      Before=kubelet.service
      
      [Service]
      Type=oneshot
      ExecStart=/sbin/swapoff -a
      ExecStart=/bin/sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
      
      [Install]
      WantedBy=multi-user.target

  # Containerd configuration
  - path: /etc/containerd/config.toml
    permissions: '0644'
    content: |
      version = 2
      [plugins]
        [plugins."io.containerd.grpc.v1.cri"]
          sandbox_image = "${image_repository}/pause:3.9"
          [plugins."io.containerd.grpc.v1.cri".containerd]
            [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]
              [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
                runtime_type = "io.containerd.runc.v2"
                [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
                  SystemdCgroup = true

  # Kubernetes sysctl settings
  - path: /etc/sysctl.d/99-kubernetes.conf
    permissions: '0644'
    content: |
      net.bridge.bridge-nf-call-iptables  = 1
      net.bridge.bridge-nf-call-ip6tables = 1
      net.ipv4.ip_forward                 = 1
      vm.swappiness                       = 0

  # Kubernetes modules
  - path: /etc/modules-load.d/kubernetes.conf
    permissions: '0644'
    content: |
      br_netfilter
      overlay

  # Docker daemon configuration
  - path: /etc/docker/daemon.json
    permissions: '0644'
    content: |
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

  # Kubernetes repositories
  - path: /etc/apt/sources.list.d/kubernetes.list
    permissions: '0644'
    content: |
      deb https://apt.kubernetes.io/ kubernetes-xenial main

  # Kubeadm configuration for master
  - path: /etc/kubernetes/kubeadm-config.yaml
    permissions: '0644'
    content: |
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
      kubernetesVersion: ${k8s_version}
      controlPlaneEndpoint: "${cluster_name}-master-1:6443"
      imageRepository: ${image_repository}
      networking:
        serviceSubnet: "10.96.0.0/12"
        podSubnet: "${pod_network_cidr}"
        dnsDomain: "cluster.local"
      ---
      apiVersion: kubelet.config.k8s.io/v1beta1
      kind: KubeletConfiguration
      cgroupDriver: systemd
      serverTLSBootstrap: true
      rotateCertificates: true

  # Kubernetes master initialization script
  - path: /opt/k8s-init.sh
    permissions: '0755'
    content: |
      #!/bin/bash
      set -e
      
      echo "Starting Kubernetes master initialization..."
      
      # Wait for system to be ready
      sleep 30
      
      # Disable swap
      sudo swapoff -a
      sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
      
      # Load kernel modules
      sudo modprobe br_netfilter
      sudo modprobe overlay
      
      # Apply sysctl settings
      sudo sysctl --system
      
      # Enable and start services
      sudo systemctl enable disable-swap
      sudo systemctl start disable-swap
      sudo systemctl enable containerd
      sudo systemctl start containerd
      sudo systemctl enable docker
      sudo systemctl start docker
      
      # Add Kubernetes repository
      curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
      echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
      
      # Update package list
      sudo apt update
      
      # Install Kubernetes packages
      sudo apt install -y kubelet=${k8s_version}-1.1.1 kubeadm=${k8s_version}-1.1.1 kubectl=${k8s_version}-1.1.1
      
      # Hold Kubernetes packages
      sudo apt-mark hold kubelet kubeadm kubectl
      
      # Initialize Kubernetes cluster (only on first master)
      if [ "${is_master}" = "true" ]; then
        echo "Initializing Kubernetes cluster..."
        sudo kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml --upload-certs
        
        # Configure kubectl for user
        mkdir -p $HOME/.kube
        sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
        sudo chown $(id -u):$(id -g) $HOME/.kube/config
        
        # Install network plugin
        if [ "${network_plugin}" = "flannel" ]; then
          kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
        elif [ "${network_plugin}" = "calico" ]; then
          kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/tigera-operator.yaml
          kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/custom-resources.yaml
        fi
        
        # Generate join command for workers
        sudo kubeadm token create --print-join-command > /tmp/join-command.sh
        chmod +x /tmp/join-command.sh
        
        echo "Kubernetes master initialized successfully!"
        echo "Join command saved to /tmp/join-command.sh"
      fi
      
      # Enable and start kubelet
      sudo systemctl enable kubelet
      sudo systemctl start kubelet

  # Environment variables
  - path: /etc/environment
    permissions: '0644'
    append: true
    content: |
      KUBECONFIG=/etc/kubernetes/admin.conf

# Runcmd - Commands to run after boot
runcmd:
  # Add user to docker group
  - usermod -aG docker ${username}
  
  # Make init script executable
  - chmod +x /opt/k8s-init.sh
  
  # Run Kubernetes initialization in background
  - nohup /opt/k8s-init.sh > /var/log/k8s-init.log 2>&1 &

# Power state
power_state:
  mode: reboot

# Final message
final_message: "Kubernetes master node setup completed after $UPTIME seconds"