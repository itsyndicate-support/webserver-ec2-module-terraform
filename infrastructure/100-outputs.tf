# Display dns information

output "http_ip" {
  value = [for instance in aws_instance.http : instance.private_ip]
}

output "db_ip" {
  value =  [for instance in aws_instance.db : instance.private_ip]
}

output "vpc_cidr" {
  value = aws_vpc.terraform.cidr_block
}

output "http_subnet_cidr" {
  value = aws_subnet.http.cidr_block
}

output "db_subnet_cidr" {
  value = aws_subnet.db.cidr_block
}
output "db_instance_public_ips" {
  value = [for instance in aws_instance.db : instance.public_ip if instance.public_ip != ""]
}