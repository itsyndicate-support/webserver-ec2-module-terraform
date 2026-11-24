resource "aws_security_group" "administration" {
  name        = "administration"
  description = "Allow default administration service"
  vpc_id      = aws_vpc.terraform.id

  ingress {
    description = "SSH within the VPC"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  ingress {
    description = "ICMP within the VPC"
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = [var.vpc_cidr]
  }

  # Ubuntu package repositories require outbound HTTP.
  #trivy:ignore:AVD-AWS-0104
  egress {
    description = "HTTP for Ubuntu package repositories"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Ubuntu package repositories require outbound HTTPS.
  #trivy:ignore:AVD-AWS-0104
  egress {
    description = "HTTPS for Ubuntu package repositories"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "administration"
  }
}

resource "aws_security_group" "web" {
  name        = "web"
  description = "Allow web incgress trafic"
  vpc_id      = aws_vpc.terraform.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "MySQL to the database subnet"
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = [var.network_db["cidr"]]
  }

  tags = {
    Name = "web"
  }
}

resource "aws_security_group" "db" {
  name        = "db"
  description = "Allow db incgress trafic"
  vpc_id      = aws_vpc.terraform.id

  ingress {
    description     = "MySQL from web instances"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.web.id]
  }

  egress = []

  tags = {
    Name = "db"
  }
}
