package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProductionDesktopUsesStableWailsIdentity(t *testing.T) {
	var config struct {
		Name           string `json:"name"`
		OutputFilename string `json:"outputfilename"`
	}
	if err := json.Unmarshal(readRepositoryFile(t, "desktop", "wails.json"), &config); err != nil {
		t.Fatalf("decode desktop Wails configuration: %v", err)
	}
	if config.Name != "go-schedule" || config.OutputFilename != "gosched-gui" {
		t.Fatalf("desktop identity = (%q, %q), want (go-schedule, gosched-gui)", config.Name, config.OutputFilename)
	}

	desktopEntry := string(readRepositoryFile(t, "brand", "platform", "linux", "go-schedule.desktop"))
	if !strings.Contains(desktopEntry, "Exec=gosched-gui") {
		t.Error("Linux desktop entry does not preserve the stable gosched-gui launcher")
	}
}

func TestReleaseWorkflowBuildsProductionWailsPayload(t *testing.T) {
	workflow := string(readRepositoryFile(t, ".github", "workflows", "release.yml"))
	for _, required := range []string{
		"go-version-file: desktop/go.mod",
		"desktop/go.sum",
		"cache-dependency-path: desktop/frontend/package-lock.json",
		"working-directory: desktop/frontend",
		"npm ci",
		"github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build",
		"desktop/build/bin/gosched-gui",
		"desktop/build/bin/go-schedule.app",
		`app="$stage/gosched-gui.app"`,
		"go build -ldflags \"$BASE\" -o \"$stage/goschedd",
		"go build -ldflags \"$BASE\" -o \"$stage/gosched",
		"build/windows/verify_wxs.ps1 -StageDir",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("release workflow is missing production Wails payload fragment %q", required)
		}
	}
	for _, retired := range []string{
		"Install Fyne build dependencies",
		"./cmd/gosched-gui",
		"goversioninfo",
	} {
		if strings.Contains(workflow, retired) {
			t.Errorf("release workflow retains retired Fyne fragment %q", retired)
		}
	}
}

func TestWindowsDesktopPackageIncludesAttributionAndGuidance(t *testing.T) {
	wxs := string(readRepositoryFile(t, "build", "windows", "goschedule.wxs"))
	for _, required := range []string{
		`Source="$(StageDir)\README.md"`,
		`Source="$(StageDir)\LICENSE"`,
		`Source="$(StageDir)\CHANGELOG.md"`,
	} {
		if !strings.Contains(wxs, required) {
			t.Errorf("Windows desktop package is missing required payload %q", required)
		}
	}

	ci := string(readRepositoryFile(t, ".github", "workflows", "ci.yml"))
	if !strings.Contains(ci, "Copy-Item README.md, LICENSE, CHANGELOG.md -Destination $stage") {
		t.Error("Windows CI staging does not provide required documentation payloads to WiX")
	}
	release := string(readRepositoryFile(t, ".github", "workflows", "release.yml"))
	if !strings.Contains(release, `cp README.md LICENSE CHANGELOG.md "$stage/"`) {
		t.Error("release staging does not provide required documentation payloads to WiX")
	}
}

func TestCIMaintainsOnlyProductionWailsDesktop(t *testing.T) {
	workflow := string(readRepositoryFile(t, ".github", "workflows", "ci.yml"))
	for _, required := range []string{
		"wails-desktop:",
		"wails-desktop-browser-contract:",
		"Windows MSI compiled & silent contract",
		"go-version-file: desktop/go.mod",
		"desktop/build/bin/gosched-gui.exe",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("CI is missing production desktop fragment %q", required)
		}
	}
	for _, retired := range []string{"wails-proof:", "wails-browser-contract:", "GUI build & test (cgo)", "Install Fyne build dependencies"} {
		if strings.Contains(workflow, retired) {
			t.Errorf("CI retains retired parallel desktop fragment %q", retired)
		}
	}
}

func TestRetiredDesktopImplementationIsAbsent(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, retiredPath := range []string{"gui", filepath.Join("cmd", "gosched-gui"), filepath.Join("experiments", "wails-foundation")} {
		if _, err := os.Stat(filepath.Join(root, retiredPath)); err == nil {
			t.Errorf("retired desktop path still exists: %s", retiredPath)
		} else if !os.IsNotExist(err) {
			t.Fatalf("inspect retired path %s: %v", retiredPath, err)
		}
	}

	rootModule := string(readRepositoryFile(t, "go.mod"))
	if strings.Contains(strings.ToLower(rootModule), "fyne") {
		t.Error("root go.mod retains a Fyne dependency")
	}
}

func TestCurrentGuidanceNamesWailsAsProductionDesktop(t *testing.T) {
	for _, parts := range [][]string{
		{"README.md"},
		{"desktop", "README.md"},
		{"docs", "build-autopilot.md"},
		{"docs", "install.md"},
		{"docs", "cli.md"},
	} {
		content := string(readRepositoryFile(t, parts...))
		if !strings.Contains(content, "Wails") {
			t.Errorf("current guidance %s does not identify Wails", filepath.Join(parts...))
		}
	}

	rootReadme := string(readRepositoryFile(t, "README.md"))
	if strings.Contains(rootReadme, "gosched-gui (Fyne GUI)") || strings.Contains(rootReadme, "gui/        Fyne") {
		t.Error("README retains Fyne as the current production desktop")
	}
}
