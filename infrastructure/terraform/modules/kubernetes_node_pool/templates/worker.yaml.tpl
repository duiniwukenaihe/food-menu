#cloud-config
# =============================================================================
# Cloud-Init Configuration for Kubernetes Worker Nodes
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

  # Kubernetes worker join script
  - path: /opt/k8s-join.sh
    permissions: '0755'
    content: |
      #!/bin/bash
      set -e
      
      echo "Starting Kubernetes worker setup..."
      
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
      
      # Wait for master to be ready
      echo "Waiting for master node to be ready..."
      sleep 60
      
      # Try to get join command from master (assuming SSH key access)
      MASTER_IP="${cluster_name}-master-1"
      JOIN_COMMAND=""
      
      # Try to get join command from master
      for i in {1..10}; do
        if ssh -o StrictHostKeyChecking=no -i /home/${username}/.ssh/id_rsa ${username}@${MASTER_IP} "test -f /tmp/join-command.sh" 2>/dev/null; then
          JOIN_COMMAND=$(ssh -o StrictHostKeyChecking=no -i /home/${username}/.ssh/id_rsa ${username}@${MASTER_IP} "cat /tmp/join-command.sh")
          break
        fi
        echo "Attempt $i: Master not ready yet, waiting 30 seconds..."
        sleep 30
      done
      
      if [ -n "$JOIN_COMMAND" ]; then
        echo "Joining Kubernetes cluster..."
        sudo $JOIN_COMMAND
        
        echo "Worker node joined successfully!"
      else
        echo "Could not get join command from master. Manual join required."
        echo "Please run the following on the master to get the join command:"
        echo "sudo kubeadm token create --print-join-command"
      fi
      
      # Enable and start kubelet
      sudo systemctl enable kubelet
      sudo systemctl start kubelet

  # Manual join instructions
  - path: /tmp/join-instructions.txt
    permissions: '0644'
    content: |
      KUBERNETES WORKER NODE JOIN INSTRUCTIONS
      ========================================
      
      If the automatic join failed, you can join this worker manually:
      
      1. SSH to the master node:
         ssh ${username}@${cluster_name}-master-1
      
      2. Get the join command:
         sudo kubeadm token create --print-join-command
      
      3. Copy the output and run it on this worker node with sudo
      
      4. Verify the node joined:
         kubectl get nodes

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
  
  # Make join script executable
  - chmod +x /opt/k8s-join.sh
  
  # Run Kubernetes join in background
  - nohup /opt/k8s-join.sh > /var/log/k8s-join.log 2>&1 &

# Power state
power_state:
  mode: reboot

# Final message
final_message: "Kubernetes worker node setup completed after $UPTIME seconds"