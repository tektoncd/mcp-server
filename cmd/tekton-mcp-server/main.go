package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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
// adding test
const ManagedByLabelKey = "app.kubernetes.io/managed-by"

func main() {
	var transport string
	var sseAddr string
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&sseAddr, "address", "", "Address to bind the SSE server to")
	flag.Parse()

	if sseAddr == "" && transport == "sse" {
		slog.Error("-address is required when transport is set to 'sse'")
		os.Exit(1)
	}

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
	tools.Add(ctx, s)
	resources.Add(ctx, s)

	slog.Info("Starting the server.")

	errC := make(chan error, 1)

	switch transport {
	case "sse":
		sseServer := server.NewSSEServer(s, server.WithSSEContextFunc(func(_ context.Context, r *http.Request) context.Context { return ctx }))
		go func() {
			errC <- sseServer.Start(sseAddr)
		}()
		slog.Info("Tekton MCP Server is listening at " + sseAddr)
	case "stdio":
		stdioServer := server.NewStdioServer(s)
		go func() {
			in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)
			errC <- stdioServer.Listen(ctx, in, out)
		}()
		_, _ = fmt.Fprintf(os.Stderr, "Tekton MCP Server running on stdio\n")
	default:
		slog.Error(fmt.Sprintf("Invalid transport %q; must be sse or stdio", transport))
		os.Exit(1)
	}

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
