package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	stepactioninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1beta1/stepaction"
)

func Add(_ context.Context, s *mcp.Server) {
	for _, template := range []*mcp.ResourceTemplate{
		{Name: "Pipeline", URITemplate: "tekton://pipeline/{namespace}/{name}"},
		{Name: "PipelineRun", URITemplate: "tekton://pipelinerun/{namespace}/{name}"},
		{Name: "Task", URITemplate: "tekton://task/{namespace}/{name}"},
		{Name: "TaskRun", URITemplate: "tekton://taskrun/{namespace}/{name}"},
		{Name: "StepAction", URITemplate: "tekton://stepaction/{namespace}/{name}"},
	} {
		s.AddResourceTemplate(template, resourceHandler)
	}
}

func resourceHandler(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := request.Params.URI
	parsed := strings.Split(uri, "/")
	resourceType := parsed[2]
	namespace := parsed[3]
	name := parsed[4]

	var jsonData []byte
	var err error

	slog.Info(fmt.Sprintf("Resource: %s, %s/%s", resourceType, namespace, name))

	switch resourceType {
	case "pipelinerun":
		jsonData, err = getPipelineRun(ctx, namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get PipelineRun %s/%s: %w", namespace, name, err)
		}
	case "taskrun":
		jsonData, err = getTaskRun(ctx, namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get TaskRun %s/%s: %w", namespace, name, err)
		}
	case "pipeline":
		jsonData, err = getPipeline(ctx, namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get Pipeline %s/%s: %w", namespace, name, err)
		}
	case "task":
		jsonData, err = getTask(ctx, namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get Task %s/%s: %w", namespace, name, err)
		}
	case "stepaction":
		jsonData, err = getStepAction(ctx, namespace, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get StepAction %s/%s: %w", namespace, name, err)
		}
	}

	contents := &mcp.ResourceContents{
		URI:      uri,
		MIMEType: "application/json;type=" + resourceType,
		Text:     string(jsonData),
	}

	return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{contents}}, nil
}

func getPipelineRun(ctx context.Context, namespace string, name string) ([]byte, error) {
	pipelineRunInformer := pipelineruninformer.Get(ctx)
	pipelineRun, err := pipelineRunInformer.Lister().PipelineRuns(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get PipelineRun %s/%s: %w", namespace, name, err)
	}
	slog.Info(fmt.Sprintf("PipelineRun %s: %v", pipelineRun.Name, pipelineRun))

	jsonData, err := json.Marshal(pipelineRun)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}
	return jsonData, nil
}

func getTaskRun(ctx context.Context, namespace string, name string) ([]byte, error) {
	taskRunInformer := taskruninformer.Get(ctx)
	taskRun, err := taskRunInformer.Lister().TaskRuns(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get TaskRun %s/%s: %w", namespace, name, err)
	}

	jsonData, err := json.Marshal(taskRun)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}
	return jsonData, nil
}

func getPipeline(ctx context.Context, namespace string, name string) ([]byte, error) {
	pipelineInformer := pipelineinformer.Get(ctx)
	pipeline, err := pipelineInformer.Lister().Pipelines(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pipeline %s/%s: %w", namespace, name, err)
	}

	jsonData, err := json.Marshal(pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}
	return jsonData, nil
}

func getTask(ctx context.Context, namespace string, name string) ([]byte, error) {
	taskInformer := taskinformer.Get(ctx)
	task, err := taskInformer.Lister().Tasks(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get Task %s/%s: %w", namespace, name, err)
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}
	return jsonData, nil
}

func getStepAction(ctx context.Context, namespace string, name string) ([]byte, error) {
	stepActionInformer := stepactioninformer.Get(ctx)
	stepAction, err := stepActionInformer.Lister().StepActions(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get StepAction %s/%s: %w", namespace, name, err)
	}

	jsonData, err := json.Marshal(stepAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to JSON: %w", err)
	}
	return jsonData, nil
}
