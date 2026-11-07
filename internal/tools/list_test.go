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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
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

func TestListSummary(t *testing.T) {
	start := metav1.Now()
	data := test.Data{
		PipelineRuns: []*v1.PipelineRun{{
			ObjectMeta: metav1.ObjectMeta{Name: "pr", Namespace: "ns1"},
			Spec:       v1.PipelineRunSpec{PipelineRef: &v1.PipelineRef{Name: "build"}},
			Status: v1.PipelineRunStatus{
				Status: duckv1.Status{Conditions: duckv1.Conditions{{
					Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse, Reason: "Failed",
				}}},
				PipelineRunStatusFields: v1.PipelineRunStatusFields{
					StartTime:    &start,
					PipelineSpec: &v1.PipelineSpec{Description: "large spec that summary drops"},
				},
			},
		}},
		TaskRuns: []*v1.TaskRun{{
			ObjectMeta: metav1.ObjectMeta{Name: "pr-clone", Namespace: "ns1", Labels: map[string]string{
				"tekton.dev/pipelineRun": "pr", "tekton.dev/pipelineTask": "clone",
			}},
			Spec: v1.TaskRunSpec{TaskRef: &v1.TaskRef{Name: "git-clone"}},
			Status: v1.TaskRunStatus{Status: duckv1.Status{Conditions: duckv1.Conditions{{
				Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue, Reason: "Succeeded",
			}}}},
		}},
	}
	ctx, _ := ttesting.SetupFakeContext(t)
	_, _ = test.SeedTestData(t, ctx, data)
	ss, cs := newSession(t, ctx)
	defer ss.Close()
	defer cs.Close()

	call := func(tool, output string) *mcp.CallToolResult {
		t.Helper()
		res, err := cs.CallTool(t.Context(), &mcp.CallToolParams{
			Name:      tool,
			Arguments: map[string]string{"namespace": "ns1", "output": output},
		})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	decode := func(res *mcp.CallToolResult) []map[string]any {
		t.Helper()
		var out []map[string]any
		if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &out); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(out) != 1 {
			t.Fatalf("expected 1 item, got %d", len(out))
		}
		return out
	}

	pr := decode(call("list_pipelineruns", "summary"))[0]
	want := map[string]string{"name": "pr", "namespace": "ns1", "ref": "build", "status": "False", "reason": "Failed"}
	for k, v := range want {
		if pr[k] != v {
			t.Errorf("pipelinerun summary %s = %v, want %q", k, pr[k], v)
		}
	}
	if pr["startTime"] == nil {
		t.Error("pipelinerun summary has no startTime")
	}
	for _, k := range []string{"metadata", "spec"} {
		if _, ok := pr[k]; ok {
			t.Errorf("pipelinerun summary contains full field %q", k)
		}
	}

	tr := decode(call("list_taskruns", "summary"))[0]
	want = map[string]string{"name": "pr-clone", "ref": "git-clone", "pipelineRun": "pr", "pipelineTask": "clone", "status": "True", "reason": "Succeeded"}
	for k, v := range want {
		if tr[k] != v {
			t.Errorf("taskrun summary %s = %v, want %q", k, tr[k], v)
		}
	}

	if full := decode(call("list_pipelineruns", "full"))[0]; full["spec"] == nil {
		t.Error("full output lost the spec")
	}

	if res := call("list_pipelineruns", "yaml"); !res.IsError {
		t.Error("invalid output value was accepted")
	}
}
