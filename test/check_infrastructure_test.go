package test

import (
    "testing"

    awsSdk "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    ec2 "github.com/aws/aws-sdk-go/service/ec2"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestCheckInfrastructure(t *testing.T) {
    t.Parallel()

    tf := &terraform.Options{TerraformDir: "../infrastructure"}
    terraform.Init(t, tf)

    assert.Greater(t, len(terraform.OutputList(t, tf, "instance_ids")), 0, "EC2 instances must be created")
    assert.Greater(t, len(terraform.OutputList(t, tf, "db_instance_ids")), 0, "DB instances must be created")

    sess := session.Must(session.NewSession(&awsSdk.Config{Region: awsSdk.String("us-east-1")}))
    ec2Client := ec2.New(sess)

    vpcID := terraform.Output(t, tf, "vpc_id")
    vpcOut, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{VpcIds: []*string{awsSdk.String(vpcID)}})
    if err != nil || len(vpcOut.Vpcs) == 0 {
        t.Fatal("VPC not found or describe failed")
    }
    assert.Equal(t, "192.168.0.0/16", *vpcOut.Vpcs[0].CidrBlock)

    for _, id := range terraform.OutputList(t, tf, "db_instance_ids") {
        res, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{awsSdk.String(id)}})
        if err != nil || len(res.Reservations[0].Instances) == 0 {
            t.Fatalf("EC2 DB Instance %s not found or describe failed", id)
        }
        assert.Nil(t, res.Reservations[0].Instances[0].PublicIpAddress, "DB instance must NOT have public IP")
    }
}
