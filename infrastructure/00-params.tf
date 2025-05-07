# parms file for aws ec2 cloud

#### VPC Network
variable "vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

#### HTTP PARAMS
variable "network_http" {
  type = map(string)
  default = {
    subnet_name = "subnet_http"
    cidr        = "192.168.1.0/24"
  }
}

# Set number of instance
variable "http_instance_names" {
  type    = set(string)
  default = ["instance-http-1", "instance-http-2"]
}

#### DB PARAMS
variable "network_db" {
  type = map(string)
  default = {
    subnet_name = "subnet_db"
    cidr        = "192.168.2.0/24"
  }
}

# Set number of instance
variable "db_instance_names" {
  type    = set(string)
  default = ["instance-db-1", "instance-db-2", "instance-db-3"]
}

# Set number of instance
variable "public_key" {
  type    = string
  description = "SSH public key to login into EC2 instance"
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQChY1YyccpHHpZzKMz3WsxLpzg1QIaoi6K/67y/Vy1/VPkliVIJE0hydQW1Mp3HElgucZM9kKv86zbfrzzCcQHS7PttTFeIEhgl4r4Gr79ZSPVU3Mn6gzHpdpkj7SCUD9a0R5SSmOT8BAcZZDsDWVGgcbUuEkR4zBCurDw+VKfcEpUTRSgwEuh1bYyQd6FDE1vLztTP7GWc4c59wvNhYBDpYbBtlGOazmWioTCtvdY/zdm56G3KcnV+Q4fYApTlo81KgdpAGGz7rfJ8DoQWkOnIWz9xt3TDMCIhOOIuwd0Le6qUt2eBiTgUC6Brsi8arVGB1Ar7GZKtBOE3N3I5xUntpDaWsawAmXaU63Q6cumQEpc7sMcKRUbTWhKV1NhbfzG31XwcI1HhWkxYVm4Yr4Zr86DPwQm5hepCy3+gUUFNaCvHPDLoR/MU3bChsC4/w1Nn1PgnQ/1tOn8lRYJYCNrUW/4kABnxo6escRSyKkIlr6rGOdq0coc+wEpQhS1BG7E5qgujtQhRQgLcRYPQfmEULWitpgjxbOUwKYnMWo0XZrxO7T6VzIBuoi1nXmPMlnfGJOuAseCdcgw9mLlNLyZKYBn5TSqsxnhM+JQdObSASAU2am+XZlHuX0RHuqxNb/6RPDzwwhjTWmnJZVfUGtzfdKLkjnWMh1QSCzCdMMHaGw== volodymyrvoitovych11@gmail.com"
}

variable "public_key_name" {
  type    = string
  description = "SSH public key name "
  default = "main-key"
}
variable "env" {
  type = string
  default = "production"
}