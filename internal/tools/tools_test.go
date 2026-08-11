package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tektoncd/mcp-server/internal/resources"
	"github.com/tektoncd/mcp-server/internal/version"
)

func newSession(t *testing.T, ctx context.Context) (*mcp.ServerSession, *mcp.ClientSession) {
	t.Helper()

	ct, st := mcp.NewInMemoryTransports()
	s := mcp.NewServer(&mcp.Implementation{Name: "Tekton", Version: version.Version}, nil)
	if err := Add(ctx, s); err != nil {
		t.Fatal(err)
	}
	resources.Add(ctx, s)
	c := mcp.NewClient(&mcp.Implementation{Name: "TektonClient", Version: version.Version}, nil)

	ss, err := s.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := c.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}

	return ss, cs
}

func TestToolRegistration(t *testing.T) {
	ss, cs := newSession(t, context.Background())
	defer ss.Close()
	defer cs.Close()

	expected := map[string]bool{
		"start_pipeline": true, "start_task": true,
		"restart_pipelinerun": true, "restart_taskrun": true,
		"get_taskrun_logs":  true,
		"list_pipelineruns": true, "list_pipelines": true, "list_taskruns": true,
		"list_tasks": true, "list_stepactions": true,
		"create_pipeline": true, "create_task": true,
		"create_pipelinerun": true, "create_taskrun": true,
		"get_pipeline": true, "get_task": true,
		"get_pipelinerun": true, "get_taskrun": true,
		"update_pipeline": true, "update_task": true, "patch_pipeline": true,
		"delete_pipeline": true, "delete_task": true,
		"delete_pipelinerun": true, "delete_taskrun": true,
		"delete_all_pipelineruns": true,
		"list_artifacthub_tasks":  true, "list_artifacthub_pipelines": true,
		"install_artifacthub_task": true, "install_artifacthub_pipeline": true,
		"trigger_artifacthub_task": true, "trigger_artifacthub_pipeline": true,
	}

	listed, err := cs.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range listed.Tools {
		if !expected[tool.Name] {
			t.Errorf("unexpected tool %q", tool.Name)
		}
		delete(expected, tool.Name)

		schema := tool.InputSchema.(map[string]any)
		required, _ := schema["required"].([]any)
		for _, field := range required {
			if field == "namespace" {
				t.Errorf("tool %q requires its defaulted namespace", tool.Name)
			}
		}
	}
	for name := range expected {
		t.Errorf("missing tool %q", name)
	}
}
