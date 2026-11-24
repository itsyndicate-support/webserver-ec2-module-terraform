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
  subnet_id                   = aws_subnet.db.id
  associate_public_ip_address = false
  user_data                   = file("scripts/first-boot-db.sh")
  user_data_replace_on_change = true

  depends_on = [aws_route_table_association.db]
  metadata_options {
    http_tokens = "required"
  }

  root_block_device {
    encrypted = true
  }
  tags = {
    Name = each.key
  }
}
