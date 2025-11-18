#!/bin/bash
# Test script for Kubernetes Terraform configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test function
test_terraform() {
    log_info "Testing Terraform configuration..."
    
    # Check if terraform is installed
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed"
        return 1
    fi
    
    # Initialize terraform
    log_info "Initializing Terraform..."
    terraform init
    
    # Validate configuration
    log_info "Validating configuration..."
    terraform validate
    
    # Format check
    log_info "Checking formatting..."
    if terraform fmt -check -diff; then
        log_success "All files are properly formatted"
    else
        log_warning "Some files need formatting"
        terraform fmt
    fi
    
    # Plan check
    log_info "Creating execution plan..."
    if terraform plan -var-file="terraform.tfvars.example" -detailed-exitcode; then
        log_success "Plan created successfully"
    elif [ $? -eq 1 ]; then
        log_error "Plan failed"
        return 1
    else
        log_success "Plan created with changes"
    fi
    
    log_success "All tests passed!"
}

# Main execution
cd "$(dirname "$0")"
test_terraform