package test

import (
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/stretchr/testify/require"
)

func TestDatabaseBootstrap(t *testing.T) {
	var cfg testConfig
	terraform.OutputStructContext(t, t.Context(), &terraform.Options{
		TerraformDir: "../infrastructure",
	}, "test_config", &cfg)

	require.NotEmpty(t, cfg.Region)
	require.NotEmpty(t, cfg.Instances["db"])

	pending := make(map[string]string)
	for name, id := range cfg.Instances["db"] {
		pending[name] = id
	}

	deadline := time.Now().Add(8 * time.Minute)
	for len(pending) > 0 {
		for name, id := range pending {
			var console string
			awsRead(t, cfg.Region, &console,
				"get-console-output", "--instance-id", id,
				"--latest", "--query", "Output",
				"--cli-connect-timeout", "10", "--cli-read-timeout", "20")

			if strings.Contains(console, "MARIADB_BOOTSTRAP_FAILED") {
				t.Fatalf("MariaDB bootstrap failed on %s (%s); inspect its console output", name, id)
			}
			if strings.Contains(console, "MARIADB_BOOTSTRAP_READY") {
				t.Logf("%s (%s): MariaDB started and SELECT 1 succeeded during bootstrap", name, id)
				delete(pending, name)
			}
		}

		if len(pending) == 0 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("No MariaDB readiness signal within 8 minutes for: %v", pending)
		}
		select {
		case <-time.After(10 * time.Second):
		case <-t.Context().Done():
			t.Fatal("Cancelled while waiting for MariaDB bootstrap")
		}
	}
}
