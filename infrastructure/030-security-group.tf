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
    description = "Allow SSH from anywhere (temporary)"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    #tfsec:ignore:aws-vpc-no-public-ingress-sgr #In prod will be explicit CIDR(not 0.0.0.0/0)
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow icmp
  ingress {
    description = "Allow ICMP from anywhere"
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    #tfsec:ignore:aws-vpc-no-public-ingress-sgr #In prod will be explicit CIDR(not 0.0.0.0/0)
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Open access to public network
  egress {
    description = "Allow all egress"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    #tfsec:ignore:aws-vpc-no-public-egress-sgr #In prod will be explicit CIDR(not 0.0.0.0/0)
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
  ingress {
    description = "Allow HTTP from internet"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    #tfsec:ignore:aws-vpc-no-public-ingress-sgr #Web server must be reachable over HTTP from public internet
    cidr_blocks = ["0.0.0.0/0"]
  }

  # https port
  ingress {
    description = "Allow HTTPS from internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    #tfsec:ignore:aws-vpc-no-public-ingress-sgr #Web server must be reachable over HTTPS from public internet
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Open access to public network
  egress {
    description = "Allow all egress"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    #tfsec:ignore:aws-vpc-no-public-egress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Open database port
resource "aws_security_group" "db" {
  name        = "db"
  description = "Allow DB ingress traffic"
  vpc_id      = aws_vpc.terraform.id
  tags = {
    Name = "db"
  }

  # db port
  ingress {
    description     = "Allow MySQL from internet (not secure!)"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.administration.id, aws_security_group.web.id]

  }

  # Open access to public network
  egress {
    description     = "Allow all egress"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [aws_security_group.administration.id, aws_security_group.web.id]
  }
}
