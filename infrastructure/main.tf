provider "aws" {
  region = "eu-north-1"
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.54.1"
    }
  }

  backend "s3" {
    encrypt      = true
    bucket       = "backend-58490149"
    region       = "eu-north-1"
    key          = "terraform.tfstate"
    use_lockfile = true
  }
  required_version = "~> 1.10"
}