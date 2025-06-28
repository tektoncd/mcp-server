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

	request := &mcp.CallToolParamsFor[startParams]{
		Arguments: startParams{Name: "hello-world", Namespace: "default"},
	}
	response, err := handlerStartPipeline(ctx, nil, request)
	if err != nil {
		t.Fatal(err)
	}
	content, _ := response.Content[0].(*mcp.TextContent)
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

	request := &mcp.CallToolParamsFor[startParams]{
		Arguments: startParams{Name: "hello-world", Namespace: "default"},
	}
	response, err := handlerStartTask(ctx, nil, request)
	if err != nil {
		t.Fatal(err)
	}
	content, _ := response.Content[0].(*mcp.TextContent)
	if !strings.HasPrefix(content.Text, "Starting task") {
		t.Fatalf("invalid resposne: %v", content)
	}
}
