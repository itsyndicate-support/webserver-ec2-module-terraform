#trivy:ignore:AWS-0104
resource "aws_security_group" "administration" {
  name        = "administration"
  description = "Administrative access and outbound connectivity"
  vpc_id      = aws_vpc.terraform.id

  ingress {
    description = "SSH from the administrator IP"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.admin_cidr]
  }

  ingress {
    description = "ICMP echo requests"
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Traffic within the VPC"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    description = "HTTP access to package repositories"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "HTTPS access to package repositories"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.environment}-administration"
    Environment = var.environment
  }
}

resource "aws_security_group" "web" {
  name        = "web"
  description = "Public HTTP and HTTPS access"
  vpc_id      = aws_vpc.terraform.id

  ingress {
    description = "Public HTTP access"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Public HTTPS access"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress = []

  tags = {
    Name        = "${var.environment}-web"
    Environment = var.environment
  }
}

resource "aws_security_group" "db" {
  name        = "db"
  description = "MariaDB access from web instances"
  vpc_id      = aws_vpc.terraform.id

  ingress {
    description     = "MariaDB access from the web security group"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.web.id]
  }

  egress = []

  tags = {
    Name        = "${var.environment}-db"
    Environment = var.environment
  }
}
