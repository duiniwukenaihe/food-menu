#cloud-config
hostname: ubuntu-template
manage_etc_hosts: true
fqdn: ubuntu-template.localdomain

# SSH configuration
ssh_pwauth: false
ssh_genkeytypes: [ed25519, rsa]

# Package updates
package_update: true
package_upgrade: true

# Install QEMU guest agent
packages:
  - qemu-guest-agent
  - cloud-init
  - netplan.io
  - openssh-server

# Enable and start services
runcmd:
  - systemctl enable qemu-guest-agent
  - systemctl start qemu-guest-agent
  - systemctl enable ssh
  - systemctl start ssh

# User configuration (will be customized when creating VMs from template)
users:
  - name: ubuntu
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: true
    ssh_authorized_keys:
      - ${ssh_public_key}

# Network configuration will be handled by network-data
network:
  config: disabled

# Cloud-init modules
cloud_init_modules:
  - migrator
  - bootcmd
  - write-files
  - growpart
  - resizefs
  - set_hostname
  - update_hostname
  - update_etc_hosts
  - rsyslog
  - users-groups
  - ssh

cloud_config_modules:
  - mounts
  - locale
  - set-passwords
  - grub-dpkg
  - apt-pipelining
  - apt-configure
  - package-update-upgrade-install
  - fan
  - landscape
  - timezone
  - resolv_conf
  - ca-certs
  - runcmd

cloud_final_modules:
  - scripts-vendor
  - scripts-per-once
  - scripts-per-boot
  - scripts-per-instance
  - scripts-user
  - ssh-authkey-fingerprints
  - keys-to-console
  - phone-home
  - final-message
  - power-state-change