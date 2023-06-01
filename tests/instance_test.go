package test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/aws"
)

func TestTests(t *testing.T) {
    t.Parallel()

    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{

		TerraformDir: "../infrastructure",
	})

	HttpName1 := terraform.Output(t, terraformOptions, "http_instance_name1")
	assert.Equal(t, "instance-http-1", HttpName1, "Instance name does not match.")

	HttpName2 := terraform.Output(t, terraformOptions, "http_instance_name2")
	assert.Equal(t, "instance-http-2", HttpName2, "Instance name does not match.")

	DbName1 := terraform.Output(t, terraformOptions, "db_instance_name1")
	assert.Equal(t, "instance-db-1", DbName1, "Instance name does not match.")

	DbName2 := terraform.Output(t, terraformOptions, "db_instance_name2")
	assert.Equal(t, "instance-db-2", DbName2, "Instance name does not match.")

	DbName3 := terraform.Output(t, terraformOptions, "db_instance_name3")
	assert.Equal(t, "instance-db-3", DbName3, "Instance name does not match.")

	VpcCidr := terraform.Output(t, terraformOptions, "vpc_cidr")
    assert.Equal(t, "192.168.0.0/16", VpcCidr, "Cidr block does not match.")

    HttpSubnetCidr := terraform.Output(t, terraformOptions, "http_subnet_cidr")
    assert.Equal(t, "192.168.1.0/24", HttpSubnetCidr, "Http subnet cidr block does not match.")

    DbSubnetCidr := terraform.Output(t, terraformOptions, "db_subnet_cidr")
    assert.Equal(t, "192.168.2.0/24", DbSubnetCidr, "Db subnet cidr block does not match.")

	db1_id := terraform.Output(t, terraformOptions, "db1_id")
    db1_public_ip := aws.GetPublicIpOfEc2Instance(t, db1_id, "eu-central-1")
    assert.Equal(t, "", db1_public_ip)

    db2_id := terraform.Output(t, terraformOptions, "db2_id")
    db2_public_ip := aws.GetPublicIpOfEc2Instance(t, db2_id, "eu-central-1")
    assert.Equal(t, "", db2_public_ip)

    db3_id := terraform.Output(t, terraformOptions, "db3_id")
    db3_public_ip := aws.GetPublicIpOfEc2Instance(t, db3_id, "eu-central-1")
    assert.Equal(t, "", db3_public_ip)
}