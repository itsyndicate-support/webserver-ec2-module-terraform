package tests

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/stretchr/testify/require"
)

func TestInfrastructure(t *testing.T) {
	ctx := t.Context()
	region := os.Getenv("AWS_REGION")
	require.NotEmpty(t, region, "AWS_REGION must be set")

	tf := &terraform.Options{
		TerraformDir: "../infrastructure",
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	require.NoError(t, err)
	client := ec2.NewFromConfig(cfg)

	var vpcID, webGroupID string

	for _, group := range []struct {
		output string
		count  int
		cidr   string
	}{
		{"http_ip", 2, "192.168.1.0/24"},
		{"db_ip", 3, "192.168.2.0/24"},
	} {
		instances := terraform.OutputMapContext(t, ctx, tf, group.output)
		require.Len(t, instances, group.count, group.output)

		for id := range instances {
			result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
				InstanceIds: []string{id},
			})
			require.NoError(t, err)
			require.Len(t, result.Reservations, 1)
			require.Len(t, result.Reservations[0].Instances, 1)

			instance := result.Reservations[0].Instances[0]
			require.NotNil(t, instance.State)
			require.Equal(t, types.InstanceStateNameRunning, instance.State.Name, id)

			if vpcID == "" {
				vpcID = aws.ToString(instance.VpcId)
				require.NotEmpty(t, vpcID)

				vpcs, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
					VpcIds: []string{vpcID},
				})
				require.NoError(t, err)
				require.Len(t, vpcs.Vpcs, 1)
				require.Equal(t, "192.168.0.0/16", aws.ToString(vpcs.Vpcs[0].CidrBlock))

				webGroups, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
					Filters: []types.Filter{
						{Name: aws.String("vpc-id"), Values: []string{vpcID}},
						{Name: aws.String("group-name"), Values: []string{"web"}},
					},
				})
				require.NoError(t, err)
				require.Len(t, webGroups.SecurityGroups, 1)
				webGroupID = aws.ToString(webGroups.SecurityGroups[0].GroupId)
				require.NotEmpty(t, webGroupID)
			}

			require.Equal(t, vpcID, aws.ToString(instance.VpcId), id)

			subnets, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
				SubnetIds: []string{aws.ToString(instance.SubnetId)},
			})
			require.NoError(t, err)
			require.Len(t, subnets.Subnets, 1)
			require.Equal(t, group.cidr, aws.ToString(subnets.Subnets[0].CidrBlock), id)
			require.Equal(t, vpcID, aws.ToString(subnets.Subnets[0].VpcId), id)

			if group.output != "db_ip" {
				continue
			}

			var securityGroupIDs []string
			for _, sg := range instance.SecurityGroups {
				securityGroupIDs = append(securityGroupIDs, aws.ToString(sg.GroupId))
			}
			require.NotEmpty(t, securityGroupIDs, id)

			groups, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
				GroupIds: securityGroupIDs,
			})
			require.NoError(t, err)
			require.Len(t, groups.SecurityGroups, len(securityGroupIDs))

			mysqlAccessFound := false

			for _, sg := range groups.SecurityGroups {
				for _, rule := range sg.IpPermissions {
					protocol := aws.ToString(rule.IpProtocol)
					allowsMySQL := protocol == "-1" ||
						((protocol == "tcp" || protocol == "6") &&
							aws.ToInt32(rule.FromPort) <= 3306 &&
							aws.ToInt32(rule.ToPort) >= 3306)

					if !allowsMySQL {
						continue
					}

					mysqlAccessFound = true
					require.Empty(t, rule.IpRanges, "DB %s allows CIDR access to port 3306", id)
					require.Empty(t, rule.Ipv6Ranges, "DB %s allows IPv6 access to port 3306", id)
					require.Empty(t, rule.PrefixListIds, "DB %s allows prefix-list access to port 3306", id)
					require.NotEmpty(t, rule.UserIdGroupPairs)

					for _, source := range rule.UserIdGroupPairs {
						require.Equal(t, webGroupID, aws.ToString(source.GroupId),
							"DB %s must allow port 3306 only from the web security group", id)
					}
				}
			}

			require.True(t, mysqlAccessFound, "DB %s has no access rule for port 3306", id)
			t.Logf("DB %s: port 3306 is restricted to the web security group", id)
		}
	}
}
