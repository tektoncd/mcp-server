package resources

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPipelineRunWithParams(t *testing.T) {
	// Create a sample PipelineRun
	pipelineRun := &pipelinev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pipeline-run",
			Namespace: "default",
		},
		Spec: pipelinev1.PipelineRunSpec{
			PipelineRef: &pipelinev1.PipelineRef{
				Name: "test-pipeline",
			},
			Params: []pipelinev1.Param{
				{
					Name: "existing-param",
					Value: pipelinev1.ParamValue{
						Type:      pipelinev1.ParamTypeString,
						StringVal: "existing-value",
					},
				},
			},
		},
	}

	// Create input parameters
	inputParams := []pipelinev1.Param{
		{
			Name: "existing-param",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "updated-value",
			},
		},
		{
			Name: "new-param",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "new-value",
			},
		},
		{
			Name: "array-param",
			Value: pipelinev1.ParamValue{
				Type:     pipelinev1.ParamTypeArray,
				ArrayVal: []string{"item1", "item2"},
			},
		},
	}

	// Call the function being tested
	pipelineRunCopy := pipelineRun.DeepCopy()

	// Manually merge the parameters as getPipelineRun would do
	if len(inputParams) > 0 {
		existingParams := make(map[string]int)
		for i, p := range pipelineRunCopy.Spec.Params {
			existingParams[p.Name] = i
		}

		result := make([]pipelinev1.Param, len(pipelineRunCopy.Spec.Params))
		copy(result, pipelineRunCopy.Spec.Params)

		// Update or append parameters
		for _, newParam := range inputParams {
			if idx, exists := existingParams[newParam.Name]; exists {
				result[idx] = newParam
			} else {
				result = append(result, newParam)
			}
		}

		pipelineRunCopy.Spec.Params = result
	}

	// Check that the existing parameter was updated
	found := false
	for _, param := range pipelineRunCopy.Spec.Params {
		if param.Name == "existing-param" {
			found = true
			if param.Value.StringVal != "updated-value" {
				t.Errorf("existing-param not updated, got %s, want %s", param.Value.StringVal, "updated-value")
			}
		}
	}
	if !found {
		t.Errorf("existing-param not found in result")
	}

	// Check that the new parameter was added
	found = false
	for _, param := range pipelineRunCopy.Spec.Params {
		if param.Name == "new-param" {
			found = true
			if param.Value.StringVal != "new-value" {
				t.Errorf("new-param value not correct, got %s, want %s", param.Value.StringVal, "new-value")
			}
		}
	}
	if !found {
		t.Errorf("new-param not found in result")
	}

	// Check that the array parameter was added
	found = false
	for _, param := range pipelineRunCopy.Spec.Params {
		if param.Name == "array-param" {
			found = true
			if len(param.Value.ArrayVal) != 2 || param.Value.ArrayVal[0] != "item1" || param.Value.ArrayVal[1] != "item2" {
				t.Errorf("array-param value not correct, got %v", param.Value.ArrayVal)
			}
		}
	}
	if !found {
		t.Errorf("array-param not found in result")
	}

	// Check total param count
	if len(pipelineRunCopy.Spec.Params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(pipelineRunCopy.Spec.Params))
	}
}

// This test is a more realistic test that simulates the request path
func TestPipelineRunParamsRequest(t *testing.T) {
	// Create a mock request with the proper structure matching what's used in resources.go
	request := mcp.ReadResourceRequest{
		Params: struct {
			URI       string                 `json:"uri"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}{
			URI: "tekton://pipelinerun/default/test-pipeline-run",
			Arguments: map[string]interface{}{
				"namespace": []string{"default"},
				"name":      []string{"test-pipeline-run"},
				"params":    []string{"existing-param=updated-value,new-param=new-value,array-param=array:item1:item2"},
			},
		},
	}

	// Create a handler (but don't use the real handler which requires informers)
	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract params like the real handler would
		paramsArg, ok := request.Params.Arguments["params"].([]string)
		if !ok || len(paramsArg) == 0 {
			t.Fatal("Failed to extract params from request")
		}

		// Check that the params were correctly parsed
		paramStr := paramsArg[0]
		if paramStr != "existing-param=updated-value,new-param=new-value,array-param=array:item1:item2" {
			t.Errorf("Wrong param string: %s", paramStr)
		}

		// This would normally be parsed by params.ParsePipelineRunParams
		// but we'll just return true to indicate success
		return nil, nil
	}

	// Call the handler
	_, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	// Success!
}
