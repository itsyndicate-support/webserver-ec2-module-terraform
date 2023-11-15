terraform {
  required_version = ">= 0.12"
  backend "s3" {
    bucket = "syndicate-tfstate"
    key    = "terraform/state"
    region = "eu-central-1"
  }
}
