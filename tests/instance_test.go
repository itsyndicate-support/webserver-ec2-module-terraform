package main

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestHTTPInstances(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../infrastructure",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	instanceNames := []string{"instance-http-1", "instance-http-2"} 

	for _, instanceName := range instanceNames {
		// Перевірка створення інстансу EC2
		instance := aws.GetEc2InstanceById(t, instanceName, "us-east-1")
		assert.NotNil(t, instance)

		// Перевірка наявності присвоєного плаваючого IP
		eip := aws.GetEipForEc2InstanceId(t, instance.ID, "us-east-1")
		assert.NotNil(t, eip)
	}
}
