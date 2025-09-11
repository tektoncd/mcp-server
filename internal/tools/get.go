package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
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
	Namespace string `json:"namespace"`
	Output    string `json:"output"`
}

func getPipeline() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[getPipelineParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Pipeline to get"
	scheme.Properties["namespace"].Description = "Namespace of the Pipeline"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"get_pipeline",
		"Get a specific Pipeline by name",
		handlerGetPipeline,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerGetPipeline(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[getPipelineParams],
) (*mcp.CallToolResultFor[string], error) {
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
	Namespace string `json:"namespace"`
	Output    string `json:"output"`
}

func getTask() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[getTaskParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Task to get"
	scheme.Properties["namespace"].Description = "Namespace of the Task"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"get_task",
		"Get a specific Task by name",
		handlerGetTask,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerGetTask(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[getTaskParams],
) (*mcp.CallToolResultFor[string], error) {
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
	Namespace string `json:"namespace"`
	Output    string `json:"output"`
}

func getPipelineRun() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[getPipelineRunParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the PipelineRun to get"
	scheme.Properties["namespace"].Description = "Namespace of the PipelineRun"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"get_pipelinerun",
		"Get a specific PipelineRun by name",
		handlerGetPipelineRun,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerGetPipelineRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[getPipelineRunParams],
) (*mcp.CallToolResultFor[string], error) {
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
	Namespace string `json:"namespace"`
	Output    string `json:"output"`
}

func getTaskRun() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[getTaskRunParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the TaskRun to get"
	scheme.Properties["namespace"].Description = "Namespace of the TaskRun"
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["output"].Description = outputFormatDescription
	scheme.Properties["output"].Default = json.RawMessage(`"yaml"`)
	scheme.Required = []string{"name"}

	return mcp.NewServerTool(
		"get_taskrun",
		"Get a specific TaskRun by name",
		handlerGetTaskRun,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerGetTaskRun(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[getTaskRunParams],
) (*mcp.CallToolResultFor[string], error) {
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
