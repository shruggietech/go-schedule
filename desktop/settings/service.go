package settings

import (
	"context"
	"errors"
	"io/fs"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/buildinfo"
)

const unavailableMessage = "Desktop settings are unavailable. Check access to the user configuration directory, then try again."
const runtimeInfoTimeout = 2 * time.Second

var productLinks = []ProductLink{
	{Key: "source", Label: "Source repository", Destination: "https://github.com/shruggietech/go-schedule"},
	{Key: "documentation", Label: "Documentation", Destination: "https://shruggietech.github.io/go-schedule/"},
	{Key: "publisher", Label: "ShruggieTech", Destination: "https://shruggie.tech"},
}

// Service composes settings snapshots and serializes preference mutations.
type Service struct {
	mu      sync.Mutex
	backend Backend
	native  Native
	deps    Dependencies
}

// NewService creates a production settings service.
func NewService(backend Backend, native Native) *Service {
	return NewServiceWithDependencies(backend, native, DefaultDependencies())
}

// NewServiceWithDependencies creates a settings service with deterministic boundaries.
func NewServiceWithDependencies(backend Backend, native Native, deps Dependencies) *Service {
	if deps.Stat == nil {
		deps.Stat = DefaultDependencies().Stat
	}
	if deps.Now == nil {
		deps.Now = DefaultDependencies().Now
	}
	if deps.Rename == nil {
		deps.Rename = DefaultDependencies().Rename
	}
	return &Service{backend: backend, native: native, deps: deps}
}

// Workspace returns desktop-local settings even when daemon paths are unavailable.
func (s *Service) Workspace(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	workspace, err := s.workspace(ctx)
	if err != nil {
		return Result{Action: "load_settings", Outcome: "unavailable", Message: unavailableMessage}
	}
	return Result{Action: "load_settings", Outcome: "accepted", Message: "Desktop settings loaded.", Workspace: &workspace}
}

// SaveAppearance persists one validated appearance and returns refreshed settings.
func (s *Service) SaveAppearance(ctx context.Context, value string) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	appearance := Appearance(value)
	if !validAppearance(appearance) {
		return Result{Action: "save_appearance", Outcome: "rejected", Message: "Choose System, Light, or Dark appearance."}
	}
	prefs, err := loadPreferences(s.deps)
	if err != nil {
		return Result{Action: "save_appearance", Outcome: "unavailable", Message: unavailableMessage}
	}
	prefs.Appearance = appearance
	if err := writePreferences(s.deps, prefs); err != nil {
		return Result{Action: "save_appearance", Outcome: "unavailable", Message: unavailableMessage}
	}
	workspace := s.workspaceWithPreferences(ctx, prefs)
	return Result{Action: "save_appearance", Outcome: "accepted", Message: "Appearance saved.", Workspace: &workspace}
}

// Restore resets only desktop-local preferences to their current defaults.
func (s *Service) Restore(ctx context.Context) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	prefs := DesktopPreferences{Version: CurrentPreferenceVersion, Appearance: AppearanceSystem, Transition: PreferenceTransition{Status: "not_required", Retired: append([]string(nil), retiredPreferenceKeys...)}}
	if err := writePreferences(s.deps, prefs); err != nil {
		return Result{Action: "restore_preferences", Outcome: "unavailable", Message: unavailableMessage}
	}
	workspace := s.workspaceWithPreferences(ctx, prefs)
	return Result{Action: "restore_preferences", Outcome: "accepted", Message: "Desktop preferences restored to system appearance.", Workspace: &workspace}
}

// CopyStoragePath resolves and copies one current allowlisted storage path.
func (s *Service) CopyStoragePath(ctx context.Context, id string) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	workspace, err := s.workspace(ctx)
	if err != nil {
		return Result{Action: "copy_storage_path", Outcome: "unavailable", Message: unavailableMessage}
	}
	for _, record := range workspace.Storage {
		if record.ID != id {
			continue
		}
		if !record.Copyable || record.Path == "" {
			return Result{Action: "copy_storage_path", Outcome: "rejected", Message: "That storage path is not currently available.", Workspace: &workspace}
		}
		if s.native == nil {
			return Result{Action: "copy_storage_path", Outcome: "unavailable", Message: "Clipboard access is unavailable.", Workspace: &workspace}
		}
		if err := s.native.ClipboardSetText(ctx, record.Path); err != nil {
			return Result{Action: "copy_storage_path", Outcome: "unavailable", Message: "The storage path could not be copied.", Workspace: &workspace}
		}
		return Result{Action: "copy_storage_path", Outcome: "accepted", Message: record.Label + " path copied.", Workspace: &workspace}
	}
	return Result{Action: "copy_storage_path", Outcome: "rejected", Message: "That storage location is not recognized.", Workspace: &workspace}
}

// OpenProductLink opens one fixed HTTPS product destination by stable key.
func (s *Service) OpenProductLink(ctx context.Context, key string) Result {
	for _, link := range productLinks {
		if link.Key != key {
			continue
		}
		destination, err := url.Parse(link.Destination)
		if err != nil || destination.Scheme != "https" || destination.Host == "" {
			return Result{Action: "open_product_link", Outcome: "rejected", Message: "That product link is not safe to open."}
		}
		if s.native == nil {
			return Result{Action: "open_product_link", Outcome: "unavailable", Message: "System browser access is unavailable."}
		}
		if err := s.native.BrowserOpenURL(ctx, link.Destination); err != nil {
			return Result{Action: "open_product_link", Outcome: "unavailable", Message: "The product link could not be opened."}
		}
		return Result{Action: "open_product_link", Outcome: "accepted", Message: link.Label + " opened in the system browser."}
	}
	return Result{Action: "open_product_link", Outcome: "rejected", Message: "That product link is not recognized."}
}

func (s *Service) workspace(ctx context.Context) (Workspace, error) {
	prefs, err := loadPreferences(s.deps)
	if err != nil {
		return Workspace{}, err
	}
	return s.workspaceWithPreferences(ctx, prefs), nil
}

func (s *Service) workspaceWithPreferences(ctx context.Context, prefs DesktopPreferences) Workspace {
	runtimeInfo := server.RuntimeInfoResponse{}
	daemonAvailable := false
	if s.backend != nil {
		runtimeCtx, cancel := context.WithTimeout(ctx, runtimeInfoTimeout)
		var err error
		runtimeInfo, err = s.backend.RuntimeInfo(runtimeCtx)
		cancel()
		daemonAvailable = err == nil
	}
	return Workspace{Preferences: prefs, PreferencePath: s.deps.Paths.Preferences, Storage: s.storage(runtimeInfo, daemonAvailable), Product: ProductInformation{Name: "go-schedule", Version: buildinfo.Version, Publisher: "ShruggieTech", Links: append([]ProductLink(nil), productLinks...)}, DaemonAvailable: daemonAvailable, LoadedAt: loadedAt(s.deps.Now())}
}

func (s *Service) storage(info server.RuntimeInfoResponse, daemonAvailable bool) []StorageRecord {
	preserve := "Preserved by normal software removal"
	wipe := "No built-in data wipe on this platform"
	if s.deps.Paths.GOOS == "windows" {
		wipe = "Removed by an explicit data wipe"
	}
	records := []StorageRecord{}
	addDaemon := func(id, label, path string) {
		owner, scope, explicit := "go-schedule", "machine", wipe
		if daemonAvailable && !pathWithin(s.deps.Paths.MachineData, path) {
			owner, scope, explicit = "external", "external", "Preserved because it is outside the application-owned data root"
		}
		if !daemonAvailable {
			path = ""
		}
		records = append(records, s.record(id, label, path, owner, scope, preserve, explicit))
	}
	addDaemon("machine-data", "Machine data", info.DataDir)
	addDaemon("task-database", "Task database", info.DatabasePath)
	addDaemon("configuration", "Configuration", info.ConfigPath)
	addDaemon("logs", "Logs", info.LogPath)
	addDaemon("runtime-state", "Runtime state", info.LockPath)
	records = append(records,
		s.record("desktop-data", "Desktop application data", s.deps.Paths.ApplicationData, "desktop", "user", preserve, wipe),
		s.record("desktop-preferences", "Desktop preferences", s.deps.Paths.Preferences, "desktop", "user", preserve, wipe),
	)
	executableDir := ""
	if filepath.IsAbs(s.deps.Paths.Executable) {
		executableDir = filepath.Dir(s.deps.Paths.Executable)
	}
	records = append(records, s.record("executable-directory", "Executable directory", executableDir, "operating-system", "runtime", "Installer-owned files are removed; development binaries are unaffected", "Installer-owned files are removed; development binaries are unaffected"))
	if executableDir != "" {
		docs := s.record("installed-documentation", "Installed documentation", filepath.Join(executableDir, "docs"), "operating-system", "runtime", "Removed with the application", "Removed with the application")
		if docs.Existence == "present" {
			records = append(records, docs)
		}
	}
	if s.deps.Paths.GOOS == "windows" {
		records = append(records, s.record("maintenance-evidence", "Maintenance evidence", s.deps.Paths.Maintenance, "go-schedule", "machine", "Not affected by software-only removal", "Removed after a complete wipe; retained only to report incomplete cleanup"))
	}
	return records
}

func (s *Service) record(id, label, path, owner, scope, normal, wipe string) StorageRecord {
	record := StorageRecord{ID: id, Label: label, Owner: owner, Scope: scope, Existence: "unavailable", NormalRemoval: normal, ExplicitWipe: wipe}
	if path == "" || !filepath.IsAbs(path) {
		return record
	}
	record.Path = filepath.Clean(path)
	record.Copyable = true
	_, err := s.deps.Stat(record.Path)
	switch {
	case err == nil:
		record.Existence = "present"
	case errors.Is(err, fs.ErrNotExist):
		record.Existence = "absent"
	default:
		record.Existence = "unavailable"
	}
	return record
}

func pathWithin(root, candidate string) bool {
	if !filepath.IsAbs(root) || !filepath.IsAbs(candidate) {
		return false
	}
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(candidate))
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
