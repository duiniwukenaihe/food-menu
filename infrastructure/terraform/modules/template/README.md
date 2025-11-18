# Template Module

This module will create a VM template from a cloud-init image that serves as the base for all Kubernetes nodes.

## Planned Features

- Download cloud image (Ubuntu/Debian/Rocky Linux)
- Create VM from cloud image
- Configure cloud-init
- Convert VM to template
- Add SSH keys and basic configuration

## To Be Implemented

Module files to be created:
- `main.tf` - Resource definitions
- `variables.tf` - Input variables
- `outputs.tf` - Output values
- `versions.tf` - Version constraints (if needed)
