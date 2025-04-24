package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func startPipeline() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("start_pipeline",
			mcp.WithDescription("Start a Pipeline"),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("Name or Reference of the Pipeline to sart"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the Pipeline is located"),
				mcp.DefaultString("default"),
			),
			// TODO add "parameters" objects
		),
		Handler: handlerStartPipeline,
	}
}

func handlerStartPipeline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}
	namespace, ok := request.Params.Arguments["namespace"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}

	pipelineInformer := pipelineinformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	if _, err := pipelineInformer.Lister().Pipelines(namespace).Get(name); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Pipeline %s/%s: %v", namespace, name, err)), nil
	}

	pr := &v1.PipelineRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tekton.dev/v1",
			Kind:       "PipelineRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: name + "-",
		},
		Spec: v1.PipelineRunSpec{
			PipelineRef: &v1.PipelineRef{
				Name: name,
			},
		},
	}

	if _, err := pipelineclientset.TektonV1().PipelineRuns(namespace).Create(ctx, pr, metav1.CreateOptions{}); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create PipelineRun %s/%s: %v", namespace, name, err)), nil
	}

	return result(fmt.Sprintf("Starting pipeline %s in namespace %s", name, namespace)), nil
}

func startTask() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("start_task",
			mcp.WithDescription("Start a Task"),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("Name or Reference of the Task to sart"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the Task is located"),
				mcp.DefaultString("default"),
			),
			// TODO add "parameters" objects
		),
		Handler: handlerStartTask,
	}
}

func handlerStartTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}
	namespace, ok := request.Params.Arguments["namespace"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}

	taskInformer := taskinformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	if _, err := taskInformer.Lister().Tasks(namespace).Get(name); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Task %s/%s: %v", namespace, name, err)), nil
	}

	pr := &v1.TaskRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tekton.dev/v1",
			Kind:       "TaskRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: name + "-",
		},
		Spec: v1.TaskRunSpec{
			TaskRef: &v1.TaskRef{
				Name: name,
			},
		},
	}

	if _, err := pipelineclientset.TektonV1().TaskRuns(namespace).Create(ctx, pr, metav1.CreateOptions{}); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create TaskRun %s/%s: %v", namespace, name, err)), nil
	}

	return result(fmt.Sprintf("Starting task %s in namespace %s", name, namespace)), nil
}
