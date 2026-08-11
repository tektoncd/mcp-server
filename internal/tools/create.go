package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type createPipelineParams struct {
	Namespace string `json:"namespace,omitempty"`
	Yaml      string `json:"yaml"`
}

func createPipeline() (serverTool, error) {
	scheme, err := jsonschema.For[createPipelineParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["namespace"].Description = "Namespace where the Pipeline will be created"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["yaml"].Description = "YAML definition of the Pipeline"
	scheme.Required = []string{"yaml"}

	return newServerTool(
		"create_pipeline",
		"Create a new Pipeline from YAML definition",
		handlerCreatePipeline,
		scheme,
	), nil
}

func handlerCreatePipeline(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[createPipelineParams],
) (*mcp.CallToolResult, error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	yamlStr := params.Arguments.Yaml

	if yamlStr == "" {
		return result("Error: YAML definition is required"), nil
	}

	var pipeline pipelinev1.Pipeline
	if err := yaml.Unmarshal([]byte(yamlStr), &pipeline); err != nil {
		return result(fmt.Sprintf("Error parsing YAML: %v", err)), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	created, err := pipelineClient.TektonV1().Pipelines(namespace).Create(ctx, &pipeline, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating Pipeline: %v", err)), nil
	}

	return result(fmt.Sprintf("Pipeline '%s' created successfully in namespace '%s'", created.Name, namespace)), nil
}

type createTaskParams struct {
	Namespace string `json:"namespace,omitempty"`
	Yaml      string `json:"yaml"`
}

func createTask() (serverTool, error) {
	scheme, err := jsonschema.For[createTaskParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["namespace"].Description = "Namespace where the Task will be created"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["yaml"].Description = "YAML definition of the Task"
	scheme.Required = []string{"yaml"}

	return newServerTool(
		"create_task",
		"Create a new Task from YAML definition",
		handlerCreateTask,
		scheme,
	), nil
}

func handlerCreateTask(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[createTaskParams],
) (*mcp.CallToolResult, error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	yamlStr := params.Arguments.Yaml

	if yamlStr == "" {
		return result("Error: YAML definition is required"), nil
	}

	var task pipelinev1.Task
	if err := yaml.Unmarshal([]byte(yamlStr), &task); err != nil {
		return result(fmt.Sprintf("Error parsing YAML: %v", err)), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	created, err := pipelineClient.TektonV1().Tasks(namespace).Create(ctx, &task, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating Task: %v", err)), nil
	}

	return result(fmt.Sprintf("Task '%s' created successfully in namespace '%s'", created.Name, namespace)), nil
}

type createPipelineRunParams struct {
	Namespace    string `json:"namespace,omitempty"`
	Yaml         string `json:"yaml,omitempty"`
	GenerateName string `json:"generateName,omitempty"`
}

func createPipelineRun() (serverTool, error) {
	scheme, err := jsonschema.For[createPipelineRunParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["namespace"].Description = "Namespace where the PipelineRun will be created"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["yaml"].Description = "YAML definition of the PipelineRun"
	scheme.Properties["generateName"].Description = "Generate name prefix for the PipelineRun (alternative to fixed name)"

	return newServerTool(
		"create_pipelinerun",
		"Create a new PipelineRun from YAML definition or generate from Pipeline",
		handlerCreatePipelineRun,
		scheme,
	), nil
}

func handlerCreatePipelineRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[createPipelineRunParams],
) (*mcp.CallToolResult, error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	yamlStr := params.Arguments.Yaml
	generateName := params.Arguments.GenerateName

	if yamlStr == "" && generateName == "" {
		return result("Error: Either YAML definition or generateName is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)

	if yamlStr != "" {
		var pipelineRun pipelinev1.PipelineRun
		if err := yaml.Unmarshal([]byte(yamlStr), &pipelineRun); err != nil {
			return result(fmt.Sprintf("Error parsing YAML: %v", err)), nil
		}

		created, err := pipelineClient.TektonV1().PipelineRuns(namespace).Create(ctx, &pipelineRun, metav1.CreateOptions{})
		if err != nil {
			return result(fmt.Sprintf("Error creating PipelineRun: %v", err)), nil
		}
		return result(fmt.Sprintf("PipelineRun '%s' created successfully in namespace '%s'", created.Name, namespace)), nil
	}

	pipelineRun := &pipelinev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: generateName,
		},
	}

	created, err := pipelineClient.TektonV1().PipelineRuns(namespace).Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating PipelineRun: %v", err)), nil
	}

	return result(fmt.Sprintf("PipelineRun '%s' created successfully in namespace '%s'", created.Name, namespace)), nil
}

type createTaskRunParams struct {
	Namespace    string `json:"namespace,omitempty"`
	Yaml         string `json:"yaml,omitempty"`
	GenerateName string `json:"generateName,omitempty"`
}

func createTaskRun() (serverTool, error) {
	scheme, err := jsonschema.For[createTaskRunParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["namespace"].Description = "Namespace where the TaskRun will be created"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["yaml"].Description = "YAML definition of the TaskRun"
	scheme.Properties["generateName"].Description = "Generate name prefix for the TaskRun (alternative to fixed name)"

	return newServerTool(
		"create_taskrun",
		"Create a new TaskRun from YAML definition",
		handlerCreateTaskRun,
		scheme,
	), nil
}

func handlerCreateTaskRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[createTaskRunParams],
) (*mcp.CallToolResult, error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	yamlStr := params.Arguments.Yaml
	generateName := params.Arguments.GenerateName

	if yamlStr == "" && generateName == "" {
		return result("Error: Either YAML definition or generateName is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)

	if yamlStr != "" {
		var taskRun pipelinev1.TaskRun
		if err := yaml.Unmarshal([]byte(yamlStr), &taskRun); err != nil {
			return result(fmt.Sprintf("Error parsing YAML: %v", err)), nil
		}

		created, err := pipelineClient.TektonV1().TaskRuns(namespace).Create(ctx, &taskRun, metav1.CreateOptions{})
		if err != nil {
			return result(fmt.Sprintf("Error creating TaskRun: %v", err)), nil
		}
		return result(fmt.Sprintf("TaskRun '%s' created successfully in namespace '%s'", created.Name, namespace)), nil
	}

	taskRun := &pipelinev1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: generateName,
		},
	}

	created, err := pipelineClient.TektonV1().TaskRuns(namespace).Create(ctx, taskRun, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating TaskRun: %v", err)), nil
	}

	return result(fmt.Sprintf("TaskRun '%s' created successfully in namespace '%s'", created.Name, namespace)), nil
}
