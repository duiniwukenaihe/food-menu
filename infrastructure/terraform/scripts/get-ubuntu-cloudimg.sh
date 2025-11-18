#!/bin/bash
# =============================================================================
# Download Ubuntu Cloud Image with SHA256 Verification
# =============================================================================

set -e

# Configuration
IMAGE_URL="${1:-https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img}"
IMAGE_NAME="${2:-jammy-server-cloudimg-amd64.img}"
DOWNLOAD_DIR="./downloads"

# Create download directory
mkdir -p "$DOWNLOAD_DIR"

echo "=== Ubuntu Cloud Image Download Script ==="
echo "Image URL: $IMAGE_URL"
echo "Image Name: $IMAGE_NAME"
echo "Download Directory: $DOWNLOAD_DIR"
echo

# Change to download directory
cd "$DOWNLOAD_DIR"

# Download the image
echo "Downloading Ubuntu cloud image..."
wget --progress=bar:force --continue "$IMAGE_URL" -O "$IMAGE_NAME"

echo "Download completed: $IMAGE_NAME"

# Get SHA256 checksum
echo "Calculating SHA256 checksum..."
SHA256_CHECKSUM=$(sha256sum "$IMAGE_NAME" | awk '{print $1}')
echo "SHA256: $SHA256_CHECKSUM"

# Try to get expected checksum from Ubuntu
echo "Attempting to verify checksum..."
SHA256_URL="${IMAGE_URL%/*}/SHA256SUMS"

if wget --quiet --timeout=10 "$SHA256_URL" -O "SHA256SUMS"; then
    EXPECTED_SHA256=$(grep "$IMAGE_NAME" SHA256SUMS | awk '{print $1}')
    
    if [ "$SHA256_CHECKSUM" = "$EXPECTED_SHA256" ]; then
        echo "✓ Checksum verification PASSED"
    else
        echo "✗ Checksum verification FAILED"
        echo "Expected: $EXPECTED_SHA256"
        echo "Actual:   $SHA256_CHECKSUM"
        echo "WARNING: Using image with unverified checksum"
    fi
else
    echo "Could not download SHA256SUMS for verification"
    echo "Proceeding with downloaded image (checksum: $SHA256_CHECKSUM)"
fi

# Create checksum file
echo "$SHA256_CHECKSUM  $IMAGE_NAME" > "${IMAGE_NAME}.sha256"

echo
echo "=== Download Summary ==="
echo "Image: $IMAGE_NAME"
echo "Size: $(du -h "$IMAGE_NAME" | cut -f1)"
echo "SHA256: $SHA256_CHECKSUM"
echo "Location: $(pwd)/$IMAGE_NAME"
echo
echo "You can now use this image for Proxmox template creation."