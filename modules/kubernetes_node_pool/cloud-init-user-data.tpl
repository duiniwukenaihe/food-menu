#cloud-config
hostname: ${hostname}
fqdn: ${hostname}.${searchdomain}
manage_etc_hosts: true
preserve_hostname: false

users:
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    lock_passwd: ${lock_passwd}
%{ if ssh_password != "" ~}
    passwd: ${ssh_password}
%{ endif ~}
%{ if length(ssh_public_keys) > 0 ~}
    ssh_authorized_keys:
%{ for key in ssh_public_keys ~}
      - ${key}
%{ endfor ~}
%{ endif ~}

package_update: true
package_upgrade: false

packages:
  - qemu-guest-agent
  - curl
  - apt-transport-https
  - ca-certificates
  - software-properties-common

runcmd:
  - systemctl enable qemu-guest-agent
  - systemctl start qemu-guest-agent
  - hostnamectl set-hostname ${hostname}
  - echo "preserve_hostname: true" >> /etc/cloud/cloud.cfg

power_state:
  mode: reboot
  timeout: 30
  condition: true
