# <h1 align="center">Automated environment provisioning</a>

### The main concept

The `main concept` of this approach is writing a `dynamic infrastructure pipeline` to achive fully automated environment creation. 

This particular task creates subtasks for implementing:

- `refactoring Terraform files` for being able to:
  - use the same resource definitions for several environments
  - use different remote backend configuration of S3 for each environment  
- `refactoring CircleCI pipeline configuration` for being able to work on each created branch (I decided to create a pipeline wich can handle `any desired number of envs`)
- `changes to the tests` for being able to work as expected in different envs

### Changes to the Terraform configuration files

To achive `multienv support`, `env` environment variable can be implemented to define the environment is currently using:

```
variable "env" {
  type = string
}
```
It's a common practice to name the git branch using the environment name. So in such case `env` var can be passed via CircleCI [predefined variable](https://circleci.com/docs/variables/#built-in-environment-variables) during the pipeline execution as the `TF_VAR_env`.

For creating unique names for the Terraform resources, `env` is used as a prefix, for example http instance Terraform resource:

```
resource "aws_instance" "http" {
  for_each      = var.http_instance_names
  ...
  tags = {
    Name = "${var.env}-${each.key}"
  }
}
```

### Dynamic backend for each env

For being able to `store the Terraform statefile` by `different keys of path`, I decided to implement the pre-templating using `envsubst`. For this part, I added the `envsubst` to my job image of the pipeline. The following `backend.tf.tpl` template is used:

```
terraform {
  required_version = ">= 0.12"
  backend "s3" {
    bucket = "syndicate-tfstate"
    key    = "terraform/$TF_VAR_env-state"
    region = "eu-central-1"
  }
}
```

`TF_VAR_env` variable is configured using `TF_VAR_env=$CIRCLE_BRANCH` command before the substitution. 

The result of running  `envsubst < backend.tf.tpl > backend-$TF_VAR_env.tf` will be used as the backend configuration Terraform executions.

### Changes to the Terratests

For making my Terratest work fine with the dynamic resource names, i defined an additional function to set the `env` prefixes to the resource names for running checks:

```
func addEnvPrefix(env string, list []string) []string {
	var result []string
	for _, item := range list {
		result = append(result, env+"-"+item)
	}
	return result
}
```
So now writing the list of the resource names looks like:

```
...
tfvarsFilePath := "../terratest/infrastructure/terraform.tfvars"
    env := os.Getenv("CIRCLE_BRANCH")

    vars := map[string]interface{}{
    "http_instance_names": addEnvPrefix(env, terraform.GetVariableAsListFromVarFile(t, tfvarsFilePath, "http_instance_names")),
    "db_instance_names":   addEnvPrefix(env, terraform.GetVariableAsListFromVarFile(t, tfvarsFilePath, "db_instance_names")),
    ...
    }
...
```

The comparing part looks the same:

```
...
// Verify EC2 instance names
    assert.ElementsMatch(t, vars["http_instance_names"], terraform.OutputList(t, terraformOptions, "http_instance_names"), "HTTP instance names do not match")
    assert.ElementsMatch(t, vars["db_instance_names"], terraform.OutputList(t, terraformOptions, "db_instance_names"), "DB instance names do not match")
...
```
Other parts of testing don't require any changes.