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

func TestStart(t *testing.T) {
	data := test.Data{
		Pipelines: []*v1.Pipeline{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
					Namespace: "default",
				},
			},
		},
		Tasks: []*v1.Task{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-world",
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
		tool     string
		response string
	}{
		{
			tool:     "start_pipeline",
			response: "Starting pipeline",
		},
		{
			tool:     "start_task",
			response: "Starting task",
		},
	}

	for _, test := range tests {
		t.Run(test.tool, func(t *testing.T) {
			response, err := cs.CallTool(ctx, &mcp.CallToolParams{
				Name:      test.tool,
				Arguments: map[string]string{"name": "hello-world", "namespace": "default"},
			})
			if err != nil {
				t.Fatal(err)
			}
			// TODO: check response.IsError
			content, _ := response.Content[0].(*mcp.TextContent)
			if !strings.HasPrefix(content.Text, test.response) {
				t.Fatalf("invalid response: %v", content)
			}
		})
	}
}
