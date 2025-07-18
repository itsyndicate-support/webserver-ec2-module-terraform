# Secure GitHub Actions Pipeline for Infrastructure Deployment

## How It Works

The workflows are located in `.github/workflows/` and include:

### `ci.yaml`

Runs on every push and pull request to `main` branch. It has next steps:

1. **tfsec**: Scans Terraform code for security misconfigurations.
2. **tflint**: Lints Terraform files to ensure best practices and catch potential issues.
3. **Terraform FMT Check**: Checks that code is formatted properly.
4. **Terraform Validate**: Validates Terraform syntax and logic.
5. **Terraform Plan**: Generates an execution plan (if all above steps succeed).

> If any step fails, the workflow stops immediately.

### `cd-apply.yaml`

Triggered manually using `workflow_dispatch`. It performs:

1. **Terraform Apply**: Applies infrastructure using `terraform apply -auto-approve`.
2. **Terratest**: Runs integration tests to verify resources were correctly created.

### `cd-destroy.yaml`

Triggered manually using `workflow_dispatch`. It performs:

1. **Terraform Destroy**: Destroys infrastructure using `terraform destroy -auto-approve`.

## Terratest Implementation

The test suite (in Go) performs the following checks after apply:

- Confirms EC2 instances and RDS database are created.
- Validates VPC and subnet CIDR blocks.
- Asserts that the database has no public IP and is not publicly accessible.

Tests are located in the `test/` directory and automatically run after infrastructure is applied.

## Challenges & Solutions

- **Rate-limiting on TFLint plugin install**: GitHub API limits were exceeded during CI. Resolved by using a `GITHUB_TOKEN` during `tflint --init`.
- **Post-apply testing**: Added Terratest and ensured it runs only after a successful apply.
- **State Locking**: Ensured workflows are separated cleanly to prevent conflicts when using S3 backend with locking enabled.

