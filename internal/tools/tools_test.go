package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tektoncd/mcp-server/internal/resources"
	"github.com/tektoncd/mcp-server/internal/version"
)

func newSession(t *testing.T, ctx context.Context) (*mcp.ServerSession, *mcp.ClientSession) {
	t.Helper()

	ct, st := mcp.NewInMemoryTransports()
	s := mcp.NewServer("Tekton", version.Version, nil)
	if err := Add(ctx, s); err != nil {
		t.Fatal(err)
	}
	resources.Add(ctx, s)
	c := mcp.NewClient("TektonClient", version.Version, nil)

	ss, err := s.Connect(ctx, st)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := c.Connect(ctx, ct)
	if err != nil {
		t.Fatal(err)
	}

	return ss, cs
}
