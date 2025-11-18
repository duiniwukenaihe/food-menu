# Worker Pool Module

This module will create and configure Kubernetes worker nodes.

## Planned Features

- Clone VMs from template
- Configure network settings
- Join workers to Kubernetes cluster
- Configure node labels and taints
- Support for auto-scaling (optional)
- Output node information

## To Be Implemented

Module files to be created:
- `main.tf` - Resource definitions
- `variables.tf` - Input variables
- `outputs.tf` - Output values (node IPs, names, etc.)
- `cloud-init.tf` - Cloud-init configuration (optional)
