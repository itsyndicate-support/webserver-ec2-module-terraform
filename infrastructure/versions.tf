
terraform {
  required_version = ">= 1.11"

  backend "s3" {
    bucket       = "m0rtyni-terraform-state-014468369410-eu-north-1-an"
    key          = "webserver-ec2/terraform.tfstate"
    region       = "eu-north-1"
    encrypt      = true
    use_lockfile = true
  }
}
