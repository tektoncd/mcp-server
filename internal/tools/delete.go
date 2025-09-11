package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deletePipelineParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func deletePipeline() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[deletePipelineParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Pipeline to delete"
	scheme.Properties["namespace"].Description = "Namespace of the Pipeline"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"delete_pipeline",
		"Delete a Pipeline",
		handlerDeletePipeline,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerDeletePipeline(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[deletePipelineParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	if name == "" {
		return result("Error: Pipeline name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	err := pipelineClient.TektonV1().Pipelines(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error deleting Pipeline: %v", err)), nil
	}

	return result(fmt.Sprintf("Pipeline '%s' deleted successfully from namespace '%s'", name, namespace)), nil
}

type deleteTaskParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func deleteTask() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[deleteTaskParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Task to delete"
	scheme.Properties["namespace"].Description = "Namespace of the Task"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"delete_task",
		"Delete a Task",
		handlerDeleteTask,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerDeleteTask(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[deleteTaskParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	if name == "" {
		return result("Error: Task name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	err := pipelineClient.TektonV1().Tasks(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error deleting Task: %v", err)), nil
	}

	return result(fmt.Sprintf("Task '%s' deleted successfully from namespace '%s'", name, namespace)), nil
}

type deletePipelineRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func deletePipelineRun() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[deletePipelineRunParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the PipelineRun to delete"
	scheme.Properties["namespace"].Description = "Namespace of the PipelineRun"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"delete_pipelinerun",
		"Delete a PipelineRun",
		handlerDeletePipelineRun,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerDeletePipelineRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[deletePipelineRunParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	if name == "" {
		return result("Error: PipelineRun name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	err := pipelineClient.TektonV1().PipelineRuns(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error deleting PipelineRun: %v", err)), nil
	}

	return result(fmt.Sprintf("PipelineRun '%s' deleted successfully from namespace '%s'", name, namespace)), nil
}

type deleteTaskRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func deleteTaskRun() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[deleteTaskRunParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the TaskRun to delete"
	scheme.Properties["namespace"].Description = "Namespace of the TaskRun"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"delete_taskrun",
		"Delete a TaskRun",
		handlerDeleteTaskRun,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerDeleteTaskRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[deleteTaskRunParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	if name == "" {
		return result("Error: TaskRun name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	err := pipelineClient.TektonV1().TaskRuns(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error deleting TaskRun: %v", err)), nil
	}

	return result(fmt.Sprintf("TaskRun '%s' deleted successfully from namespace '%s'", name, namespace)), nil
}

type deleteAllPipelineRunsParams struct {
	Namespace     string `json:"namespace"`
	LabelSelector string `json:"labelSelector"`
	FieldSelector string `json:"fieldSelector"`
}

func deleteAllPipelineRuns() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[deleteAllPipelineRunsParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["namespace"].Description = "Namespace to delete PipelineRuns from"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["labelSelector"].Description = "Label selector to filter PipelineRuns to delete"
	scheme.Properties["fieldSelector"].Description = "Field selector to filter PipelineRuns to delete"

	return mcp.NewServerTool(
		"delete_all_pipelineruns",
		"Delete multiple PipelineRuns based on selectors",
		handlerDeleteAllPipelineRuns,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerDeleteAllPipelineRuns(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[deleteAllPipelineRunsParams],
) (*mcp.CallToolResultFor[string], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	labelSelector := params.Arguments.LabelSelector
	fieldSelector := params.Arguments.FieldSelector

	pipelineClient := pipelineclient.Get(ctx)

	deleteOptions := metav1.DeleteOptions{}
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
	}

	err := pipelineClient.TektonV1().PipelineRuns(namespace).DeleteCollection(ctx, deleteOptions, listOptions)
	if err != nil {
		return result(fmt.Sprintf("Error deleting PipelineRuns: %v", err)), nil
	}

	return result(fmt.Sprintf("PipelineRuns deleted successfully from namespace '%s' with selectors", namespace)), nil
}
