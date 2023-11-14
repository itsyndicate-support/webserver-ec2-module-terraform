# Define ssh to config in instance

# Create default ssh publique key
resource "aws_key_pair" "user_key" {
  key_name   = "user-key"
  public_key = file("${var.pub_key_path}${var.pub_key_filename}")
}

