package tools

import (
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpdateOperations(t *testing.T) {
	data := test.Data{
		Pipelines: []*v1.Pipeline{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipeline",
					Namespace: "default",
				},
				Spec: v1.PipelineSpec{
					Description: "Original description",
				},
			},
		},
		Tasks: []*v1.Task{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-task",
					Namespace: "default",
				},
				Spec: v1.TaskSpec{
					Description: "Original task description",
				},
			},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	updatedPipelineYAML := `
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: test-pipeline
spec:
  description: Updated pipeline description
  tasks:
  - name: updated-task
    taskSpec:
      steps:
      - image: alpine
        script: echo updated`

	updatedTaskYAML := `
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: test-task
spec:
  description: Updated task description
  steps:
  - image: alpine
    script: echo updated`

	tests := []struct {
		name     string
		tool     string
		args     map[string]interface{}
		expected string
	}{
		{
			name: "update_pipeline",
			tool: "update_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"yaml":      updatedPipelineYAML,
			},
			expected: "Pipeline 'test-pipeline' updated successfully",
		},
		{
			name: "update_task",
			tool: "update_task",
			args: map[string]interface{}{
				"name":      "test-task",
				"namespace": "default",
				"yaml":      updatedTaskYAML,
			},
			expected: "Task 'test-task' updated successfully",
		},
		{
			name: "patch_pipeline",
			tool: "patch_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"patch":     `[{"op": "replace", "path": "/spec/description", "value": "Patched description"}]`,
			},
			expected: "Pipeline 'test-pipeline' patched successfully",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := cs.CallTool(ctx, &mcp.CallToolParams{
				Name:      test.tool,
				Arguments: test.args,
			})
			if err != nil {
				t.Fatal(err)
			}

			content, ok := response.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatal("Expected text content")
			}

			if !strings.Contains(content.Text, test.expected) {
				t.Fatalf("Expected response to contain '%s', got '%s'", test.expected, content.Text)
			}
		})
	}
}

func TestUpdateOperationsErrors(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	tests := []struct {
		name     string
		tool     string
		args     map[string]interface{}
		expected string
	}{
		{
			name: "update_pipeline_not_found",
			tool: "update_pipeline",
			args: map[string]interface{}{
				"name":      "nonexistent",
				"namespace": "default",
				"yaml":      "apiVersion: tekton.dev/v1\nkind: Pipeline\nmetadata:\n  name: test",
			},
			expected: "Error getting existing Pipeline",
		},
		{
			name: "update_pipeline_invalid_yaml",
			tool: "update_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"yaml":      "invalid yaml content",
			},
			expected: "Error parsing YAML",
		},
		{
			name: "update_pipeline_missing_yaml",
			tool: "update_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
			},
			expected: "Error: Name and YAML definition are required",
		},
		{
			name: "patch_pipeline_invalid_patch",
			tool: "patch_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"patch":     "invalid json patch",
			},
			expected: "Error patching Pipeline",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := cs.CallTool(ctx, &mcp.CallToolParams{
				Name:      test.tool,
				Arguments: test.args,
			})
			if err != nil {
				t.Fatal(err)
			}

			content, ok := response.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatal("Expected text content")
			}

			if !strings.Contains(content.Text, test.expected) {
				t.Fatalf("Expected error message to contain '%s', got '%s'", test.expected, content.Text)
			}
		})
	}
}
