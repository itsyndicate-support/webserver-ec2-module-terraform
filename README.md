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
# Report

In this task I created CircleCI CI pipeline and FluxCD CD pipeline. My implementation consists of:
1) CircleCI CI pipeline with 2 workflows: just-test and test-and-pr.
2) FluxCD CD pipeline that uses [terraform controller](https://github.com/flux-iac/tofu-controller) to automatically deploy infrastructure in AWS

## CI
Two implemented workflows `just-test` and `test-and-pr` do next things:
1) `just-test` workflow tests all commits for secure terraform code with [tfsec](https://github.com/aquasecurity/tfsec) and for it's ability to correctly deploy infrastructure using terragrunt and [terratest](https://pkg.go.dev/github.com/gruntwork-io/terratest@v0.46.15) golang package.
2) `test-and-pr` does the same thing as `just-test` but it triggers when commit message has `create pr` in it and automatically generates PR into main branch.

Example of successfull workflows:
![workflow triggering examples](/img/workflows.png)

![created PR](/img/auto_pr.png)

## CD
CD pipeline consists of default FluxCD config files in default [path](/kubernetes/fluxcd/repositories/infra-repo/clusters/dev-cluster/flux-system/) and separate FluxCD [specific k8s manifests](/kubernetes/fluxcd/repositories/infra-repo/apps/tf-app/) to monitor state of `/infrastructure` directory in this repository.
1) `gitrepository.yaml` configures which repository and on which branch to monitor
2) `terraform.yaml` uses `gitrepository.yaml` to monitor `/infrastructure` folder for changes and on initial deployment to k8s cluster perform `terraform apply` with terraform files in specified directory. Also this manifest will automatically destroy all terraform infrastructre on deletion from cluster using `destroyResourcesOnDeletion` spec.

Bootstraped FluxCD in k8s cluster:
![list of pods in k8s cluster](/img/FluxCD_bootstraped.png)

Result of terraform controller deploying infrastructure:
![k8s resources](/img/k8s_resources.png)

![aws infrastructure](/img/deployed_infrastructure.png)