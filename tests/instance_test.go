package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/testing"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/gruntwork-io/terratest/modules/aws"
)

func TestTests(t *testing.T) {
    t.Parallel()

    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir: "/home/circleci/project/infrastructure",
    })

    testInstanceExistence(t, terraformOptions)
    testCidrBlocks(t, terraformOptions)
    testDbInstances(t, terraformOptions)
}

func testInstanceExistence(t *testing.T, terraformOptions *terraform.Options) {
    instanceIDs := []string{"instance_id1", "instance_id2"}

    for _, instanceID := range instanceIDs {
        exists := terraform.ResourceExists(t, terraformOptions, instanceID)
        assert.True(t, exists, "Instance does not exist: %s", instanceID)
    }

    dbInstances := []string{"db1_id", "db2_id", "db3_id"}

    for _, dbInstance := range dbInstances {
        exists := terraform.ResourceExists(t, terraformOptions, dbInstance)
        assert.True(t, exists, "Database instance does not exist: %s", dbInstance)
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

    for _, tt := range dbInstances {
        db_id := terraform.Output(t, terraformOptions, tt)
        db_public_ip := aws.GetPublicIpOfEc2Instance(t, db_id, "eu-central-1")
        assert.Equal(t, "", db_public_ip)
    }
}
