package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestInfrastructure(t *testing.T) {
	t.Parallel()

	tf := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
	})

	outputs := terraform.OutputAll(t, tf)

	validateInstancesCreated(t, outputs)
	validateVPCAndSubnets(t, outputs)
	validateDBIsPrivate(t, outputs)
}

func validateInstancesCreated(t *testing.T, outputs map[string]interface{}) {
	httpIPs := outputs["http_private_ip"].(map[string]interface{})
	dbIPs := outputs["db_private_ip"].(map[string]interface{})

	assert.NotEmpty(t, httpIPs, "Expected at least one HTTP instance")
	assert.NotEmpty(t, dbIPs, "Expected at least one DB instance")
}

func validateVPCAndSubnets(t *testing.T, outputs map[string]interface{}) {
	vpcCIDR := outputs["vpc_cidr"].(string)
	httpSubnet := outputs["http_subnet_cidrs"].(string)
	dbSubnet := outputs["db_subnet_cidrs"].(string)

	assert.Equal(t, "192.168.0.0/16", vpcCIDR, "Unexpected VPC CIDR block")
	assert.Equal(t, "192.168.1.0/24", httpSubnet, "Unexpected HTTP subnet CIDR")
	assert.Equal(t, "192.168.2.0/24", dbSubnet, "Unexpected DB subnet CIDR")
}

func validateDBIsPrivate(t *testing.T, outputs map[string]interface{}) {
	dbPublicIPs := outputs["db_public_ip"].(map[string]interface{})

	for instanceID, val := range dbPublicIPs {
		ip := val.(string)
		assert.Equal(t, "", ip, "DB instance %s must nop have public ip", instanceID)
	}
}
