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

output "instance_ids" {
  value = values(aws_instance.http)[*].id
}

output "vpc_id" {
  value = aws_vpc.terraform.id
}

output "db_instance_ids" {
  value = values(aws_instance.db)[*].id
}