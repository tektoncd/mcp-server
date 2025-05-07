package tools

import (
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHandlerStartPipeline(t *testing.T) {
	data := test.Data{
		Pipelines: []*v1.Pipeline{
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

	request := newCallToolRequest(map[string]any{"name": "hello-world", "namespace": "default"})
	response, err := handlerStartPipeline(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	content, _ := mcp.AsTextContent(response.Content[0])
	if !strings.HasPrefix(content.Text, "Starting pipeline") {
		t.Fatalf("invalid resposne: %v", content)
	}
}

func TestHandlerStartTask(t *testing.T) {
	data := test.Data{
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

	request := newCallToolRequest(map[string]any{"name": "hello-world", "namespace": "default"})
	response, err := handlerStartTask(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	content, _ := mcp.AsTextContent(response.Content[0])
	if !strings.HasPrefix(content.Text, "Starting task") {
		t.Fatalf("invalid resposne: %v", content)
	}
}
