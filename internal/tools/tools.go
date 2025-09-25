package tools

import (
	"context"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultNamespace = "default"

func Add(_ context.Context, s *mcp.Server) error {
	// Start tools
	startPipelineTool, err := startPipeline()
	if err != nil {
		return err
	}
	startTaskTool, err := startTask()
	if err != nil {
		return err
	}

	// Restart tools
	restartPipelineRunTool, err := restartPipelineRun()
	if err != nil {
		return err
	}
	restartTaskRunTool, err := restartTaskRun()
	if err != nil {
		return err
	}

	// Log tools
	getTaskRunLogsTool, err := getTaskRunLogs()
	if err != nil {
		return err
	}

	// Create tools
	createPipelineTool, err := createPipeline()
	if err != nil {
		return err
	}
	createTaskTool, err := createTask()
	if err != nil {
		return err
	}
	createPipelineRunTool, err := createPipelineRun()
	if err != nil {
		return err
	}
	createTaskRunTool, err := createTaskRun()
	if err != nil {
		return err
	}

	// Update tools
	updatePipelineTool, err := updatePipeline()
	if err != nil {
		return err
	}
	updateTaskTool, err := updateTask()
	if err != nil {
		return err
	}
	patchPipelineTool, err := patchPipeline()
	if err != nil {
		return err
	}

	// Delete tools
	deletePipelineTool, err := deletePipeline()
	if err != nil {
		return err
	}
	deleteTaskTool, err := deleteTask()
	if err != nil {
		return err
	}
	deletePipelineRunTool, err := deletePipelineRun()
	if err != nil {
		return err
	}
	deleteTaskRunTool, err := deleteTaskRun()
	if err != nil {
		return err
	}
	deleteAllPipelineRunsTool, err := deleteAllPipelineRuns()
	if err != nil {
		return err
	}

	// Get tools
	getPipelineTool, err := getPipeline()
	if err != nil {
		return err
	}
	getTaskTool, err := getTask()
	if err != nil {
		return err
	}
	getPipelineRunTool, err := getPipelineRun()
	if err != nil {
		return err
	}
	getTaskRunTool, err := getTaskRun()
	if err != nil {
		return err
	}

	// Artifact Hub tools
	listArtifactHubTasksTool := listArtifactHubTasks()
	listArtifactHubPipelinesTool := listArtifactHubPipelines()
	installArtifactHubTaskTool, err := installArtifactHubTask()
	if err != nil {
		return err
	}
	installArtifactHubPipelineTool, err := installArtifactHubPipeline()
	if err != nil {
		return err
	}
	triggerArtifactHubTaskTool := triggerArtifactHubTask()
	triggerArtifactHubPipelineTool := triggerArtifactHubPipeline()

	s.AddTools(
		// Existing tools
		startPipelineTool,
		startTaskTool,
		restartPipelineRunTool,
		restartTaskRunTool,
		getTaskRunLogsTool,
		listPipelineRuns(),
		listPipelines(),
		listTaskRuns(),
		listTasks(),
		listStepactions(),

		// Create operations
		createPipelineTool,
		createTaskTool,
		createPipelineRunTool,
		createTaskRunTool,

		// Read/Get operations
		getPipelineTool,
		getTaskTool,
		getPipelineRunTool,
		getTaskRunTool,

		// Update operations
		updatePipelineTool,
		updateTaskTool,
		patchPipelineTool,

		// Delete operations
		deletePipelineTool,
		deleteTaskTool,
		deletePipelineRunTool,
		deleteTaskRunTool,
		deleteAllPipelineRunsTool,

		// Artifact Hub operations
		listArtifactHubTasksTool,
		listArtifactHubPipelinesTool,
		installArtifactHubTaskTool,
		installArtifactHubPipelineTool,
		triggerArtifactHubTaskTool,
		triggerArtifactHubPipelineTool,
	)
	return nil
}

func result(s string) *mcp.CallToolResultFor[string] {
	return &mcp.CallToolResultFor[string]{
		Content: []mcp.Content{&mcp.TextContent{Text: s}},
	}
}
