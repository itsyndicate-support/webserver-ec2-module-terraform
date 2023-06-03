# Security group configuration

# Default administration port
resource "aws_security_group" "administration" {
  name        = "administration"
  description = "Allow default administration service"
  vpc_id      = aws_vpc.terraform.id
  tags = {
    Name = "administration"
  }

  # Open ssh port
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["91.214.136.165/32"]
  }

  # Allow icmp
  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["91.214.136.165/32"]
  }

  # Open access to public network
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["91.214.136.165/32"]
  }
}

# Open web port
resource "aws_security_group" "web" {
  name        = "web"
  description = "Allow web incgress trafic"
  vpc_id      = aws_vpc.terraform.id
  tags = {
    Name = "${var.env}-web"
  }

  # http port
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["91.214.136.165/32"]
  }

  # https port
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["91.214.136.165/32"]
  }

  # Open access to public network
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["91.214.136.165/32"]
  }
}

# Open database port
resource "aws_security_group" "db" {
  name        = "db"
  description = "Allow db incgress trafic"
  vpc_id      = aws_vpc.terraform.id
  tags = {
    Name = "${var.env}-db"
  }

  # db port
  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = ["91.214.136.165/32"]
  }

  # Open access to public network
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["91.214.136.165/32"]
  }
}

