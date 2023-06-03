// Display DNS information

// HTTP instance private IP addresses
output "http_ip" {
  value = {
    for instance in aws_instance.http :
    instance.id => instance.private_ip
  }
}

// DB instance private IP addresses
output "db_ip" {
  value = {
    for instance in aws_instance.db :
    instance.id => instance.private_ip
  }
}

// HTTP instance names
output "http_instance_name1" {
  value = aws_instance.http["instance-http-1"].tags.Name
}

output "http_instance_name2" {
  value = aws_instance.http["instance-http-2"].tags.Name
}

// DB instance names
output "db_instance_name1" {
  value = aws_instance.db["instance-db-1"].tags.Name
}

output "db_instance_name2" {
  value = aws_instance.db["instance-db-2"].tags.Name
}

output "db_instance_name3" {
  value = aws_instance.db["instance-db-3"].tags.Name
}

// VPC CIDR block
output "vpc_cidr" {
  value = aws_vpc.terraform.cidr_block
}

// HTTP subnet CIDR block
output "http_subnet_cidr" {
  value = aws_subnet.http.cidr_block
}

// DB subnet CIDR block
output "db_subnet_cidr" {
  value = aws_subnet.db.cidr_block
}

// DB instance IDs
output "db1_id" {
  value = aws_instance.db["instance-db-1"].id
}

output "db2_id" {
  value = aws_instance.db["instance-db-2"].id
}

output "db3_id" {
  value = aws_instance.db["instance-db-3"].id
}

// Http instance IDs
output "http1_id" {
  value = aws_instance.http["instance-http-1"].id
}

output "http2_id" {
  value = aws_instance.http["instance-http-2"].id
}
