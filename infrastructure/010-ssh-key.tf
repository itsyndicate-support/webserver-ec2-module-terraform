# Define ssh to config in instance

# Create default ssh publique key
resource "aws_key_pair" "user_key" {
  key_name   = "user-key"
  public_key = file("~/.ssh/spacelink.pub")
}

