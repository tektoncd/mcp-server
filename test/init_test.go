//go:build e2e
// +build e2e

/*
Copyright 2024 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	knativetest "knative.dev/pkg/test"
)

func init() {
	flag.StringVar(&mcpServerURL, "mcp-server-url", "", "URL of the MCP server to test against")
	flag.BoolVar(&deployMCPServer, "deploy-mcp-server", true, "Whether to deploy the MCP server as part of the test")
}

// TestMain initializes the test environment
func TestMain(m *testing.M) {
	flag.Parse()
	c := m.Run()
	fmt.Fprintf(os.Stderr, "Using kubeconfig at `%s` with cluster `%s`\n", knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster)
	os.Exit(c)
}