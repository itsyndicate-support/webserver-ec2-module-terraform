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

### Add new environment

1. In ALL workflows (cd, ci, destroy) need to add to environment new option
2. In infrastructure/backend add file with new terraform backend
3. In infrastructure/vars add file with new variables