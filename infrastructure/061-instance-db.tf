#### INSTANCE DB ####

# Create instance
resource "aws_instance" "db" {
  for_each      = var.db_instance_names
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"
  key_name      = aws_key_pair.user_key.key_name
  subnet_id     = aws_subnet.db.id
  user_data     = file("scripts/first-boot-db.sh")

  vpc_security_group_ids = [
    aws_security_group.administration.id,
    aws_security_group.db.id,
  ]

  root_block_device {
    volume_type = "gp3"
    volume_size = 8
    encrypted   = true
  }

  metadata_options {
    http_tokens   = "required"
    http_endpoint = "enabled"
  }

  tags = {
    Name        = each.key
    Environment = var.env
  }
}

