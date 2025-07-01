package tools

import (
	"context"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func Add(_ context.Context, s *mcp.Server) error {
	startPipelineTool, err := startPipeline()
	if err != nil {
		return err
	}
	startTaskTool, err := startTask()
	if err != nil {
		return err
	}
	restartPipelineRunTool, err := restartPipelineRun()
	if err != nil {
		return err
	}
	restartTaskRunTool, err := restartTaskRun()
	if err != nil {
		return err
	}
	getTaskRunLogsTool, err := getTaskRunLogs()
	if err != nil {
		return err
	}

	s.AddTools(
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
	)
	return nil
}

func result(s string) *mcp.CallToolResultFor[string] {
	return &mcp.CallToolResultFor[string]{
		Content: []mcp.Content{&mcp.TextContent{Text: s}},
	}
}
