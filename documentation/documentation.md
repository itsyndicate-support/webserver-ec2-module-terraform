# Simple GitHub Actions Pipeline for Infrastructure Deployment

## How It Works

The workflows are located in `.github/workflows/` and include:

### `validate.yaml`

Runs on every push and pull request еto `main` branch. It has next steps:

1. **Setup Dependencies**: Checks out the code, installs Terraform, sets up AWS credentials.
2. **Terraform FMT Check**: Validates Terraform formatting.
3. **Terraform Validate**: Ensures Terraform configuration is syntactically valid.
4. **Terraform Plan**: Generates an execution plan (read-only).

### `apply.yaml`

Triggered manually using `workflow_dispatch`. It performs:

1. **Terraform Apply**: Applies infrastructure changes with `-auto-approve`.
2. **Terraform Destroy (Optional)**: If `destroy` input is set to true, destroys infrastructure.

## Challenges & Solutions

- **Terraform State Locking**: Terraform Destroy could 
not destroy infra because there was no perms to tfstate and locks. 
So S3 backend was used

