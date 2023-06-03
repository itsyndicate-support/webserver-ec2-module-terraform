package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTests(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "/home/circleci/project/infrastructure",
	}

	testInstanceExistence(t, terraformOptions)
	testCidrBlocks(t, terraformOptions)
	testDbInstances(t, terraformOptions)
}

func testInstanceExistence(t *testing.T, terraformOptions *terraform.Options) {
	instanceIDs := []string{"http1_id", "http2_id"}

	for _, instanceID := range instanceIDs {
		_, ok := terraform.OutputMap(t, terraformOptions, instanceID)
		assert.True(t, ok, "Instance does not exist: %s", instanceID)
	}

	dbInstances := []string{"db1_id", "db2_id", "db3_id"}

	for _, dbInstance := range dbInstances {
		_, ok := terraform.OutputMap(t, terraformOptions, dbInstance)
		assert.True(t, ok, "Database instance does not exist: %s", dbInstance)
	}
}

func testCidrBlocks(t *testing.T, terraformOptions *terraform.Options) {
	cidrBlocks := []struct {
		output   string
		expected string
	}{
		{"vpc_cidr", "192.168.0.0/16"},
		{"http_subnet_cidr", "192.168.1.0/24"},
		{"db_subnet_cidr", "192.168.2.0/24"},
	}

	for _, tt := range cidrBlocks {
		actualCidr := terraform.Output(t, terraformOptions, tt.output)
		assert.Equal(t, tt.expected, actualCidr, "Cidr block does not match.")
	}
}

func testDbInstances(t *testing.T, terraformOptions *terraform.Options) {
	dbInstances := []string{"db1_id", "db2_id", "db3_id"}

	for _, dbInstance := range dbInstances {
		dbID := terraform.Output(t, terraformOptions, dbInstance)
		dbPublicIP := aws.GetPublicIpOfEc2Instance(t, dbID, "eu-central-1")
		assert.Equal(t, "", dbPublicIP)
	}
}
