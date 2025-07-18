# parms file for aws ec2 cloud

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "administrator_ip" {
  type    = string
  default = "141.138.97.234/32"
}

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
    az          = "us-east-1a"
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
    az          = "us-east-1a"
  }
}

# Set number of instance
variable "db_instance_names" {
  type    = set(string)
  default = ["instance-db-1", "instance-db-2", "instance-db-3"]
}

# Set number of instance
variable "public_key" {
  type        = string
  description = "SSH public key to login into EC2 instance"
  default     = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQyBc2y0+lf5P5r+TqOqihPoJsp0she5qcSXoGJ+XAov59LGfSeTabaDKXC9R2Jw04xNNy3B6vQvKVXpqZWyz7kRumKDhu2rJiS3zjWBP1wcGyNIFdi2frzqdLeQG4EN8WjPuQOh86N1QhAbhnVD4WnHgNRVwH62ll101JHpcYdNoHRnBAzAlsUXx1/l3QajywR29MxW03fgaiLyJRgngzWGdnsdkmrcYde2ZceUEXEgXHITXogUqVgy2E3nr/TcrMehwlHk1XN1b29QzHA5p+rOINMyWFBdNeVfpA1VL5NAKgaiABRZeMUsPkaLlEt09TfCWuw9TlXNOJRc44FcRZ key-us-east-1"
}

