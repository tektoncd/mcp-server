package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createTask() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewToolWithRawSchema(
			"create_task",
			"Create a task",
			json.RawMessage(getSchemaForType("github.com/tektoncd/pipeline/pkg/apis/pipeline/v1.Task"))),
		Handler: handlerCreateTask,
	}
}

func handlerCreateTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(request.Params.Arguments["task"])
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal request params", err), nil
	}

	var task v1.Task
	if err = json.Unmarshal(data, &task); err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to unmarshal request params", err), nil
	}

	pipelineclientset := pipelineclient.Get(ctx)
	if _, err := pipelineclientset.TektonV1().Tasks(task.Namespace).Create(ctx, &task, metav1.CreateOptions{}); err != nil {
		return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to create Task %s/%s", task.Namespace, task.Name), err), nil
	}

	return result(fmt.Sprintf("Creating task %s in namespace %s", task.Name, task.Namespace)), nil
}
