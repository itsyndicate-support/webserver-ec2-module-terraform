# Here it needs KMS key to encrypt this log group
#tfsec:ignore:aws-cloudwatch-log-group-customer-key
resource "aws_cloudwatch_log_group" "vpc_logs" {
  name              = "/vpc/flow-logs"
  retention_in_days = 14
}

resource "aws_flow_log" "vpc_flow_logs" {
  log_destination_type = "cloud-watch-logs" # or "s3"
  traffic_type         = "ALL"              # ACCEPT, REJECT, or ALL
  vpc_id               = aws_vpc.terraform.id

  log_destination = aws_cloudwatch_log_group.vpc_logs.arn
  iam_role_arn    = aws_iam_role.vpc_flow_logs_role.arn

  tags = {
    Name = "vpc-flow-logs"
  }
}
