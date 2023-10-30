package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
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
	httpInstances := terraform.OutputList(t, terraformOptions, "http_id")
	dbInstances := terraform.OutputList(t, terraformOptions, "db_id")

	assert.Equal(t, 2, len(httpInstances), "Expected 2 HTTP instances to be created")
	assert.NotEmpty(t, dbInstances, "Expected a database instance to be created")

	// Get the public IP address of the EC2 instance.
	publicIP := aws.GetPublicIpsOfEc2Instances(t, dbInstances, "us-east-1")

	// Ensure that the public IP is empty (i.e., the instance doesn't have a public IP).
	assert.Empty(t, publicIP, "Expected the EC2 instance to not have a public IP")

	// Check that two subnets are created.
	http_subnet := terraform.Output(t, terraformOptions, "http_subnet")
	db_subnet := terraform.Output(t, terraformOptions, "db_subnet")
	vpc := terraform.Output(t, terraformOptions, "vpc")

	assert.NotEmpty(t, http_subnet, "Expected the HTTP subnet to be created")
	assert.NotEmpty(t, db_subnet, "Expected the DB subnet to be created")
	assert.NotEmpty(t, vpc, "Expected the VPC to be created")
}
