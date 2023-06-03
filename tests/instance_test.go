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
        TerraformDir: "/home/circleci/project/infrastructure",
    })

    testInstanceNames(t, terraformOptions)
    testCidrBlocks(t, terraformOptions)
    testDbInstances(t, terraformOptions)
}

func testInstanceCreation(t *testing.T, terraformOptions *terraform.Options) {
    // Перевірка створення інстансів EC2
    ec2ResourceType := "aws_instance"
    ec2Count := terraform.GetResourceCount(t, terraformOptions, ec2ResourceType)
    assert.Greater(t, ec2Count, 0, "No EC2 instances were created.")

    // Перевірка створення баз даних
    dbResourceType := "aws_db_instance"
    dbCount := terraform.GetResourceCount(t, terraformOptions, dbResourceType)
    assert.Greater(t, dbCount, 0, "No database instances were created.")
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
