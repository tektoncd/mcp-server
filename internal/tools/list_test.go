package tools

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHandlerListStepActions(t *testing.T) {
	data := test.Data{
		StepActions: []*v1beta1.StepAction{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo-sa1", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar-sa2", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	tests := []struct {
		name          string
		request       mcp.CallToolRequest
		expectedNames []string
	}{
		{
			name:          "No namespace, no filters",
			request:       newCallToolRequest(map[string]any{}),
			expectedNames: []string{"foo-sa1", "bar-sa2"},
		},
		{
			name:          "With namespace",
			request:       newCallToolRequest(map[string]any{"namespace": "ns1"}),
			expectedNames: []string{"foo-sa1"},
		},
		{
			name:          "With label selector",
			request:       newCallToolRequest(map[string]any{"label-selector": "env=dev"}),
			expectedNames: []string{"bar-sa2"},
		},
		{
			name:          "With prefix filter",
			request:       newCallToolRequest(map[string]any{"prefix": "foo"}),
			expectedNames: []string{"foo-sa1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := handlerListStepaction(ctx, test.request)
			if err != nil {
				t.Fatal(err)
			}

			actual := res.Content[0].(mcp.TextContent)
			var stepactions []*v1beta1.StepAction
			if err = json.Unmarshal([]byte(actual.Text), &stepactions); err != nil {
				t.Fatalf("failed to unmarshal stepactions: %v", err)
			}
			for _, stepaction := range stepactions {
				if !slices.Contains(test.expectedNames, stepaction.GetName()) {
					t.Fatalf("response contained unexpected result: %v", stepaction)
				}
			}
		})
	}
}

func TestHandlerListTaskRuns(t *testing.T) {
	data := test.Data{
		TaskRuns: []*v1.TaskRun{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo-tr1", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar-tr2", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	tests := []struct {
		name          string
		request       mcp.CallToolRequest
		expectedNames []string
	}{
		{
			name:          "No namespace, no filters",
			request:       newCallToolRequest(map[string]any{}),
			expectedNames: []string{"foo-tr1", "bar-tr2"},
		},
		{
			name:          "With namespace",
			request:       newCallToolRequest(map[string]any{"namespace": "ns1"}),
			expectedNames: []string{"foo-tr1"},
		},
		{
			name:          "With label selector",
			request:       newCallToolRequest(map[string]any{"label-selector": "env=dev"}),
			expectedNames: []string{"bar-tr2"},
		},
		{
			name:          "With prefix filter",
			request:       newCallToolRequest(map[string]any{"prefix": "foo"}),
			expectedNames: []string{"foo-tr1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := handlerListTaskRun(ctx, test.request)
			if err != nil {
				t.Fatal(err)
			}

			actual := res.Content[0].(mcp.TextContent)
			var stepactions []*v1.TaskRun
			if err = json.Unmarshal([]byte(actual.Text), &stepactions); err != nil {
				t.Fatalf("failed to unmarshal stepactions: %v", err)
			}
			for _, stepaction := range stepactions {
				if !slices.Contains(test.expectedNames, stepaction.GetName()) {
					t.Fatalf("response contained unexpected result: %v", stepaction)
				}
			}
		})
	}
}

func TestHandlerListTasks(t *testing.T) {
	data := test.Data{
		Tasks: []*v1.Task{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo-task1", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar-task2", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	tests := []struct {
		name          string
		request       mcp.CallToolRequest
		expectedNames []string
	}{
		{
			name:          "No namespace, no filters",
			request:       newCallToolRequest(map[string]any{}),
			expectedNames: []string{"foo-task1", "bar-task2"},
		},
		{
			name:          "With namespace",
			request:       newCallToolRequest(map[string]any{"namespace": "ns1"}),
			expectedNames: []string{"foo-task1"},
		},
		{
			name:          "With label selector",
			request:       newCallToolRequest(map[string]any{"label-selector": "env=dev"}),
			expectedNames: []string{"bar-task2"},
		},
		{
			name:          "With prefix filter",
			request:       newCallToolRequest(map[string]any{"prefix": "foo"}),
			expectedNames: []string{"foo-task1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := handlerListTask(ctx, test.request)
			if err != nil {
				t.Fatal(err)
			}

			actual := res.Content[0].(mcp.TextContent)
			var tasks []*v1.Task
			if err = json.Unmarshal([]byte(actual.Text), &tasks); err != nil {
				t.Fatalf("failed to unmarshal tasks: %v", err)
			}
			for _, task := range tasks {
				if !slices.Contains(test.expectedNames, task.GetName()) {
					t.Fatalf("response contained unexpected result: %v", task)
				}
			}
		})
	}
}

func TestHandlerListPipelines(t *testing.T) {
	data := test.Data{
		Pipelines: []*v1.Pipeline{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo-pipeline1", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar-pipeline2", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	tests := []struct {
		name          string
		request       mcp.CallToolRequest
		expectedNames []string
	}{
		{
			name:          "No namespace, no filters",
			request:       newCallToolRequest(map[string]any{}),
			expectedNames: []string{"foo-pipeline1", "bar-pipeline2"},
		},
		{
			name:          "With namespace",
			request:       newCallToolRequest(map[string]any{"namespace": "ns1"}),
			expectedNames: []string{"foo-pipeline1"},
		},
		{
			name:          "With label selector",
			request:       newCallToolRequest(map[string]any{"label-selector": "env=dev"}),
			expectedNames: []string{"bar-pipeline2"},
		},
		{
			name:          "With prefix filter",
			request:       newCallToolRequest(map[string]any{"prefix": "foo"}),
			expectedNames: []string{"foo-pipeline1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := handlerListPipeline(ctx, test.request)
			if err != nil {
				t.Fatal(err)
			}

			actual := res.Content[0].(mcp.TextContent)
			var pipelines []*v1.Pipeline
			if err = json.Unmarshal([]byte(actual.Text), &pipelines); err != nil {
				t.Fatalf("failed to unmarshal pipelines: %v", err)
			}
			for _, pipeline := range pipelines {
				if !slices.Contains(test.expectedNames, pipeline.GetName()) {
					t.Fatalf("response contained unexpected result: %v", pipeline)
				}
			}
		})
	}
}

func TestHandlerListPipelineRuns(t *testing.T) {
	data := test.Data{
		PipelineRuns: []*v1.PipelineRun{
			{ObjectMeta: metav1.ObjectMeta{Name: "foo-pr1", Namespace: "ns1", Labels: map[string]string{"env": "prod"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "bar-pr2", Namespace: "ns2", Labels: map[string]string{"env": "dev"}}},
		},
	}

	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)

	tests := []struct {
		name          string
		request       mcp.CallToolRequest
		expectedNames []string
	}{
		{
			name:          "No namespace, no filters",
			request:       newCallToolRequest(map[string]any{}),
			expectedNames: []string{"foo-pr1", "bar-pr2"},
		},
		{
			name:          "With namespace",
			request:       newCallToolRequest(map[string]any{"namespace": "ns1"}),
			expectedNames: []string{"foo-pr1"},
		},
		{
			name:          "With label selector",
			request:       newCallToolRequest(map[string]any{"label-selector": "env=dev"}),
			expectedNames: []string{"bar-pr2"},
		},
		{
			name:          "With prefix filter",
			request:       newCallToolRequest(map[string]any{"prefix": "foo"}),
			expectedNames: []string{"foo-pr1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := handlerListPipelineRun(ctx, test.request)
			if err != nil {
				t.Fatal(err)
			}

			actual := res.Content[0].(mcp.TextContent)
			var prs []*v1.PipelineRun
			if err = json.Unmarshal([]byte(actual.Text), &prs); err != nil {
				t.Fatalf("failed to unmarshal pipelineruns: %v", err)
			}
			for _, pr := range prs {
				if !slices.Contains(test.expectedNames, pr.GetName()) {
					t.Fatalf("response contained unexpected result: %v", pr)
				}
			}
		})
	}
}
