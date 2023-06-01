provider "aws" {
  region = "eu-central-1"
}

terraform {
  backend "s3" {
    bucket = "state-bucket-circleci"
    key    = "terraform.tfstate"
    region = "eu-central-1"
  }
}
