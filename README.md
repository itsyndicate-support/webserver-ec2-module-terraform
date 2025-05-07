# Project Setup Instructions
This document provides instructions to bootstrap Flux in the "dev" branch, and create new environment with "test" branch.

## Prerequisites
### Ensure you have the following tools installed on your local machine:

- Git
- Flux CLI
- GitHub token. Create [here](https://github.com/settings/tokens)
## Steps to Create "dev" Branch
Clone the Repository:

```bash
git clone https://github.com/phosphor-v/webserver-ec2-module-terraform
cd webserver-ec2-module-terraform/
```

**Create the `dev` Branch:**

```bash
git checkout -b dev
```


**Push the `dev` Branch to Remote:**

```bash
git push origin dev
```
## Bootstrap Flux in the `dev` Branch:

**Create github token env:**
```bash
export GITHUB_TOKEN=<gh-token>
```
Replace `<gh-token>` with your [GitHub token](https://github.com/settings/tokens).

**Ensure you are in the root directory of your repository and execute the following command:**

```bash
flux bootstrap github \
  --token-auth \
  --owner=<github-username> \
  --repository=webserver-ec2-module-terraform \
  --branch=dev \
  --personal \
  --private=false \
  --path=kubernetes/fluxcd/repositories/infra-repo/clusters/dev-cluster
```
Replace `<github-username>` with your actual GitHub username.

**Install tofu-controller:**
```bash
kubectl apply -f https://raw.githubusercontent.com/flux-iac/tofu-controller/main/docs/release.yaml
```

**Set-up AWS credentials**
```bash
kubectl create secret generic aws-credentials \
  --from-literal=AWS_ACCESS_KEY_ID=<acces_key> \
  --from-literal=<secret_key> -n flux-system

kubectl create secret generic terraform-s3-backend \
  --from-literal=access_key=<acces_key> \
  --from-literal=secret_key=<secret_key> -n flux-system
```
Replace `<acces_key>` and `<secret_key>`  with your own

## Steps to Create "test" Branch with Empty Commit
**(Optional) Checkout branch with required source. Typically `main`**
```bash
git checkout main
```

**Create the "test" Branch:**
```bash
git checkout -b test
```

**Create an Empty Commit:**

```bash
git commit --allow-empty -m "trigger"
```

**Push the "test" Branch to Remote:**

```bash
git push origin test
```

## Results
Terratest configured to run tests with branch name as variable.


After running terratest, pipeline will create new flux environment configuration for this branch and push it to `dev` branch. Check for existing environment is also included to prevent errors.


Fluxcd will automatically add it to cluster and create tf resources.

### Some pictures of deployment:
![image](https://github.com/phosphor-v/webserver-ec2-module-terraform/assets/113942255/bcc54d91-0815-4cc9-bee9-634195ec6b2d)
![image](https://github.com/phosphor-v/webserver-ec2-module-terraform/assets/113942255/e92564d3-7b19-4143-8efe-d684d9414e3b)
