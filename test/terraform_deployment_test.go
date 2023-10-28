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
		TerraformDir: "webserver-ec2-module-terraform",

		// Variables to pass to your Terraform code.
		Vars: map[string]interface{}{
			"vpc_cidr": "192.168.0.0/16",
			"network_http": map[string]interface{}{
				"subnet_name": "subnet_http",
				"cidr":        "192.168.1.0/24",
			},
			"network_db": map[string]interface{}{
				"subnet_name": "subnet_db",
				"cidr":        "192.168.2.0/24",
			},
		},

		// Variables to pass when running 'terraform init'.
		VarFiles: []string{"webserver-ec2-module-terraform/infrastructure"},
	}

	// Check that the EC2 instances are created.
	httpInstances := terraform.OutputList(t, terraformOptions, "http_id")
	dbInstance := terraform.Output(t, terraformOptions, "db_id")

	assert.Equal(t, 2, len(httpInstances), "Expected 2 HTTP instances to be created")
	assert.NotEmpty(t, dbInstance, "Expected a database instance to be created")

	// Get the public IP address of the EC2 instance.
	publicIP := aws.GetPublicIpOfEc2Instance(t, dbInstance, "us-east-1")

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
