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
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["176.37.225.82/32"]
  }

  # Allow icmp
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["176.37.225.82/32", "192.168.0.0/16"]
  }

  # Open access to public network
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    #tfsec:ignore:aws-ec2-no-public-egress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Open web port
resource "aws_security_group" "web" {
  name        = "web"
  description = "Allow web incgress trafic"
  vpc_id      = aws_vpc.terraform.id
  tags = {
    Name = "web"
  }

  # http port
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    #tfsec:ignore:aws-ec2-no-public-ingress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }

  # https port
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    #tfsec:ignore:aws-ec2-no-public-ingress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Open access to public network
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    #tfsec:ignore:aws-ec2-no-public-egress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Open database port
resource "aws_security_group" "db" {
  name        = "db"
  description = "Allow db incgress trafic"
  vpc_id      = aws_vpc.terraform.id
  tags = {
    Name = "db"
  }

  # db port
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = ["192.168.0.0/16"]
  }

  # Open access to public network
  #tfsec:ignore:aws-ec2-add-description-to-security-group-rule
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    #tfsec:ignore:aws-ec2-no-public-egress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }
}

