terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
  
  default_tags {
    tags = var.tags
  }
}

# Data sources
data "aws_ami" "ubuntu" {
  most_recent = true
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
  
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  
  owners = ["099720109477"] # Canonical
}

# Calculate availability zones if not provided
locals {
  availability_zones = length(var.availability_zones) > 0 ? var.availability_zones : slice(data.aws_availability_zones.available.names, 0, 2)
}

data "aws_availability_zones" "available" {
  state = "available"
}

# VPC Module
module "vpc" {
  source = "./modules/k8s-nodes"
  
  cluster_name           = var.cluster_name
  vpc_cidr              = var.vpc_cidr
  public_subnet_cidrs   = var.public_subnet_cidrs
  private_subnet_cidrs  = var.private_subnet_cidrs
  availability_zones    = local.availability_zones
  tags                  = var.tags
  
  create_vpc = true
}

# Kubernetes Nodes Module
module "k8s_nodes" {
  source = "./modules/k8s-nodes"
  
  cluster_name            = var.cluster_name
  vpc_id                 = module.vpc.vpc_id
  public_subnet_ids      = module.vpc.public_subnet_ids
  private_subnet_ids     = module.vpc.private_subnet_ids
  availability_zones      = local.availability_zones
  
  master_instance_type   = var.master_instance_type
  worker_instance_type   = var.worker_instance_type
  master_count           = var.master_count
  worker_count           = var.worker_count
  
  ssh_key_name           = var.ssh_key_name
  ssh_public_key         = var.ssh_public_key
  
  network_plugin         = var.network_plugin
  pod_network_cidr       = var.pod_network_cidr
  kubernetes_version     = var.kubernetes_version
  
  enable_monitoring      = var.enable_monitoring
  enable_ebs_optimization = var.enable_ebs_optimization
  root_volume_size       = var.root_volume_size
  root_volume_type       = var.root_volume_type
  data_volume_size       = var.data_volume_size
  data_volume_type       = var.data_volume_type
  
  ami_id                 = data.aws_ami.ubuntu.id
  tags                   = var.tags
  
  create_vpc = false
}