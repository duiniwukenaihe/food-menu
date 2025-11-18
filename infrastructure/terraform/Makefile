# Kubernetes Cluster Management Makefile

.PHONY: help init plan apply destroy status clean fmt validate lint

# Default target
help: ## Show this help message
	@echo "Kubernetes Cluster Management Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Terraform commands
init: ## Initialize Terraform
	@echo "Initializing Terraform..."
	terraform init

plan: ## Show execution plan
	@echo "Planning Terraform deployment..."
	terraform plan -var-file="terraform.tfvars"

apply: ## Apply Terraform configuration
	@echo "Applying Terraform configuration..."
	terraform apply -var-file="terraform.tfvars"

destroy: ## Destroy Terraform resources
	@echo "Destroying Terraform resources..."
	terraform destroy -var-file="terraform.tfvars"

status: ## Show cluster status
	@echo "=== Cluster Information ==="
	terraform output cluster_name 2>/dev/null || echo "Cluster not deployed"
	terraform output kubernetes_version 2>/dev/null || echo "Version not available"
	terraform output network_plugin 2>/dev/null || echo "Network plugin not available"
	terraform output cluster_endpoint 2>/dev/null || echo "Endpoint not available"
	@echo ""
	@echo "=== Master Nodes ==="
	terraform output -json master_public_ips 2>/dev/null | jq -r '.value[]' | nl || echo "No master nodes"
	@echo ""
	@echo "=== Worker Nodes ==="
	terraform output -json worker_public_ips 2>/dev/null | jq -r '.value[]' | nl || echo "No worker nodes"

# Development commands
fmt: ## Format Terraform files
	@echo "Formatting Terraform files..."
	terraform fmt

validate: ## Validate Terraform configuration
	@echo "Validating Terraform configuration..."
	terraform validate

lint: ## Run security and best practices checks
	@echo "Running security checks..."
	@if command -v tfsec >/dev/null 2>&1; then \
		tfsec .; \
	else \
		echo "tfsec not installed. Install with: curl -sfL https://raw.githubusercontent.com/aquasecurity/tfsec/master/scripts/install_linux.sh | sh -s -- -b /usr/local/bin"; \
	fi

# Cluster management
ssh-master: ## SSH to the first master node
	@echo "Connecting to master node..."
	@MASTER_IP=$$(terraform output -json master_public_ips 2>/dev/null | jq -r '.value[0]'); \
	if [ "$$MASTER_IP" = "null" ] || [ -z "$$MASTER_IP" ]; then \
		echo "Master node not found. Deploy the cluster first."; \
		exit 1; \
	fi; \
	ssh-keygen -f "$$HOME/.ssh/k8s-cluster-key" -P "" 2>/dev/null || true; \
	ssh -i "$$HOME/.ssh/k8s-cluster-key" ubuntu@$$MASTER_IP

get-kubeconfig: ## Get kubeconfig from master node
	@echo "Getting kubeconfig from master..."
	@MASTER_IP=$$(terraform output -json master_public_ips 2>/dev/null | jq -r '.value[0]'); \
	if [ "$$MASTER_IP" = "null" ] || [ -z "$$MASTER_IP" ]; then \
		echo "Master node not found. Deploy the cluster first."; \
		exit 1; \
	fi; \
	mkdir -p $$HOME/.kube; \
	scp -i "$$HOME/.ssh/k8s-cluster-key" ubuntu@$$MASTER_IP:/home/ubuntu/.kube/config "$$HOME/.kube/config-$$(terraform output -raw cluster_name 2>/dev/null || echo "k8s-cluster")" 2>/dev/null; \
	sed -i.bak "s/127.0.0.1:6443/$$MASTER_IP:6443/g" "$$HOME/.kube/config-$$(terraform output -raw cluster_name 2>/dev/null || echo "k8s-cluster")" 2>/dev/null; \
	echo "Kubeconfig saved to $$HOME/.kube/config-$$(terraform output -raw cluster_name 2>/dev/null || echo "k8s-cluster")"

# Scaling operations
scale-workers: ## Scale worker nodes (usage: make scale-workers N=3)
	@if [ -z "$(N)" ]; then \
		echo "Usage: make scale-workers N=<number_of_workers>"; \
		exit 1; \
	fi
	@echo "Scaling workers to $(N)..."
	@sed -i "s/worker_count = .*/worker_count = $(N)/" terraform.tfvars
	$(MAKE) plan apply

scale-masters: ## Scale master nodes (usage: make scale-masters N=3)
	@if [ -z "$(N)" ]; then \
		echo "Usage: make scale-masters N=<number_of_masters>"; \
		exit 1; \
	fi
	@echo "Scaling masters to $(N)..."
	@sed -i "s/master_count = .*/master_count = $(N)/" terraform.tfvars
	$(MAKE) plan apply

# Backup and restore
backup: ## Create cluster backup
	@echo "Creating cluster backup..."
	@BACKUP_DIR="$$HOME/k8s-backups/$$(date +%Y%m%d-%H%M%S)"; \
	mkdir -p $$BACKUP_DIR; \
	if [ -f "terraform.tfstate" ]; then cp terraform.tfstate $$BACKUP_DIR/; fi; \
	$(MAKE) get-kubeconfig; \
	cp "$$HOME/.kube/config-$$(terraform output -raw cluster_name 2>/dev/null || echo "k8s-cluster")" $$BACKUP_DIR/ 2>/dev/null || true; \
	echo "Backup created at $$BACKUP_DIR"

# Cleanup
clean: ## Clean up temporary files
	@echo "Cleaning up temporary files..."
	@rm -f *.tfplan
	@rm -f *.backup
	@rm -f .terraform.lock.hcl
	@rm -rf .terraform/

full-clean: clean ## Clean up everything including state (DANGEROUS)
	@echo "WARNING: This will remove all state files!"
	@read -p "Are you sure? (y/N): " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		rm -f *.tfstate*; \
		rm -f crash.log*; \
		echo "All state files removed."; \
	else \
		echo "Cancelled."; \
	fi

# Development setup
setup-dev: ## Set up development environment
	@echo "Setting up development environment..."
	@if ! command -v terraform >/dev/null 2>&1; then \
		echo "Installing Terraform..."; \
		curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -; \
		sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $$(lsb_release -cs) main"; \
		sudo apt-get update && sudo apt-get install terraform; \
	fi
	@if ! command -v jq >/dev/null 2>&1; then \
		echo "Installing jq..."; \
		sudo apt-get install -y jq; \
	fi
	@if ! command -v tfsec >/dev/null 2>&1; then \
		echo "Installing tfsec..."; \
		curl -sfL https://raw.githubusercontent.com/aquasecurity/tfsec/master/scripts/install_linux.sh | sh -s -- -b /usr/local/bin; \
	fi
	$(MAKE) init

# Quick deployment
quick-deploy: init fmt validate apply ## Quick deployment with all checks

# Quick destroy
quick-destroy: status destroy ## Quick destroy with status check

# Show outputs
outputs: ## Show all Terraform outputs
	@terraform output

# Show sensitive outputs
sensitive-outputs: ## Show sensitive outputs (tokens, etc.)
	@terraform output kubeadm_join_command
	@terraform output kubeadm_join_token
	@terraform output kubeadm_ca_cert_hash

# Install dependencies
deps: ## Install required dependencies
	@echo "Installing dependencies..."
	@terraform init -upgrade

# Check for updates
check-updates: ## Check for provider updates
	@echo "Checking for provider updates..."
	@terraform providers

# Security scan
security-scan: lint ## Run comprehensive security scan
	@if command -v checkov >/dev/null 2>&1; then \
		echo "Running Checkov..."; \
		checkov -d .; \
	else \
		echo "Checkov not installed. Install with: pip install checkov"; \
	fi
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
