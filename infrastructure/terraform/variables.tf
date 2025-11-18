variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "k8s-cluster"
}

variable "region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24"]
}

variable "master_instance_type" {
  description = "EC2 instance type for master nodes"
  type        = string
  default     = "t3.medium"
}

variable "worker_instance_type" {
  description = "EC2 instance type for worker nodes"
  type        = string
  default     = "t3.large"
}

variable "master_count" {
  description = "Number of master nodes"
  type        = number
  default     = 1
}

variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
}

variable "ssh_key_name" {
  description = "SSH key pair name for EC2 instances"
  type        = string
  default     = "k8s-keypair"
}

variable "ssh_public_key" {
  description = "SSH public key content"
  type        = string
  default     = ""
}

variable "network_plugin" {
  description = "Network plugin for Kubernetes (flannel, cilium, calico)"
  type        = string
  default     = "flannel"
  
  validation {
    condition = contains(["flannel", "cilium", "calico"], var.network_plugin)
    error_message = "The network_plugin must be one of: flannel, cilium, calico."
  }
}

variable "pod_network_cidr" {
  description = "Pod network CIDR"
  type        = string
  default     = "10.244.0.0/16"
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.30.0"
}

variable "tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    "Project"     = "kubernetes-cluster"
    "Environment" = "development"
    "ManagedBy"   = "terraform"
  }
}

variable "enable_monitoring" {
  description = "Enable CloudWatch monitoring for instances"
  type        = bool
  default     = true
}

variable "enable_ebs_optimization" {
  description = "Enable EBS optimization for instances"
  type        = bool
  default     = true
}

variable "root_volume_size" {
  description = "Root volume size in GB"
  type        = number
  default     = 30
}

variable "root_volume_type" {
  description = "Root volume type"
  type        = string
  default     = "gp3"
}

variable "data_volume_size" {
  description = "Data volume size in GB for worker nodes"
  type        = number
  default     = 50
}

variable "data_volume_type" {
  description = "Data volume type for worker nodes"
  type        = string
  default     = "gp3"
}

variable "availability_zones" {
  description = "Availability zones for subnets"
  type        = list(string)
  default     = []
}