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

func TestDeleteOperations(t *testing.T) {
	data := test.Data{
		Pipelines: []*v1.Pipeline{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipeline",
					Namespace: "default",
				},
			},
		},
		Tasks: []*v1.Task{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-task",
					Namespace: "default",
				},
			},
		},
		PipelineRuns: []*v1.PipelineRun{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipelinerun",
					Namespace: "default",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipelinerun-2",
					Namespace: "default",
					Labels: map[string]string{
						"app": "test",
					},
				},
			},
		},
		TaskRuns: []*v1.TaskRun{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-taskrun",
					Namespace: "default",
				},
			},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
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
			name: "delete_pipeline",
			tool: "delete_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
			},
			expected: "Pipeline 'test-pipeline' deleted successfully",
		},
		{
			name: "delete_task",
			tool: "delete_task",
			args: map[string]interface{}{
				"name":      "test-task",
				"namespace": "default",
			},
			expected: "Task 'test-task' deleted successfully",
		},
		{
			name: "delete_pipelinerun",
			tool: "delete_pipelinerun",
			args: map[string]interface{}{
				"name":      "test-pipelinerun",
				"namespace": "default",
			},
			expected: "PipelineRun 'test-pipelinerun' deleted successfully",
		},
		{
			name: "delete_taskrun",
			tool: "delete_taskrun",
			args: map[string]interface{}{
				"name":      "test-taskrun",
				"namespace": "default",
			},
			expected: "TaskRun 'test-taskrun' deleted successfully",
		},
		{
			name: "delete_all_pipelineruns_with_label",
			tool: "delete_all_pipelineruns",
			args: map[string]interface{}{
				"namespace":     "default",
				"labelSelector": "app=test",
			},
			expected: "PipelineRuns deleted successfully from namespace 'default' with selectors",
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

func TestDeleteOperationsErrors(t *testing.T) {
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
			name: "delete_pipeline_not_found",
			tool: "delete_pipeline",
			args: map[string]interface{}{
				"name":      "nonexistent",
				"namespace": "default",
			},
			expected: "Error deleting Pipeline",
		},
		{
			name: "delete_pipeline_missing_name",
			tool: "delete_pipeline",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: Pipeline name is required",
		},
		{
			name: "delete_task_not_found",
			tool: "delete_task",
			args: map[string]interface{}{
				"name":      "nonexistent",
				"namespace": "default",
			},
			expected: "Error deleting Task",
		},
		{
			name: "delete_pipelinerun_missing_name",
			tool: "delete_pipelinerun",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: PipelineRun name is required",
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
