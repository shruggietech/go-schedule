package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

type daemonIdentityClientFake struct {
	manifest server.ManifestResponse
	rename   string
	reset    string
	err      error
}

func (f *daemonIdentityClientFake) Manifest(context.Context) (server.ManifestResponse, error) {
	return f.manifest, f.err
}
func (f *daemonIdentityClientFake) RenameDaemon(_ context.Context, name string) (server.ManifestResponse, error) {
	f.rename = name
	return f.manifest, f.err
}
func (f *daemonIdentityClientFake) ResetDaemonIdentity(_ context.Context, confirmation string) (server.ManifestResponse, error) {
	f.reset = confirmation
	return f.manifest, f.err
}

func TestDaemonCommandsExposeManifestRenameAndReset(t *testing.T) {
	manifest := server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Workshop", ProductVersion: "1.2.0", OperatingMode: "local_only", Capabilities: []string{"tasks"}, Platform: server.ManifestPlatform{OS: "windows", Architecture: "amd64"}}
	fake := &daemonIdentityClientFake{manifest: manifest}
	for _, test := range []struct {
		args []string
		want string
	}{
		{args: []string{"manifest"}, want: "Installation ID: daemon-1"},
		{args: []string{"rename", "Workshop scheduler"}, want: "Daemon renamed to Workshop"},
		{args: []string{"reset-identity", "--confirm", "daemon-1"}, want: "Installation identity reset to daemon-1"},
	} {
		var output bytes.Buffer
		cmd := newDaemonCmdWithClient(fake)
		cmd.SetOut(&output)
		cmd.SetArgs(test.args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("%v: %v", test.args, err)
		}
		if !strings.Contains(output.String(), test.want) {
			t.Fatalf("%v output = %q", test.args, output.String())
		}
	}
	if fake.rename != "Workshop scheduler" || fake.reset != "daemon-1" {
		t.Fatalf("rename=%q reset=%q", fake.rename, fake.reset)
	}
}

func TestDaemonManifestSupportsJSONAndValidation(t *testing.T) {
	fake := &daemonIdentityClientFake{manifest: server.ManifestResponse{InstallationID: "daemon-1"}}
	oldJSON := jsonOut
	jsonOut = true
	defer func() { jsonOut = oldJSON }()
	var output bytes.Buffer
	cmd := newDaemonCmdWithClient(fake)
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"manifest"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"installation_id": "daemon-1"`) {
		t.Fatalf("JSON output = %q", output.String())
	}
	jsonOut = false
	for _, args := range [][]string{{"rename", ""}, {"reset-identity"}} {
		cmd := newDaemonCmdWithClient(fake)
		cmd.SetArgs(args)
		if err := cmd.Execute(); !errors.Is(err, errUsage) {
			t.Fatalf("%v error = %v", args, err)
		}
	}
}
