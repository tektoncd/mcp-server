#!/usr/bin/env bash

# Copyright 2024 The Tekton Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script runs the end-to-end tests for the Tekton MCP Server

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/e2e-common.sh

# Default settings
SKIP_MCP_DEPLOY=${SKIP_MCP_DEPLOY:-"false"}
MCP_SERVER_IMAGE=${MCP_SERVER_IMAGE:-""}
E2E_GO_TEST_TIMEOUT=${E2E_GO_TEST_TIMEOUT:-"20m"}
SKIP_CLEANUP=${SKIP_CLEANUP:-"false"}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --skip-mcp-deploy)
      SKIP_MCP_DEPLOY="true"
      shift
      ;;
    --mcp-server-image)
      MCP_SERVER_IMAGE="$2"
      shift 2
      ;;
    --skip-cleanup)
      SKIP_CLEANUP="true"
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

header "Setting up test environment"

# Ensure we have a cluster to test against
if [[ -z "${KUBECONFIG:-}" ]]; then
  echo "KUBECONFIG not set, using default"
  export KUBECONFIG="${HOME}/.kube/config"
fi

# Check if kubectl can connect to cluster
if ! kubectl cluster-info &> /dev/null; then
  echo "ERROR: Cannot connect to Kubernetes cluster"
  echo "Please ensure you have a running cluster and KUBECONFIG is set correctly"
  exit 1
fi

# Check if Tekton Pipelines is installed
if ! kubectl get crd pipelines.tekton.dev &> /dev/null; then
  echo "ERROR: Tekton Pipelines CRDs not found"
  echo "Please install Tekton Pipelines first:"
  echo "  kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml"
  exit 1
fi

header "Deploying MCP Server"

if [[ "${SKIP_MCP_DEPLOY}" == "false" ]]; then
  echo "Building and deploying MCP server..."
  
  # If no image specified, build with ko
  if [[ -z "${MCP_SERVER_IMAGE}" ]]; then
    if ! command -v ko &> /dev/null; then
      echo "ERROR: ko is not installed. Please install ko or provide --mcp-server-image"
      exit 1
    fi
    
    # Set KO_DOCKER_REPO if not already set
    if [[ -z "${KO_DOCKER_REPO:-}" ]]; then
      echo "ERROR: KO_DOCKER_REPO not set. Please set it to your image registry"
      echo "Example: export KO_DOCKER_REPO=gcr.io/my-project"
      exit 1
    fi
    
    echo "Building MCP server image with ko..."
    ko apply -B -f config/
  else
    echo "Using pre-built image: ${MCP_SERVER_IMAGE}"
    # Replace the image in the deployment and apply all configs
    for file in config/*.yaml; do
      if [[ "$file" == "config/300-deployment.yaml" ]]; then
        sed "s|image: ko://.*|image: ${MCP_SERVER_IMAGE}|" "$file" | kubectl apply -f -
      else
        kubectl apply -f "$file"
      fi
    done
  fi
  
  # Wait for deployment to be ready
  echo "Waiting for MCP server deployment to be ready..."
  kubectl wait --for=condition=available --timeout=120s deployment/tekton-mcp-server -n tekton-mcp-server || {
    echo "ERROR: MCP server deployment failed to become ready"
    echo "Deployment status:"
    kubectl get deployment tekton-mcp-server -n tekton-mcp-server
    echo "Pod status:"
    kubectl get pods -n tekton-mcp-server -l app.kubernetes.io/name=tekton-mcp-server
    echo "Pod logs:"
    kubectl logs -n tekton-mcp-server -l app.kubernetes.io/name=tekton-mcp-server --tail=50
    exit 1
  }
  
  # Setup port-forward for local testing
  echo "Setting up port-forward to MCP server..."
  kubectl port-forward -n tekton-mcp-server service/tekton-mcp-server 3000:3000 &
  PORT_FORWARD_PID=$!
  sleep 5  # Give port-forward time to establish
  
  # Test that MCP server is responding
  if ! curl -s http://localhost:3000/health &> /dev/null; then
    echo "WARNING: MCP server health check failed, but continuing with tests"
  fi
  
  export MCP_SERVER_URL="http://localhost:3000"
else
  echo "Skipping MCP server deployment (--skip-mcp-deploy was set)"
  if [[ -z "${MCP_SERVER_URL:-}" ]]; then
    export MCP_SERVER_URL="http://localhost:3000"
    echo "Using default MCP_SERVER_URL: ${MCP_SERVER_URL}"
  fi
fi

header "Running E2E Tests"

# Run the Go tests
echo "Running E2E tests with timeout ${E2E_GO_TEST_TIMEOUT}..."
failed=0

go test -v -count=1 -tags=e2e -timeout="${E2E_GO_TEST_TIMEOUT}" ./test \
  -mcp-server-url="${MCP_SERVER_URL}" \
  -deploy-mcp-server=false \
  --kubeconfig="${KUBECONFIG}" || failed=1

header "Cleanup"

if [[ "${SKIP_CLEANUP}" == "false" ]]; then
  if [[ "${SKIP_MCP_DEPLOY}" == "false" ]]; then
    echo "Cleaning up MCP server deployment..."
    
    # Kill port-forward if it's running
    if [[ -n "${PORT_FORWARD_PID:-}" ]]; then
      kill ${PORT_FORWARD_PID} 2>/dev/null || true
    fi
    
    # Delete the deployment and all resources
    kubectl delete -f config/ --ignore-not-found=true
  fi
  
  # Clean up test namespaces
  echo "Cleaning up test namespaces..."
  kubectl delete namespaces -l tekton-mcp-test=true --ignore-not-found=true
else
  echo "Skipping cleanup (--skip-cleanup was set)"
fi

if [[ ${failed} -eq 0 ]]; then
  success "E2E tests passed!"
else
  echo "E2E tests failed!"
  exit 1
fi