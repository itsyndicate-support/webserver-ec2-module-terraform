# Display dns information

output "http_ip" {
  value = {
    for instance in aws_instance.http :
    instance.id => instance.private_ip
  }
}

output "db_ip" {
  value = {
    for instance in aws_instance.db :
    instance.id => instance.private_ip
  }
}

output "vpc_id" {
  value = aws_vpc.terraform.id
}

output "http_subnet_id" {
  value = aws_subnet.http.id
}

output "db_subnet_id" {
  value = aws_subnet.db.id
}

output "db_security_group_id" {
  value = aws_security_group.db.id
}

output "db_route_table_id" {
  value = aws_route_table.private.id
}
