package resources

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestPipelineRunResourcesTemplate(t *testing.T) {
	ctx := context.Background()
	template, _ := pipelineRunResources(ctx)

	// Convert to JSON to check the template string value
	jsonData, err := json.Marshal(template)
	assert.NoError(t, err)

	// Verify JSON contains the expected template string
	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, "tekton://pipelinerun/{namespace}/{name}")
	assert.Equal(t, "PipelineRun", template.Name)
}

func TestPipelineRunResourcesHandlerMissingNamespace(t *testing.T) {
	ctx := context.Background()
	handler := pipelineRunHandler(ctx)

	request := mcp.ReadResourceRequest{
		Params: struct {
			URI       string                 `json:"uri"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}{
			URI: "tekton://pipelinerun/",
			Arguments: map[string]interface{}{
				"name": []string{"test-pipeline-run"},
			},
		},
	}

	_, err := handler(ctx, request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")
}

func TestPipelineRunResourcesHandlerMissingName(t *testing.T) {
	ctx := context.Background()
	handler := pipelineRunHandler(ctx)

	request := mcp.ReadResourceRequest{
		Params: struct {
			URI       string                 `json:"uri"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}{
			URI: "tekton://pipelinerun/",
			Arguments: map[string]interface{}{
				"namespace": []string{"test-namespace"},
			},
		},
	}

	_, err := handler(ctx, request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}
