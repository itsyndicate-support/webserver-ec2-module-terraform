# Multi-environment Terraform deployment

I updated the existing GitHub Actions workflow to support `test` and `production`. Both environments use the same Terraform configuration, with separate variable files, resources and S3 state files.

Each environment creates one VPC, two subnets, two Apache EC2 instances, three MariaDB EC2 instances and two Elastic IP addresses in `eu-north-1`. I kept the original instance counts and connected the database bootstrap script. Resource tags identify the environment, and SSH key-pair names include the environment name.

![Infrastructure diagram](./img/03-multiple-network.png)

## Workflows

| Workflow | Trigger | Purpose |
|---|---|---|
| `.github/workflows/terraform-checks.yml` | Push and pull request | Run formatting, validation, TFLint and Trivy checks |
| `.github/workflows/terraform.yml` | Manual | Run plan, apply or destroy for the selected environment |

Terraform state is stored in S3 with native state locking. Each environment has its own concurrency group, so runs targeting the same environment cannot overlap.

GitHub Actions authenticates through OIDC and receives temporary AWS credentials. AWS operations run from `main` in the repository specified by `DEPLOY_REPOSITORY`.

## Setup

Create a private S3 bucket with versioning and encryption enabled. Create a GitHub OIDC provider and an IAM deployment role.

The role uses `AmazonEC2FullAccess` and an inline S3 policy. The S3 policy grants `s3:ListBucket` on the bucket, `s3:GetObject` and `s3:PutObject` on both state files, and `s3:GetObject`, `s3:PutObject` and `s3:DeleteObject` on their `.tflock` files.

Configure these repository variables under **Settings → Secrets and variables → Actions → Variables**:

| Variable | Value |
|---|---|
| `AWS_REGION` | `eu-north-1` |
| `AWS_ROLE_ARN` | Deployment role ARN |
| `TF_STATE_BUCKET` | `wai-ai-terraform-state-week11` |
| `DEPLOY_REPOSITORY` | `wai-AI/webserver-ec2-module-terraform` |
| `SSH_PUBLIC_KEY` | Full SSH public key |
| `ADMIN_CIDR` | Administrator's public IPv4 address with `/32` |

Create `test` and `production` under **Settings → Environments**. Add `TF_STATE_KEY` as an environment variable in each:

| Environment | Variable file | `TF_STATE_KEY` |
|---|---|---|
| `test` | `infrastructure/environments/test.tfvars` | `webserver-ec2/test/terraform.tfstate` |
| `production` | `infrastructure/environments/production.tfvars` | `webserver-ec2/production/terraform.tfstate` |

For both GitHub Environments, select **Selected branches and tags** and allow the `main` branch.

The IAM trust policy must allow these OIDC subjects for this repository, with audience `sts.amazonaws.com`:

```text
repo:wai-AI@122587908/webserver-ec2-module-terraform@1449424385:environment:test
repo:wai-AI@122587908/webserver-ec2-module-terraform@1449424385:environment:production
```

For a fork, use its repository name, owner and repository IDs, deployment role and state bucket.

The workflows use Terraform `1.16.0`. The provider lock file is committed to the repository. Environment variable files contain configuration values only; credentials are not committed.

## Usage

Pushes and pull requests run the configuration checks automatically.

For AWS operations, open **Actions → Terraform AWS → Run workflow**. Select `main`, choose the target environment and select an operation:

| Operation | Confirmation | Result |
|---|---|---|
| `plan` | Leave empty | Show the proposed changes |
| `apply` | Leave empty | Generate and apply a fresh plan, then run HTTP checks and Terratest |
| `destroy` | Enter `destroy` | Generate and apply a destroy plan for the selected environment |

The selected environment determines the GitHub Environment, variable file and S3 state key. Apply and destroy use saved plans generated in the same workflow run.

The `http_urls` output contains the public Apache URLs. The HTTP checks retry while the servers finish bootstrapping. Terratest checks instance counts, running state, VPC and subnet configuration, and database security-group rules.

## Add a new environment

To add another environment, for example `staging`:

1. Create `infrastructure/environments/staging.tfvars` with:

   ```hcl
   environment = "staging"
   ```

   Add any environment-specific overrides to the same file.

2. Add `staging` to the allowed values in the `environment` validation in `infrastructure/00-params.tf`.

3. Add `staging` to the environment input options in `.github/workflows/terraform.yml`.

4. Create the GitHub Environment `staging`. Set its `TF_STATE_KEY` to `webserver-ec2/staging/terraform.tfstate` and allow deployments from `main`.

5. Add its environment subject to the IAM trust policy. Grant access to the new state object and its `.tflock` object in the S3 policy.

6. Update Terratest expectations if instance counts or network CIDRs differ from the current configuration.

7. Merge the changes into `main`. Run `plan` for `staging`, review the changes, then run `apply`.

## Verification

Run `plan` and `apply` for both `test` and `production`. Check that the HTTP checks and Terratest pass for each environment.

Confirm that the environments have different VPC IDs and separate S3 state files. Run a follow-up plan for each environment; with unchanged configuration and AMI selection, it should report `No changes`.

To check cleanup isolation, destroy `test`, then run a plan for `production` and check that its resources remain unchanged. Keep the workflow run links for the implementation report.

## Cleanup

Select the environment, choose `destroy` and enter `destroy` in the confirmation field. Terraform removes the resources managed by that environment's state.

The state bucket, IAM role and OIDC provider were created separately and remain after destroy.

Automatic cleanup could be added with a scheduled workflow for the temporary `test` environment. It should use the test state key and the same `terraform-test` concurrency group. Production cleanup should remain a manual operation.

## Self-hosted runner

The deployment workflow supports `ubuntu-latest` and `terraform-runner`.

See [Runner evaluation](docs/runner-report.md) for the performance measurements, cost comparison and findings.

### Installation

1. Create a separate EC2 instance with Ubuntu 24.04 x86_64. The tested configuration used `c7i-flex.large`, 4 GiB RAM and a 20 GB gp3 disk.

   Allow outbound access to GitHub, AWS APIs and package repositories. Restrict inbound SSH to the administrator's IP address.

2. Connect over SSH and install the required packages:

```bash
sudo apt-get update
sudo apt-get install -y git curl jq unzip ca-certificates
```

3. Install AWS CLI:

```bash
curl -fsSL https://awscli.amazonaws.com/v2/install.sh \
  -o /tmp/awscli-install-runner.sh
sudo bash /tmp/awscli-install-runner.sh --system
aws --version
```

4. In the repository, open **Settings → Actions → Runners → New self-hosted runner**. Select **Linux x64** and follow GitHub's download, checksum verification and extraction commands.

   From the extracted runner directory, install its dependencies:

```bash
sudo ./bin/installdependencies.sh
```

5. Run the configuration command generated by GitHub as the `ubuntu` user. Use the default runner group, add the custom label `terraform-runner` and keep the default work directory.

   Use a fresh registration token from GitHub. Do not commit the token or runner credentials.

6. Install and start the service from the runner directory:

```bash
sudo ./svc.sh install ubuntu
sudo ./svc.sh start
sudo ./svc.sh status
```

7. Check that GitHub shows the runner as **Idle**.

The workflow installs Terraform, TFLint, Trivy and Go as required. AWS access uses the existing OIDC configuration.

The lab also used 4 GiB swap. To reproduce it, follow the [AWS swap-file instructions](https://repost.aws/knowledge-center/ec2-memory-swap-file). Swap helps with memory pressure but does not replace sufficient RAM.

### Running the pipeline

Open **Actions → Terraform AWS → Run workflow**:

- Branch: `main`
- Runner: `terraform-runner`
- Environment: `test` or `production`
- Operation: `plan`, `apply` or `destroy`
- Confirmation: enter `destroy` only for deletion

Keep self-hosted deployment runs limited to trusted workflows. Pull-request checks use GitHub-hosted runners.
