package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/platform"
)

func TestV13PackageDefaultsRemainOptIn(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	candidateDir := filepath.Join(t.TempDir(), "v1.3 Package Shape ü")
	if err := os.MkdirAll(candidateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(candidateDir, "goschedd")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command("go", "build", "-o", binary, "./cmd/goschedd")
	build.Dir = root
	platform.HideConsole(build)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build package-shaped daemon: %v\n%s", err, output)
	}

	stateDir := filepath.Join(t.TempDir(), "retained state")
	cfg := config.Default()
	cfg.DataDir = stateDir
	cfg.AdminGroup = v13AccessibleAdminGroup(t)
	ipcDir := filepath.Join("/tmp", fmt.Sprintf("s072-%d-%d", os.Getpid(), time.Now().UnixNano()))
	cfg.IPCPath = filepath.Join(ipcDir, "goschedd.sock")
	if runtime.GOOS == "windows" {
		cfg.IPCPath = fmt.Sprintf(`\\.\pipe\goschedd-s072-%d-%d`, os.Getpid(), time.Now().UnixNano())
	} else {
		t.Cleanup(func() { _ = os.RemoveAll(ipcDir) })
	}
	configPath := filepath.Join(stateDir, "config.json")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	for start := 1; start <= 2; start++ {
		runV13CandidateAndInspect(t, binary, configPath, cfg.IPCPath, start)
	}
}

func v13AccessibleAdminGroup(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		return "Users"
	}
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("resolve current user for restricted IPC: %v", err)
	}
	group, err := user.LookupGroupId(currentUser.Gid)
	if err != nil {
		t.Fatalf("resolve current group for restricted IPC: %v", err)
	}
	if group.Name == "" {
		t.Fatal("current group has no name for restricted IPC")
	}
	return group.Name
}

func runV13CandidateAndInspect(t *testing.T, binary, configPath, endpoint string, start int) {
	t.Helper()
	var output bytes.Buffer
	command := exec.Command(binary, "--config", configPath)
	command.Stdout = &output
	command.Stderr = &output
	platform.HideConsole(command)
	if err := command.Start(); err != nil {
		t.Fatalf("start package-shaped daemon (launch %d): %v", start, err)
	}
	stopped := false
	defer func() {
		if !stopped && command.Process != nil {
			_ = command.Process.Kill()
			_ = command.Wait()
		}
	}()

	api := client.New(endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := waitForV13Health(ctx, api); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		stopped = true
		t.Fatalf("package-shaped daemon never became healthy (launch %d): %v\n%s", start, err, output.String())
	}

	channels, err := api.ListNotificationChannels(ctx)
	if err != nil || len(channels) != 0 {
		t.Fatalf("launch %d notification channels=%+v err=%v", start, channels, err)
	}
	deliveries, err := api.ListNotificationDeliveries(ctx, domain.NotificationDeliveryFilter{})
	if err != nil || len(deliveries) != 0 {
		t.Fatalf("launch %d notification deliveries=%+v err=%v", start, deliveries, err)
	}
	status, err := api.MCPHTTPStatus(ctx)
	if err != nil {
		t.Fatalf("launch %d MCP HTTP status: %v", start, err)
	}
	if status.Enabled || status.Endpoint != "" || len(status.AllowedOrigins) != 0 || status.CredentialFingerprint != "" || status.EnabledAt != nil || status.ClientName != "" || status.LastAccessedAt != nil || status.RequestCount != 0 {
		t.Fatalf("launch %d exposed MCP HTTP state before enablement: %+v", start, status)
	}

	if start == 1 {
		created, err := api.CreateTask(ctx, server.TaskCreateRequest{Name: "S072 offline draft"})
		if err != nil {
			t.Fatalf("create task through local IPC: %v", err)
		}
		if created.Task.Name != "S072 offline draft" {
			t.Fatalf("created task=%+v", created.Task)
		}
	}
	tasks, err := api.ListTasks(ctx, "", "")
	if err != nil || len(tasks) != 1 || tasks[0].Name != "S072 offline draft" {
		t.Fatalf("launch %d retained tasks=%+v err=%v", start, tasks, err)
	}

	if err := command.Process.Kill(); err != nil {
		t.Fatalf("stop package-shaped daemon (launch %d): %v", start, err)
	}
	if err := command.Wait(); err == nil {
		t.Fatalf("launch %d unexpectedly reported a clean exit after forced stop", start)
	}
	stopped = true
	if !bytes.Contains(output.Bytes(), []byte(`"access_mode":"restricted"`)) {
		t.Fatalf("launch %d did not report restricted IPC: %s", start, output.String())
	}
}

func waitForV13Health(ctx context.Context, api *client.Client) error {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		attemptCtx, cancel := context.WithTimeout(ctx, time.Second)
		_, err := api.Health(attemptCtx)
		cancel()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
