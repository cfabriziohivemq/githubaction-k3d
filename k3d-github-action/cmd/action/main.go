package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cfabriziohivemq/k3d-github-action/internal/k3d"
)

func main() {
	// Parse flags for GitHub Action inputs
	clusterName := flag.String("name", "test-cluster", "Name of the k3d cluster")
	k3sVersion := flag.String("version", "latest", "K3s version to use")
	ports := flag.String("ports", "", "Comma-separated list of ports to expose (e.g., '8080:80,8443:443')")
	agents := flag.Int("agents", 0, "Number of agent nodes")
	timeout := flag.Duration("timeout", 5*time.Minute, "Timeout for operations")
	flag.Parse()

	// Create a cancelable context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Handle termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Received termination signal, cleaning up...")
		cancel()
	}()

	// Set up cluster options
	opts := &k3d.ClusterOptions{
		Name:       *clusterName,
		K3sVersion: *k3sVersion,
		Ports:      strings.Split(*ports, ","),
		AgentNodes: *agents,
	}

	// Create the cluster
	log.Printf("Creating k3d cluster: %s", opts.Name)
	cluster, err := k3d.CreateCluster(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	// Output for GitHub Actions
	fmt.Printf("::set-output name=cluster-name::%s\n", cluster.Name)
	fmt.Printf("::set-output name=kubeconfig-path::%s\n", cluster.KubeconfigPath)

	// For modern GitHub Actions (>=v1.12.0)
	fmt.Printf("cluster-name=%s\n", cluster.Name)
	fmt.Printf("kubeconfig-path=%s\n", cluster.KubeconfigPath)

	log.Printf("✅ Successfully created cluster: %s", cluster.Name)
	log.Printf("📄 Kubeconfig: %s", cluster.KubeconfigPath)
}
