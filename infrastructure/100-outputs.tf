# Display dns information

output "http_private_ip" {
  value = {
    for instance in aws_instance.http :
    instance.id => instance.private_ip
  }
}

output "db_private_ip" {
  value = {
    for instance in aws_instance.db :
    instance.id => instance.private_ip
  }
}

output "db_public_ip" {
  value = {
    for instance in aws_instance.db :
    instance.id => instance.public_ip
  }
}

# ID VPC
output "vpc_id" {
  value = aws_vpc.terraform.id
}

# CIDR of VPC
output "vpc_cidr" {
  value = aws_vpc.terraform.cidr_block
}

# CIDR of HTTP subnets
output "http_subnet_cidrs" {
  value = aws_subnet.http.cidr_block
}

# CIDR of DB subnets
output "db_subnet_cidrs" {
  value = aws_subnet.db.cidr_block
}
