package k3d

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	//"regexp"
	//"strings"

	"github.com/cfabriziohivemq/k3d-github-action/internal/exec"
)

// ClusterOptions defines the configuration for a k3d cluster
type ClusterOptions struct {
	Name       string
	K3sVersion string
	Ports      []string
	AgentNodes int
	// Add more options as needed
}

// Cluster represents a created k3d cluster
type Cluster struct {
	Name           string
	KubeconfigPath string
}

// CreateCluster creates a new k3d cluster using the CLI
func CreateCluster(ctx context.Context, opts *ClusterOptions) (*Cluster, error) {
	// Start with base command
	args := []string{"cluster", "create", opts.Name}

	// Add version if specified
	if opts.K3sVersion != "latest" && opts.K3sVersion != "" {
		args = append(args, "--image", fmt.Sprintf("rancher/k3s:%s", opts.K3sVersion))
	}

	// Add port mappings
	for _, port := range opts.Ports {
		if port == "" {
			continue
		}
		args = append(args, "-p", port)
	}

	// Add agent nodes if specified
	if opts.AgentNodes > 0 {
		args = append(args, "--agents", fmt.Sprintf("%d", opts.AgentNodes))
	}

	// Wait for the cluster to be ready
	args = append(args, "--wait")

	// Execute the command
	output, err := exec.RunCommand(ctx, "k3d", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w\nOutput: %s", err, output)
	}

	// Get kubeconfig path
	kubeconfigPath, err := getKubeconfig(ctx, opts.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	return &Cluster{
		Name:           opts.Name,
		KubeconfigPath: kubeconfigPath,
	}, nil
}

// DeleteCluster deletes a k3d cluster
func DeleteCluster(ctx context.Context, name string) error {
	output, err := exec.RunCommand(ctx, "k3d", "cluster", "delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete cluster: %w\nOutput: %s", err, output)
	}
	return nil
}

// getKubeconfig gets the path to the kubeconfig for a cluster
func getKubeconfig(ctx context.Context, clusterName string) (string, error) {
	// First check if k3d can write the kubeconfig
	tempKubeconfigDir, err := os.MkdirTemp("", "k3d-kubeconfig-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	tempKubeconfig := filepath.Join(tempKubeconfigDir, "kubeconfig")

	output, err := exec.RunCommand(ctx, "k3d", "kubeconfig", "write", clusterName, "--output", tempKubeconfig)
	if err != nil {
		os.RemoveAll(tempKubeconfigDir) // Clean up
		return "", fmt.Errorf("failed to write kubeconfig: %w\nOutput: %s", err, output)
	}

	return tempKubeconfig, nil
}
