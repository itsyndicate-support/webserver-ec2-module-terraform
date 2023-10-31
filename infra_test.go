package test

import (
	"testing"
	"net/http"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMatchTfvarsAndOutputs(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../infrastructure",
		VarsFiles:    []string{"../infrastructure/terraform.tfvars"},
	}

	// Read variables from the terraform.tfvars file
	terraformVars := terraform.ReadVarFiles(t, terraformOptions.VarsFiles)
	vpcCidr := terraformVars["vpc_cidr"].(string)
	httpSubnetCidr := terraformVars["network_http"].(map[string]interface{})["cidr"].(string)
	dbSubnetCidr := terraformVars["network_db"].(map[string]interface{})["cidr"].(string)


	// Fetch the VPC CIDR block from Terraform output
	actualVpcCidr := terraform.Output(t, terraformOptions, "vpc_cidr")
	actualHttpSubnetCidr := terraform.Output(t, terraformOptions, "http_subnet_cidr")
	actualDbSubnetCidr := terraform.Output(t, terraformOptions, "db_subnet_cidr")

	// Check if the values match
	assert.Equal(t, vpcCidr, actualVpcCidr)
	assert.Equal(t, httpSubnetCidr, actualHttpSubnetCidr)
	assert.Equal(t, dbSubnetCidr, actualDbSubnetCidr)



	httpInstanceNames := terraformVars["http_instance_names"].([]interface{})
	dbInstanceNames := terraformVars["db_instance_names"].([]interface{})

	httpIPs := terraform.OutputMap(t, terraformOptions, "http_ip")
	dbIPs := terraform.OutputMap(t, terraformOptions, "db_ip")

	// Compare the counts
	assert.Equal(t, len(httpInstanceNames), len(httpIPs))
	assert.Equal(t, len(dbInstanceNames), len(dbIPs))


	httpPublicIPs := terraform.OutputList(t, terraformOptions, "http_public_ips")

	// Test HTTP connection for each instance
	for _, ip := range httpPublicIPs {
		resp, err := http.Get("http://" + ip + "/")
		defer resp.Body.Close()

		// Check if the HTTP request was successful
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}


	dbPublicIPs := terraform.OutputList(t, terraformOptions, "db_public_ips")

	// Check if the list of DB instance public IP addresses is empty
	assert.Empty(t, dbPublicIPs)
}