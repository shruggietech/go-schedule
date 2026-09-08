package mcpobserve

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestServerDiscoveryAndEveryApprovedResource(t *testing.T) {
	ctx := context.Background()
	reader := &fakeReader{health: server.HealthResponse{Status: "ok", Version: "v1.3.0"}}
	for i := 0; i < 101; i++ {
		reader.tasks = append(reader.tasks, server.TaskObservationResponse{ID: fmt.Sprintf("task-%03d", i), State: domain.TaskActive})
	}
	observeServer := NewServer(reader, "test")
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := observeServer.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	observeClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, &mcp.ClientOptions{Capabilities: &mcp.ClientCapabilities{}})
	clientSession, err := observeClient.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()
	if got := clientSession.InitializeResult().ProtocolVersion; got != "2026-07-28" {
		t.Fatalf("protocol version = %q", got)
	}
	if clientSession.InitializeResult().Capabilities.Tools != nil {
		t.Fatal("server advertised tools")
	}
	resources, err := clientSession.ListResources(ctx, nil)
	if err != nil || len(resources.Resources) != 5 {
		t.Fatalf("resources = %+v, %v", resources, err)
	}
	templates, err := clientSession.ListResourceTemplates(ctx, nil)
	if err != nil || len(templates.ResourceTemplates) != 4 {
		t.Fatalf("templates = %+v, %v", templates, err)
	}
	for _, resource := range resources.Resources {
		result, readErr := clientSession.ReadResource(ctx, &mcp.ReadResourceParams{URI: resource.URI})
		if readErr != nil || len(result.Contents) != 1 {
			t.Fatalf("read %s = %+v, %v", resource.URI, result, readErr)
		}
		if resource.URI == tasksURI {
			var first Envelope
			if err := json.Unmarshal([]byte(result.Contents[0].Text), &first); err != nil || first.Page.NextURI == "" {
				t.Fatalf("task first page = %+v, %v", first, err)
			}
			continued, err := clientSession.ReadResource(ctx, &mcp.ReadResourceParams{URI: first.Page.NextURI})
			if err != nil || len(continued.Contents) != 1 {
				t.Fatalf("continued tasks = %+v, %v", continued, err)
			}
			var next Envelope
			if err := json.Unmarshal([]byte(continued.Contents[0].Text), &next); err != nil || next.Page.Count != 1 {
				t.Fatalf("task continuation = %+v, %v", next, err)
			}
		}
	}
}
