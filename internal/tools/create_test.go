package tools

import (
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
)

func TestCreateOperations(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)

	// No need to seed data for create operations
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	pipelineYAML := `
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: test-pipeline
spec:
  tasks:
  - name: hello
    taskSpec:
      steps:
      - image: alpine
        script: echo hello`

	taskYAML := `
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: test-task
spec:
  steps:
  - image: alpine
    script: echo hello`

	pipelineRunYAML := `
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: test-pipelinerun
spec:
  pipelineRef:
    name: test-pipeline`

	taskRunYAML := `
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  name: test-taskrun
spec:
  taskRef:
    name: test-task`

	tests := []struct {
		name     string
		tool     string
		args     map[string]interface{}
		expected string
	}{
		{
			name: "create_pipeline",
			tool: "create_pipeline",
			args: map[string]interface{}{
				"namespace": "default",
				"yaml":      pipelineYAML,
			},
			expected: "Pipeline 'test-pipeline' created successfully",
		},
		{
			name: "create_task",
			tool: "create_task",
			args: map[string]interface{}{
				"namespace": "default",
				"yaml":      taskYAML,
			},
			expected: "Task 'test-task' created successfully",
		},
		{
			name: "create_pipelinerun_with_yaml",
			tool: "create_pipelinerun",
			args: map[string]interface{}{
				"namespace": "default",
				"yaml":      pipelineRunYAML,
			},
			expected: "PipelineRun 'test-pipelinerun' created successfully",
		},
		{
			name: "create_taskrun_with_yaml",
			tool: "create_taskrun",
			args: map[string]interface{}{
				"namespace": "default",
				"yaml":      taskRunYAML,
			},
			expected: "TaskRun 'test-taskrun' created successfully",
		},
		{
			name: "create_pipelinerun_with_generateName",
			tool: "create_pipelinerun",
			args: map[string]interface{}{
				"namespace":    "default",
				"generateName": "generated-pr-",
			},
			expected: "PipelineRun '' created successfully", // Empty name since generateName doesn't work in test
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

func TestCreateOperationsErrors(t *testing.T) {
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
			name: "create_pipeline_invalid_yaml",
			tool: "create_pipeline",
			args: map[string]interface{}{
				"namespace": "default",
				"yaml":      "invalid yaml content",
			},
			expected: "Error parsing YAML",
		},
		{
			name: "create_pipeline_missing_yaml",
			tool: "create_pipeline",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: YAML definition is required",
		},
		{
			name: "create_pipelinerun_missing_both",
			tool: "create_pipelinerun",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: Either YAML definition or generateName is required",
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
