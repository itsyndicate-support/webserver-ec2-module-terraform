# Display dns information

output "http_ip" {
  value = [for instance in aws_instance.http : instance.public_ip]
}

output "db_ip" {
  value = [for instance in aws_instance.db : instance.public_ip]
}

output "http_subnet" {
  value = aws_subnet.http.id
}

output "db_subnet" {
  value = aws_subnet.db.id
}

output "vpc" {
  value = aws_vpc.terraform.id
}

