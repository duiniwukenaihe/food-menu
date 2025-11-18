version: 2
ethernets:
  eth0:
    dhcp4: false
    addresses:
      - ${ip}
    gateway4: ${gateway}
    nameservers:
      addresses:
        - ${nameserver}
      search:
        - ${searchdomain}
