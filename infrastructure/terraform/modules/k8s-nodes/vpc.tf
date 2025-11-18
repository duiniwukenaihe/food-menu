# VPC Resources (only if create_vpc is true)
resource "aws_vpc" "main" {
  count = var.create_vpc ? 1 : 0
  
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-vpc"
  })
}

resource "aws_internet_gateway" "main" {
  count = var.create_vpc ? 1 : 0
  
  vpc_id = aws_vpc.main[0].id
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-igw"
  })
}

resource "aws_eip" "nat" {
  count = var.create_vpc ? length(var.public_subnet_cidrs) : 0
  
  domain = "vpc"
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-nat-${count.index}"
  })
  
  depends_on = [aws_internet_gateway.main]
}

resource "aws_nat_gateway" "main" {
  count = var.create_vpc ? length(var.public_subnet_cidrs) : 0
  
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-nat-${count.index}"
  })
  
  depends_on = [aws_internet_gateway.main]
}

# Public Subnets
resource "aws_subnet" "public" {
  count = var.create_vpc ? length(var.public_subnet_cidrs) : 0
  
  vpc_id                  = aws_vpc.main[0].id
  cidr_block              = var.public_subnet_cidrs[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-public-${count.index}"
    Type = "Public"
  })
}

# Private Subnets
resource "aws_subnet" "private" {
  count = var.create_vpc ? length(var.private_subnet_cidrs) : 0
  
  vpc_id            = aws_vpc.main[0].id
  cidr_block        = var.private_subnet_cidrs[count.index]
  availability_zone = var.availability_zones[count.index]
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-private-${count.index}"
    Type = "Private"
  })
}

# Route Tables
resource "aws_route_table" "public" {
  count = var.create_vpc ? 1 : 0
  
  vpc_id = aws_vpc.main[0].id
  
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main[0].id
  }
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-public-rt"
  })
}

resource "aws_route_table" "private" {
  count = var.create_vpc ? length(var.private_subnet_cidrs) : 0
  
  vpc_id = aws_vpc.main[0].id
  
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }
  
  tags = merge(var.tags, {
    Name = "${var.cluster_name}-private-rt-${count.index}"
  })
}

# Route Table Associations
resource "aws_route_table_association" "public" {
  count = var.create_vpc ? length(var.public_subnet_cidrs) : 0
  
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public[0].id
}

resource "aws_route_table_association" "private" {
  count = var.create_vpc ? length(var.private_subnet_cidrs) : 0
  
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}