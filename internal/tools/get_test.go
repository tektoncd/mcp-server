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

func TestGetOperations(t *testing.T) {
	data := test.Data{
		Pipelines: []*v1.Pipeline{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipeline",
					Namespace: "default",
				},
				Spec: v1.PipelineSpec{
					Description: "Test pipeline for get operation",
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
					Description: "Test task for get operation",
				},
			},
		},
		PipelineRuns: []*v1.PipelineRun{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pipelinerun",
					Namespace: "default",
				},
				Spec: v1.PipelineRunSpec{
					PipelineRef: &v1.PipelineRef{
						Name: "test-pipeline",
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
				Spec: v1.TaskRunSpec{
					TaskRef: &v1.TaskRef{
						Name: "test-task",
					},
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
		contains []string
	}{
		{
			name: "get_pipeline_yaml",
			tool: "get_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"output":    "yaml",
			},
			contains: []string{"name: test-pipeline", "Test pipeline for get operation"},
		},
		{
			name: "get_pipeline_json",
			tool: "get_pipeline",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"output":    "json",
			},
			contains: []string{`"name": "test-pipeline"`, `"description": "Test pipeline for get operation"`},
		},
		{
			name: "get_task_yaml",
			tool: "get_task",
			args: map[string]interface{}{
				"name":      "test-task",
				"namespace": "default",
				"output":    "yaml",
			},
			contains: []string{"name: test-task", "Test task for get operation"},
		},
		{
			name: "get_pipelinerun_yaml",
			tool: "get_pipelinerun",
			args: map[string]interface{}{
				"name":      "test-pipelinerun",
				"namespace": "default",
				"output":    "yaml",
			},
			contains: []string{"name: test-pipelinerun", "pipelineRef:"},
		},
		{
			name: "get_taskrun_yaml",
			tool: "get_taskrun",
			args: map[string]interface{}{
				"name":      "test-taskrun",
				"namespace": "default",
				"output":    "yaml",
			},
			contains: []string{"name: test-taskrun", "taskRef:"},
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

			for _, expected := range test.contains {
				if !strings.Contains(content.Text, expected) {
					t.Fatalf("Expected response to contain '%s', got '%s'", expected, content.Text)
				}
			}
		})
	}
}

func TestGetOperationsErrors(t *testing.T) {
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
			name: "get_pipeline_not_found",
			tool: "get_pipeline",
			args: map[string]interface{}{
				"name":      "nonexistent",
				"namespace": "default",
			},
			expected: "Error getting Pipeline",
		},
		{
			name: "get_task_missing_name",
			tool: "get_task",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: Task name is required",
		},
		{
			name: "get_pipelinerun_not_found",
			tool: "get_pipelinerun",
			args: map[string]interface{}{
				"name":      "nonexistent",
				"namespace": "default",
			},
			expected: "Error getting PipelineRun",
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
