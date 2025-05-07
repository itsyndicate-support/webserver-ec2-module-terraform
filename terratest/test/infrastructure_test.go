package test

import (
	"os"
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestTerraformModule(t *testing.T) {
    t.Parallel()
    env := os.Getenv("CIRCLE_BRANCH")
    opts := &terraform.Options{
        TerraformDir: "../terragrunt/infra-module",
        TerraformBinary : "terragrunt",
        Vars: map[string]interface{}{
			"env": env,
			"public_key_name": env,
		},
    }

    defer terraform.TgDestroyAll(t, opts)
    terraform.TgApplyAll(t, opts)
    httpInstances := terraform.OutputList(t, opts, "http_ip")
	dbInstances := terraform.OutputList(t, opts, "db_ip")

	assert.Equal(t, 2, len(httpInstances), "Expected 2 HTTP instances to be created")
	assert.Equal(t, 3, len(dbInstances), "Expected 3 DB instances to be created")

    vpcCidr := terraform.Output(t, opts, "vpc_cidr")
    assert.Equal(t, vpcCidr, "192.168.0.0/16", "Expected VPC CIDR block to be 192.168.0.0/16")
    HttpSubnetCidr := terraform.Output(t, opts, "http_subnet_cidr")
    assert.Equal(t, "192.168.1.0/24", HttpSubnetCidr, "Expected Http subnet cidr block to be 192.168.1.0/24")

    DbSubnetCidr := terraform.Output(t, opts, "db_subnet_cidr")
    assert.Equal(t, "192.168.2.0/24", DbSubnetCidr, "Expected Http subnet cidr block to be 192.168.2.0/24")
    dbPublicIPs := terraform.OutputList(t, opts, "db_instance_public_ips")
    assert.Empty(t, dbPublicIPs, "Database should not be accessible from the internet")
}
