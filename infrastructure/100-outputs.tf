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


#--------------- Terratests |---------------

output "first_web_server_ssh_key" {
  value = aws_instance.http["instance-http-1"].key_name
}

output "second_web_server_ssh_key" {
  value = aws_instance.http["instance-http-2"].key_name
}

output "first_db_server_ssh_key" {
  value = aws_instance.db["instance-db-1"].key_name
}

output "second_db_server_ssh_key" {
  value = aws_instance.db["instance-db-2"].key_name
}

output "third_db_server_ssh_key" {
  value = aws_instance.db["instance-db-3"].key_name
}

output "vpc_cidr_block" {
  value = aws_vpc.terraform.cidr_block
}

output "http_subnet_cidr_block" {
  value = aws_subnet.http.cidr_block
}

output "db_subnet_cidr_block" {
  value = aws_subnet.db.cidr_block
}

output "first_db_server_public_ip" {
  value = aws_instance.db["instance-db-1"].public_ip
}

output "second_db_server_public_ip" {
  value = aws_instance.db["instance-db-2"].public_ip
}

output "third_db_server_public_ip" {
  value = aws_instance.db["instance-db-3"].public_ip
}

output "first_db_server_id" {
  value = aws_instance.db["instance-db-1"].id
}

output "second_db_server_id" {
  value = aws_instance.db["instance-db-2"].id
}

output "third_db_server_id" {
  value = aws_instance.db["instance-db-3"].id
}
