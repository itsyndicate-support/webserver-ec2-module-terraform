variable "prefix" {
  description = "Prefix used for resource naming."
  type        = string
  default = "my-"
}

variable "aws_region" {
  description = "AWS region to create the VPC, assuming zones `a` and `b` exists."
  type        = string
  default = "us-east-1"
}

variable "instance_type" {
  type = string
  default = "t3.small"
}

variable "path_to_private_key" {
  type = string
  default = "./my-org-runner.2025-07-25.private-key.pem"
}

variable "github_app" {
  description = "GitHub for API usages."

  type = object({
    id         = string
    key_base64 = string
  })

  default = {
    id         = "1664700"
    key_base64 = "default-base64-key"
  }
}

variable "ami_filter_info" {
  type = object({
    name = string
    owner = string
  })

  default = {
    name = "github-runner-al2023-x86_64-202507271355"
    owner = "104280348236"
  }
}
