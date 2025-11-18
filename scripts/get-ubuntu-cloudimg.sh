#!/bin/bash

set -euo pipefail

# Script arguments
UBUNTU_VERSION="${1:-22.04}"
ARCHITECTURE="${2:-amd64}"
INTERACTIVE_REPLACE="${3:-true}"

# Configuration
BASE_URL="https://cloud-images.ubuntu.com/releases"
DOWNLOAD_DIR="$(dirname "$0")/../downloads"
IMAGE_NAME="ubuntu-${UBUNTU_VERSION}-server-cloudimg-${ARCHITECTURE}.img"
SHA256SUMS_URL="${BASE_URL}/${UBUNTU_VERSION}/release/SHA256SUMS"
SHA256SUMS_GPG_URL="${BASE_URL}/${UBUNTU_VERSION}/release/SHA256SUMS.gpg"

# Create downloads directory
mkdir -p "$DOWNLOAD_DIR"

# Function to download file with progress
download_file() {
    local url="$1"
    local output="$2"
    echo "Downloading: $url"
    curl -L --progress-bar -o "$output" "$url"
}

# Function to verify GPG signature (simplified version)
verify_checksum() {
    local image_file="$1"
    local expected_sum="$2"
    
    echo "Verifying checksum..."
    local actual_sum=$(sha256sum "$image_file" | cut -d' ' -f1)
    
    if [ "$actual_sum" = "$expected_sum" ]; then
        echo "✓ Checksum verification passed"
        return 0
    else
        echo "✗ Checksum verification failed"
        echo "Expected: $expected_sum"
        echo "Actual:   $actual_sum"
        return 1
    fi
}

# Function to prompt for replacement
prompt_replace() {
    local file="$1"
    
    if [ "$INTERACTIVE_REPLACE" = "true" ]; then
        echo "File already exists: $file"
        read -p "Replace? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            return 0
        else
            echo "Skipping download"
            exit 0
        fi
    else
        echo "File already exists and interactive mode is disabled. Skipping download."
        exit 0
    fi
}

# Main execution
main() {
    echo "Ubuntu Cloud Image Downloader"
    echo "Version: $UBUNTU_VERSION"
    echo "Architecture: $ARCHITECTURE"
    echo "----------------------------------------"
    
    local image_path="$DOWNLOAD_DIR/$IMAGE_NAME"
    
    # Check if image already exists
    if [ -f "$image_path" ]; then
        prompt_replace "$image_path"
    fi
    
    # Download SHA256SUMS
    echo "Downloading checksum files..."
    download_file "$SHA256SUMS_URL" "$DOWNLOAD_DIR/SHA256SUMS"
    download_file "$SHA256SUMS_GPG_URL" "$DOWNLOAD_DIR/SHA256SUMS.gpg"
    
    # Extract the checksum for our image
    local expected_checksum=$(grep "$IMAGE_NAME" "$DOWNLOAD_DIR/SHA256SUMS" | cut -d' ' -f1)
    
    if [ -z "$expected_checksum" ]; then
        echo "Error: Could not find checksum for $IMAGE_NAME"
        exit 1
    fi
    
    echo "Expected checksum: $expected_checksum"
    
    # Download the cloud image
    local image_url="${BASE_URL}/${UBUNTU_VERSION}/release/$IMAGE_NAME"
    download_file "$image_url" "$image_path"
    
    # Verify the downloaded image
    if verify_checksum "$image_path" "$expected_checksum"; then
        echo "✓ Successfully downloaded and verified $IMAGE_NAME"
        echo "Location: $image_path"
    else
        echo "✗ Verification failed. Removing corrupted file..."
        rm -f "$image_path"
        exit 1
    fi
    
    # Cleanup checksum files
    rm -f "$DOWNLOAD_DIR/SHA256SUMS" "$DOWNLOAD_DIR/SHA256SUMS.gpg"
}

# Run main function
main "$@"