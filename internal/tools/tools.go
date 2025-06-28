package tools

import (
	"context"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func Add(_ context.Context, s *mcp.Server) {
	s.AddTools(
		startPipeline(),
		startTask(),
		restartPipelineRun(),
		restartTaskRun(),
		getTaskRunLogs(),
		listPipelineRuns(),
		listPipelines(),
		listTaskRuns(),
		listTasks(),
		listStepactions(),
	)
}

func result(s string) *mcp.CallToolResultFor[string] {
	return &mcp.CallToolResultFor[string]{
		Content: []mcp.Content{&mcp.TextContent{Text: s}},
	}
}
