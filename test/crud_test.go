//go:build e2e
// +build e2e

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

	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestPipelineCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c, namespace := setup(ctx, t)

	pipelineName := generateName("test-pipeline")

	t.Run("Create_Pipeline", func(t *testing.T) {
		pipelineYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: %s
  namespace: %s
spec:
  params:
  - name: message
    type: string
    default: "Hello World"
  tasks:
  - name: echo-task
    taskSpec:
      params:
      - name: message
        type: string
      steps:
      - name: echo
        image: mirror.gcr.io/ubuntu
        script: |
          #!/bin/bash
          echo "$(params.message)"
    params:
    - name: message
      value: "$(params.message)"
`, pipelineName, namespace)

		// Create pipeline via MCP
		result, err := c.MCPSession.CallTool(ctx, "create_pipeline", map[string]interface{}{
			"namespace": namespace,
			"yaml":      pipelineYAML,
		})
		if err != nil {
			t.Fatalf("Failed to create pipeline via MCP: %v", err)
		}

		// Verify the result
		if !isSuccessResult(result) {
			t.Fatalf("Create pipeline failed: %v", result)
		}

		// Verify via Kubernetes API
		pipeline, err := c.V1PipelineClient.Get(ctx, pipelineName, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to get pipeline via K8s API: %v", err)
		}
		if pipeline.Name != pipelineName {
			t.Errorf("Pipeline name mismatch: got %s, want %s", pipeline.Name, pipelineName)
		}
	})

	t.Run("Get_Pipeline", func(t *testing.T) {
		// Get pipeline via MCP
		result, err := c.MCPSession.CallTool(ctx, "get_pipeline", map[string]interface{}{
			"name":      pipelineName,
			"namespace": namespace,
			"output":    "yaml",
		})
		if err != nil {
			t.Fatalf("Failed to get pipeline via MCP: %v", err)
		}

		// Verify we got YAML back
		yamlContent := getResultContent(result)
		if yamlContent == "" {
			t.Fatal("No YAML content returned from get_pipeline")
		}

		// Parse the YAML to verify it's valid
		var pipelineObj map[string]interface{}
		if err := yaml.Unmarshal([]byte(yamlContent), &pipelineObj); err != nil {
			t.Fatalf("Failed to parse returned YAML: %v", err)
		}

		// Verify the name matches
		if metadata, ok := pipelineObj["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				if name != pipelineName {
					t.Errorf("Pipeline name mismatch in YAML: got %s, want %s", name, pipelineName)
				}
			}
		}
	})

	t.Run("List_Pipelines", func(t *testing.T) {
		// List pipelines via MCP
		result, err := c.MCPSession.CallTool(ctx, "list_pipelines", map[string]interface{}{
			"namespace": namespace,
			"prefix":    "test-pipeline",
		})
		if err != nil {
			t.Fatalf("Failed to list pipelines via MCP: %v", err)
		}

		// Verify our pipeline is in the list
		content := getResultContent(result)
		if content == "" {
			t.Fatal("No content returned from list_pipelines")
		}

		// Check that our pipeline name appears in the list
		if !contains(content, pipelineName) {
			t.Errorf("Pipeline %s not found in list", pipelineName)
		}
	})

	t.Run("Update_Pipeline", func(t *testing.T) {
		updatedYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: %s
  namespace: %s
  labels:
    updated: "true"
spec:
  params:
  - name: message
    type: string
    default: "Updated Hello World"
  tasks:
  - name: echo-task
    taskSpec:
      params:
      - name: message
        type: string
      steps:
      - name: echo
        image: mirror.gcr.io/ubuntu
        script: |
          #!/bin/bash
          echo "Updated: $(params.message)"
    params:
    - name: message
      value: "$(params.message)"
`, pipelineName, namespace)

		// Update pipeline via MCP
		result, err := c.MCPSession.CallTool(ctx, "update_pipeline", map[string]interface{}{
			"name":      pipelineName,
			"namespace": namespace,
			"yaml":      updatedYAML,
		})
		if err != nil {
			t.Fatalf("Failed to update pipeline via MCP: %v", err)
		}

		if !isSuccessResult(result) {
			t.Fatalf("Update pipeline failed: %v", result)
		}

		// Verify the update via Kubernetes API
		pipeline, err := c.V1PipelineClient.Get(ctx, pipelineName, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to get updated pipeline: %v", err)
		}

		// Check that the label was added
		if val, ok := pipeline.Labels["updated"]; !ok || val != "true" {
			t.Error("Pipeline was not updated with new label")
		}
	})

	t.Run("Delete_Pipeline", func(t *testing.T) {
		// Delete pipeline via MCP
		result, err := c.MCPSession.CallTool(ctx, "delete_pipeline", map[string]interface{}{
			"name":      pipelineName,
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to delete pipeline via MCP: %v", err)
		}

		if !isSuccessResult(result) {
			t.Fatalf("Delete pipeline failed: %v", result)
		}

		// Verify deletion via Kubernetes API
		_, err = c.V1PipelineClient.Get(ctx, pipelineName, metav1.GetOptions{})
		if err == nil {
			t.Error("Pipeline still exists after deletion")
		}
	})
}

func TestTaskCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c, namespace := setup(ctx, t)

	taskName := generateName("test-task")

	t.Run("Create_Task", func(t *testing.T) {
		taskYAML := fmt.Sprintf(`
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: %s
  namespace: %s
spec:
  params:
  - name: greeting
    type: string
    default: "Hello"
  steps:
  - name: greet
    image: mirror.gcr.io/ubuntu
    script: |
      #!/bin/bash
      echo "$(params.greeting) from Task"
`, taskName, namespace)

		// Create task via MCP
		result, err := c.MCPSession.CallTool(ctx, "create_task", map[string]interface{}{
			"namespace": namespace,
			"yaml":      taskYAML,
		})
		if err != nil {
			t.Fatalf("Failed to create task via MCP: %v", err)
		}

		if !isSuccessResult(result) {
			t.Fatalf("Create task failed: %v", result)
		}

		// Verify via Kubernetes API
		task, err := c.V1TaskClient.Get(ctx, taskName, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to get task via K8s API: %v", err)
		}
		if task.Name != taskName {
			t.Errorf("Task name mismatch: got %s, want %s", task.Name, taskName)
		}
	})

	t.Run("Get_Task", func(t *testing.T) {
		// Get task via MCP with JSON output
		result, err := c.MCPSession.CallTool(ctx, "get_task", map[string]interface{}{
			"name":      taskName,
			"namespace": namespace,
			"output":    "json",
		})
		if err != nil {
			t.Fatalf("Failed to get task via MCP: %v", err)
		}

		content := getResultContent(result)
		if content == "" {
			t.Fatal("No content returned from get_task")
		}

		// Verify it contains the task name
		if !contains(content, taskName) {
			t.Errorf("Task name %s not found in response", taskName)
		}
	})

	t.Run("List_Tasks", func(t *testing.T) {
		// List tasks via MCP
		result, err := c.MCPSession.CallTool(ctx, "list_tasks", map[string]interface{}{
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to list tasks via MCP: %v", err)
		}

		content := getResultContent(result)
		if !contains(content, taskName) {
			t.Errorf("Task %s not found in list", taskName)
		}
	})

	t.Run("Delete_Task", func(t *testing.T) {
		// Delete task via MCP
		result, err := c.MCPSession.CallTool(ctx, "delete_task", map[string]interface{}{
			"name":      taskName,
			"namespace": namespace,
		})
		if err != nil {
			t.Fatalf("Failed to delete task via MCP: %v", err)
		}

		if !isSuccessResult(result) {
			t.Fatalf("Delete task failed: %v", result)
		}

		// Verify deletion
		_, err = c.V1TaskClient.Get(ctx, taskName, metav1.GetOptions{})
		if err == nil {
			t.Error("Task still exists after deletion")
		}
	})
}

// Helper functions

func isSuccessResult(result *mcp.CallToolResultFor[any]) bool {
	if result == nil {
		return false
	}
	// Check if there's an error in the result
	if len(result.Content) > 0 {
		// Success typically has content
		return true
	}
	return false
}

func getResultContent(result *mcp.CallToolResultFor[any]) string {
	if result == nil || len(result.Content) == 0 {
		return ""
	}
	
	// MCP results typically have text content
	if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
		return textContent.Text
	}
	
	return fmt.Sprintf("%v", result.Content[0])
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}