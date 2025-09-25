package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tektoncd/mcp-server/internal/artifacthub"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonClientSet "github.com/tektoncd/pipeline/pkg/client/injection/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

type artifactHubSearchParams struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type artifactHubInstallParams struct {
	PackageID string `json:"packageId"`
	Version   string `json:"version"`
	Namespace string `json:"namespace"`
}

func listArtifactHubTasks() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_artifacthub_tasks",
		"List Tekton tasks from Artifact Hub with search options",
		handlerListArtifactHubTasks,
	)
}

func listArtifactHubPipelines() *mcp.ServerTool {
	return mcp.NewServerTool(
		"list_artifacthub_pipelines",
		"List Tekton pipelines from Artifact Hub with search options",
		handlerListArtifactHubPipelines,
	)
}

func installArtifactHubTask() (*mcp.ServerTool, error) {
	return mcp.NewServerTool(
		"install_artifacthub_task",
		"Install and trigger a Tekton task from Artifact Hub to the cluster",
		handlerInstallArtifactHubTask,
	), nil
}

func installArtifactHubPipeline() (*mcp.ServerTool, error) {
	return mcp.NewServerTool(
		"install_artifacthub_pipeline",
		"Install and trigger a Tekton pipeline from Artifact Hub to the cluster",
		handlerInstallArtifactHubPipeline,
	), nil
}

func handlerListArtifactHubTasks(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[artifactHubSearchParams],
) (*mcp.CallToolResultFor[string], error) {
	// Set default limit if not provided
	if params.Arguments.Limit <= 0 {
		params.Arguments.Limit = 20
	}

	client := artifacthub.NewClient()

	resp, err := client.SearchTektonTasks(ctx, params.Arguments.Query, params.Arguments.Limit)
	if err != nil {
		return result(fmt.Sprintf("Error searching Artifact Hub for tasks: %v", err)), nil
	}

	if len(resp.Packages) == 0 {
		return result("No Tekton tasks found on Artifact Hub"), nil
	}

	// Format results for display
	var output strings.Builder
	output.WriteString(fmt.Sprintf("Found %d Tekton tasks on Artifact Hub:\n\n", len(resp.Packages)))

	for i, pkg := range resp.Packages {
		output.WriteString(fmt.Sprintf("%d. **%s** (v%s)\n", i+1, pkg.DisplayName, pkg.Version))
		output.WriteString(fmt.Sprintf("   ID: %s\n", pkg.PackageID))
		if pkg.Description != "" {
			output.WriteString(fmt.Sprintf("   Description: %s\n", pkg.Description))
		}
		if len(pkg.Keywords) > 0 {
			output.WriteString(fmt.Sprintf("   Keywords: %s\n", strings.Join(pkg.Keywords, ", ")))
		}
		output.WriteString(fmt.Sprintf("   Repository: %s\n", pkg.Repository.DisplayName))
		if pkg.HomeURL != "" {
			output.WriteString(fmt.Sprintf("   Homepage: %s\n", pkg.HomeURL))
		}
		output.WriteString("\n")
	}

	return result(output.String()), nil
}

func handlerListArtifactHubPipelines(
	ctx context.Context,
	cc *mcp.ServerSession,
	request *mcp.CallToolParamsFor[artifactHubSearchParams],
) (*mcp.CallToolResultFor[string], error) {
	// Set default limit if not provided
	if request.Arguments.Limit <= 0 {
		request.Arguments.Limit = 20
	}

	client := artifacthub.NewClient()

	resp, err := client.SearchTektonPipelines(ctx, request.Arguments.Query, request.Arguments.Limit)
	if err != nil {
		return result(fmt.Sprintf("Error searching Artifact Hub for pipelines: %v", err)), nil
	}

	if len(resp.Packages) == 0 {
		return result("No Tekton pipelines found on Artifact Hub"), nil
	}

	// Format results for display
	var output strings.Builder
	output.WriteString(fmt.Sprintf("Found %d Tekton pipelines on Artifact Hub:\n\n", len(resp.Packages)))

	for i, pkg := range resp.Packages {
		output.WriteString(fmt.Sprintf("%d. **%s** (v%s)\n", i+1, pkg.DisplayName, pkg.Version))
		output.WriteString(fmt.Sprintf("   ID: %s\n", pkg.PackageID))
		if pkg.Description != "" {
			output.WriteString(fmt.Sprintf("   Description: %s\n", pkg.Description))
		}
		if len(pkg.Keywords) > 0 {
			output.WriteString(fmt.Sprintf("   Keywords: %s\n", strings.Join(pkg.Keywords, ", ")))
		}
		output.WriteString(fmt.Sprintf("   Repository: %s\n", pkg.Repository.DisplayName))
		if pkg.HomeURL != "" {
			output.WriteString(fmt.Sprintf("   Homepage: %s\n", pkg.HomeURL))
		}
		output.WriteString("\n")
	}

	return result(output.String()), nil
}

func handlerInstallArtifactHubTask(
	ctx context.Context,
	cc *mcp.ServerSession,
	request *mcp.CallToolParamsFor[artifactHubInstallParams],
) (*mcp.CallToolResultFor[string], error) {
	if request.Arguments.PackageID == "" {
		return result("Error: packageId parameter is required"), nil
	}

	if request.Arguments.Namespace == "" {
		request.Arguments.Namespace = defaultNamespace
	}

	client := artifacthub.NewClient()

	// Get package details
	pkg, err := client.GetPackage(ctx, request.Arguments.PackageID)
	if err != nil {
		return result(fmt.Sprintf("Error getting package from Artifact Hub: %v", err)), nil
	}

	if pkg.ContentURL == "" {
		return result("Error: Package does not have a content URL"), nil
	}

	// Get package content (YAML definition)
	content, err := client.GetPackageContent(ctx, pkg.ContentURL)
	if err != nil {
		return result(fmt.Sprintf("Error getting package content: %v", err)), nil
	}

	// Parse and apply the task to the cluster
	tektonClient := tektonClientSet.Get(ctx)

	// Parse the YAML content
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, _, err := decoder.Decode([]byte(content), nil, nil)
	if err != nil {
		return result(fmt.Sprintf("Error parsing task YAML: %v", err)), nil
	}

	task, ok := obj.(*v1.Task)
	if !ok {
		return result("Error: Content is not a valid Tekton Task"), nil
	}

	// Set namespace
	task.Namespace = request.Arguments.Namespace

	// Clear resource version for new creation
	task.ResourceVersion = ""

	// Create the task in the cluster
	createdTask, err := tektonClient.TektonV1().Tasks(request.Arguments.Namespace).Create(ctx, task, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating task in cluster: %v", err)), nil
	}

	return result(fmt.Sprintf("Successfully installed Tekton task '%s' (v%s) to namespace '%s' as '%s'",
		pkg.DisplayName, pkg.Version, request.Arguments.Namespace, createdTask.Name)), nil
}

func handlerInstallArtifactHubPipeline(
	ctx context.Context,
	cc *mcp.ServerSession,
	request *mcp.CallToolParamsFor[artifactHubInstallParams],
) (*mcp.CallToolResultFor[string], error) {
	if request.Arguments.PackageID == "" {
		return result("Error: packageId parameter is required"), nil
	}

	if request.Arguments.Namespace == "" {
		request.Arguments.Namespace = defaultNamespace
	}

	client := artifacthub.NewClient()

	// Get package details
	pkg, err := client.GetPackage(ctx, request.Arguments.PackageID)
	if err != nil {
		return result(fmt.Sprintf("Error getting package from Artifact Hub: %v", err)), nil
	}

	if pkg.ContentURL == "" {
		return result("Error: Package does not have a content URL"), nil
	}

	// Get package content (YAML definition)
	content, err := client.GetPackageContent(ctx, pkg.ContentURL)
	if err != nil {
		return result(fmt.Sprintf("Error getting package content: %v", err)), nil
	}

	// Parse and apply the pipeline to the cluster
	tektonClient := tektonClientSet.Get(ctx)

	// Parse the YAML content
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, _, err := decoder.Decode([]byte(content), nil, nil)
	if err != nil {
		return result(fmt.Sprintf("Error parsing pipeline YAML: %v", err)), nil
	}

	pipeline, ok := obj.(*v1.Pipeline)
	if !ok {
		return result("Error: Content is not a valid Tekton Pipeline"), nil
	}

	// Set namespace
	pipeline.Namespace = request.Arguments.Namespace

	// Clear resource version for new creation
	pipeline.ResourceVersion = ""

	// Create the pipeline in the cluster
	createdPipeline, err := tektonClient.TektonV1().Pipelines(request.Arguments.Namespace).Create(ctx, pipeline, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating pipeline in cluster: %v", err)), nil
	}

	return result(fmt.Sprintf("Successfully installed Tekton pipeline '%s' (v%s) to namespace '%s' as '%s'",
		pkg.DisplayName, pkg.Version, request.Arguments.Namespace, createdPipeline.Name)), nil
}

// Helper function to trigger a task run from an installed task
func triggerArtifactHubTask() *mcp.ServerTool {
	return mcp.NewServerTool(
		"trigger_artifacthub_task",
		"Trigger a Tekton task that was installed from Artifact Hub",
		handlerTriggerArtifactHubTask,
	)
}

// Helper function to trigger a pipeline run from an installed pipeline
func triggerArtifactHubPipeline() *mcp.ServerTool {
	return mcp.NewServerTool(
		"trigger_artifacthub_pipeline",
		"Trigger a Tekton pipeline that was installed from Artifact Hub",
		handlerTriggerArtifactHubPipeline,
	)
}

type triggerParams struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Params    map[string]interface{} `json:"params"`
}

func handlerTriggerArtifactHubTask(
	ctx context.Context,
	cc *mcp.ServerSession,
	request *mcp.CallToolParamsFor[triggerParams],
) (*mcp.CallToolResultFor[string], error) {
	if request.Arguments.Name == "" {
		return result("Error: name parameter is required"), nil
	}

	if request.Arguments.Namespace == "" {
		request.Arguments.Namespace = defaultNamespace
	}

	tektonClient := tektonClientSet.Get(ctx)

	// Create TaskRun
	taskRun := &v1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: request.Arguments.Name + "-run-",
			Namespace:    request.Arguments.Namespace,
		},
		Spec: v1.TaskRunSpec{
			TaskRef: &v1.TaskRef{
				Name: request.Arguments.Name,
			},
		},
	}

	// Add parameters if provided
	if len(request.Arguments.Params) > 0 {
		for key, value := range request.Arguments.Params {
			param := v1.Param{
				Name: key,
			}
			// Convert value to string
			switch v := value.(type) {
			case string:
				param.Value.StringVal = v
			case int:
				param.Value.StringVal = strconv.Itoa(v)
			case float64:
				param.Value.StringVal = strconv.FormatFloat(v, 'f', -1, 64)
			case bool:
				param.Value.StringVal = strconv.FormatBool(v)
			default:
				param.Value.StringVal = fmt.Sprintf("%v", v)
			}
			taskRun.Spec.Params = append(taskRun.Spec.Params, param)
		}
	}

	createdTaskRun, err := tektonClient.TektonV1().TaskRuns(request.Arguments.Namespace).Create(ctx, taskRun, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating TaskRun: %v", err)), nil
	}

	return result(fmt.Sprintf("Successfully triggered TaskRun '%s' for task '%s' in namespace '%s'",
		createdTaskRun.Name, request.Arguments.Name, request.Arguments.Namespace)), nil
}

func handlerTriggerArtifactHubPipeline(
	ctx context.Context,
	cc *mcp.ServerSession,
	request *mcp.CallToolParamsFor[triggerParams],
) (*mcp.CallToolResultFor[string], error) {
	if request.Arguments.Name == "" {
		return result("Error: name parameter is required"), nil
	}

	if request.Arguments.Namespace == "" {
		request.Arguments.Namespace = defaultNamespace
	}

	tektonClient := tektonClientSet.Get(ctx)

	// Create PipelineRun
	pipelineRun := &v1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: request.Arguments.Name + "-run-",
			Namespace:    request.Arguments.Namespace,
		},
		Spec: v1.PipelineRunSpec{
			PipelineRef: &v1.PipelineRef{
				Name: request.Arguments.Name,
			},
		},
	}

	// Add parameters if provided
	if len(request.Arguments.Params) > 0 {
		for key, value := range request.Arguments.Params {
			param := v1.Param{
				Name: key,
			}
			// Convert value to string
			switch v := value.(type) {
			case string:
				param.Value.StringVal = v
			case int:
				param.Value.StringVal = strconv.Itoa(v)
			case float64:
				param.Value.StringVal = strconv.FormatFloat(v, 'f', -1, 64)
			case bool:
				param.Value.StringVal = strconv.FormatBool(v)
			default:
				param.Value.StringVal = fmt.Sprintf("%v", v)
			}
			pipelineRun.Spec.Params = append(pipelineRun.Spec.Params, param)
		}
	}

	createdPipelineRun, err := tektonClient.TektonV1().PipelineRuns(request.Arguments.Namespace).Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		return result(fmt.Sprintf("Error creating PipelineRun: %v", err)), nil
	}

	return result(fmt.Sprintf("Successfully triggered PipelineRun '%s' for pipeline '%s' in namespace '%s'",
		createdPipelineRun.Name, request.Arguments.Name, request.Arguments.Namespace)), nil
}
