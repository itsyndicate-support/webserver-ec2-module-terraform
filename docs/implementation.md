I added two GitHub Actions workflows to manage the infrastructure in this repository.

The deployment creates two Apache servers and three MariaDB servers in `eu-north-1`, using one VPC and two subnets. I kept the original instance counts and bootstrap scripts. I added public IPv4 addresses for package installation, restricted SSH access to the configured administrator IP, and limited port 3306 to instances using the web security group.

## Workflows

| Workflow | Trigger | Purpose |
|---|---|---|
| `terraform-checks.yml` | Push and pull request | Check formatting and validate the Terraform configuration |
| `terraform.yml` | Manual | Run plan, apply or destroy |

Terraform state is stored in S3 with native state locking. A shared concurrency group prevents overlapping runs of the AWS workflow.

GitHub Actions authenticates through OIDC and receives temporary AWS credentials.

## Setup

Create a private S3 bucket with versioning and encryption enabled. Then create a GitHub OIDC provider and an IAM role that trusts the fork’s `main` branch.

The deployment role uses `AmazonEC2FullAccess` and an inline policy granting access to the S3 state object and its lock file.

Configure these repository variables under **Settings → Secrets and variables → Actions → Variables**:

| Variable | Value |
|---|---|
| `AWS_REGION` | `eu-north-1` |
| `AWS_ROLE_ARN` | Deployment role ARN |
| `TF_STATE_BUCKET` | State bucket name |
| `TF_STATE_KEY` | `webserver-ec2/lab/terraform.tfstate` |
| `DEPLOY_REPOSITORY` | Fork owner and repository name |
| `SSH_PUBLIC_KEY` | Full SSH public key |
| `ADMIN_CIDR` | Administrator’s public IPv4 address with `/32` |

The workflows use Terraform `1.16.0`. The provider lock file is committed to the repository.

## Usage

Pushes and pull requests run the checks automatically.

For AWS operations, open **Actions → Terraform AWS → Run workflow**, select `main`, and choose an operation:

| Operation | Confirmation |
|---|---|
| `plan` | Leave empty |
| `apply` | Leave empty |
| `destroy` | Enter `destroy` |

Apply generates a fresh plan and applies it in the same run. Destroy generates and applies a destroy plan.

The state bucket, IAM role and OIDC provider were created separately and remain after destroy.

## Test results

- Formatting and validation checks passed in GitHub Actions.
- Both Apache servers returned `HTTP 200 OK`.
- All three DB servers completed cloud-init, reported MariaDB as active, and successfully executed `SELECT VERSION()`.
- A follow-up Terraform plan completed successfully.
- The manual destroy workflow completed successfully.

## Security checks and post-apply tests

Both workflows run Terraform formatting and validation checks, TFLint, and Trivy. In the deployment workflow, these checks must pass before Terraform can generate a plan. If a check fails, subsequent deployment steps are skipped.

The security changes enable EC2 root volume encryption, require IMDSv2, and limit outbound internet traffic to HTTP and HTTPS.

Terratest runs automatically after a successful apply. It reads the existing Terraform outputs and queries AWS to check that:

- Two HTTP instances and three database instances are running.
- The VPC CIDR is `192.168.0.0/16`, with HTTP instances in `192.168.1.0/24` and database instances in `192.168.2.0/24`.
- Every database instance allows inbound TCP/3306 only from the web security group. All attached security groups are checked, including rules covering multiple ports or all protocols.

The tests use the same temporary AWS credentials as Terraform. To run them, select `apply` in the manual Terraform AWS workflow.

