package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"knative.dev/pkg/signals"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Tekton",
		"0.0.1", // FIXME get this from internal package
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	ctx := signals.NewContext()

	slog.Info("Starting the server.")
	// Start the stdio server
	stdioServer := server.NewStdioServer(s)
	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output tekton-mcp string
	_, _ = fmt.Fprintf(os.Stderr, "Tekton MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		slog.Info("Shutting down server...")
	case err := <-errC:
		if err != nil {
			slog.Error(fmt.Sprintf("Error running server: %v", err))
			os.Exit(1)
		}
	}
}
