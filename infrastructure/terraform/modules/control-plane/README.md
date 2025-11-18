# Control Plane Module

This module will create and configure Kubernetes control plane (master) nodes.

## Planned Features

- Clone VMs from template
- Configure network settings
- Initialize Kubernetes control plane
- Set up high availability (for multiple masters)
- Configure load balancer for API server (optional)
- Output cluster join tokens and endpoints

## To Be Implemented

Module files to be created:
- `main.tf` - Resource definitions
- `variables.tf` - Input variables
- `outputs.tf` - Output values (node IPs, join tokens, etc.)
- `cloud-init.tf` - Cloud-init configuration (optional)
