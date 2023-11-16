package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformDeployment(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../infrastructure",
	}

	// Check vpc for creation.
	vpc := terraform.Output(t, terraformOptions, "vpc")

	assert.NotEmpty(t, vpc, "Expected the VPC to be created")

	// Check subnets for creation.
	httpSubnet := terraform.Output(t, terraformOptions, "httpSubnet")
	dbSubnet := terraform.Output(t, terraformOptions, "dbSubnet")

	assert.NotEmpty(t, httpSubnet, "Expected the HTTP subnet to be created")
	assert.NotEmpty(t, dbSubnet, "Expected the DB subnet to be created")

	// Check EC2 instances for creation.
	httpInstances := terraform.OutputList(t, terraformOptions, "http_ip")
	dbInstances := terraform.OutputList(t, terraformOptions, "db_ip")

	assert.Equal(t, 2, len(httpInstances), "Expected 2 HTTP instances to be created")
	assert.Equal(t, 2, len(dbInstances), "Expected 2 DB instances to be created")
}