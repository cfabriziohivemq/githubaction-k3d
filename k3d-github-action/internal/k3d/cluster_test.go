package k3d

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestCreateCluster(t *testing.T) {
	// Skip if we're in CI and can't create real clusters
	if _, err := exec.LookPath("k3d"); err != nil {
		t.Skip("k3d not found in PATH, skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := &ClusterOptions{
		Name:       "test-cluster-unit-test",
		K3sVersion: "latest",
	}

	cluster, err := CreateCluster(ctx, opts)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// Clean up after the test
	defer DeleteCluster(context.Background(), opts.Name)

	if cluster.Name != opts.Name {
		t.Errorf("Expected cluster name %s, got %s", opts.Name, cluster.Name)
	}

	if cluster.KubeconfigPath == "" {
		t.Error("Expected kubeconfig path, got empty string")
	}
}
