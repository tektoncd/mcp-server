package tools

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvertToParamValue_String(t *testing.T) {
	result := convertToParamValue("hello")

	if result.Type != v1.ParamTypeString {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeString, result.Type)
	}
	if result.StringVal != "hello" {
		t.Errorf("Expected StringVal 'hello', got %q", result.StringVal)
	}
}

func TestConvertToParamValue_Int(t *testing.T) {
	result := convertToParamValue(42)

	if result.Type != v1.ParamTypeString {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeString, result.Type)
	}
	if result.StringVal != "42" {
		t.Errorf("Expected StringVal '42', got %q", result.StringVal)
	}
}

func TestConvertToParamValue_Float(t *testing.T) {
	result := convertToParamValue(3.14)

	if result.Type != v1.ParamTypeString {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeString, result.Type)
	}
	if result.StringVal != "3.14" {
		t.Errorf("Expected StringVal '3.14', got %q", result.StringVal)
	}
}

func TestConvertToParamValue_Bool(t *testing.T) {
	result := convertToParamValue(true)

	if result.Type != v1.ParamTypeString {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeString, result.Type)
	}
	if result.StringVal != "true" {
		t.Errorf("Expected StringVal 'true', got %q", result.StringVal)
	}
}

func TestConvertToParamValue_StringSlice(t *testing.T) {
	input := []string{"a", "b", "c"}
	result := convertToParamValue(input)

	if result.Type != v1.ParamTypeArray {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeArray, result.Type)
	}
	if diff := cmp.Diff(input, result.ArrayVal); diff != "" {
		t.Errorf("ArrayVal mismatch (-want +got):\n%s", diff)
	}
}

func TestConvertToParamValue_InterfaceSlice(t *testing.T) {
	input := []interface{}{"a", 1, true}
	result := convertToParamValue(input)

	if result.Type != v1.ParamTypeArray {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeArray, result.Type)
	}
	expected := []string{"a", "1", "true"}
	if diff := cmp.Diff(expected, result.ArrayVal); diff != "" {
		t.Errorf("ArrayVal mismatch (-want +got):\n%s", diff)
	}
}

func TestConvertToParamValue_StringMap(t *testing.T) {
	input := map[string]string{"key1": "value1", "key2": "value2"}
	result := convertToParamValue(input)

	if result.Type != v1.ParamTypeObject {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeObject, result.Type)
	}
	if diff := cmp.Diff(input, result.ObjectVal); diff != "" {
		t.Errorf("ObjectVal mismatch (-want +got):\n%s", diff)
	}
}

func TestConvertToParamValue_InterfaceMap(t *testing.T) {
	input := map[string]interface{}{"str": "value", "num": 42, "bool": false}
	result := convertToParamValue(input)

	if result.Type != v1.ParamTypeObject {
		t.Errorf("Expected type %s, got %s", v1.ParamTypeObject, result.Type)
	}
	expected := map[string]string{"str": "value", "num": "42", "bool": "false"}
	if diff := cmp.Diff(expected, result.ObjectVal); diff != "" {
		t.Errorf("ObjectVal mismatch (-want +got):\n%s", diff)
	}
}

func TestConvertToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"other", struct{ Name string }{"test"}, "{test}"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := convertToString(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTriggerArtifactHubTask(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		expected string
	}{
		{
			name: "trigger with string param",
			args: map[string]interface{}{
				"name":      "test-task",
				"namespace": "default",
				"params": map[string]interface{}{
					"string-param": "value1",
				},
			},
			expected: "Successfully triggered TaskRun",
		},
		{
			name: "trigger with array param",
			args: map[string]interface{}{
				"name":      "test-task",
				"namespace": "default",
				"params": map[string]interface{}{
					"array-param": []interface{}{"a", "b", "c"},
				},
			},
			expected: "Successfully triggered TaskRun",
		},
		{
			name: "trigger without params",
			args: map[string]interface{}{
				"name":      "test-task",
				"namespace": "default",
			},
			expected: "Successfully triggered TaskRun",
		},
		{
			name: "trigger without name",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: name parameter is required",
		},
		{
			name: "trigger with default namespace",
			args: map[string]interface{}{
				"name": "test-task",
			},
			expected: "Successfully triggered TaskRun",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh context and session for each test to avoid "already exists" errors
			ctx, _ := ttesting.SetupFakeContext(t)

			data := test.Data{
				Tasks: []*v1.Task{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-task",
							Namespace: "default",
						},
						Spec: v1.TaskSpec{
							Params: []v1.ParamSpec{
								{Name: "string-param", Type: v1.ParamTypeString},
								{Name: "array-param", Type: v1.ParamTypeArray},
							},
							Steps: []v1.Step{
								{
									Name:   "echo",
									Image:  "alpine",
									Script: "echo hello",
								},
							},
						},
					},
				},
			}
			_, _ = test.SeedTestData(t, ctx, data)

			ss, cs := newSession(t, ctx)
			defer ss.Close()
			defer cs.Close()

			response, err := cs.CallTool(ctx, &mcp.CallToolParams{
				Name:      "trigger_artifacthub_task",
				Arguments: tc.args,
			})
			if err != nil {
				t.Fatal(err)
			}

			content, ok := response.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatal("Expected text content")
			}

			if !strings.Contains(content.Text, tc.expected) {
				t.Errorf("Expected response to contain %q, got %q", tc.expected, content.Text)
			}
		})
	}
}

func TestTriggerArtifactHubPipeline(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		expected string
	}{
		{
			name: "trigger with string param",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"params": map[string]interface{}{
					"string-param": "value1",
				},
			},
			expected: "Successfully triggered PipelineRun",
		},
		{
			name: "trigger with array param",
			args: map[string]interface{}{
				"name":      "test-pipeline",
				"namespace": "default",
				"params": map[string]interface{}{
					"array-param": []interface{}{"x", "y", "z"},
				},
			},
			expected: "Successfully triggered PipelineRun",
		},
		{
			name: "trigger without name",
			args: map[string]interface{}{
				"namespace": "default",
			},
			expected: "Error: name parameter is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh context and session for each test to avoid "already exists" errors
			ctx, _ := ttesting.SetupFakeContext(t)

			data := test.Data{
				Pipelines: []*v1.Pipeline{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-pipeline",
							Namespace: "default",
						},
						Spec: v1.PipelineSpec{
							Params: []v1.ParamSpec{
								{Name: "string-param", Type: v1.ParamTypeString},
								{Name: "array-param", Type: v1.ParamTypeArray},
							},
							Tasks: []v1.PipelineTask{
								{
									Name: "task1",
									TaskSpec: &v1.EmbeddedTask{
										TaskSpec: v1.TaskSpec{
											Steps: []v1.Step{
												{
													Name:   "echo",
													Image:  "alpine",
													Script: "echo hello",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			_, _ = test.SeedTestData(t, ctx, data)

			ss, cs := newSession(t, ctx)
			defer ss.Close()
			defer cs.Close()

			response, err := cs.CallTool(ctx, &mcp.CallToolParams{
				Name:      "trigger_artifacthub_pipeline",
				Arguments: tc.args,
			})
			if err != nil {
				t.Fatal(err)
			}

			content, ok := response.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatal("Expected text content")
			}

			if !strings.Contains(content.Text, tc.expected) {
				t.Errorf("Expected response to contain %q, got %q", tc.expected, content.Text)
			}
		})
	}
}

func TestListArtifactHubTasks(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	// Note: This test verifies the tool exists and can be called.
	// Actual HTTP calls to Artifact Hub are not made in unit tests.
	// The client tests cover the HTTP behavior with mocks.

	// Verify the tool is registered
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, tool := range tools.Tools {
		if tool.Name == "list_artifacthub_tasks" {
			found = true
			break
		}
	}

	if !found {
		t.Error("list_artifacthub_tasks tool not found")
	}
}

func TestListArtifactHubPipelines(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	// Verify the tool is registered
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, tool := range tools.Tools {
		if tool.Name == "list_artifacthub_pipelines" {
			found = true
			break
		}
	}

	if !found {
		t.Error("list_artifacthub_pipelines tool not found")
	}
}

func TestInstallArtifactHubTask(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	// Verify the tool is registered
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, tool := range tools.Tools {
		if tool.Name == "install_artifacthub_task" {
			found = true
			break
		}
	}

	if !found {
		t.Error("install_artifacthub_task tool not found")
	}
}

func TestInstallArtifactHubPipeline(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	// Verify the tool is registered
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, tool := range tools.Tools {
		if tool.Name == "install_artifacthub_pipeline" {
			found = true
			break
		}
	}

	if !found {
		t.Error("install_artifacthub_pipeline tool not found")
	}
}

func TestInstallArtifactHubTask_MissingPackageID(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	response, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "install_artifacthub_task",
		Arguments: map[string]interface{}{
			"namespace": "default",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	content, ok := response.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Expected text content")
	}

	if !strings.Contains(content.Text, "Error: packageId parameter is required") {
		t.Errorf("Expected error about missing packageId, got %q", content.Text)
	}
}

func TestInstallArtifactHubPipeline_MissingPackageID(t *testing.T) {
	ctx, _ := ttesting.SetupFakeContext(t)
	data := test.Data{}
	_, _ = test.SeedTestData(t, ctx, data)

	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	response, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "install_artifacthub_pipeline",
		Arguments: map[string]interface{}{
			"namespace": "default",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	content, ok := response.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Expected text content")
	}

	if !strings.Contains(content.Text, "Error: packageId parameter is required") {
		t.Errorf("Expected error about missing packageId, got %q", content.Text)
	}
}
