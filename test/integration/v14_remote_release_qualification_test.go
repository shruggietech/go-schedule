package integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/platform"
	"github.com/shruggietech/go-schedule/internal/remoteenroll"
)

func TestV14RemoteReleaseJourney(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	packageDir := filepath.Join(t.TempDir(), "v1.4 Package Shape ü")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	candidate := buildV14Daemon(t, root, packageDir, "goschedd-candidate")
	replacement := buildV14Daemon(t, root, packageDir, "goschedd-replacement")

	stateDir := filepath.Join(t.TempDir(), "retained remote state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	endpoint := v14IPCEndpoint(t)
	address := v14AvailableAddress(t)
	configPath := filepath.Join(stateDir, "daemon-config.json")
	cfg := config.Default()
	cfg.DataDir = stateDir
	cfg.IPCPath = endpoint
	cfg.AdminGroup = v13AccessibleAdminGroup(t)
	writeV14Config(t, configPath, cfg)

	local, output, stop := startV14Daemon(t, candidate, configPath, endpoint)
	manifest, err := local.Manifest(context.Background())
	if err != nil {
		stop()
		t.Fatal(err)
	}
	if connection, err := net.DialTimeout("tcp", address, 250*time.Millisecond); err == nil {
		_ = connection.Close()
		stop()
		t.Fatal("default configuration unexpectedly exposed the remote listener")
	}
	stop()
	if strings.Contains(output.String(), "remote HTTPS listener ready") {
		t.Fatal("default start reported a remote listener")
	}

	certificatePath, keyPath, certificatePEM := writeV14Certificate(t, stateDir)
	cfg.Remote = config.RemoteConfig{Enabled: true, BindAddress: address, CertificateFile: certificatePath, PrivateKeyFile: keyPath}
	writeV14Config(t, configPath, cfg)
	local, output, stop = startV14Daemon(t, candidate, configPath, endpoint)
	defer stop()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := waitForV14Remote(ctx, address); err != nil {
		t.Fatalf("remote listener did not become reachable: %v\n%s", err, output.String())
	}

	manageSecret, err := local.CreatePairing(ctx, server.PairingCreateRequest{DisplayName: "S080 manager", Kind: domain.ActorKindCLI, Capability: domain.CapabilityManage})
	if err != nil {
		t.Fatal(err)
	}
	managerCredential, pairedManifest, err := remoteenroll.Exchange(ctx, remoteenroll.Draft{Address: "https://" + address, DaemonID: manageSecret.DaemonID, PairingID: manageSecret.ID, Phrase: manageSecret.Phrase, CertificatePEM: certificatePEM, DisplayName: manageSecret.DisplayName, Kind: manageSecret.Kind, Capability: manageSecret.Capability})
	if err != nil || pairedManifest.InstallationID != manifest.InstallationID {
		t.Fatalf("manager pairing manifest=%+v err=%v", pairedManifest, err)
	}
	manager, err := client.NewRemote("https://"+address, certificatePEM, managerCredential.Token, manifest.InstallationID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manager.VerifyIdentity(ctx); err != nil {
		t.Fatal(err)
	}
	created, err := manager.CreateTask(ctx, server.TaskCreateRequest{Name: "S080 retained task"})
	if err != nil || created.Task.Name != "S080 retained task" {
		t.Fatalf("created=%+v err=%v", created, err)
	}

	observeSecret, err := local.CreatePairing(ctx, server.PairingCreateRequest{DisplayName: "S080 observer", Kind: domain.ActorKindDesktop, Capability: domain.CapabilityObserve})
	if err != nil {
		t.Fatal(err)
	}
	observerCredential, _, err := remoteenroll.Exchange(ctx, remoteenroll.Draft{Address: "https://" + address, DaemonID: observeSecret.DaemonID, PairingID: observeSecret.ID, Phrase: observeSecret.Phrase, CertificatePEM: certificatePEM, DisplayName: observeSecret.DisplayName, Kind: observeSecret.Kind, Capability: observeSecret.Capability})
	if err != nil {
		t.Fatal(err)
	}
	observer, err := client.NewRemote("https://"+address, certificatePEM, observerCredential.Token, manifest.InstallationID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := observer.VerifyIdentity(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := observer.CreateTask(ctx, server.TaskCreateRequest{Name: "must be denied"}); err == nil {
		t.Fatal("observer mutation succeeded")
	} else {
		var status *client.StatusError
		if !errors.As(err, &status) || status.Code != server.CodeForbidden {
			t.Fatalf("observer mutation err=%v", err)
		}
	}
	events, err := local.ListAudit(ctx, domain.AuditQuery{Operation: "tasks.create", Limit: 10})
	if err != nil || len(events) != 2 {
		t.Fatalf("task creation audits=%+v err=%v", events, err)
	}
	results := map[domain.AuditResult]string{}
	for _, event := range events {
		results[event.Result] = event.ActorID
	}
	if results[domain.AuditResultSucceeded] != managerCredential.Actor.ID || results[domain.AuditResultDenied] != observerCredential.Actor.ID {
		t.Fatalf("audit attribution=%+v", results)
	}
	if _, err := local.RevokeCredential(ctx, managerCredential.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.ListTasks(ctx, "", ""); err == nil || !strings.Contains(err.Error(), "credential revoked") {
		t.Fatalf("revoked request err=%v", err)
	}
	if _, err := local.Health(ctx); err != nil {
		t.Fatalf("local IPC after revocation: %v", err)
	}
	stop()

	cfg.Remote = config.RemoteConfig{}
	writeV14Config(t, configPath, cfg)
	local, disabledOutput, stopDisabled := startV14Daemon(t, replacement, configPath, endpoint)
	defer stopDisabled()
	restored, err := local.Manifest(ctx)
	if err != nil || restored.InstallationID != manifest.InstallationID {
		t.Fatalf("replacement manifest=%+v err=%v", restored, err)
	}
	tasks, err := local.ListTasks(ctx, "", "")
	if err != nil || len(tasks) != 1 || tasks[0].Name != "S080 retained task" {
		t.Fatalf("replacement tasks=%+v err=%v", tasks, err)
	}
	if connection, err := net.DialTimeout("tcp", address, 250*time.Millisecond); err == nil {
		_ = connection.Close()
		t.Fatal("disabled replacement unexpectedly exposed the remote listener")
	}
	stopDisabled()
	for name, data := range map[string]string{"enabled output": output.String(), "disabled output": disabledOutput.String()} {
		for _, secret := range []string{manageSecret.Phrase, managerCredential.Token, observeSecret.Phrase, observerCredential.Token} {
			if secret != "" && strings.Contains(data, secret) {
				t.Fatalf("%s disclosed protected value", name)
			}
		}
	}
}

func TestV14RemoteReleaseDocumentationContract(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	documents := map[string][]string{
		"docs/remote-access.md": {
			"## Operator runbook",
			"Private-network HTTPS is the recommended direct deployment.",
			"SSH-tunneled HTTPS is the recommended path",
			"Direct public HTTPS is an advanced supported deployment, never a zero-configuration recommendation.",
			"certificate issuance and renewal",
			"gosched service install --config",
			"gosched profile pair",
			"gosched credential revoke",
			"To disable network access",
			"For an upgrade",
			"Back up the daemon database, configuration, certificate, and private key",
			"A mutation that loses its response is never replayed automatically",
		},
		"docs/cli.md":             {"install [--config FILE]", "records exactly `--config <absolute-path>`"},
		"docs/INSTALL-linux.md":   {"/var/lib/goschedule/config.json", "installation alone never opens one"},
		"docs/INSTALL-macos.md":   {"/Library/Application Support/goschedule/config.json", "installation alone never opens one"},
		"docs/INSTALL-windows.md": {"C:\\ProgramData\\goschedule\\config.json", "MSI installation and upgrade never create a remote configuration or open a network listener"},
	}
	for relative, required := range documents {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		for _, text := range required {
			if !bytes.Contains(data, []byte(text)) {
				t.Errorf("%s missing %q", relative, text)
			}
		}
	}
}

func buildV14Daemon(t *testing.T, root, directory, name string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(directory, name)
	command := exec.Command("go", "build", "-o", path, "./cmd/goschedd")
	command.Dir = root
	platform.HideConsole(command)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build %s: %v\n%s", name, err, output)
	}
	return path
}

func v14IPCEndpoint(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(`\\.\pipe\goschedd-s080-%d-%d`, os.Getpid(), time.Now().UnixNano())
	}
	directory := filepath.Join(os.TempDir(), fmt.Sprintf("s080-%d-%d", os.Getpid(), time.Now().UnixNano()))
	t.Cleanup(func() { _ = os.RemoveAll(directory) })
	return filepath.Join(directory, "goschedd.sock")
}

func v14AvailableAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	return address
}

func writeV14Config(t *testing.T, path string, cfg config.Config) {
	t.Helper()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeV14Certificate(t *testing.T, directory string) (string, string, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 120))
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{SerialNumber: serial, Subject: pkix.Name{CommonName: "S080 loopback"}, NotBefore: time.Now().Add(-time.Minute), NotAfter: time.Now().Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	certificatePath := filepath.Join(directory, "server.crt")
	keyPath := filepath.Join(directory, "server.key")
	if err := os.WriteFile(certificatePath, certificate, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyPath, key, 0o600); err != nil {
		t.Fatal(err)
	}
	return certificatePath, keyPath, string(certificate)
}

func startV14Daemon(t *testing.T, binary, configPath, endpoint string) (*client.Client, *bytes.Buffer, func()) {
	t.Helper()
	var output bytes.Buffer
	command := exec.Command(binary, "--config", configPath)
	command.Stdout = &output
	command.Stderr = &output
	platform.HideConsole(command)
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	stopped := false
	stop := func() {
		if stopped {
			return
		}
		stopped = true
		if command.Process != nil {
			_ = command.Process.Kill()
			_ = command.Wait()
		}
	}
	local := client.New(endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := waitForV13Health(ctx, local); err != nil {
		stop()
		t.Fatalf("daemon never became healthy: %v\n%s", err, output.String())
	}
	return local, &output, stop
}

func waitForV14Remote(ctx context.Context, address string) error {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		connection, err := net.DialTimeout("tcp", address, 250*time.Millisecond)
		if err == nil {
			_ = connection.Close()
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
