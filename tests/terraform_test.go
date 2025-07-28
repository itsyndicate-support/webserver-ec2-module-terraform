package test

import (
	"fmt"
	"net"
	"testing"
	"time"

	terratestAws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var awsRegion string = "us-east-1"

func TestTerraform(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../infrastructure",
	})

	CheckEc2Instances(t, terraformOptions)
	CheckCidrs(t, terraformOptions)
	CheckDbNotAccessable(t, terraformOptions)
}

func CheckEc2Instances(t *testing.T, terraformOptions *terraform.Options){
	instanceHttpIDs := terraform.OutputList(t, terraformOptions, "instance_http_ids")
	instanceDbIDs := terraform.OutputList(t, terraformOptions, "instance_db_ids")
	vpcId := terraform.Output(t, terraformOptions, "vpc_id")
	totalCount := len(instanceHttpIDs) + len(instanceDbIDs)

	filters := map[string][]string{
		"instance-state-name": {"running"},
		"vpc-id": {vpcId},
	}

	ids := terratestAws.GetEc2InstanceIdsByFilters(t, awsRegion, filters)

	assert.Equal(t, totalCount, len(ids))
}

func CheckCidrs(t *testing.T, terraformOptions *terraform.Options){
	vpcCidr := terraform.Output(t, terraformOptions, "vpc_cidr")
	subnetHttpCird := terraform.Output(t, terraformOptions, "subnet_http_cidr")
	subnetDbCird := terraform.Output(t, terraformOptions, "subnet_db_cidr")

	assert.Equal(t, "192.168.0.0/16", vpcCidr)
	assert.Equal(t, "192.168.1.0/24", subnetHttpCird)
	assert.Equal(t, "192.168.2.0/24", subnetDbCird)
}

func CheckDbNotAccessable(t *testing.T, terraformOptions *terraform.Options){
	instanceIDs := terraform.OutputList(t, terraformOptions, "instance_db_ids")
	publicIps := terratestAws.GetPublicIpsOfEc2Instances(t, instanceIDs, awsRegion)

	ipsCount := 0
	for key, value := range publicIps {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		if value != "" {
			ipsCount++
		}
	}

	noPublicIps := 0 == ipsCount
	assert.Equal(t, true, noPublicIps)
	
	if noPublicIps{
		return
	}

	// Check connection to MySQL
	timeout := 5 * time.Second
	for name, ip := range publicIps{
		address := fmt.Sprintf("%s:3306", ip)
		
		conn, err := net.DialTimeout("tcp", address, timeout)
		if conn != nil {
			conn.Close()
		}

		assert.Error(t, err, "Instance %s is accessible at %s!", name, address)
	}
}