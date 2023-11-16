
output "httpIP" {
  value = {
    for instance in aws_instance.http :
    instance.id => instance.private_ip
  }
}

output "dbIP" {
  value = {
    for instance in aws_instance.db :
    instance.id => instance.private_ip
  }
}
output "httpSubnet" {
  value = aws_subnet.http.id
}

output "dbSubnet" {
  value = aws_subnet.db.id
}

output "vpc" {
  value = aws_vpc.terraform.id
}