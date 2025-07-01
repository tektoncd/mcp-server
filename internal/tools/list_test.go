package tools

import (
	"encoding/json"
	"slices"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestList(t *testing.T) {
	data := test.Data{
		StepActions: []*v1beta1.StepAction{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
		Tasks: []*v1.Task{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
		TaskRuns: []*v1.TaskRun{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
		Pipelines: []*v1.Pipeline{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
		PipelineRuns: []*v1.PipelineRun{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	tools := []string{
		"list_stepactions",
		"list_tasks",
		"list_taskruns",
		"list_pipelines",
		"list_pipelineruns",
	}

	tests := []struct {
		name      string
		arguments map[string]string
		expected  []string
	}{
		{
			name:      "No namespace, no filters",
			arguments: map[string]string{},
			expected:  []string{"foo", "bar"},
		},
		{
			name:      "With namespace",
			arguments: map[string]string{"namespace": "ns1"},
			expected:  []string{"foo"},
		},
		{
			name:      "With label selector",
			arguments: map[string]string{"labelSelector": "env=dev"},
			expected:  []string{"bar"},
		},
		{
			name:      "With prefix filter",
			arguments: map[string]string{"prefix": "foo"},
			expected:  []string{"foo"},
		},
	}

	for _, tool := range tools {
		t.Run(tool, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					response, err := cs.CallTool(t.Context(), &mcp.CallToolParams{
						Name:      tool,
						Arguments: test.arguments,
					})
					if err != nil {
						t.Fatal(err)
					}

					actual := response.Content[0].(*mcp.TextContent)
					var objects []map[string]any
					if err = json.Unmarshal([]byte(actual.Text), &objects); err != nil {
						t.Fatalf("failed to unmarshal objects: %v", err)
					}
					for _, object := range objects {
						metadata := object["metadata"].(map[string]any)
						if !slices.Contains(test.expected, metadata["name"].(string)) {
							t.Fatalf("response contained unexpected result: %v", object)
						}
					}
				})
			}
		})
	}
}
