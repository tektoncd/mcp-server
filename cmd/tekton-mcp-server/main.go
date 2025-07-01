package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tektoncd/mcp-server/internal/resources"
	"github.com/tektoncd/mcp-server/internal/tools"
	"github.com/tektoncd/mcp-server/internal/version"
	"k8s.io/client-go/tools/clientcmd"
	filteredinformerfactory "knative.dev/pkg/client/injection/kube/informers/factory/filtered"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"
)

// ManagedByLabelKey is the label key used to mark what is managing this resource
const ManagedByLabelKey = "app.kubernetes.io/managed-by"

func main() {
	var transport string
	var httpAddr string
	flag.StringVar(&transport, "transport", "http", "Transport type (stdio or http)")
	flag.StringVar(&httpAddr, "address", ":3000", "Address to bind the HTTP server to")
	flag.Parse()

	if httpAddr == "" && transport == "http" {
		slog.Error("-address is required when transport is set to 'http'")
		os.Exit(1)
	}

	// Create MCP server
	s := mcp.NewServer("Tekton", version.Version, nil)

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

	if err = tools.Add(ctx, s); err != nil {
		slog.Error(fmt.Sprintf("unable to add tools: %v", err))
		os.Exit(1)
	}

	resources.Add(ctx, s)

	slog.Info("Starting the server.")

	errC := make(chan error, 1)

	switch transport {
	case "http":

		streamableHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server { return s }, nil)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			streamableHandler.ServeHTTP(w, r.WithContext(ctx))
		})
		server := &http.Server{
			Addr:              httpAddr,
			Handler:           handler,
			ReadHeaderTimeout: 3 * time.Second,
		}

		go func() {
			errC <- server.ListenAndServe()
		}()
		slog.Info("Tekton MCP Server is listening at " + httpAddr)
	case "stdio":
		go func() {
			errC <- s.Run(ctx, mcp.NewStdioTransport())
		}()
		_, _ = fmt.Fprintf(os.Stderr, "Tekton MCP Server running on stdio\n")
	default:
		slog.Error(fmt.Sprintf("Invalid transport %q; must be http or stdio", transport))
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
