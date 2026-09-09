package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/platform"
)

func TestPackagedMCPCommandDiscovery(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	packageDir := filepath.Join(t.TempDir(), "Package Shape ü")
	binary := filepath.Join(packageDir, "gosched")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	build := exec.Command("go", "build", "-o", binary, "./cmd/gosched")
	build.Dir = root
	platform.HideConsole(build)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build packaged command: %v\n%s", err, output)
	}
	command := exec.Command(binary, "mcp", "serve")
	platform.HideConsole(command)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
