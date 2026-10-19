package clientsecret

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type readBackend struct {
	mode          string
	service, user string
}

func (*readBackend) Set(string, string, string) error { return nil }
func (*readBackend) Delete(string, string) error      { return nil }
func (b *readBackend) Get(service, user string) (string, error) {
	b.service, b.user = service, user
	if b.mode == "block" {
		if os.WriteFile(user, []byte("ready"), 0600) != nil {
			return "", errors.New("fixture readiness failed")
		}
		for {
			time.Sleep(time.Hour)
		}
	}
	if b.mode == "error" {
		return "", errors.New("PRIVATE raw native failure")
	}
	return "private-fixture-token", nil
}
func TestReadHelperProtocolScopesNativeReadAndScrubsFailure(t *testing.T) {
	for _, mode := range []string{"success", "error"} {
		backend := &readBackend{mode: mode}
		var output bytes.Buffer
		code := serveReadHelper(strings.NewReader(`{"daemon_id":"daemon","credential_id":"credential"}`), &output, backend)
		if code != 0 || backend.service != servicePrefix+"daemon" || backend.user != "credential" {
			t.Fatalf("code=%d backend=%+v", code, backend)
		}
		if strings.Contains(output.String(), "PRIVATE") {
			t.Fatal("raw native failure disclosed")
		}
		if mode == "error" && strings.Contains(output.String(), "private-fixture-token") {
			t.Fatal("token disclosed on failure")
		}
	}
	if handled, _ := RunReadHelper([]string{"health"}, nil, nil); handled {
		t.Fatal("ordinary CLI intercepted")
	}
	for _, input := range []string{`{}`, `{"daemon_id":"d","credential_id":"c","unknown":true}`, `{"daemon_id":"d","credential_id":"c"} {}`} {
		if handled, code := RunReadHelper([]string{readHelperArgument}, strings.NewReader(input), io.Discard); !handled || code == 0 {
			t.Fatal("invalid private protocol accepted")
		}
	}
	if _, err := NewWithBackend(&readBackend{}).LoadContext(context.Background(), "d", "c"); err == nil {
		t.Fatal("blocking legacy backend silently used")
	}
}
func TestReadHelperProcess(t *testing.T) {
	mode := os.Args[len(os.Args)-1]
	if mode != "s107-read-success" && mode != "s107-read-block" && mode != "s107-read-error" {
		return
	}
	mode = strings.TrimPrefix(mode, "s107-read-")
	os.Exit(serveReadHelper(os.Stdin, os.Stdout, &readBackend{mode: mode}))
}
func TestReadCommandUsesPrivatePipesAndJoinsCanceledHelper(t *testing.T) {
	arguments := []string{"-test.run=^TestReadHelperProcess$", "--", "s107-read-success"}
	token, err := runReadCommand(context.Background(), os.Args[0], arguments, readRequest{DaemonID: "d", CredentialID: "c"})
	if err != nil || token != "private-fixture-token" {
		t.Fatalf("successful helper failed: %v", err)
	}
	arguments[2] = "s107-read-error"
	if token, err := runReadCommand(context.Background(), os.Args[0], arguments, readRequest{DaemonID: "d", CredentialID: "c"}); err == nil || token != "" || strings.Contains(err.Error(), "PRIVATE") {
		t.Fatal("helper failure was not scrubbed")
	}
	arguments[2] = "s107-read-block"
	ready := filepath.Join(t.TempDir(), "ready")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := runReadCommand(ctx, os.Args[0], arguments, readRequest{DaemonID: "d", CredentialID: ready})
		done <- err
	}()
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			cancel()
			<-done
			t.Fatal("private helper did not become ready")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("canceled native helper was not joined")
	}
}
func TestReadResponseLimitRejectsOversizedSecret(t *testing.T) {
	var buffer boundedResponse
	encoded, _ := json.Marshal(readResponse{Token: strings.Repeat("x", readResponseLimit)})
	if _, err := io.Copy(&buffer, bytes.NewReader(encoded)); err == nil || buffer.Len() != 0 {
		t.Fatal("oversized secret response accepted")
	}
}
