package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/userprefs"
)

func runPreferenceCommand(t *testing.T, deps Dependencies, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	deps.Out = &out
	deps.Err = &bytes.Buffer{}
	root := newRootWithDependencies(deps)
	if command, _, _ := root.Find([]string{"preferences"}); command == root {
		root.AddCommand(newPreferencesCmd())
	}
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func TestPreferencesHeadlessScopedUpdatesAndExplicitReset(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop", "preferences.json")
	deps := Dependencies{PreferencePath: path, Secrets: noPreferenceSecrets{}}
	out, err := runPreferenceCommand(t, deps, "--json", "--profile", "unavailable-do-not-contact", "preferences", "set", "--appearance", "dark", "--popup-enabled", "--popup-daemons", "daemon-a")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, `"transition"`) {
		t.Fatalf("obsolete transition metadata: %s", out)
	}
	var result preferenceOutput
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	if result.Path != path || result.User == "" || result.Preferences.Appearance != userprefs.AppearanceDark || !result.Preferences.Popups.Enabled {
		t.Fatalf("result=%+v", result)
	}
	if _, err := runPreferenceCommand(t, deps, "preferences", "set", "--popup-enabled=false", "--popup-daemons="); err != nil {
		t.Fatal(err)
	}
	prefs, err := userprefs.NewStore(userprefs.Options{Path: path}).Load()
	if err != nil || prefs.Appearance != userprefs.AppearanceDark || prefs.Popups.Enabled || len(prefs.Popups.DaemonIDs) != 0 {
		t.Fatalf("state=%+v %v", prefs, err)
	}
	if _, err := runPreferenceCommand(t, deps, "preferences", "reset"); !errors.Is(err, errUsage) {
		t.Fatalf("missing reset scope: %v", err)
	}
	if _, err := runPreferenceCommand(t, deps, "preferences", "reset", "--scope", "popups"); err != nil {
		t.Fatal(err)
	}
	prefs, _ = userprefs.NewStore(userprefs.Options{Path: path}).Load()
	if prefs.Appearance != userprefs.AppearanceDark || prefs.Popups.Enabled {
		t.Fatalf("reset touched unrelated appearance: %+v", prefs)
	}
	out, err = runPreferenceCommand(t, deps, "preferences", "show")
	if strings.Contains(out, "Legacy preference") || strings.Contains(out, "Retired legacy") {
		t.Fatalf("obsolete labels: %s", out)
	}
	if err != nil || !strings.Contains(out, path) || !strings.Contains(out, "dark") {
		t.Fatalf("human output=%q err=%v", out, err)
	}
}

func TestPreferenceMalformedInputDoesNotInitializeIsolatedStore(t *testing.T) {
	for _, args := range [][]string{
		{"set", "--appearance", "sepia"},
		{"set", "--popup-conditions="},
		{"set", "--popup-severities", "critical"},
		{"set", "--popup-daemons", "duplicate,duplicate"},
		{"set"},
		{"reset", "--scope", "unknown"},
		{"show", "unexpected"},
	} {
		path := filepath.Join(t.TempDir(), "preferences.json")
		_, err := runPreferenceCommand(t, Dependencies{PreferencePath: path, Secrets: noPreferenceSecrets{}}, append([]string{"preferences"}, args...)...)
		if err == nil {
			t.Fatalf("accepted %v", args)
		}
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("invalid command %v initialized state: %v", args, err)
		}
	}
}

type noPreferenceSecrets struct{}

func (noPreferenceSecrets) Probe() error { panic("preference command accessed native credentials") }
func (noPreferenceSecrets) Save(string, string, string) error {
	panic("preference command accessed native credentials")
}
func (noPreferenceSecrets) Load(string, string) (string, error) {
	panic("preference command accessed native credentials")
}
func (noPreferenceSecrets) Delete(string, string) error {
	panic("preference command accessed native credentials")
}

type preferenceProfiles struct{ active string }

func (s *preferenceProfiles) Load() (clientprofile.Collection, error) {
	return clientprofile.Collection{Version: 1, ActiveDesktopProfileID: s.active, Profiles: []clientprofile.Profile{{ID: "profile-one", Label: "Example"}}}, nil
}
func (s *preferenceProfiles) Add(p clientprofile.Profile) (clientprofile.Profile, error) {
	return p, errors.New("unexpected add")
}
func (s *preferenceProfiles) Rename(string, string) (clientprofile.Profile, error) {
	return clientprofile.Profile{}, errors.New("unexpected rename")
}
func (s *preferenceProfiles) SetActive(id string) error { s.active = id; return nil }
func (s *preferenceProfiles) RemoveWith(string, func(clientprofile.Profile) error) (clientprofile.Profile, error) {
	return clientprofile.Profile{}, errors.New("unexpected removal")
}
func (s *preferenceProfiles) Path() string { return "injected-profiles.json" }

func TestSavedDesktopDefaultIsLocalPersistenceOnly(t *testing.T) {
	profiles := &preferenceProfiles{}
	deps := Dependencies{Profiles: profiles, Secrets: noPreferenceSecrets{}}
	out, err := runPreferenceCommand(t, deps, "--json", "--profile", "unavailable-do-not-contact", "profile", "default", "set", "Example")
	if err != nil || profiles.active != "profile-one" {
		t.Fatalf("out=%q active=%q err=%v", out, profiles.active, err)
	}
	if !strings.Contains(out, `"profile_id": "profile-one"`) || !strings.Contains(out, "injected-profiles.json") {
		t.Fatalf("output=%q", out)
	}
	if _, err := runPreferenceCommand(t, deps, "profile", "default", "set", "missing"); !errors.Is(err, errUsage) || profiles.active != "profile-one" {
		t.Fatalf("missing profile changed default: %q %v", profiles.active, err)
	}
	if _, err := runPreferenceCommand(t, deps, "profile", "default", "reset"); err != nil || profiles.active != "" {
		t.Fatalf("reset=%q %v", profiles.active, err)
	}
	out, err = runPreferenceCommand(t, deps, "profile", "default", "show")
	if err != nil || !strings.Contains(out, "This computer") {
		t.Fatalf("out=%q err=%v", out, err)
	}
}

func (noPreferenceSecrets) LoadContext(context.Context, string, string) (string, error) {
	return "", errors.New("unexpected credential read")
}
