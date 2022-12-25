package test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/gruntwork-io/terratest/modules/aws"
)

func TestMytests(t *testing.T) {
    t.Parallel()

    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{

		TerraformDir: "/home/circleci/project/infrastructure",
	})

    terraform.Init(t, terraformOptions)
    first_web_server_ssh_key := terraform.Output(t, terraformOptions, "first_web_server_ssh_key")
    assert.Equal(t, "gl-Frankfurt", first_web_server_ssh_key)

    second_web_server_ssh_key := terraform.Output(t, terraformOptions, "second_web_server_ssh_key")
    assert.Equal(t, "gl-Frankfurt", second_web_server_ssh_key)

    first_db_server_ssh_key := terraform.Output(t, terraformOptions, "first_db_server_ssh_key")
    assert.Equal(t, "gl-Frankfurt", first_db_server_ssh_key)

    second_db_server_ssh_key := terraform.Output(t, terraformOptions, "second_db_server_ssh_key")
    assert.Equal(t, "gl-Frankfurt", second_db_server_ssh_key)

    third_db_server_ssh_key := terraform.Output(t, terraformOptions, "third_db_server_ssh_key")
    assert.Equal(t, "gl-Frankfurt", third_db_server_ssh_key)

    vpc_cidr_block := terraform.Output(t, terraformOptions, "vpc_cidr_block")
    assert.Equal(t, "192.168.0.0/16", vpc_cidr_block)

    http_subnet_cidr_block := terraform.Output(t, terraformOptions, "http_subnet_cidr_block")
    assert.Equal(t, "192.168.1.0/24", http_subnet_cidr_block)

    db_subnet_cidr_block := terraform.Output(t, terraformOptions, "db_subnet_cidr_block")
    assert.Equal(t, "192.168.2.0/24", db_subnet_cidr_block)

    first_db_serverID := terraform.Output(t, terraformOptions, "first_db_server_id")
    first_db_server_piblic_ip := aws.GetPublicIpOfEc2Instance(t, first_db_serverID, "eu-central-1")
    assert.Equal(t, "", first_db_server_piblic_ip)

    second_db_serverID := terraform.Output(t, terraformOptions, "second_db_server_id")
    second_db_server_piblic_ip := aws.GetPublicIpOfEc2Instance(t, second_db_serverID, "eu-central-1")
    assert.Equal(t, "", second_db_server_piblic_ip)

    third_db_serverID := terraform.Output(t, terraformOptions, "third_db_server_id")
    third_db_server_piblic_ip := aws.GetPublicIpOfEc2Instance(t, third_db_serverID, "eu-central-1")
    assert.Equal(t, "", third_db_server_piblic_ip)

}
