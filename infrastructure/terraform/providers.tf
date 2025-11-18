# =============================================================================
# Proxmox Provider Configuration
# =============================================================================

provider "proxmox" {
  # API Configuration
  pm_api_url = var.proxmox_api_url
  
  # Authentication - Support both password and API token
  pm_api_token_id = var.proxmox_api_token_id != "" ? var.proxmox_api_token_id : null
  pm_api_token_secret = var.proxmox_api_token_secret != "" ? var.proxmox_api_token_secret : null
  
  # Fallback to username/password if no API token provided
  pm_user = var.proxmox_api_token_id == "" ? var.proxmox_user : null
  pm_password = var.proxmox_api_token_id == "" ? var.proxmox_password : null
  
  # TLS Configuration
  pm_tls_insecure = var.proxmox_insecure
  
  # Parallel operations
  pm_parallel = 4
  pm_timeout = 600
}

# =============================================================================
# SSH Provider Configuration (for Proxmox host operations)
# =============================================================================

provider "ssh" {
  host = var.proxmox_host
  user = var.proxmox_ssh_user
  private_key = file(var.proxmox_ssh_key)
  
  # Connection settings
  port = 22
  timeout = "60s"
  
  # Host key verification (disable for testing, enable for production)
  host_key = "accept-new"
}