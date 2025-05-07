# Flux part:
### Create cluster
```bash
kind create cluster --name production-cluster
```
### Add github token
***Create token [here](https://github.com/settings/tokens)***
```bash
export GITHUB_TOKEN=<gh-token>
```
### bootstrap flux
```bash
flux bootstrap github \
  --owner=YOUR_GITHUB_USERNAME \
  --repository=YOUR_FORKED_REPO \
  --branch=main \
  --path=kubernetes/fluxcd/clusters/production-cluster
```
### Add AWS secrets
```bash
kubectl create secret generic aws-credentials \
  --from-literal=aws_access_key_id=YOUR_ACCESS_KEY \
  --from-literal=aws_secret_access_key=YOUR_SECRET_KEY -n flux-system
```
### install tofu controller
```bash
kubectl apply -f https://raw.githubusercontent.com/flux-iac/tofu-controller/main/docs/release.yaml
```
### Configure flux 
***MODIFY*** `kubernetes/fluxcd/repositories/infra-repo/apps/tf-app/gitrepository.yaml` before apply
```bash
kubectl -n flux-system apply -f kubernetes/fluxcd/repositories/infra-repo/apps/tf-app/gitrepository.yaml
kubectl -n flux-system apply -f kubernetes/fluxcd/repositories/infra-repo/apps/tf-app/terraform.yaml
```
# Circle part:
## Create token
```bash
circleci namespace create <name> --org-id <your-organization-id>
```
```bash
circleci runner resource-class create <namespace>/<resource-class> <description> --generate-token
```
## Install runner to k8s

### Add helm image
```bash
helm repo add container-agent https://packagecloud.io/circleci/container-agent/helm
helm repo update
```

### Create namespace
```bash
kubectl create namespace circleci
```

### Update values.yaml:
```yaml
agent:
  resourceClasses:
    namespace/my-rc:
      token: <resource_class_token>
```
### Install container agent
```bash
helm install container-agent container-agent/container-agent -n circleci -f values.yaml
```