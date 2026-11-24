# Task 1 — Terraform infrastructure pipeline

## Implementation

The workflow is defined in .github/workflows/terraform.yml.

The checks job runs terraform fmt, initializes the provider, and runs
terraform validate. The terraform job depends on successful checks.

Pushes to main and week11-* branches run terraform plan.
Pull requests to main run configuration checks.
Manual runs provide plan, apply, and destroy choices.

The workflow uses Terraform 1.16.2. AWS credentials are supplied through
the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY repository secrets.

Terraform state is stored in S3 in us-east-1. The backend uses encryption,
bucket versioning, and S3 lock files. AWS operations share a concurrency
group to coordinate runs against the same state.

## Usage

Open Actions, select Terraform CI, and choose Run workflow.
Select the project branch and choose plan, apply, or destroy.

An apply run generates a saved plan and applies it in the same job.
A destroy run removes the infrastructure tracked in the project state.

After all infrastructure work is complete, remove the versions and delete
markers from the state bucket, then delete the bucket.

## Challenges and resolutions

The initial provider lock file caused a checksum error on Linux.
Checksums for darwin_arm64 and linux_amd64 were added.

The original eu-north-1 region was unavailable to the lab user.
The provider region was changed to us-east-1.

Creating the state bucket initially failed with AccessDenied.
A bucket-scoped IAM policy was added for setup, state access, and cleanup.

## Verification

Formatting, validation, backend initialization, and planning have passed
in GitHub Actions.

For deployment verification, inspect the manual apply run for
"Apply complete", then inspect the created EC2 instances and network.
For cleanup verification, inspect the manual destroy run for
"Destroy complete".
