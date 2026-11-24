package test

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Region  string `json:"region"`
	VPCID   string `json:"vpc_id"`
	VPCCIDR string `json:"vpc_cidr"`
	WebSGID string `json:"web_sg_id"`
	Subnets map[string]struct {
		ID   string `json:"id"`
		CIDR string `json:"cidr"`
	} `json:"subnets"`
	Instances map[string]map[string]string `json:"instances"`
}

type instance struct {
	InstanceID, VpcID, SubnetID, PrivateIPAddress string
	PublicIPAddress, PublicDNSName                string
	State                                         struct{ Name string }
	Tags                                          []struct{ Key, Value string }
	SecurityGroups                                []struct{ GroupID string }
	NetworkInterfaces                             []struct {
		Association   struct{ PublicIP string }
		IPv6Addresses []struct{ IPv6Address string }
	}
}

func awsRead(t *testing.T, region string, result any, args ...string) {
	t.Helper()
	args = append([]string{
		"ec2", "--region", region,
		"--output", "json", "--no-cli-pager",
	}, args...)

	cmd := exec.CommandContext(t.Context(), "aws", args...)
	out, err := cmd.Output()
	if failure, ok := err.(*exec.ExitError); ok {
		t.Logf("AWS CLI: %s", failure.Stderr)
	}
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(out, result))
}

func TestPostApply(t *testing.T) {
	var cfg testConfig
	terraform.OutputStructContext(t, t.Context(), &terraform.Options{
		TerraformDir: "../infrastructure",
	}, "test_config", &cfg)

	require.NotEmpty(t, cfg.Region)
	require.NotEmpty(t, cfg.VPCID)
	require.NotEmpty(t, cfg.WebSGID)

	instances := make(map[string]instance)
	for _, role := range []string{"http", "db"} {
		require.NotEmpty(t, cfg.Instances[role])
		for _, id := range cfg.Instances[role] {
			var actual instance
			awsRead(t, cfg.Region, &actual,
				"describe-instances", "--instance-ids", id,
				"--query", "Reservations[0].Instances[0]")
			instances[id] = actual
		}
	}

	t.Run("Network", func(t *testing.T) {
		var vpc struct{ VpcID, CidrBlock, State string }
		awsRead(t, cfg.Region, &vpc,
			"describe-vpcs", "--vpc-ids", cfg.VPCID,
			"--query", "Vpcs[0]")

		require.Equal(t, cfg.VPCID, vpc.VpcID)
		require.Equal(t, cfg.VPCCIDR, vpc.CidrBlock)
		require.Equal(t, "available", vpc.State)

		for _, role := range []string{"http", "db"} {
			expected := cfg.Subnets[role]
			require.NotEmpty(t, expected.ID)

			var subnet struct{ SubnetID, VpcID, CidrBlock, State string }
			awsRead(t, cfg.Region, &subnet,
				"describe-subnets", "--subnet-ids", expected.ID,
				"--query", "Subnets[0]")

			require.Equal(t, expected.ID, subnet.SubnetID)
			require.Equal(t, cfg.VPCID, subnet.VpcID)
			require.Equal(t, expected.CIDR, subnet.CidrBlock)
			require.Equal(t, "available", subnet.State)
		}
	})

	t.Run("EC2Hosts", func(t *testing.T) {
		for _, role := range []string{"http", "db"} {
			for name, id := range cfg.Instances[role] {
				actual := instances[id]
				require.Equal(t, id, actual.InstanceID, name)
				require.Equal(t, "running", actual.State.Name, name)
				require.Equal(t, cfg.VPCID, actual.VpcID, name)
				require.Equal(t, cfg.Subnets[role].ID, actual.SubnetID, name)
				require.NotEmpty(t, actual.PrivateIPAddress, name)

				var nameTag string
				for _, tag := range actual.Tags {
					if tag.Key == "Name" {
						nameTag = tag.Value
					}
				}
				require.Equal(t, name, nameTag)
			}
		}
	})

	t.Run("DatabasePrivate", func(t *testing.T) {
		for name, id := range cfg.Instances["db"] {
			actual := instances[id]
			require.Empty(t, actual.PublicIPAddress, name)
			require.Empty(t, actual.PublicDNSName, name)

			for _, nic := range actual.NetworkInterfaces {
				require.Empty(t, nic.Association.PublicIP, name)
				require.Empty(t, nic.IPv6Addresses, name)
			}

			require.NotEmpty(t, actual.SecurityGroups, name)
			foundDatabaseRule := false

			for _, attached := range actual.SecurityGroups {
				var group struct {
					IpPermissions []struct {
						IpProtocol       string
						FromPort, ToPort int
						IpRanges         []struct{ CidrIP string }
						IPv6Ranges       []struct{ CidrIPv6 string }
						PrefixListIDs    []struct{ PrefixListID string }
						UserIDGroupPairs []struct{ GroupID string }
					}
				}

				awsRead(t, cfg.Region, &group,
					"describe-security-groups", "--group-ids", attached.GroupID,
					"--query", "SecurityGroups[0]")

				for _, rule := range group.IpPermissions {
					allowsSQL := rule.IpProtocol == "-1" ||
						((rule.IpProtocol == "tcp" || rule.IpProtocol == "6") &&
							rule.FromPort <= 3306 && rule.ToPort >= 3306)

					if !allowsSQL {
						continue
					}

					foundDatabaseRule = true
					require.Empty(t, rule.IpRanges, name)
					require.Empty(t, rule.IPv6Ranges, name)
					require.Empty(t, rule.PrefixListIDs, name)
					require.NotEmpty(t, rule.UserIDGroupPairs, name)

					for _, source := range rule.UserIDGroupPairs {
						require.Equal(t, cfg.WebSGID, source.GroupID, name)
					}
				}
			}

			require.True(t, foundDatabaseRule,
				"%s must allow SQL from the web security group", name)
		}
	})
}
