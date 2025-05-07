package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tektoncd/mcp-server/internal/params"
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
				mcp.Description("Name or Reference of the Pipeline to start"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the Pipeline is located"),
				mcp.DefaultString("default"),
			),
			mcp.WithString("params",
				mcp.Description("Parameters for the Pipeline in the format: key1=value1,key2=value2,key3=array:val1:val2:val3,key4=object:k1=v1:k2=v2"),
			),
		),
		Handler: handlerStartPipeline,
	}
}

func handlerStartPipeline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name must be a string"), nil
	}
	namespace, ok := request.Params.Arguments["namespace"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}

	// Extract parameters if provided
	var pipelineParams []v1.Param
	if paramsStr, ok := request.Params.Arguments["params"].(string); ok && paramsStr != "" {
		var err error
		pipelineParams, err = params.ParsePipelineRunParams(paramsStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse parameters: %v", err)), nil
		}
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
			Params: pipelineParams,
		},
	}

	createdPr, err := pipelineclientset.TektonV1().PipelineRuns(namespace).Create(ctx, pr, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create PipelineRun %s/%s: %v", namespace, name, err)), nil
	}

	// Include parameter information in the result message
	resultMsg := fmt.Sprintf("Starting pipeline %s in namespace %s", name, namespace)
	if len(pipelineParams) > 0 {
		paramDescs := make([]string, 0, len(pipelineParams))
		for _, p := range pipelineParams {
			var valueStr string
			switch p.Value.Type {
			case v1.ParamTypeArray:
				valueStr = fmt.Sprintf("[%s]", strings.Join(p.Value.ArrayVal, ", "))
			case v1.ParamTypeObject:
				pairs := make([]string, 0, len(p.Value.ObjectVal))
				for k, v := range p.Value.ObjectVal {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
				}
				valueStr = fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
			default:
				valueStr = p.Value.StringVal
			}
			paramDescs = append(paramDescs, fmt.Sprintf("%s=%s", p.Name, valueStr))
		}
		resultMsg += fmt.Sprintf(" with parameters: %s", strings.Join(paramDescs, ", "))
	}

	resultMsg += fmt.Sprintf("\nCreated PipelineRun: %s", createdPr.Name)

	return result(resultMsg), nil
}

func startTask() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("start_task",
			mcp.WithDescription("Start a Task"),
			mcp.WithString("name", mcp.Required(),
				mcp.Description("Name or Reference of the Task to start"),
			),
			mcp.WithString("namespace",
				mcp.Description("Namespace where the Task is located"),
				mcp.DefaultString("default"),
			),
			mcp.WithString("params",
				mcp.Description("Parameters for the Task in the format: key1=value1,key2=value2,key3=array:val1:val2:val3,key4=object:k1=v1:k2=v2"),
			),
		),
		Handler: handlerStartTask,
	}
}

func handlerStartTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name must be a string"), nil
	}
	namespace, ok := request.Params.Arguments["namespace"].(string)
	if !ok {
		return mcp.NewToolResultError("namespace must be a string"), nil
	}

	// Extract parameters if provided
	var taskParams []v1.Param
	if paramsStr, ok := request.Params.Arguments["params"].(string); ok && paramsStr != "" {
		var err error
		taskParams, err = params.ParsePipelineRunParams(paramsStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to parse parameters: %v", err)), nil
		}
	}

	taskInformer := taskinformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	if _, err := taskInformer.Lister().Tasks(namespace).Get(name); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get Task %s/%s: %v", namespace, name, err)), nil
	}

	tr := &v1.TaskRun{
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
			Params: taskParams,
		},
	}

	createdTr, err := pipelineclientset.TektonV1().TaskRuns(namespace).Create(ctx, tr, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create TaskRun %s/%s: %v", namespace, name, err)), nil
	}

	// Include parameter information in the result message
	resultMsg := fmt.Sprintf("Starting task %s in namespace %s", name, namespace)
	if len(taskParams) > 0 {
		paramDescs := make([]string, 0, len(taskParams))
		for _, p := range taskParams {
			var valueStr string
			switch p.Value.Type {
			case v1.ParamTypeArray:
				valueStr = fmt.Sprintf("[%s]", strings.Join(p.Value.ArrayVal, ", "))
			case v1.ParamTypeObject:
				pairs := make([]string, 0, len(p.Value.ObjectVal))
				for k, v := range p.Value.ObjectVal {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
				}
				valueStr = fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
			default:
				valueStr = p.Value.StringVal
			}
			paramDescs = append(paramDescs, fmt.Sprintf("%s=%s", p.Name, valueStr))
		}
		resultMsg += fmt.Sprintf(" with parameters: %s", strings.Join(paramDescs, ", "))
	}

	resultMsg += fmt.Sprintf("\nCreated TaskRun: %s", createdTr.Name)

	return result(resultMsg), nil
}
