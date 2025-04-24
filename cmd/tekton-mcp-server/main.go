package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tektoncd/mcp-server/internal/resources"
	"github.com/tektoncd/mcp-server/internal/tools"
	"k8s.io/client-go/tools/clientcmd"
	filteredinformerfactory "knative.dev/pkg/client/injection/kube/informers/factory/filtered"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"
)

// ManagedByLabelKey is the label key used to mark what is managing this resource
const ManagedByLabelKey = "app.kubernetes.io/managed-by"

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

	// Load kubernetes configuration
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	cfg, err := kubeConfig.ClientConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get Kubernetes config: %v", err))
		os.Exit(1)
	}

	// Start informers through knative injection functions (in context)
	ctx = filteredinformerfactory.WithSelectors(ctx, ManagedByLabelKey)
	// slog.Info("Registering %d informer factories", len(injection.Default.GetInformerFactories()))
	// slog.Info("Registering %d informers", len(injection.Default.GetInformers()))
	ctx, startInformers := injection.EnableInjectionOrDie(ctx, cfg)

	// Start the injection clients and informers.
	startInformers()

	slog.Info("Adding tools and resources to the server.")
	tools.Add(s)
	resources.Add(s)

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
