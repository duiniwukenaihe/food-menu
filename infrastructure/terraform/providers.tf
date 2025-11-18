provider "proxmox" {
  endpoint = var.proxmox_endpoint

  # Authentication options (use environment variables or variables):
  # - For API token authentication:
  #   Set PROXMOX_VE_API_TOKEN environment variable or use api_token variable
  # - For username/password authentication:
  #   Set PROXMOX_VE_USERNAME and PROXMOX_VE_PASSWORD environment variables
  #   or use username/password variables
  username  = var.proxmox_username
  password  = var.proxmox_password
  api_token = var.proxmox_api_token

  # Skip TLS verification if using self-signed certificates
  # In production, consider using proper TLS certificates and set this to false
  insecure = var.proxmox_insecure

  # SSH connection configuration for Proxmox host operations
  ssh {
    agent    = false
    username = var.proxmox_ssh_username
    # Private key can be provided via environment variable PROXMOX_VE_SSH_PRIVATE_KEY
    # or via the private_key parameter below
    private_key = var.proxmox_ssh_private_key != "" ? file(var.proxmox_ssh_private_key) : null
  }
}
