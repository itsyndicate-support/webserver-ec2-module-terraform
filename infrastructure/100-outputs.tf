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

output "http_urls" {
  value = {
    for name, eip in aws_eip.public_http :
    name => "http://${eip.public_ip}"
  }
}
