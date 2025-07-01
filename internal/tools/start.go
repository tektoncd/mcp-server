package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type startParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func startPipeline() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[startParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name or referece of the pipeline"
	scheme.Properties["namespace"].Description = "Namespace of the pipeline"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)

	return mcp.NewServerTool(
		"start_pipeline",
		"Start a Pipeline",
		handlerStartPipeline,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerStartPipeline(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[startParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace

	pipelineInformer := pipelineinformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	if _, err := pipelineInformer.Lister().Pipelines(namespace).Get(name); err != nil {
		return nil, fmt.Errorf("failed to get Pipeline %s/%s: %w", namespace, name, err)
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
		return nil, fmt.Errorf("failed to create PipelineRun %s/%s: %w", namespace, name, err)
	}

	return result(fmt.Sprintf("Starting pipeline %s in namespace %s", name, namespace)), nil
}

func startTask() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[startParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name or referece of the task"
	scheme.Properties["namespace"].Description = "Namespace of the task"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)

	return mcp.NewServerTool(
		"start_task",
		"Start a Task",
		handlerStartTask,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerStartTask(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[startParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	taskInformer := taskinformer.Get(ctx)
	pipelineclientset := pipelineclient.Get(ctx)

	if _, err := taskInformer.Lister().Tasks(namespace).Get(name); err != nil {
		return nil, fmt.Errorf("failed to get Task %s/%s: %w", namespace, name, err)
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
		return nil, fmt.Errorf("failed to create TaskRun %s/%s: %w", namespace, name, err)
	}

	return result(fmt.Sprintf("Starting task %s in namespace %s", name, namespace)), nil
}
