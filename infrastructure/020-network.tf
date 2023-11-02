# Network configuration

# VPC creation
#tfsec:ignore:aws-ec2-require-vpc-flow-logs-for-all-vpcs
resource "aws_vpc" "terraform" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  tags = {
    Name = "${var.environment}-vpc-http"
  }
}

# http subnet configuration
resource "aws_subnet" "http" {
  vpc_id     = aws_vpc.terraform.id
  cidr_block = var.network_http["cidr"]
  #tfsec:ignore:aws-ec2-no-public-ip-subnet
  map_public_ip_on_launch = true
  tags = {
    Name = "${var.environment}-subnet-http"
  }
  depends_on = [aws_internet_gateway.gw]
}

# db subnet configuration
resource "aws_subnet" "db" {
  vpc_id     = aws_vpc.terraform.id
  cidr_block = var.network_db["cidr"]
  tags = {
    Name = "${var.environment}-subnet-db"
  }
  depends_on = [aws_internet_gateway.gw]
}

# External gateway configuration
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.terraform.id
  tags = {
    Name = "${var.environment}-internet-gateway"
  }
}

