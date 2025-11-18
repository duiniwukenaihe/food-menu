# Environment Variables Reference

This document lists all environment variables that can be used with the Terraform configuration.

## Required Environment Variables

At minimum, you need to set **one of the following authentication methods**:

### Option 1: Username and Password Authentication

```bash
export PROXMOX_VE_USERNAME="root@pam"
export PROXMOX_VE_PASSWORD="your-password-here"
```

### Option 2: API Token Authentication (Recommended)

```bash
export PROXMOX_VE_API_TOKEN="root@pam!terraform=12345678-1234-1234-1234-123456789abc"
```

To create an API token in Proxmox:
1. Navigate to Datacenter → Permissions → API Tokens
2. Click "Add"
3. Select user (e.g., root@pam)
4. Enter a Token ID (e.g., "terraform")
5. Uncheck "Privilege Separation" for full access
6. Click "Add"
7. Copy the displayed token secret (you won't be able to see it again)

## Optional Environment Variables

### SSH Configuration

```bash
# SSH private key for Proxmox host operations (alternative to file path)
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)"

# Or just specify the path in terraform.tfvars
# proxmox_ssh_private_key = "~/.ssh/id_rsa"
```

### Proxmox Connection

```bash
# Override Proxmox endpoint (default: https://192.168.0.200:8006)
export TF_VAR_proxmox_endpoint="https://your-proxmox-ip:8006"

# Override Proxmox node name (default: pve)
export TF_VAR_proxmox_node_name="pve1"
```

## Complete Example

### Development Environment

```bash
#!/bin/bash
# save as: env-dev.sh

# Proxmox authentication (choose one method)
export PROXMOX_VE_USERNAME="root@pam"
export PROXMOX_VE_PASSWORD="your-dev-password"

# SSH key for Proxmox host
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)"

# Optional: Override defaults
export TF_VAR_environment="dev"
export TF_VAR_project_name="k8s-dev"
```

### Production Environment

```bash
#!/bin/bash
# save as: env-prod.sh

# Use API token for production (more secure)
export PROXMOX_VE_API_TOKEN="root@pam!terraform=your-production-token-here"

# SSH key for Proxmox host
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa_prod)"

# Production overrides
export TF_VAR_environment="prod"
export TF_VAR_project_name="k8s-prod"
export TF_VAR_control_plane_count="3"
export TF_VAR_worker_count="5"
```

## Usage

Source the environment file before running Terraform commands:

```bash
# Source the environment variables
source env-dev.sh

# Or for production
source env-prod.sh

# Then run Terraform commands
terraform plan
terraform apply
```

## Security Best Practices

1. **Never commit environment files with secrets** to version control
   - Add `env-*.sh` to `.gitignore`
   - Use `*.env` or `*.sh` patterns

2. **Use API tokens instead of passwords** when possible
   - API tokens can be easily revoked
   - They can have limited scopes and expiration dates

3. **Restrict file permissions** for files containing secrets
   ```bash
   chmod 600 env-*.sh
   ```

4. **Use a secrets manager** for production environments
   - HashiCorp Vault
   - AWS Secrets Manager
   - Azure Key Vault
   - etc.

5. **Rotate credentials regularly**
   - Change API tokens periodically
   - Update SSH keys as part of security policy

## Verifying Environment Variables

Check if your environment variables are set correctly:

```bash
# Check Proxmox authentication
echo "Username: $PROXMOX_VE_USERNAME"
echo "Password set: $([ -n "$PROXMOX_VE_PASSWORD" ] && echo "Yes" || echo "No")"
echo "API Token set: $([ -n "$PROXMOX_VE_API_TOKEN" ] && echo "Yes" || echo "No")"

# Check SSH key
echo "SSH Key set: $([ -n "$PROXMOX_VE_SSH_PRIVATE_KEY" ] && echo "Yes" || echo "No")"

# Check Terraform variables
env | grep TF_VAR_
```

## Troubleshooting

### "Error: error creating Proxmox client"

- Verify your credentials are correct
- Check that PROXMOX_VE_PASSWORD or PROXMOX_VE_API_TOKEN is set
- Ensure the Proxmox endpoint is accessible

### "SSH authentication failed"

- Verify the SSH private key is correct
- Check that the corresponding public key is in Proxmox host's `~/.ssh/authorized_keys`
- Test SSH connection manually:
  ```bash
  ssh -i ~/.ssh/id_rsa root@192.168.0.200
  ```

### Variables not being picked up

- Ensure environment variables are exported (not just set)
- Check variable names are correct (case-sensitive)
- For Terraform variables, use `TF_VAR_` prefix
- Source the environment file in the same shell session where you run Terraform

## Reference Links

- [Proxmox Provider Documentation](https://registry.terraform.io/providers/bpg/proxmox/latest/docs)
- [Terraform Environment Variables](https://developer.hashicorp.com/terraform/cli/config/environment-variables)
- [Proxmox API Tokens](https://pve.proxmox.com/wiki/User_Management#pveum_tokens)
