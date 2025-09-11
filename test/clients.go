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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	v1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1beta1"
	"k8s.io/client-go/kubernetes"
	knativetest "knative.dev/pkg/test"
)

var (
	mcpServerURL    string
	deployMCPServer bool
)

// MCPTestClients holds all the clients needed for E2E tests
type MCPTestClients struct {
	MCPClient                *mcp.Client
	MCPSession               *mcp.ClientSession
	KubeClient               kubernetes.Interface
	V1PipelineClient         v1.PipelineInterface
	V1TaskClient             v1.TaskInterface
	V1TaskRunClient          v1.TaskRunInterface
	V1PipelineRunClient      v1.PipelineRunInterface
	V1beta1StepActionClient  v1beta1.StepActionInterface
	Namespace                string
}

// newMCPTestClients creates and initializes all test clients
func newMCPTestClients(t *testing.T, namespace string) (*MCPTestClients, error) {
	t.Helper()

	// Setup Kubernetes and Tekton clients
	cfg, err := knativetest.BuildClientConfig(knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %w", err)
	}

	tektonClient, err := versioned.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tekton client: %w", err)
	}

	// Setup MCP client
	mcpClient, mcpSession, err := setupMCPClient(t, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to setup MCP client: %w", err)
	}

	return &MCPTestClients{
		MCPClient:               mcpClient,
		MCPSession:              mcpSession,
		KubeClient:              kubeClient,
		V1PipelineClient:        tektonClient.TektonV1().Pipelines(namespace),
		V1TaskClient:            tektonClient.TektonV1().Tasks(namespace),
		V1TaskRunClient:         tektonClient.TektonV1().TaskRuns(namespace),
		V1PipelineRunClient:     tektonClient.TektonV1().PipelineRuns(namespace),
		V1beta1StepActionClient: tektonClient.TektonV1beta1().StepActions(namespace),
		Namespace:               namespace,
	}, nil
}

// setupMCPClient creates and connects an MCP client to the server
func setupMCPClient(t *testing.T, namespace string) (*mcp.Client, *mcp.ClientSession, error) {
	t.Helper()

	// Determine MCP server URL
	serverURL := mcpServerURL
	if serverURL == "" {
		if deployMCPServer {
			// If we're deploying the server, use the service URL in the MCP server namespace
			serverURL = "http://tekton-mcp-server.tekton-mcp-server.svc.cluster.local:3000"
		} else {
			// Default to localhost with port-forward
			serverURL = "http://localhost:3000"
		}
	}

	// Create MCP client
	mcpClient := mcp.NewClient("tekton-e2e-test-client", "1.0.0", nil)

	// Create HTTP transport
	transport := &httpTransport{
		url:    serverURL,
		client: &http.Client{Timeout: 30 * time.Second},
	}

	// Connect to MCP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session, err := mcpClient.Connect(ctx, transport)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MCP server at %s: %w", serverURL, err)
	}

	// Wait for initialization
	if err := session.Initialize(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize MCP session: %w", err)
	}

	return mcpClient, session, nil
}

// httpTransport implements mcp.Transport for HTTP connections
type httpTransport struct {
	url    string
	client *http.Client
}

func (h *httpTransport) Send(message []byte) error {
	resp, err := h.client.Post(h.url, "application/json", bytes.NewReader(message))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	return nil
}

func (h *httpTransport) Receive() ([]byte, error) {
	// For HTTP transport, we would typically use SSE or WebSocket
	// This is a simplified implementation
	return nil, fmt.Errorf("not implemented for basic HTTP transport")
}

func (h *httpTransport) Close() error {
	return nil
}

// Cleanup closes the MCP session
func (c *MCPTestClients) Cleanup() {
	if c.MCPSession != nil {
		c.MCPSession.Close()
	}
}