#!/bin/bash
# test-local.sh

# Configuration
CLUSTER_NAME=${1:-test-local}  # Use first argument or default to "test-local"
TEST_DIR=/tmp/k3d-test
IMAGE=cfabriziohivemq/k3d-github-action

# Create test directory
mkdir -p $TEST_DIR

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Clean existing cluster with same name if exists
echo -e "${YELLOW}Checking for existing cluster...${NC}"
if k3d cluster list 2>/dev/null | grep -q "$CLUSTER_NAME"; then
  echo -e "${YELLOW}Deleting existing cluster $CLUSTER_NAME...${NC}"
  k3d cluster delete "$CLUSTER_NAME"
fi

# Run the container
echo -e "${YELLOW}Creating cluster with container...${NC}"
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $TEST_DIR:/tmp/output \
  $IMAGE \
  --name "$CLUSTER_NAME" \
  --ports 8080:80

# Check if successful
if [ $? -ne 0 ]; then
  echo -e "${RED}Container execution failed!${NC}"
  exit 1
fi

# Verify cluster was created
echo -e "${YELLOW}Verifying cluster was created...${NC}"
if k3d cluster list | grep -q "$CLUSTER_NAME"; then
  echo -e "${GREEN}Cluster '$CLUSTER_NAME' created successfully!${NC}"
else
  echo -e "${RED}Failed to find cluster '$CLUSTER_NAME'${NC}"
  exit 1
fi

# Get kubeconfig and test access
echo -e "${YELLOW}Testing Kubernetes access...${NC}"
k3d kubeconfig get "$CLUSTER_NAME" > $TEST_DIR/kubeconfig
export KUBECONFIG=$TEST_DIR/kubeconfig

if kubectl get nodes; then
  echo -e "${GREEN}Successfully connected to cluster!${NC}"

  # Show running pods
  echo -e "${YELLOW}Pods running in the cluster:${NC}"
  kubectl get pods -A
else
  echo -e "${RED}Failed to connect to cluster${NC}"
  exit 1
fi

# Ask to delete the cluster
read -p "Do you want to delete the test cluster? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo -e "${YELLOW}Cleaning up...${NC}"
  k3d cluster delete "$CLUSTER_NAME"
  echo -e "${GREEN}Cleanup complete!${NC}"
fi
