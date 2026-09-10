package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/config"
)

func TestServiceInstallBindsValidatedAbsoluteConfiguration(t *testing.T) {
	configuration := config.Default()
	configuration.DataDir = filepath.Join(t.TempDir(), "daemon state")
	configuration.AdminGroup = ""
	path := filepath.Join(t.TempDir(), "remote config.json")
	data, err := json.Marshal(configuration)
	if err != nil {
		t.Fatal(err)
	}
	if err := osWriteFile(path, data); err != nil {
		t.Fatal(err)
	}

	var executable string
	var arguments []string
	control := func(_ string, gotExecutable string, gotArguments []string) (string, error) {
		executable = gotExecutable
		arguments = append([]string(nil), gotArguments...)
		return "install ok", nil
	}
	command := newServiceCmdWith(control, func() (string, error) { return filepath.Join(t.TempDir(), "goschedd"), nil })
	command.SetArgs([]string{"install", "--config", path})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	if executable == "" || !reflect.DeepEqual(arguments, []string{"--config", abs}) {
		t.Fatalf("executable=%q arguments=%q", executable, arguments)
	}
}

func TestServiceInstallWithoutConfigurationPreservesDefaults(t *testing.T) {
	called := false
	control := func(action, _ string, arguments []string) (string, error) {
		called = true
		if action != "install" || len(arguments) != 0 {
			t.Fatalf("action=%q arguments=%q", action, arguments)
		}
		return "install ok", nil
	}
	command := newServiceCmdWith(control, func() (string, error) { return "goschedd", nil })
	command.SetArgs([]string{"install"})
	if err := command.Execute(); err != nil || !called {
		t.Fatalf("called=%v err=%v", called, err)
	}
}

func TestServiceInstallRejectsInvalidConfigurationBeforeRegistration(t *testing.T) {
	for _, test := range []struct {
		name    string
		path    func(*testing.T) string
		message string
	}{
		{name: "missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "missing.json") }, message: "read service configuration"},
		{name: "invalid", path: func(t *testing.T) string {
			path := filepath.Join(t.TempDir(), "invalid.json")
			if err := osWriteFile(path, []byte(`{"remote":{"enabled":true}}`)); err != nil {
				t.Fatal(err)
			}
			return path
		}, message: "remote.bind_address"},
	} {
		t.Run(test.name, func(t *testing.T) {
			called := false
			control := func(string, string, []string) (string, error) { called = true; return "", errors.New("must not run") }
			command := newServiceCmdWith(control, func() (string, error) { return "goschedd", nil })
			command.SetArgs([]string{"install", "--config", test.path(t)})
			err := command.Execute()
			if err == nil || !strings.Contains(err.Error(), test.message) || called {
				t.Fatalf("called=%v err=%v", called, err)
			}
		})
	}
}

func osWriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}
