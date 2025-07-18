terraform {
  backend "s3" {
    bucket       = "w11t1-terraform-state"
    key          = "state"
    region       = "us-east-1"
    use_lockfile = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.3"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}
