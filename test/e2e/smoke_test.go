//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestTools(t *testing.T) {
	endpoint := os.Getenv("MCP_SERVER_URL")
	namespace := os.Getenv("E2E_NAMESPACE")
	if endpoint == "" || namespace == "" {
		t.Fatal("MCP_SERVER_URL and E2E_NAMESPACE are required")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	session := connect(t, ctx, endpoint)
	defer session.Close()

	name := fmt.Sprintf("mcp-smoke-%d", time.Now().UnixNano())
	created, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "create_task",
		Arguments: map[string]any{
			"namespace": namespace,
			"yaml": fmt.Sprintf(`apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: %s
spec:
  steps:
    - name: noop
      image: busybox:1.36.1
      command: ["true"]
`, name),
		},
	})
	if err != nil {
		t.Fatalf("create_task: %v", err)
	}
	if text := resultText(t, created); created.IsError || !strings.Contains(text, "created successfully") {
		t.Fatalf("create_task: %s", text)
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		listed, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: "list_tasks",
			Arguments: map[string]any{
				"namespace": namespace,
				"prefix":    name,
			},
		})
		if err == nil && !listed.IsError && containsTask(resultText(t, listed), name) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("list_tasks did not return %q: %v", name, err)
		}
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(time.Second):
		}
	}
}

func connect(t *testing.T, ctx context.Context, endpoint string) *mcp.ClientSession {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		client := mcp.NewClient("tekton-mcp-e2e", "dev", nil)
		session, err := client.Connect(ctx, mcp.NewStreamableClientTransport(endpoint, nil))
		if err == nil {
			return session
		}
		if time.Now().After(deadline) {
			t.Fatalf("connect to %s: %v", endpoint, err)
		}
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(time.Second):
		}
	}
}

func resultText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) != 1 {
		t.Fatalf("unexpected tool result: %#v", result.Content)
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("unexpected tool result type: %T", result.Content[0])
	}
	return text.Text
}

func containsTask(text, name string) bool {
	var tasks []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	}
	if json.Unmarshal([]byte(text), &tasks) != nil {
		return false
	}
	for _, task := range tasks {
		if task.Metadata.Name == name {
			return true
		}
	}
	return false
}
