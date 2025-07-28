resource "random_id" "random" {
  byte_length = 20
}

module "github-runner" {
  source  = "github-aws-runners/github-runner/aws"
  version = "6.7.1"

  aws_region = var.aws_region
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  github_app = {
    key_base64     = filebase64(var.path_to_private_key)
    id             = var.github_app.id
    webhook_secret = random_id.random.hex
  }

  webhook_lambda_zip                = "lambdas-download/webhook.zip"
  runner_binaries_syncer_lambda_zip = "lambdas-download/runner-binaries-syncer.zip"
  runners_lambda_zip                = "lambdas-download/runners.zip"
  ami_housekeeper_lambda_zip        = "lambdas-download/ami-housekeeper.zip"
  
  create_service_linked_role_spot = true
  enable_organization_runners = true
  enable_ssm_on_runners = true

  instance_types = [var.instance_type]

  enable_ephemeral_runners = true
  enable_job_queued_check = true

  delay_webhook_event   = 0
  runners_maximum_count = 10

  enable_user_data_debug_logging_runner = true

  eventbridge = {
    enable = false
  }

  ami = {
    filter = {
      name   = [var.ami_filter_info.name]
      state  = ["available"]
    }
    owners = [var.ami_filter_info.owner]
  }

  enable_userdata = false
}
