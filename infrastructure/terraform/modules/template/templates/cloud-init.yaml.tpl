#cloud-config
# =============================================================================
# Cloud-Init Configuration for Kubernetes Template
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
  - kubelet
  - kubeadm
  - kubectl

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
          sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.9"
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

  # Environment variables
  - path: /etc/environment
    permissions: '0644'
    append: true
    content: |
      KUBECONFIG=/etc/kubernetes/admin.conf

# Runcmd - Commands to run after boot
runcmd:
  # Enable and start services
  - systemctl enable disable-swap
  - systemctl start disable-swap
  - systemctl enable containerd
  - systemctl start containerd
  - systemctl enable docker
  - systemctl start docker
  
  # Load kernel modules
  - modprobe br_netfilter
  - modprobe overlay
  
  # Apply sysctl settings
  - sysctl --system
  
  # Create Kubernetes directories
  - mkdir -p /etc/kubernetes
  - mkdir -p /var/lib/kubelet
  - mkdir -p /var/lib/kubernetes
  - mkdir -p /opt/cni/bin
  
  # Add user to docker group
  - usermod -aG docker ${username}
  
  # Disable swap in fstab
  - sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
  
  # Hold Kubernetes packages
  - apt-mark hold kubelet kubeadm kubectl
  
  # Create kubeadm configuration directory
  - mkdir -p /etc/kubernetes/kubeadm
  
  # Status message
  - echo "Kubernetes template setup completed successfully!"

# Power state
power_state:
  mode: reboot

# Final message
final_message: "Kubernetes template setup completed after $UPTIME seconds"