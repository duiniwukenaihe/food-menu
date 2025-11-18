#!/bin/bash

# Simple validation script for the Ubuntu template module

set -euo pipefail

echo "Validating Ubuntu template module..."
echo "=================================="

# Check if required files exist
required_files=(
    "main.tf"
    "variables.tf"
    "outputs.tf"
    "modules.tf"
    "modules/template/main.tf"
    "modules/template/variables.tf"
    "modules/template/outputs.tf"
    "modules/template/cloud-init-user-data.tpl"
    "modules/template/cloud-init-network-data.tpl"
    "scripts/get-ubuntu-cloudimg.sh"
    "terraform.tfvars.example"
)

echo "Checking required files..."
for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✓ $file exists"
    else
        echo "✗ $file missing"
        exit 1
    fi
done

# Check if script is executable
if [ -x "scripts/get-ubuntu-cloudimg.sh" ]; then
    echo "✓ Download script is executable"
else
    echo "✗ Download script is not executable"
    exit 1
fi

# Check Terraform syntax (if terraform is available)
if command -v terraform >/dev/null 2>&1; then
    echo "Validating Terraform syntax..."
    if terraform fmt -check -recursive >/dev/null 2>&1; then
        echo "✓ Terraform formatting is valid"
    else
        echo "✗ Terraform formatting issues found"
        terraform fmt -recursive
        exit 1
    fi
    
    if terraform validate >/dev/null 2>&1; then
        echo "✓ Terraform configuration is valid"
    else
        echo "✗ Terraform validation failed"
        terraform validate
        exit 1
    fi
else
    echo "⚠ Terraform not found, skipping syntax validation"
fi

echo "=================================="
echo "✓ All validations passed!"
echo "The Ubuntu template module is ready to use."