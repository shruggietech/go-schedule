package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/platform"
)

func TestPackagedMCPCommandDiscovery(t *testing.T) {
	packageDir := filepath.Join(t.TempDir(), "Package Shape ü")
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		t.Fatal(err)
	}
	binary := preparePackagedCandidate(t, packageDir, "gosched", "gosched")
	command := exec.Command(binary, "mcp", "serve")
	platform.HideConsole(command)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "package-smoke", Version: "test"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: command}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	resources, err := session.ListResources(ctx, nil)
	if err != nil || len(resources.Resources) != 5 {
		t.Fatalf("resources=%+v err=%v", resources, err)
	}
	templates, err := session.ListResourceTemplates(ctx, nil)
	if err != nil || len(templates.ResourceTemplates) != 4 {
		t.Fatalf("templates=%+v err=%v", templates, err)
	}
	tools, err := session.ListTools(ctx, nil)
	if err != nil || len(tools.Tools) != 0 {
		t.Fatalf("tools=%+v err=%v", tools, err)
	}
}
