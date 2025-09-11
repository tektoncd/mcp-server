/*
Copyright 2024 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"knative.dev/pkg/apis"
	knativetest "knative.dev/pkg/test"
	"knative.dev/pkg/test/helpers"
)

const (
	// Default timeout for wait operations
	defaultTimeout = 5 * time.Minute
	// Default interval for polling
	defaultInterval = 1 * time.Second
)

// setup creates a test namespace and clients
func setup(ctx context.Context, t *testing.T) (*MCPTestClients, string) {
	t.Helper()

	namespace := helpers.AppendRandomString("mcp-test")

	// Create namespace
	if err := createNamespace(ctx, t, namespace); err != nil {
		t.Fatalf("Failed to create namespace %s: %v", namespace, err)
	}

	// Setup cleanup
	t.Cleanup(func() {
		t.Logf("Cleaning up namespace %s", namespace)
		deleteNamespace(ctx, t, namespace)
	})

	// Create clients
	clients, err := newMCPTestClients(t, namespace)
	if err != nil {
		t.Fatalf("Failed to create test clients: %v", err)
	}

	t.Cleanup(func() {
		clients.Cleanup()
	})

	// Deploy MCP server if needed
	if deployMCPServer {
		if err := deployMCPServerToNamespace(ctx, t, clients, namespace); err != nil {
			t.Fatalf("Failed to deploy MCP server: %v", err)
		}
	}

	return clients, namespace
}

// createNamespace creates a test namespace
func createNamespace(ctx context.Context, t *testing.T, namespace string) error {
	t.Helper()

	cfg, err := knativetest.BuildClientConfig(knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster)
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"tekton-mcp-test": "true",
			},
		},
	}

	_, err = kubeClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	return err
}

// deleteNamespace deletes a test namespace
func deleteNamespace(ctx context.Context, t *testing.T, namespace string) {
	t.Helper()

	cfg, err := knativetest.BuildClientConfig(knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster)
	if err != nil {
		t.Logf("Failed to build config for cleanup: %v", err)
		return
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Logf("Failed to create client for cleanup: %v", err)
		return
	}

	if err := kubeClient.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{}); err != nil {
		t.Logf("Failed to delete namespace %s: %v", namespace, err)
	}
}

// deployMCPServerToNamespace deploys the MCP server to the test namespace
func deployMCPServerToNamespace(ctx context.Context, t *testing.T, clients *MCPTestClients, namespace string) error {
	t.Helper()

	// This would typically use ko or kubectl to deploy the server
	// For now, we'll assume it's already deployed or running locally
	t.Logf("MCP server deployment to namespace %s (assuming pre-deployed or local)", namespace)
	
	// Wait for MCP server to be ready
	return waitForMCPServerReady(ctx, t, clients)
}

// waitForMCPServerReady waits for the MCP server to be ready to accept connections
func waitForMCPServerReady(ctx context.Context, t *testing.T, clients *MCPTestClients) error {
	t.Helper()

	return wait.PollImmediate(defaultInterval, 30*time.Second, func() (bool, error) {
		// Try to list tools as a health check
		tools, err := clients.MCPSession.CallTool(ctx, "list_pipelines", map[string]interface{}{
			"namespace": clients.Namespace,
		})
		if err != nil {
			t.Logf("MCP server not ready yet: %v", err)
			return false, nil
		}
		return tools != nil, nil
	})
}

// WaitForPipelineRunState waits for a PipelineRun to reach a certain state
func WaitForPipelineRunState(ctx context.Context, c *MCPTestClients, name string, checkFunc func(*v1.PipelineRun) bool, desc string) error {
	return wait.PollImmediate(defaultInterval, defaultTimeout, func() (bool, error) {
		pr, err := c.V1PipelineRunClient.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return checkFunc(pr), nil
	})
}

// WaitForTaskRunState waits for a TaskRun to reach a certain state  
func WaitForTaskRunState(ctx context.Context, c *MCPTestClients, name string, checkFunc func(*v1.TaskRun) bool, desc string) error {
	return wait.PollImmediate(defaultInterval, defaultTimeout, func() (bool, error) {
		tr, err := c.V1TaskRunClient.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return checkFunc(tr), nil
	})
}

// IsSuccessful checks if a PipelineRun/TaskRun completed successfully
func IsSuccessful(t *testing.T, status apis.Status) bool {
	cond := status.GetCondition(apis.ConditionSucceeded)
	if cond == nil {
		t.Logf("No succeeded condition found")
		return false
	}
	t.Logf("Condition: %v", cond)
	return cond.Status == corev1.ConditionTrue
}

// IsFailed checks if a PipelineRun/TaskRun failed
func IsFailed(t *testing.T, status apis.Status) bool {
	cond := status.GetCondition(apis.ConditionSucceeded)
	if cond == nil {
		t.Logf("No succeeded condition found")
		return false
	}
	t.Logf("Condition: %v", cond)
	return cond.Status == corev1.ConditionFalse
}

// generateName generates a unique name for test resources
func generateName(base string) string {
	return helpers.AppendRandomString(base)
}