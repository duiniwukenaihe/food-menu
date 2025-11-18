# =============================================================================
# Makefile for Kubernetes on Proxmox Terraform Project
# =============================================================================

.PHONY: help init plan apply destroy validate fmt clean download-image test

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Terraform commands
init: ## Initialize Terraform providers and modules
	@echo "Initializing Terraform..."
	terraform init

plan: ## Show Terraform execution plan
	@echo "Creating Terraform plan..."
	terraform plan

apply: ## Apply Terraform configuration
	@echo "Applying Terraform configuration..."
	terraform apply

destroy: ## Destroy all Terraform-managed resources
	@echo "Destroying Terraform resources..."
	terraform destroy

validate: ## Validate Terraform configuration
	@echo "Validating Terraform configuration..."
	terraform validate

fmt: ## Format Terraform files
	@echo "Formatting Terraform files..."
	terraform fmt

clean: ## Clean Terraform cache and temporary files
	@echo "Cleaning Terraform cache..."
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -f *.tfplan
	rm -f terraform.tfstate.backup

# Image management
download-image: ## Download Ubuntu cloud image with verification
	@echo "Downloading Ubuntu cloud image..."
	./scripts/get-ubuntu-cloudimg.sh

# Testing and validation
test: ## Run Terraform tests
	@echo "Running Terraform tests..."
	terraform test

lint: ## Run Terraform linter
	@echo "Running Terraform linter..."
	terraform fmt -check
	terraform validate

check: ## Run all checks (fmt, validate, lint)
	@echo "Running all checks..."
	$(MAKE) fmt
	$(MAKE) validate
	$(MAKE) lint

# Cluster management
plan-cluster: ## Plan cluster deployment
	@echo "Planning cluster deployment..."
	terraform plan -var-file=terraform.tfvars

deploy-cluster: ## Deploy Kubernetes cluster
	@echo "Deploying Kubernetes cluster..."
	terraform apply -var-file=terraform.tfvars

scale-up-workers: ## Scale up worker nodes (set WORKER_COUNT)
	@if [ -z "$(WORKER_COUNT)" ]; then \
		echo "Usage: make scale-up-workers WORKER_COUNT=N"; \
		exit 1; \
	fi
	@echo "Scaling workers to $(WORKER_COUNT)..."
	terraform apply -var="worker_count=$(WORKER_COUNT)" -var-file=terraform.tfvars

scale-up-masters: ## Scale up master nodes (set MASTER_COUNT)
	@if [ -z "$(MASTER_COUNT)" ]; then \
		echo "Usage: make scale-up-masters MASTER_COUNT=N"; \
		exit 1; \
	fi
	@echo "Scaling masters to $(MASTER_COUNT)..."
	terraform apply -var="master_count=$(MASTER_COUNT)" -var-file=terraform.tfvars

# Output and information
outputs: ## Show Terraform outputs
	@echo "Terraform outputs:"
	terraform output

cluster-info: ## Show cluster information
	@echo "=== Kubernetes Cluster Information ==="
	@echo "Cluster Name: $$(terraform output -raw cluster_name 2>/dev/null || echo 'Not deployed')"
	@echo "Kubernetes Version: $$(terraform output -raw kubernetes_version 2>/dev/null || echo 'Not deployed')"
	@echo "Master Nodes: $$(terraform output -json | jq -r '.master_nodes.value.count // "Not deployed"')"
	@echo "Worker Nodes: $$(terraform output -json | jq -r '.worker_nodes.value.count // "Not deployed"')"
	@echo "Network Plugin: $$(terraform output -raw network_plugin 2>/dev/null || echo 'Not deployed')"

ssh-master: ## SSH to first master node
	@MASTER_IP=$$(terraform output -json | jq -r '.master_ips.value[0] // empty'); \
	if [ -n "$$MASTER_IP" ]; then \
		echo "SSH to master node: $$MASTER_IP"; \
		ssh -i $$(terraform output -raw ssh_private_key 2>/dev/null || echo "~/.ssh/id_rsa") ubuntu@$$MASTER_IP; \
	else \
		echo "No master IP found. Cluster may not be deployed."; \
	fi

list-vms: ## List all VMs in Proxmox
	@echo "Listing VMs in Proxmox..."
	@qm list 2>/dev/null || echo "qm command not available. Please run on Proxmox host."

# Backup and restore
backup-state: ## Backup Terraform state
	@echo "Backing up Terraform state..."
	cp terraform.tfstate terraform.tfstate.backup.$$(date +%Y%m%d_%H%M%S)
	@echo "State backed up"

restore-state: ## Restore Terraform state (set BACKUP_FILE)
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "Usage: make restore-state BACKUP_FILE=terraform.tfstate.backup.YYYYMMDD_HHMMSS"; \
		exit 1; \
	fi
	@echo "Restoring Terraform state from $(BACKUP_FILE)..."
	cp $(BACKUP_FILE) terraform.tfstate
	@echo "State restored"

# Development and debugging
debug-plan: ## Show detailed Terraform plan with debug
	@echo "Creating detailed Terraform plan..."
	TF_LOG=DEBUG terraform plan -detailed-exitcode

graph: ## Generate dependency graph
	@echo "Generating Terraform graph..."
	terraform graph | dot -Tpng > dependency-graph.png
	@echo "Dependency graph saved to dependency-graph.png"

# Security and maintenance
security-scan: ## Run security scan on configuration
	@echo "Running security scan..."
	@echo "Checking for sensitive data..."
	@if grep -r "password\|secret\|key" . --include="*.tf" --include="*.tfvars" | grep -v "sensitive\|example" | head -10; then \
		echo "⚠️  Potential sensitive data found in configuration files"; \
	else \
		echo "✅ No obvious sensitive data found"; \
	fi

update-providers: ## Update Terraform providers
	@echo "Updating Terraform providers..."
	terraform init -upgrade

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v terraform-docs >/dev/null 2>&1; then \
		terraform-docs markdown . > DOCUMENTATION.md; \
		echo "Documentation generated in DOCUMENTATION.md"; \
	else \
		echo "terraform-docs not installed. Install with: go install github.com/terraform-docs/terraform-docs@latest"; \
	fi

# Quick start shortcuts
quick-start: ## Quick start: init, plan, apply
	@echo "Quick start: Initializing, planning, and applying..."
	$(MAKE) init
	$(MAKE) plan
	$(MAKE) apply

quick-destroy: ## Quick destroy: destroy without confirmation
	@echo "Quick destroy: Destroying all resources..."
	terraform destroy -auto-approve

# Environment management
dev-env: ## Set up development environment
	@echo "Setting up development environment..."
	@if [ ! -f terraform.tfvars ]; then \
		cp terraform.tfvars.example terraform.tfvars; \
		echo "Created terraform.tfvars from example. Please edit with your values."; \
	else \
		echo "terraform.tfvars already exists"; \
	fi

prod-env: ## Set up production environment variables
	@echo "Setting production environment variables..."
	@echo "export TF_VAR_environment=prod" >> .env
	@echo "export TF_VAR_master_count=3" >> .env
	@echo "export TF_VAR_worker_count=5" >> .env
	@echo "Production environment variables set in .env"

# Monitoring and health
health-check: ## Check cluster health
	@echo "Checking cluster health..."
	@MASTER_IP=$$(terraform output -json | jq -r '.master_ips.value[0] // empty'); \
	if [ -n "$$MASTER_IP" ]; then \
		echo "Checking cluster health via master: $$MASTER_IP"; \
		ssh -i $$(terraform output -raw ssh_private_key 2>/dev/null || echo "~/.ssh/id_rsa") ubuntu@$$MASTER_IP "kubectl get nodes && kubectl get pods --all-namespaces"; \
	else \
		echo "No master IP found. Cluster may not be deployed."; \
	fi

# Version information
version: ## Show Terraform and provider versions
	@echo "Terraform version:"
	@terraform version
	@echo ""
	@echo "Provider versions:"
	@terraform providers

# Clean and reset
reset: ## Complete reset: destroy, clean, and reinitialize
	@echo "Complete reset: This will destroy everything and clean up!"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	$(MAKE) destroy
	$(MAKE) clean
	rm -f terraform.tfstate
	rm -f terraform.tfstate.backup.*
	$(MAKE) init
	@echo "Reset completed"