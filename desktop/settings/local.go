package settings

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/platform"
	"github.com/shruggietech/go-schedule/internal/winuninstall"
)

// Backend supplies authoritative daemon-owned runtime paths.
type Backend interface {
	RuntimeInfo(context.Context) (server.RuntimeInfoResponse, error)
}

// Native supplies the bounded operating-system actions used by Settings.
type Native interface {
	ClipboardSetText(context.Context, string) error
	BrowserOpenURL(context.Context, string) error
}

// Paths contains resolved platform paths used to build a settings workspace.
type Paths struct {
	Preferences       string
	LegacyPreferences string
	ApplicationData   string
	MachineData       string
	Executable        string
	Maintenance       string
	GOOS              string
}

// Dependencies supplies deterministic filesystem, platform, and time seams.
type Dependencies struct {
	Paths  Paths
	Stat   func(string) (os.FileInfo, error)
	Now    func() time.Time
	Rename func(string, string) error
}

// DefaultDependencies resolves production paths and standard-library operations.
func DefaultDependencies() Dependencies {
	configDir, _ := os.UserConfigDir()
	homeDir, _ := os.UserHomeDir()
	executable, _ := os.Executable()
	applicationData := ""
	legacyRoot := ""
	if configDir != "" {
		applicationData = filepath.Join(configDir, "go-schedule", "desktop")
		legacyRoot = filepath.Join(configDir, "fyne", "tech.shruggie.goschedule")
	}
	if runtime.GOOS == "darwin" && homeDir != "" {
		legacyRoot = filepath.Join(homeDir, "Library", "Preferences", "fyne", "tech.shruggie.goschedule")
	}
	machineData := platform.DataDir()
	maintenance := ""
	if runtime.GOOS == "windows" {
		maintenance = winuninstall.CleanupResultPath(filepath.Dir(machineData))
	}
	return Dependencies{
		Paths: Paths{Preferences: joinWhenSet(applicationData, "preferences.json"), LegacyPreferences: joinWhenSet(legacyRoot, "preferences.json"), ApplicationData: applicationData, MachineData: machineData, Executable: executable, Maintenance: maintenance, GOOS: runtime.GOOS},
		Stat:  os.Stat, Now: time.Now, Rename: os.Rename,
	}
}

func joinWhenSet(root, name string) string {
	if root == "" {
		return ""
	}
	return filepath.Join(root, name)
}

// LocalBackend adapts the existing protected local daemon client.
type LocalBackend struct{ daemon Backend }

// NewLocalBackend creates a runtime-info adapter without opening a network listener.
func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }

// RuntimeInfo returns the daemon's effective storage paths.
func (b *LocalBackend) RuntimeInfo(ctx context.Context) (server.RuntimeInfoResponse, error) {
	return b.daemon.RuntimeInfo(ctx)
}
