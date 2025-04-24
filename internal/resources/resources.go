package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	pipelineinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipeline"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	taskinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/task"
	taskruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/taskrun"
	stepactioninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1beta1/stepaction"
)

func Add(ctx context.Context, s *server.MCPServer) {
	s.AddResourceTemplate(pipelineRunResources(ctx))
}

func pipelineRunResources(ctx context.Context) (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
		"tekton://pipelinerun/{namespace}/{name}",
		"PipelineRun",
	), pipelineRunHandler(ctx)
}

func pipelineRunHandler(ctx context.Context) server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		ns, ok := request.Params.Arguments["namespace"].([]string)
		if !ok || len(ns) == 0 {
			return nil, errors.New("namespace is required")
		}
		namespace := ns[0]

		n, ok := request.Params.Arguments["name"].([]string)
		if !ok || len(n) == 0 {
			return nil, errors.New("name is required")
		}
		name := n[0]

		uri := request.Params.URI
		resourceType := strings.Split(uri, "/")[2]

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

		contents := mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json;type=" + resourceType,
			Text:     string(jsonData),
		}

		return []mcp.ResourceContents{contents}, nil
	}
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
