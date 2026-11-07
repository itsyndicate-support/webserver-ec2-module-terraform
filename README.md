## How to

![infra instance with multiple network](./img/03-multiple-network.png "infra instance with multiple network")

### Create stack

```bash
terraform apply
```

This script will create:
- 1 VPC
- 2 networks
- 2 HTTP instances
- 3 DB instances

### Delete stack

```bash
terraform destroy
```

### GitHub Actions Pipeline

The repository uses GitHub Actions to automate Terraform validation, security checks, planning, and deployment.

**Workflows:**

- **Terraform CI** – runs `terraform fmt`, `terraform init`, and `terraform validate` on pushes and pull requests.
- **Terraform Pipeline** – runs Terraform validation, TFLint, and Trivy security scans before generating and uploading a Terraform Plan artifact.
- **Terraform Apply** – manually triggered workflow that creates or updates the infrastructure and runs Terratest after a successful apply.

**Security and Code Quality**

- **TFLint** checks Terraform code for best-practice issues.
- **Trivy** scans Terraform configuration for security misconfigurations.
- Terraform Plan runs only if TFLint and Trivy succeed.

**Terratest**

After Terraform Apply, Terratest verifies:
- Expected HTTP and DB EC2 instance counts and running state.
- VPC and subnet CIDR ranges.
- DB instances have no public IP addresses.
- DB Security Group does not allow public MySQL access.

**AWS Authentication**

GitHub Actions authenticates to AWS using OpenID Connect (OIDC) and an IAM role. No long-lived AWS access keys are stored in GitHub.

**Terraform Remote State**

Terraform state is stored in an Amazon S3 bucket in the `eu-north-1` region. State locking is enabled using S3 lockfiles.

**How to use the pipeline:**

1. Push changes to the `main` branch to trigger Terraform CI and Terraform Pipeline.
2. Check the Terraform Pipeline result and review the generated plan artifact.
3. Open GitHub Actions and manually run the **Terraform Apply** workflow when deployment is required.

Local Terraform commands can also be used from the `infrastructure` directory after running `terraform init`.