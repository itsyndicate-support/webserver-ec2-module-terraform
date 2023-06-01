provider "aws" {
  region = "eu-central-1"
}

terraform {
  backend "s3" {
    bucket = "my-terraform-state-bucket"
    key    = "terraform.tfstate"
    region = "eu-central-1"
  }
}
