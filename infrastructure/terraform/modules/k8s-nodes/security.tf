# Security Group for Kubernetes nodes
resource "aws_security_group" "k8s" {
  name_prefix = "${var.cluster_name}-k8s-sg"
  vpc_id      = var.create_vpc ? aws_vpc.main[0].id : var.vpc_id
  
  # SSH access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH access"
  }
  
  # Kubernetes API server
  ingress {
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = var.create_vpc ? [aws_vpc.main[0].cidr_block] : ["10.0.0.0/16"]
    description = "Kubernetes API server"
  }
  
  # etcd
  ingress {
    from_port   = 2379
    to_port     = 2380
    protocol    = "tcp"
    security_groups = [aws_security_group.k8s.id]
    description = "etcd server client API"
  }
  
  # Kubelet API
  ingress {
    from_port   = 10250
    to_port     = 10250
    protocol    = "tcp"
    security_groups = [aws_security_group.k8s.id]
    description = "Kubelet API"
  }
  
  # kube-scheduler
  ingress {
    from_port   = 10259
    to_port     = 10259
    protocol    = "tcp"
    security_groups = [aws_security_group.k8s.id]
    description = "kube-scheduler"
  }
  
  # kube-controller-manager
  ingress {
    from_port   = 10257
    to_port     = 10257
    protocol    = "tcp"
    security_groups = [aws_security_group.k8s.id]
    description = "kube-controller-manager"
  }
  
  # NodePort Services
  ingress {
    from_port   = 30000
    to_port     = 32767
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "NodePort Services"
  }
  
  # Container Network Interface (CNI)
  ingress {
    from_port   = 8285
    to_port     = 8285
    protocol    = "udp"
    security_groups = [aws_security_group.k8s.id]
    description = "Flannel VXLAN"
  }
  
  ingress {
    from_port   = 8472
    to_port     = 8472
    protocol    = "udp"
    security_groups = [aws_security_group.k8s.id]
    description = "Flannel VXLAN alternative"
  }
  
  # Cilium
  ingress {
    from_port   = 4240
    to_port     = 4240
    protocol    = "tcp"
    security_groups = [aws_security_group.k8s.id]
    description = "Cilium health check"
  }
  
  # Calico
  ingress {
    from_port   = 179
    to_port     = 179
    protocol    = "tcp"
    security_groups = [aws_security_group.k8s.id]
    description = "Calico BGP"
  }
  
  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-k8s-sg"
  })
  
  lifecycle {
    create_before_destroy = true
  }
}