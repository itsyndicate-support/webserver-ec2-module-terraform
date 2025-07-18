#### INSTANCE DB ####

# Create instance
resource "aws_instance" "db" {
  for_each      = var.db_instance_names
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"
  key_name      = aws_key_pair.user_key.key_name
  vpc_security_group_ids = [
    aws_security_group.administration.id,
    aws_security_group.db.id,
  ]
  subnet_id = aws_subnet.db.id
  user_data = file("scripts/first-boot-db.sh")

  # Enforce IMDSv2
  metadata_options {
    http_tokens   = "required"
    http_endpoint = "enabled"
  }

  # Encrypt root volume
  root_block_device {
    encrypted   = true
    volume_size = 8
    volume_type = "gp3"
  }

  tags = {
    Name = each.key
  }
}
