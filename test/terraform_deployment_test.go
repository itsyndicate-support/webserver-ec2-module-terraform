package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformDeployment(t *testing.T) {
	// Define the Terraform options.
	terraformOptions := &terraform.Options{
		// Set the path to your Terraform code.
		TerraformDir: "../infrastructure",
	}

	// Check that the EC2 instances are created.
	httpInstances := terraform.OutputList(t, terraformOptions, "http_ip")
	dbInstances := terraform.OutputList(t, terraformOptions, "db_ip")

	assert.Equal(t, 2, len(httpInstances), "Expected 2 HTTP instances to be created")
	assert.Equal(t, 2, len(dbInstances), "Expected 2 DB instances to be created")

	// Check that DB instances don't have public IP
	assert.ElementsMatch(t, dbInstances, []string{"", ""}, "Expected the DB instance to not have a public IP")

	// Check that two subnets are created.
	http_subnet := terraform.Output(t, terraformOptions, "http_subnet")
	db_subnet := terraform.Output(t, terraformOptions, "db_subnet")
	vpc := terraform.Output(t, terraformOptions, "vpc")

	assert.NotEmpty(t, http_subnet, "Expected the HTTP subnet to be created")
	assert.NotEmpty(t, db_subnet, "Expected the DB subnet to be created")
	assert.NotEmpty(t, vpc, "Expected the VPC to be created")
}
