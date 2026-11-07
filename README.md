# How to

![infra instance with multiple network](./img/03-multiple-network.png "infra instance with multiple network")

### Create stack

```
terraform apply
```

This script will create:
-   1 vpc
-   2 networks
-   2 instances http
-   3 instances db

### Delete stack

```
terraform destroy
```
### GitHub Actions Pipeline

The repository uses GitHub Actions to automate Terraform validation and deployment.

**Workflows:**

- **Terraform CI** – runs `terraform fmt` and `terraform validate` on pushes and pull requests.
- **Terraform Plan** – initializes Terraform, generates an execution plan and uploads it as an artifact.
- **Terraform Apply** – manually triggered workflow that creates or updates the infrastructure.

**AWS Authentication**

GitHub Actions authenticates to AWS using OpenID Connect (OIDC) and an IAM role. No long-lived AWS access keys are stored in GitHub.

**Terraform Remote State**

Terraform state is stored in an Amazon S3 bucket in the `eu-north-1` region. State locking is enabled using S3 lockfiles.

**How to use the pipeline:**

1. Push changes to the `main` branch to trigger Terraform CI and Plan.
2. Check the Terraform Plan workflow result and review the generated plan artifact.
3. Open GitHub Actions and manually run the **Terraform Apply** workflow when deployment is required.

Local Terraform commands can also be used from the `infrastructure` directory after running `terraform init`.
