terraform {
  backend "s3" {
    bucket       = "week11-tfstate-047472448571-dmitriy-o"
    key          = "webserver/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    use_lockfile = true
  }
}
