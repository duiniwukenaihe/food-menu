.PHONY: help init validate plan apply destroy fmt clean docs

# Default target
help:
	@echo "Available targets:"
	@echo "  init      - Initialize Terraform working directory"
	@echo "  validate  - Validate Terraform configuration"
	@echo "  fmt       - Format Terraform files"
	@echo "  plan      - Generate and show execution plan"
	@echo "  apply     - Build or change infrastructure"
	@echo "  destroy   - Destroy Terraform-managed infrastructure"
	@echo "  clean     - Remove Terraform working files"
	@echo "  docs      - Generate documentation (requires terraform-docs)"
	@echo ""
	@echo "Environment variables:"
	@echo "  PROXMOX_VE_USERNAME   - Proxmox username (e.g., root@pam)"
	@echo "  PROXMOX_VE_PASSWORD   - Proxmox password"
	@echo "  PROXMOX_VE_API_TOKEN  - Proxmox API token (alternative to username/password)"
	@echo ""
	@echo "Example usage:"
	@echo "  make init"
	@echo "  make plan"
	@echo "  make apply"

init:
	@echo "Initializing Terraform..."
	terraform init

validate: init
	@echo "Validating Terraform configuration..."
	terraform validate

fmt:
	@echo "Formatting Terraform files..."
	terraform fmt -recursive

plan: validate
	@echo "Generating execution plan..."
	terraform plan

apply: validate
	@echo "Applying configuration..."
	terraform apply

apply-auto: validate
	@echo "Applying configuration (auto-approve)..."
	terraform apply -auto-approve

destroy:
	@echo "Destroying infrastructure..."
	terraform destroy

destroy-auto:
	@echo "Destroying infrastructure (auto-approve)..."
	terraform destroy -auto-approve

clean:
	@echo "Cleaning Terraform working files..."
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate
	rm -f terraform.tfstate.backup
	rm -f *.tfplan

docs:
	@echo "Generating documentation..."
	@if command -v terraform-docs >/dev/null 2>&1; then \
		terraform-docs markdown table . > TERRAFORM.md; \
		echo "Documentation generated in TERRAFORM.md"; \
	else \
		echo "terraform-docs not found. Install it from: https://terraform-docs.io/"; \
		exit 1; \
	fi

check-env:
	@echo "Checking environment variables..."
	@if [ -z "$$PROXMOX_VE_USERNAME" ] && [ -z "$$PROXMOX_VE_API_TOKEN" ]; then \
		echo "ERROR: Either PROXMOX_VE_USERNAME or PROXMOX_VE_API_TOKEN must be set"; \
		exit 1; \
	fi
	@if [ -n "$$PROXMOX_VE_USERNAME" ] && [ -z "$$PROXMOX_VE_PASSWORD" ]; then \
		echo "WARNING: PROXMOX_VE_USERNAME is set but PROXMOX_VE_PASSWORD is not"; \
	fi
	@echo "Environment variables are configured"

show:
	@echo "Showing current state..."
	terraform show

state-list:
	@echo "Listing resources in state..."
	terraform state list

output:
	@echo "Showing outputs..."
	terraform output
