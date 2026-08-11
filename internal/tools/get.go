package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	outputFormatDescription = "Output format (json or yaml)"
	outputFormatYAML        = "yaml"
	outputFormatJSON        = "json"
)

type getPipelineParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Output    string `json:"output,omitempty"`
}

func getPipeline() (serverTool, error) {
	scheme, err := jsonschema.For[getPipelineParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Pipeline to get"
	scheme.Properties["namespace"].Description = "Namespace of the Pipeline"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return newServerTool(
		"get_pipeline",
		"Get a specific Pipeline by name",
		handlerGetPipeline,
		scheme,
	), nil
}

func handlerGetPipeline(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[getPipelineParams],
) (*mcp.CallToolResult, error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	output := params.Arguments.Output
	if output == "" {
		output = outputFormatYAML
	}

	if name == "" {
		return result("Error: Pipeline name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	pipeline, err := pipelineClient.TektonV1().Pipelines(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error getting Pipeline: %v", err)), nil
	}

	var outputStr string
	if output == outputFormatJSON {
		jsonData, err := json.MarshalIndent(pipeline, "", "  ")
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to JSON: %v", err)), nil
		}
		outputStr = string(jsonData)
	} else {
		yamlData, err := yaml.Marshal(pipeline)
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to YAML: %v", err)), nil
		}
		outputStr = string(yamlData)
	}

	return result(outputStr), nil
}

type getTaskParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Output    string `json:"output,omitempty"`
}

func getTask() (serverTool, error) {
	scheme, err := jsonschema.For[getTaskParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Task to get"
	scheme.Properties["namespace"].Description = "Namespace of the Task"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return newServerTool(
		"get_task",
		"Get a specific Task by name",
		handlerGetTask,
		scheme,
	), nil
}

func handlerGetTask(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[getTaskParams],
) (*mcp.CallToolResult, error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	output := params.Arguments.Output
	if output == "" {
		output = outputFormatYAML
	}

	if name == "" {
		return result("Error: Task name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	task, err := pipelineClient.TektonV1().Tasks(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error getting Task: %v", err)), nil
	}

	var outputStr string
	if output == outputFormatJSON {
		jsonData, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to JSON: %v", err)), nil
		}
		outputStr = string(jsonData)
	} else {
		yamlData, err := yaml.Marshal(task)
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to YAML: %v", err)), nil
		}
		outputStr = string(yamlData)
	}

	return result(outputStr), nil
}

type getPipelineRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Output    string `json:"output,omitempty"`
}

func getPipelineRun() (serverTool, error) {
	scheme, err := jsonschema.For[getPipelineRunParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the PipelineRun to get"
	scheme.Properties["namespace"].Description = "Namespace of the PipelineRun"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return newServerTool(
		"get_pipelinerun",
		"Get a specific PipelineRun by name",
		handlerGetPipelineRun,
		scheme,
	), nil
}

func handlerGetPipelineRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[getPipelineRunParams],
) (*mcp.CallToolResult, error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	output := params.Arguments.Output
	if output == "" {
		output = outputFormatYAML
	}

	if name == "" {
		return result("Error: PipelineRun name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	pipelineRun, err := pipelineClient.TektonV1().PipelineRuns(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error getting PipelineRun: %v", err)), nil
	}

	var outputStr string
	if output == outputFormatJSON {
		jsonData, err := json.MarshalIndent(pipelineRun, "", "  ")
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to JSON: %v", err)), nil
		}
		outputStr = string(jsonData)
	} else {
		yamlData, err := yaml.Marshal(pipelineRun)
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to YAML: %v", err)), nil
		}
		outputStr = string(yamlData)
	}

	return result(outputStr), nil
}

type getTaskRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Output    string `json:"output,omitempty"`
}

func getTaskRun() (serverTool, error) {
	scheme, err := jsonschema.For[getTaskRunParams](nil)
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the TaskRun to get"
	scheme.Properties["namespace"].Description = "Namespace of the TaskRun"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return newServerTool(
		"get_taskrun",
		"Get a specific TaskRun by name",
		handlerGetTaskRun,
		scheme,
	), nil
}

func handlerGetTaskRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *callToolParamsFor[getTaskRunParams],
) (*mcp.CallToolResult, error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	output := params.Arguments.Output
	if output == "" {
		output = outputFormatYAML
	}

	if name == "" {
		return result("Error: TaskRun name is required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)
	taskRun, err := pipelineClient.TektonV1().TaskRuns(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error getting TaskRun: %v", err)), nil
	}

	var outputStr string
	if output == outputFormatJSON {
		jsonData, err := json.MarshalIndent(taskRun, "", "  ")
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to JSON: %v", err)), nil
		}
		outputStr = string(jsonData)
	} else {
		yamlData, err := yaml.Marshal(taskRun)
		if err != nil {
			return result(fmt.Sprintf("Error marshaling to YAML: %v", err)), nil
		}
		outputStr = string(yamlData)
	}

	return result(outputStr), nil
}
