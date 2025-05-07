package infra_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func GetFirstThreeOctets(fullIp string) string {
	ipAddrOctets := strings.Split(fullIp, ".")
	return ipAddrOctets[0] + "." + ipAddrOctets[1] + "." + ipAddrOctets[2]
}

func TestInfra(t *testing.T) {

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    "../terragrunt/infra-module",
		TerraformBinary: "terragrunt",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Extracting test data from terragrunt apply
	instanceData := terraform.OutputMap(t, terraformOptions, "http_ip")
	instancesRunning, err := strconv.ParseBool(terraform.Output(t, terraformOptions, "http_instances_state"))
	if err != nil {
		t.Fatalf("Failed to parse bool for http_instances_state: %v", err)
	}
	dbInstanceData := terraform.OutputMap(t, terraformOptions, "db_ip")
	dbInstancesRunning, err := strconv.ParseBool(terraform.Output(t, terraformOptions, "db_instances_state"))
	if err != nil {
		t.Fatalf("Failed to parse bool for db_instances_state: %v", err)
	}
	dbPublicAccess, err := strconv.ParseBool(terraform.Output(t, terraformOptions, "db_public_access"))
	if err != nil {
		t.Fatalf("Failed to parse bool for db_public_access: %v", err)
	}

	// Testing:
	// 1) EC2 instances are correctly created
	assert.True(t, instancesRunning)
	assert.True(t, dbInstancesRunning)
	// 2) CIDR addresses for VPC and subnets
	for _, instanceIp := range instanceData {
		assert.Equal(t, "200.100.1", GetFirstThreeOctets(instanceIp))
	}
	for _, dbInstanceIp := range dbInstanceData {
		assert.Equal(t, "200.100.2", GetFirstThreeOctets(dbInstanceIp))
	}
	// 3) Testing if DB instance is accessible over internet
	assert.False(t, dbPublicAccess)
}
