package test

import (
    "os"
    "fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestInfrastructure(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../terratest/infrastructure",
    }

    // Get variables from terraform.tfvars
    tfvarsFilePath := "../terratest/infrastructure/terraform.tfvars"
    env, ok := os.LookupEnv("CIRCLE_BRANCH")

    vars := map[string]interface{}{
        "http_instance_names": addEnvPrefix(env, terraform.GetVariableAsListFromVarFile(t, tfvarsFilePath, "http_instance_names")),
		"db_instance_names":   addEnvPrefix(env, terraform.GetVariableAsListFromVarFile(t, tfvarsFilePath, "db_instance_names")),
        "vpc_cidr":            terraform.GetVariableAsStringFromVarFile(t, tfvarsFilePath, "vpc_cidr"),
        "http_subnet_cidr":    terraform.GetVariableAsMapFromVarFile(t, tfvarsFilePath, "network_http")["cidr"],
        "db_subnet_cidr":      terraform.GetVariableAsMapFromVarFile(t, tfvarsFilePath, "network_db")["cidr"],
    }

    // Verify EC2 instance names
    assert.ElementsMatch(t, vars["http_instance_names"], terraform.OutputList(t, terraformOptions, "http_instance_names"), "HTTP instance names do not match")
    assert.ElementsMatch(t, vars["db_instance_names"], terraform.OutputList(t, terraformOptions, "db_instance_names"), "DB instance names do not match")

    // Verify CIDR blocks
    // Verify VPC cidr
    assert.Equal(t, vars["vpc_cidr"], terraform.Output(t, terraformOptions, "vpc_cidr"), "VPC CIDR block does not match")

    // Verify http instances cidr
    expectedHTTPSubnetCidr := vars["http_subnet_cidr"].(string)
    actualHTTPSubnetCidr := terraform.Output(t, terraformOptions, "http_subnet_cidr")
    assert.Equal(t, expectedHTTPSubnetCidr, actualHTTPSubnetCidr, "HTTP subnet CIDR block does not match")

    // Verify db instances cidr
    expectedDBSubnetCidr := vars["db_subnet_cidr"].(string)
    actualDBSubnetCidr := terraform.Output(t, terraformOptions, "db_subnet_cidr")
    assert.Equal(t, expectedDBSubnetCidr, actualDBSubnetCidr, "DB subnet CIDR block does not match")

    // Verify if the database is accessible from the internet
    dbPublicIPs := terraform.OutputList(t, terraformOptions, "db_instance_public_ips")
    assert.Empty(t, dbPublicIPs, "Database should not be accessible from the internet")
}


func addEnvPrefix(env string, list []string) []string {
	var result []string
	for _, item := range list {
		result = append(result, env+"-"+item)
	}
	return result
}
