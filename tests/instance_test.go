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
		TerraformDir: "../infrastructure",
	}

	// Deploy the infrastructure and retrieve outputs
	terraform.InitAndApply(t, terraformOptions)
	httpIPs := terraform.OutputMap(t, terraformOptions, "http_ip")
	dbIPs := terraform.OutputMap(t, terraformOptions, "db_ip")

	// Validate the DNS information
	assert.NotEmpty(t, httpIPs)
	assert.NotEmpty(t, dbIPs)

	// Validate individual instance IP addresses
	for instanceID, privateIP := range httpIPs {
		assert.True(t, aws.InstancePrivateIPExists(t, instanceID, privateIP, "eu-central-1"))
	}

	for instanceID, privateIP := range dbIPs {
		assert.True(t, aws.InstancePrivateIPExists(t, instanceID, privateIP, "eu-central-1"))
	}
}
