package cli

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/platform"
)

func TestMCPHTTPCommandRegistersCompleteRuntimeSurface(t *testing.T) {
	root := newMCPHTTPCmd()
	want := map[string]bool{"status": true, "enable": true, "rotate": true, "disable": true}
	for _, child := range root.Commands() {
		delete(want, child.Name())
	}
	if len(want) != 0 {
		t.Fatalf("missing HTTP commands: %v", want)
	}
	enable := newMCPHTTPEnableCmd()
	enable.SetArgs([]string{"--port", "0"})
	if err := enable.Execute(); err == nil {
		t.Fatal("invalid port succeeded")
	}
}

func TestMCPHTTPCredentialOutputIsExplicitlyOneTime(t *testing.T) {
	previousJSON := jsonOut
	jsonOut = false
	t.Cleanup(func() { jsonOut = previousJSON })
	previousStdout := os.Stdout
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = write
	t.Cleanup(func() { os.Stdout = previousStdout })
	result := server.MCPHTTPCredentialResponse{MCPHTTPStatusResponse: server.MCPHTTPStatusResponse{Enabled: true, Endpoint: "http://127.0.0.1:43123/mcp", AllowedOrigins: []string{}}, Credential: "one-time-secret"}
	if err := printMCPHTTPCredential(result, "localhost MCP enabled"); err != nil {
		t.Fatal(err)
	}
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	if err := read.Close(); err != nil {
		t.Fatal(err)
	}
	if got := string(output); !strings.Contains(got, "credential (shown once): one-time-secret") || !strings.Contains(got, result.Endpoint) {
		t.Fatalf("output=%q", got)
	}
}

func TestMCPStdioSubprocessHasCleanProtocolAndStopsOnDisconnect(t *testing.T) {
	if os.Getenv("GO_SCHEDULE_MCP_HELPER") == "1" {
		os.Args = []string{"gosched", "mcp", "serve"}
		os.Exit(Execute())
	}
	command := exec.Command(os.Args[0], "-test.run=TestMCPStdioSubprocessHasCleanProtocolAndStopsOnDisconnect")
	platform.HideConsole(command)
	command.Env = append(os.Environ(), "GO_SCHEDULE_MCP_HELPER=1")
	stdin, err := command.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	var stderr strings.Builder
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	request := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}` + "\n"
	if _, err := io.WriteString(stdin, request); err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatalf("no initialize response: %v, stderr=%q", scanner.Err(), stderr.String())
	}
	var response map[string]any
	if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
		t.Fatalf("stdout is not protocol JSON: %q", scanner.Text())
	}
	result, _ := response["result"].(map[string]any)
	capabilities, _ := result["capabilities"].(map[string]any)
	if _, found := capabilities["tools"]; found {
		t.Fatalf("initialize advertised tools: %s", scanner.Text())
	}
	if err := stdin.Close(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("helper exit = %v, stderr=%q", err, stderr.String())
		}
	case <-time.After(5 * time.Second):
		_ = command.Process.Kill()
		t.Fatal("MCP subprocess did not stop after host disconnect")
	}
}
