data "aws_region" "current" {}

output "test_config" {
  description = "Resource identifiers and expected configuration for post-apply tests."

  value = {
    region    = data.aws_region.current.region
    vpc_id    = aws_vpc.terraform.id
    vpc_cidr  = var.vpc_cidr
    web_sg_id = aws_security_group.web.id

    subnets = {
      http = {
        id   = aws_subnet.http.id
        cidr = var.network_http["cidr"]
      }
      db = {
        id   = aws_subnet.db.id
        cidr = var.network_db["cidr"]
      }
    }

    instances = {
      http = { for name, instance in aws_instance.http : name => instance.id }
      db   = { for name, instance in aws_instance.db : name => instance.id }
    }
  }
}
