package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"
)

const (
	namespacePipelineDescription = "Namespace of the Pipeline"
	namespaceTaskDescription     = "Namespace of the Task"
)

type updatePipelineParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Yaml      string `json:"yaml"`
}

func updatePipeline() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[updatePipelineParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Pipeline to update"
	scheme.Properties["namespace"].Description = namespacePipelineDescription
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["yaml"].Description = "Updated YAML definition of the Pipeline"
	scheme.Required = []string{"name", "yaml"}

	return mcp.NewServerTool(
		"update_pipeline",
		"Update an existing Pipeline",
		handlerUpdatePipeline,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerUpdatePipeline(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[updatePipelineParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	yamlStr := params.Arguments.Yaml

	if name == "" || yamlStr == "" {
		return result("Error: Name and YAML definition are required"), nil
	}

	var pipeline pipelinev1.Pipeline
	if err := yaml.Unmarshal([]byte(yamlStr), &pipeline); err != nil {
		return result(fmt.Sprintf("Error parsing YAML: %v", err)), nil
	}

	pipelineClient := pipelineclient.Get(ctx)

	existing, err := pipelineClient.TektonV1().Pipelines(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error getting existing Pipeline: %v", err)), nil
	}

	pipeline.Name = name
	pipeline.Namespace = namespace
	pipeline.ResourceVersion = existing.ResourceVersion

	updated, err := pipelineClient.TektonV1().Pipelines(namespace).Update(ctx, &pipeline, metav1.UpdateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error updating Pipeline: %v", err)), nil
	}

	return result(fmt.Sprintf("Pipeline '%s' updated successfully in namespace '%s'", updated.Name, namespace)), nil
}

type updateTaskParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Yaml      string `json:"yaml"`
}

func updateTask() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[updateTaskParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Task to update"
	scheme.Properties["namespace"].Description = namespaceTaskDescription
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["yaml"].Description = "Updated YAML definition of the Task"
	scheme.Required = []string{"name", "yaml"}

	return mcp.NewServerTool(
		"update_task",
		"Update an existing Task",
		handlerUpdateTask,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerUpdateTask(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[updateTaskParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	yamlStr := params.Arguments.Yaml

	if name == "" || yamlStr == "" {
		return result("Error: Name and YAML definition are required"), nil
	}

	var task pipelinev1.Task
	if err := yaml.Unmarshal([]byte(yamlStr), &task); err != nil {
		return result(fmt.Sprintf("Error parsing YAML: %v", err)), nil
	}

	pipelineClient := pipelineclient.Get(ctx)

	existing, err := pipelineClient.TektonV1().Tasks(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error getting existing Task: %v", err)), nil
	}

	task.Name = name
	task.Namespace = namespace
	task.ResourceVersion = existing.ResourceVersion

	updated, err := pipelineClient.TektonV1().Tasks(namespace).Update(ctx, &task, metav1.UpdateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error updating Task: %v", err)), nil
	}

	return result(fmt.Sprintf("Task '%s' updated successfully in namespace '%s'", updated.Name, namespace)), nil
}

type patchPipelineParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Patch     string `json:"patch"`
}

func patchPipeline() (*mcp.ServerTool, error) {
	scheme, err := jsonschema.For[patchPipelineParams]()
	if err != nil {
		return nil, err
	}

	scheme.Properties["name"].Description = "Name of the Pipeline to patch"
	scheme.Properties["namespace"].Description = namespacePipelineDescription
	scheme.Properties["namespace"].Default = json.RawMessage(`"default"`)
	scheme.Properties["patch"].Description = "JSON patch to apply to the Pipeline"
	scheme.Required = []string{"name", "patch"}

	return mcp.NewServerTool(
		"patch_pipeline",
		"Apply a JSON patch to an existing Pipeline",
		handlerPatchPipeline,
		mcp.Input(mcp.Schema(scheme)),
	), nil
}

func handlerPatchPipeline(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[patchPipelineParams],
) (*mcp.CallToolResultFor[string], error) {
	name := params.Arguments.Name
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}
	patchStr := params.Arguments.Patch

	if name == "" || patchStr == "" {
		return result("Error: Name and patch are required"), nil
	}

	pipelineClient := pipelineclient.Get(ctx)

	patched, err := pipelineClient.TektonV1().Pipelines(namespace).Patch(
		ctx,
		name,
		types.JSONPatchType,
		[]byte(patchStr),
		metav1.PatchOptions{},
	)
	if err != nil {
		return result(fmt.Sprintf("Error patching Pipeline: %v", err)), nil
	}

	return result(fmt.Sprintf("Pipeline '%s' patched successfully in namespace '%s'", patched.Name, namespace)), nil
}
