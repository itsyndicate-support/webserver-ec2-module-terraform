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

If you want to use this pipeline in your environment, you need the following:
```
- input your S3 Bucket config in './infrastructure/main.tf';

- if your remote .tfstate file doesn't have the name 'terraform.tfstate', you need to make appropriative changes in '.circleci/config.yml' for '-backend-config' parameter;

- (optional) if you want, you can rename jobs and 'requires' fields in the workflow according to your environment in '.circleci/config.yml';

- create 'FIRST_ENVIRONMENT' variable in your CircleCi project with the value of your environment name (keep in mind that a folder with .tfstate file will have the value of this variable as the name;

- (optional) if you want to create CircleCi project variable with another name, you also need to rename it in .circleci/config.yml;

- create context 'terraform' in CircleCI with your AWS IAM user credentials;

- enjoy using this pipeline!!
```