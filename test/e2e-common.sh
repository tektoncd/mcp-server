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

# Common functions for E2E test scripts

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

# Print a header message
function header() {
  echo
  echo -e "${GREEN}======================================${NC}"
  echo -e "${GREEN}${1}${NC}"
  echo -e "${GREEN}======================================${NC}"
  echo
}

# Print a success message
function success() {
  echo -e "${GREEN}✓ ${1}${NC}"
}

# Print a warning message
function warning() {
  echo -e "${YELLOW}⚠ ${1}${NC}"
}

# Print an error message
function error() {
  echo -e "${RED}✗ ${1}${NC}"
}

# Check if a command exists
function command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Wait for a pod to be ready
function wait_for_pod() {
  local namespace=$1
  local label_selector=$2
  local timeout=${3:-120}
  
  echo "Waiting for pod with selector '${label_selector}' in namespace '${namespace}'..."
  kubectl wait --for=condition=ready pod \
    --namespace="${namespace}" \
    --selector="${label_selector}" \
    --timeout="${timeout}s"
}

# Get the logs of a pod
function get_pod_logs() {
  local namespace=$1
  local label_selector=$2
  local tail_lines=${3:-100}
  
  kubectl logs --namespace="${namespace}" \
    --selector="${label_selector}" \
    --tail="${tail_lines}"
}

# Check if a CRD exists
function crd_exists() {
  local crd_name=$1
  kubectl get crd "${crd_name}" &> /dev/null
}

# Create a namespace if it doesn't exist
function create_namespace() {
  local namespace=$1
  
  if kubectl get namespace "${namespace}" &> /dev/null; then
    echo "Namespace '${namespace}' already exists"
  else
    echo "Creating namespace '${namespace}'..."
    kubectl create namespace "${namespace}"
  fi
}

# Delete a namespace if it exists
function delete_namespace() {
  local namespace=$1
  
  if kubectl get namespace "${namespace}" &> /dev/null; then
    echo "Deleting namespace '${namespace}'..."
    kubectl delete namespace "${namespace}" --ignore-not-found=true
  fi
}

# Apply a YAML file with retries
function apply_with_retry() {
  local yaml_file=$1
  local max_retries=${2:-3}
  local retry_delay=${3:-5}
  
  for i in $(seq 1 "${max_retries}"); do
    if kubectl apply -f "${yaml_file}"; then
      success "Successfully applied ${yaml_file}"
      return 0
    else
      if [[ ${i} -lt ${max_retries} ]]; then
        warning "Failed to apply ${yaml_file}, retrying in ${retry_delay} seconds..."
        sleep "${retry_delay}"
      else
        error "Failed to apply ${yaml_file} after ${max_retries} attempts"
        return 1
      fi
    fi
  done
}

# Check cluster connectivity
function check_cluster() {
  echo "Checking cluster connectivity..."
  if kubectl cluster-info &> /dev/null; then
    success "Connected to Kubernetes cluster"
    kubectl version --short
  else
    error "Cannot connect to Kubernetes cluster"
    return 1
  fi
}

# Check required tools
function check_prerequisites() {
  local tools=("kubectl" "go")
  
  echo "Checking prerequisites..."
  for tool in "${tools[@]}"; do
    if command_exists "${tool}"; then
      success "${tool} is installed"
    else
      error "${tool} is not installed"
      return 1
    fi
  done
  
  # Check Go version
  local go_version=$(go version | awk '{print $3}' | sed 's/go//')
  local required_version="1.21"
  if [[ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" = "$required_version" ]]; then
    success "Go version ${go_version} meets minimum requirement (${required_version})"
  else
    error "Go version ${go_version} does not meet minimum requirement (${required_version})"
    return 1
  fi
}

# Export common environment variables
export TEST_NAMESPACE="${TEST_NAMESPACE:-tekton-pipelines}"
export TEST_TIMEOUT="${TEST_TIMEOUT:-300}"