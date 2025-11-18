# SSH Key Pair
resource "aws_key_pair" "k8s" {
  count = var.create_vpc && var.ssh_public_key != "" ? 1 : 0
  
  key_name   = var.ssh_key_name
  public_key = var.ssh_public_key
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-keypair"
  })
}

# Master Nodes
resource "aws_instance" "master" {
  count = var.master_count
  
  ami                    = var.ami_id
  instance_type          = var.master_instance_type
  key_name               = var.create_vpc ? aws_key_pair.k8s[0].key_name : var.ssh_key_name
  subnet_id              = var.create_vpc ? aws_subnet.public[count.index % length(aws_subnet.public)].id : var.public_subnet_ids[count.index % length(var.public_subnet_ids)]
  vpc_security_group_ids = [aws_security_group.k8s.id]
  monitoring             = var.enable_monitoring
  ebs_optimized          = var.enable_ebs_optimization
  
  root_block_device {
    volume_size           = var.root_volume_size
    volume_type           = var.root_volume_type
    delete_on_termination = true
    encrypted             = true
  }
  
  user_data = templatefile("${path.module}/../../cloud-init/master.yaml.tpl", {
    network_plugin             = var.network_plugin
    pod_network_cidr          = var.pod_network_cidr
    api_server_advertise_address = "0.0.0.0"
  })
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-master-${count.index}"
    Role = "master"
  })
  
  depends_on = [aws_security_group.k8s]
}

# Worker Nodes
resource "aws_instance" "worker" {
  count = var.worker_count
  
  ami                    = var.ami_id
  instance_type          = var.worker_instance_type
  key_name               = var.create_vpc ? aws_key_pair.k8s[0].key_name : var.ssh_key_name
  subnet_id              = var.create_vpc ? aws_subnet.public[count.index % length(aws_subnet.public)].id : var.public_subnet_ids[count.index % length(var.public_subnet_ids)]
  vpc_security_group_ids = [aws_security_group.k8s.id]
  monitoring             = var.enable_monitoring
  ebs_optimized          = var.enable_ebs_optimization
  
  root_block_device {
    volume_size           = var.root_volume_size
    volume_type           = var.root_volume_type
    delete_on_termination = true
    encrypted             = true
  }
  
  # Additional data volume for worker nodes
  ebs_block_device {
    device_name           = "/dev/sdf"
    volume_size           = var.data_volume_size
    volume_type           = var.data_volume_type
    delete_on_termination = true
    encrypted             = true
  }
  
  user_data = templatefile("${path.module}/../../cloud-init/worker.yaml.tpl", {
    join_command = "sudo kubeadm join ${aws_instance.master[0].private_ip}:6443 --token ${local.kubeadm_join_token} --discovery-token-ca-cert-hash ${local.kubeadm_ca_cert_hash}"
  })
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-worker-${count.index}"
    Role = "worker"
  })
  
  depends_on = [aws_instance.master, aws_security_group.k8s]
}

# Wait for master to be ready before creating workers
resource "null_resource" "wait_for_master" {
  count = var.master_count > 0 ? 1 : 0
  
  triggers = {
    master_ids = join(",", aws_instance.master[*].id)
  }
  
  provisioner "remote-exec" {
    connection {
      host        = aws_instance.master[0].public_ip
      type        = "ssh"
      user        = "ubuntu"
      private_key = var.create_vpc ? tls_private_key.k8s[0].private_key : file("~/.ssh/${var.ssh_key_name}")
    }
    
    inline = [
      "while [ ! -f /tmp/kubeadm-join-token ]; do echo 'Waiting for master to initialize...'; sleep 10; done",
      "echo 'Master is ready!'"
    ]
  }
}

# Local variables for join token and cert hash
locals {
  kubeadm_join_token = "abcdef.0123456789abcdef"
  kubeadm_ca_cert_hash = "sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
}