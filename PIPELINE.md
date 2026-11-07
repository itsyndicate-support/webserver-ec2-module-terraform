# Terraform CI/CD Pipeline

## Security and Code Quality

Before Terraform Plan, the pipeline runs:

- **Terraform fmt and validate** – checks formatting and configuration validity.
- **TFLint** – detects Terraform best-practice issues.
- **Trivy** – scans Terraform configuration for security misconfigurations.

Terraform Plan runs only if TFLint and Trivy succeed.

## Terraform Apply and Terratest

Terraform Apply is triggered manually using GitHub Actions workflow_dispatch.

After a successful apply, Terratest checks:

- Expected HTTP and DB EC2 instance counts and running state.
- VPC CIDR: 192.168.0.0/16.
- HTTP subnet CIDR: 192.168.1.0/24.
- DB subnet CIDR: 192.168.2.0/24.
- DB instances have no public IP addresses.
- DB Security Group does not allow public MySQL access.

Terratest uses AWS credentials configured through GitHub Actions OIDC.

## Failure Handling

GitHub Actions stops subsequent steps when a step fails.
Terraform Plan depends on successful TFLint and Trivy jobs.
Terratest runs only after Terraform Apply succeeds.