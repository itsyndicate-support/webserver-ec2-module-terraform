provider "aws" {
  region = "eu-central-1"
}

terraform {
  backend "s3" {
    bucket = "project-terraform-tfstate"
    key    = "dev/terraform.tfstate"
    region = "eu-central-1"
  }
}
