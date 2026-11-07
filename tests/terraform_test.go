package tests

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestInfrastructure(t *testing.T) {
	ctx := context.Background()
	opts := &terraform.Options{TerraformDir: "../infrastructure"}

	http := terraform.OutputMap(t, opts, "http_ip")
	db := terraform.OutputMap(t, opts, "db_ip")
	vpc := terraform.Output(t, opts, "vpc_id")
	httpSubnet := terraform.Output(t, opts, "http_subnet_id")
	dbSubnet := terraform.Output(t, opts, "db_subnet_id")
	dbSG := terraform.Output(t, opts, "db_security_group_id")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-north-1"))
	if err != nil {
		t.Fatal(err)
	}
	client := ec2.NewFromConfig(cfg)

	// Check EC2 instances.
	if len(http) != 2 || len(db) != 3 {
		t.Fatalf("Expected 2 HTTP and 3 DB instances")
	}

	ids := []string{}
	for id := range http {
		ids = append(ids, id)
	}
	for id := range db {
		ids = append(ids, id)
	}

	instances, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{InstanceIds: ids})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range instances.Reservations {
		for _, i := range r.Instances {
			id := aws.ToString(i.InstanceId)

			if i.State == nil || i.State.Name != "running" {
				t.Errorf("Instance %s is not running", id)
			}
			if _, isDB := db[id]; isDB && aws.ToString(i.PublicIpAddress) != "" {
				t.Errorf("DB instance %s has a public IP", id)
			}
		}
	}

	// Check VPC and subnet CIDRs.
	vpcs, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{VpcIds: []string{vpc}})
	if err != nil {
		t.Fatal(err)
	}
	if len(vpcs.Vpcs) != 1 || aws.ToString(vpcs.Vpcs[0].CidrBlock) != "192.168.0.0/16" {
		t.Error("Incorrect VPC CIDR")
	}

	subnets, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{SubnetIds: []string{httpSubnet, dbSubnet}})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range subnets.Subnets {
		cidr := aws.ToString(s.CidrBlock)
		if aws.ToString(s.SubnetId) == httpSubnet && cidr != "192.168.1.0/24" {
			t.Error("Incorrect HTTP subnet CIDR")
		}
		if aws.ToString(s.SubnetId) == dbSubnet && cidr != "192.168.2.0/24" {
			t.Error("Incorrect DB subnet CIDR")
		}
	}

	// Check that the database is not accessible from the internet.
	groups, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{GroupIds: []string{dbSG}})
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range groups.SecurityGroups[0].IpPermissions {
		protocol := aws.ToString(p.IpProtocol)
		allowsDB := protocol == "-1" ||
			(protocol == "tcp" && p.FromPort != nil && p.ToPort != nil &&
				*p.FromPort <= 3306 && *p.ToPort >= 3306)

		if allowsDB {
			for _, ip := range p.IpRanges {
				if aws.ToString(ip.CidrIp) == "0.0.0.0/0" {
					t.Error("Database is accessible from the internet")
				}
			}
			for _, ip := range p.Ipv6Ranges {
				if aws.ToString(ip.CidrIpv6) == "::/0" {
					t.Error("Database is accessible from the internet via IPv6")
				}
			}
		}
	}
}
