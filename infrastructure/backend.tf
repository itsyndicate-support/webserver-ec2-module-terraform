terraform {
  backend "s3" {
    encrypt      = true
    bucket       = "backend-58490149"
    region       = "eu-north-1"
    key          = "terraform.tfstate"
    use_lockfile = true
  }
}
