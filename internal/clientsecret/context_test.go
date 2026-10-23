package clientsecret

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"

	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/platform"
)

type readBackend struct {
	mode          string
	service, user string
	calls         atomic.Int32
}

func (*readBackend) Set(string, string, string) error { return nil }
func (*readBackend) Delete(string, string) error      { return nil }
func (b *readBackend) Get(service, user string) (string, error) {
	b.calls.Add(1)
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
		return "private-fixture-token", errors.New("PRIVATE raw native failure")
	}
	if b.mode == "large" {
		return strings.Repeat("x", readResponseLimit+1), nil
	}
	return "private-fixture-token", nil
}
func framed(t *testing.T, input string) *bytes.Buffer {
	t.Helper()
	var buffer bytes.Buffer
	if err := writeFrame(&buffer, []byte(input)); err != nil {
		t.Fatal(err)
	}
	return &buffer
}
func TestReadHelperProtocolScopesNativeReadAndScrubsFailure(t *testing.T) {
	for _, mode := range []string{"success", "error"} {
		backend := &readBackend{mode: mode}
		var output bytes.Buffer
		code := serveReadHelper(framed(t, `{"daemon_id":"daemon","credential_id":"credential"}`), &output, backend)
		if code != 0 || backend.service != servicePrefix+"daemon" || backend.user != "credential" {
			t.Fatalf("code=%d backend=%+v", code, backend)
		}
		data, err := readFrame(&output, readResponseLimit)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte("PRIVATE")) || mode == "error" && bytes.Contains(data, []byte("private-fixture-token")) {
			t.Fatal("native failure disclosed token or raw error")
		}
	}
	for _, input := range []string{`{}`, `{"daemon_id":"d","credential_id":"c","unknown":true}`, `{"daemon_id":"d","credential_id":"c"} {}`} {
		backend := &readBackend{}
		if serveReadHelper(framed(t, input), io.Discard, backend) == 0 || backend.calls.Load() != 0 {
			t.Fatal("invalid private request accepted")
		}
	}
	if _, err := NewWithBackend(&readBackend{}).LoadContext(context.Background(), "d", "c"); err == nil {
		t.Fatal("blocking legacy backend silently used")
	}
}
func TestReadHelperRejectsDirectEntryAndUntrustedLivePeer(t *testing.T) {
	if handled, _ := RunReadHelper([]string{"health"}); handled {
		t.Fatal("ordinary CLI intercepted")
	}
	if handled, code := RunReadHelper([]string{readHelperArgument}); !handled || code == 0 {
		t.Fatal("direct standard-stream entry accepted")
	}
	for _, pid := range []int{0, -1, os.Getpid()} {
		if authenticateReadParent(pid) == nil {
			t.Fatalf("untrusted parent PID %d accepted", pid)
		}
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if sameReadExecutable(executable, executable) != nil {
		t.Fatal("own executable identity rejected")
	}
	unrelated := filepath.Join(t.TempDir(), filepath.Base(executable))
	if err := os.WriteFile(unrelated, []byte("unrelated process image"), 0600); err != nil {
		t.Fatal(err)
	}
	if sameReadExecutable(executable, unrelated) == nil {
		t.Fatal("same-name unrelated executable authenticated")
	}
	listener, endpoint, closeChannel, err := listenReadChannel()
	if err != nil {
		t.Fatal(err)
	}
	defer closeChannel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	// Windows may discard a pipe client that closes before Accept completes.
	// Establish both real endpoints before running peer authentication so an
	// immediate, correct rejection cannot race the fixture setup.
	type dialResult struct {
		connection net.Conn
		err        error
	}
	dialed := make(chan dialResult, 1)
	go func() {
		connection, err := dialReadChannel(ctx, endpoint)
		dialed <- dialResult{connection, err}
	}()
	connection, acceptErr := listener.Accept()
	peer := <-dialed
	if peer.connection != nil {
		defer func() { _ = peer.connection.Close() }()
	}
	if connection != nil {
		defer func() { _ = connection.Close() }()
	}
	if acceptErr != nil || peer.err != nil {
		t.Fatalf("establish native peer fixture: accept=%v dial=%v", acceptErr, peer.err)
	}
	backend := &readBackend{}
	done := make(chan int, 1)
	go func() { done <- serveAuthenticatedConnection(ctx, peer.connection, backend) }()
	// This endpoint belongs to the test itself, not its actual parent. Even valid
	// request metadata must never cause a credential read from this forged peer.
	_ = writeFrame(connection, []byte(`{"daemon_id":"d","credential_id":"c"}`))
	_ = connection.Close()
	if code := <-done; code == 0 || backend.calls.Load() != 0 {
		t.Fatal("forged parent channel reached credential storage")
	}
}
func TestReadChannelRejectsOtherClientsBeforeSendingAnyRequest(t *testing.T) {
	listener, endpoint, closeChannel, err := listenReadChannel()
	if err != nil {
		t.Fatal(err)
	}
	defer closeChannel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	done := make(chan error, 1)
	go func() { _, err := acceptReadHelper(ctx, listener, os.Getpid()+1000000); done <- err }()
	connection, err := dialReadChannel(ctx, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	var data [1]byte
	n, _ := connection.Read(data[:])
	_ = connection.Close()
	cancel()
	if n != 0 {
		t.Fatal("request sent to an unrelated client")
	}
	if <-done == nil {
		t.Fatal("unrelated client accepted")
	}
}
func TestReadHelperProcess(t *testing.T) {
	if len(os.Args) < 3 {
		return
	}
	mode := os.Args[len(os.Args)-2]
	if !strings.HasPrefix(mode, "s107-read-") {
		return
	}
	if mode == "s107-read-expire" {
		os.Exit(runAuthenticatedHelper(os.Args[len(os.Args)-1], &readBackend{mode: "block"}, 500*time.Millisecond))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	os.Exit(serveAuthenticatedHelper(ctx, os.Args[len(os.Args)-1], &readBackend{mode: strings.TrimPrefix(mode, "s107-read-")}))
}
func TestReadCommandUsesAuthenticatedChannelAndJoinsCanceledHelper(t *testing.T) {
	arguments := []string{"-test.run=^TestReadHelperProcess$", "--", "s107-read-success"}
	token, err := runReadCommand(context.Background(), os.Args[0], arguments, readRequest{DaemonID: "d", CredentialID: "c"})
	if err != nil || token != "private-fixture-token" {
		t.Fatalf("authenticated helper failed: %v", err)
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
			t.Fatal("authenticated helper did not become ready")
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
func TestReadFrameRejectsOversizedOrIncompleteInput(t *testing.T) {
	var input bytes.Buffer
	if err := binary.Write(&input, binary.BigEndian, uint32(readResponseLimit+1)); err != nil {
		t.Fatal(err)
	}
	if _, err := readFrame(&input, readResponseLimit); err == nil {
		t.Fatal("oversized frame accepted")
	}
	if _, err := readFrame(strings.NewReader("short"), readResponseLimit); err == nil {
		t.Fatal("incomplete frame accepted")
	}
	var output bytes.Buffer
	if serveReadHelper(framed(t, `{"daemon_id":"d","credential_id":"c"}`), &output, &readBackend{mode: "large"}) == 0 || output.Len() != 0 {
		t.Fatal("oversized response disclosed")
	}
}

func TestPackagedNativeHelperRejectsUnrelatedParentAndDirectStreams(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	executable := filepath.Join(t.TempDir(), "gosched"+extension)
	buildCtx, buildCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	err = buildPackagedCLI(buildCtx, root, executable)
	buildCancel()
	if err != nil {
		t.Fatal(err)
	}
	// Compilation is setup, and must not consume the authentication budget.
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	var stdout, stderr bytes.Buffer
	direct := exec.CommandContext(ctx, executable, readHelperArgument)
	platform.HideConsole(direct)
	direct.Stdout = &stdout
	direct.Stderr = &stderr
	if direct.Run() == nil || stdout.Len() != 0 {
		t.Fatal("direct private flag accepted standard-stream credential requests")
	}
	listener, endpoint, closeChannel, err := listenReadChannel()
	if err != nil {
		t.Fatal(err)
	}
	defer closeChannel()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	command := exec.CommandContext(ctx, executable, readHelperArgument, endpoint)
	platform.HideConsole(command)
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	exited := make(chan struct{})
	go func() { defer close(exited); finished <- command.Wait(); _ = listener.Close() }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-exited:
		case <-time.After(2 * time.Second):
			t.Error("owned rejection helper did not exit")
		}
	})
	connection, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var probe [1]byte
	n, readErr := connection.Read(probe[:])
	_ = connection.Close()
	if timeout, ok := readErr.(interface{ Timeout() bool }); ok && timeout.Timeout() {
		t.Fatal("unrelated process image reached request reading before authentication")
	}
	if err := <-finished; err == nil || n != 0 || stdout.Len() != 0 {
		t.Fatal("unrelated parent or standard streams received a credential response")
	}
}

func buildPackagedCLI(ctx context.Context, root, executable string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("build isolated CLI setup: %w", err)
	}
	goPath, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("build isolated CLI setup: find Go: %w", err)
	}
	stdout, stderr := executable+".build.stdout", executable+".build.stderr"
	build, err := clifixture.StartProcess(clifixture.Command{Executable: goPath, Arguments: []string{"build", "-o", executable, "./cmd/gosched"}, Directory: root, StdoutPath: stdout, StderrPath: stderr, Environment: os.Environ()})
	if err != nil {
		return fmt.Errorf("build isolated CLI setup: start owned build: %w", err)
	}
	runErr := build.Wait(ctx)
	cleanupErr := build.Stop()
	output, outputErr := os.ReadFile(stdout)
	diagnostic, diagnosticErr := os.ReadFile(stderr)
	if err := errors.Join(runErr, cleanupErr, outputErr, diagnosticErr, ctx.Err()); err != nil {
		return fmt.Errorf("build isolated CLI setup: %w: %s %s", err, output, diagnostic)
	}
	return nil
}

func TestPackagedCLIBuildReportsSetupCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	executable := filepath.Join(t.TempDir(), "unused")
	err := buildPackagedCLI(ctx, t.TempDir(), executable)
	if !errors.Is(err, context.Canceled) || !strings.Contains(err.Error(), "CLI setup") {
		t.Fatalf("setup cancellation was not distinct: %v", err)
	}
	if _, err := os.Stat(executable); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("canceled build unexpectedly produced a fixture: %v", err)
	}
}

func TestAuthenticatedHelperDeadlineExpiresWithoutParentCancellation(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := runReadCommand(ctx, os.Args[0], []string{"-test.run=^TestReadHelperProcess$", "--", "s107-read-expire"}, readRequest{DaemonID: "d", CredentialID: ready})
	if err == nil || ctx.Err() != nil {
		t.Fatalf("helper did not enforce its own lifetime: %v", err)
	}
	if _, err := os.Stat(ready); err != nil {
		t.Fatal("helper did not enter the fixture read")
	}
}
