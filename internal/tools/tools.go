package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Add(_ context.Context, s *server.MCPServer) {
	s.AddTools(
		startPipeline(),
		startTask(),
		listPipelineRuns(),
		listPipelines(),
		listTaskRuns(),
		listTasks(),
		listStepactions(),
	)
}

func result(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(s),
		},
	}
}
