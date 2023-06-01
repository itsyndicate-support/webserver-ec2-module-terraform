package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestInfrastructure(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "/home/circleci/project/infrastructure",
	}

	// Check the existence of EC2 instances
	instanceIDs := []string{
		terraform.Output(t, terraformOptions, "http_instance_id1"),
		terraform.Output(t, terraformOptions, "http_instance_id2"),
		terraform.Output(t, terraformOptions, "db_instance_id1"),
		terraform.Output(t, terraformOptions, "db_instance_id2"),
		terraform.Output(t, terraformOptions, "db_instance_id3"),
	}

	for _, instanceID := range instanceIDs {
		assert.True(t, aws.InstanceExists(t, instanceID, "eu-central-1"))
	}

	// Check the existence of VPC and subnets with the correct CIDR blocks
	vpcCidr := terraform.Output(t, terraformOptions, "vpc_cidr")
	assert.Equal(t, "192.168.0.0/16", vpcCidr)

	httpSubnetCidr := terraform.Output(t, terraformOptions, "http_subnet_cidr")
	assert.Equal(t, "192.168.1.0/24", httpSubnetCidr)

	dbSubnetCidr := terraform.Output(t, terraformOptions, "db_subnet_cidr")
	assert.Equal(t, "192.168.2.0/24", dbSubnetCidr)

	// Check the absence of access to the database from the internet
	dbInstanceIDs := []string{
		terraform.Output(t, terraformOptions, "db_instance_id1"),
		terraform.Output(t, terraformOptions, "db_instance_id2"),
		terraform.Output(t, terraformOptions, "db_instance_id3"),
	}

	for _, dbInstanceID := range dbInstanceIDs {
		publicIP := aws.GetPublicIpOfEc2Instance(t, dbInstanceID, "eu-central-1")
		assert.Empty(t, publicIP)
	}
}
