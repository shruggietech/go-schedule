package watcher

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	recoveryInterval = 2 * time.Second
	maxHealthReason  = 512
)

// Store supplies the current durable watcher definitions.
type Store interface {
	ListFilesystemWatchers() ([]domain.FilesystemWatcher, error)
}

// Dispatcher submits one watcher-originated request to the scheduler.
type Dispatcher interface {
	FireFilesystemWatcher(string) error
}

// HealthReporter receives deduplicated runtime health transitions.
type HealthReporter func(domain.FilesystemWatcher, domain.WatcherHealth)

type observerFactory func() (observer, error)

type fileSnapshot struct {
	size    int64
	modTime time.Time
}

type pendingCandidate struct {
	watcherID   string
	path        string
	generation  uint64
	dueAt       time.Time
	snapshot    *fileSnapshot
	stableSince time.Time
}

// Manager owns one native observer and one generation-scoped event loop.
type Manager struct {
	store      Store
	dispatcher Dispatcher
	clock      clock.Clock
	log        *slog.Logger
	factory    observerFactory
	reporter   HealthReporter
	reload     chan struct{}
	ready      chan struct{}
	readyOnce  sync.Once

	mu      sync.RWMutex
	health  map[string]domain.WatcherHealth
	configs map[string]domain.FilesystemWatcher
}

// New creates a filesystem watcher manager using the native observer.
func New(store Store, dispatcher Dispatcher, clk clock.Clock, log *slog.Logger) *Manager {
	return newManager(store, dispatcher, clk, log, newFSObserver)
}

func newManager(store Store, dispatcher Dispatcher, clk clock.Clock, log *slog.Logger, factory observerFactory) *Manager {
	if log == nil {
		log = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}
	return &Manager{store: store, dispatcher: dispatcher, clock: clk, log: log, factory: factory, reload: make(chan struct{}, 1), ready: make(chan struct{}), health: map[string]domain.WatcherHealth{}, configs: map[string]domain.FilesystemWatcher{}}
}

// SetHealthReporter registers the transition callback used by live events.
func (m *Manager) SetHealthReporter(reporter HealthReporter) { m.reporter = reporter }

// Reload asks the event loop to replace all registrations. Requests are coalesced.
func (m *Manager) Reload() {
	select {
	case m.reload <- struct{}{}:
	default:
	}
}

// Ready closes after the first registration attempt completes.
func (m *Manager) Ready() <-chan struct{} { return m.ready }

// Health returns the current runtime state for id.
func (m *Manager) Health(id string) domain.WatcherHealth {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if health, ok := m.health[id]; ok {
		return health
	}
	return domain.WatcherHealth{State: domain.WatcherDegraded, Reason: "watcher runtime has not loaded this definition", ChangedAt: m.clock.Now().UTC()}
}

// Run owns native observation until ctx is cancelled.
func (m *Manager) Run(ctx context.Context) error {
	state := runtimeState{manager: m, pending: map[string]pendingCandidate{}}
	state.rebuild()
	m.readyOnce.Do(func() { close(m.ready) })
	defer state.closeObserver()
	for {
		timer, timerC := state.nextTimer()
		var events <-chan fsnotify.Event
		var errs <-chan error
		if state.observer != nil {
			events = state.observer.Events()
			errs = state.observer.Errors()
		}
		select {
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			return nil
		case <-m.reload:
			if timer != nil {
				timer.Stop()
			}
			state.rebuild()
		case event, ok := <-events:
			if timer != nil {
				timer.Stop()
			}
			if !ok {
				state.degradeAll("filesystem observer stopped")
				continue
			}
			state.handleEvent(event)
		case err, ok := <-errs:
			if timer != nil {
				timer.Stop()
			}
			if !ok {
				state.degradeAll("filesystem observer stopped")
				continue
			}
			reason := "filesystem observer failed"
			if errors.Is(err, fsnotify.ErrEventOverflow) {
				reason = "filesystem observer queue overflowed"
			}
			state.degradeAll(reason)
		case <-timerC:
			state.processDue()
		}
	}
}

type runtimeState struct {
	manager       *Manager
	observer      observer
	definitions   map[string]domain.FilesystemWatcher
	configs       map[string]domain.FilesystemWatcher
	observed      map[string]bool
	pending       map[string]pendingCandidate
	generation    uint64
	recoveryDue   time.Time
	needsRecovery bool
}

func (s *runtimeState) rebuild() {
	s.closeObserver()
	s.generation++
	s.pending = map[string]pendingCandidate{}
	s.definitions = map[string]domain.FilesystemWatcher{}
	s.configs = map[string]domain.FilesystemWatcher{}
	s.observed = map[string]bool{}
	s.needsRecovery = false
	s.recoveryDue = time.Time{}
	definitions, err := s.manager.store.ListFilesystemWatchers()
	if err != nil {
		s.manager.log.Error("watcher: load definitions", "err", err)
		s.scheduleRecovery()
		return
	}
	current := make(map[string]domain.FilesystemWatcher, len(definitions))
	for _, definition := range definitions {
		current[definition.ID] = definition
		s.definitions[definition.ID] = definition
		if !definition.Enabled {
			s.manager.setHealth(definition, domain.WatcherDisabled, "watcher is disabled")
		}
	}
	s.manager.replaceConfigs(current)
	if !hasEnabled(definitions) {
		return
	}
	obs, err := s.manager.factory()
	if err != nil {
		for _, definition := range definitions {
			if definition.Enabled {
				s.manager.setHealth(definition, domain.WatcherDegraded, "filesystem observer could not start")
			}
		}
		s.scheduleRecovery()
		return
	}
	s.observer = obs
	for _, definition := range definitions {
		if !definition.Enabled {
			continue
		}
		if err := s.register(definition, s.observed); err != nil {
			s.manager.setHealth(definition, domain.WatcherDegraded, registrationReason(err))
			s.scheduleRecovery()
			continue
		}
		s.configs[definition.ID] = definition
		s.manager.setHealth(definition, domain.WatcherActive, "")
	}
}

func (s *runtimeState) register(definition domain.FilesystemWatcher, added map[string]bool) error {
	root := definition.Path
	if definition.Kind == domain.WatcherFile {
		root = filepath.Dir(root)
	}
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errUnsupportedLink
	}
	if !info.IsDir() {
		return errNotDirectory
	}
	if hasLinkComponent(root) {
		return errUnsupportedLink
	}
	if !definition.Recursive || definition.Kind == domain.WatcherFile {
		return addObservation(s.observer, root, added)
	}
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			return nil
		}
		if path != root && entry.Type()&os.ModeSymlink != 0 {
			return filepath.SkipDir
		}
		return addObservation(s.observer, path, added)
	})
}

func hasLinkComponent(path string) bool {
	for current := filepath.Clean(path); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return false
		}
	}
}

var (
	errUnsupportedLink = errors.New("unsupported linked root")
	errNotDirectory    = errors.New("observation root is not a directory")
)

func addObservation(obs observer, path string, added map[string]bool) error {
	path = filepath.Clean(path)
	if added[path] {
		return nil
	}
	if err := obs.Add(path); err != nil {
		return err
	}
	added[path] = true
	return nil
}

func (s *runtimeState) handleEvent(event fsnotify.Event) {
	path := filepath.Clean(event.Name)
	if event.Has(fsnotify.Create) {
		s.addCreatedDirectory(path)
	}
	if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
		s.forgetObserved(path)
		if s.degradeAffectedRoots(path, "configured observation root was removed or replaced") {
			return
		}
	}
	if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) && !event.Has(fsnotify.Rename) {
		return
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return
	}
	now := s.manager.clock.Now()
	for id, definition := range s.configs {
		if !matches(definition, path) {
			continue
		}
		key := id + "\x00" + normalizedPath(path)
		s.pending[key] = pendingCandidate{watcherID: id, path: path, generation: s.generation, dueAt: now.Add(definition.Debounce)}
	}
}

func (s *runtimeState) addCreatedDirectory(path string) {
	info, err := os.Lstat(path)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || s.observer == nil {
		return
	}
	for id, watched := range s.configs {
		if watched.Kind != domain.WatcherDirectory || !watched.Recursive || !within(path, watched.Path) {
			continue
		}
		err := filepath.WalkDir(path, func(candidate string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !entry.IsDir() {
				return nil
			}
			if candidate != path && entry.Type()&os.ModeSymlink != 0 {
				return filepath.SkipDir
			}
			return addObservation(s.observer, candidate, s.observed)
		})
		if err != nil {
			s.manager.setHealth(watched, domain.WatcherDegraded, registrationReason(err))
			delete(s.configs, id)
			s.discardPending(id)
			s.scheduleRecovery()
		}
	}
}

func (s *runtimeState) processDue() {
	now := s.manager.clock.Now()
	if s.needsRecovery && !s.recoveryDue.After(now) {
		s.recoverDefinitions()
		return
	}
	for key, candidate := range s.pending {
		if candidate.dueAt.After(now) {
			continue
		}
		definition, ok := s.configs[candidate.watcherID]
		if !ok || candidate.generation != s.generation {
			delete(s.pending, key)
			continue
		}
		info, err := os.Lstat(candidate.path)
		if err != nil || !info.Mode().IsRegular() {
			delete(s.pending, key)
			continue
		}
		current := fileSnapshot{size: info.Size(), modTime: info.ModTime()}
		if candidate.snapshot == nil || *candidate.snapshot != current {
			candidate.snapshot = &current
			candidate.stableSince = now
			candidate.dueAt = now.Add(definition.Stability)
			s.pending[key] = candidate
			continue
		}
		if now.Sub(candidate.stableSince) < definition.Stability {
			candidate.dueAt = candidate.stableSince.Add(definition.Stability)
			s.pending[key] = candidate
			continue
		}
		delete(s.pending, key)
		if err := s.manager.dispatcher.FireFilesystemWatcher(candidate.watcherID); err != nil {
			s.manager.log.Warn("watcher: task dispatch rejected", "watcher", candidate.watcherID, "err", err)
		}
	}
}

func (s *runtimeState) recoverDefinitions() {
	if s.observer == nil {
		s.rebuild()
		return
	}
	s.needsRecovery = false
	for id, definition := range s.definitions {
		if !definition.Enabled {
			continue
		}
		if _, active := s.configs[id]; active {
			continue
		}
		if err := s.register(definition, s.observed); err != nil {
			s.manager.setHealth(definition, domain.WatcherDegraded, registrationReason(err))
			s.scheduleRecovery()
			continue
		}
		s.configs[id] = definition
		s.manager.setHealth(definition, domain.WatcherActive, "")
	}
}

func (s *runtimeState) discardPending(watcherID string) {
	for key, candidate := range s.pending {
		if candidate.watcherID == watcherID {
			delete(s.pending, key)
		}
	}
}

func (s *runtimeState) nextTimer() (*clock.Timer, <-chan time.Time) {
	var earliest time.Time
	for _, candidate := range s.pending {
		if earliest.IsZero() || candidate.dueAt.Before(earliest) {
			earliest = candidate.dueAt
		}
	}
	if s.needsRecovery && (earliest.IsZero() || s.recoveryDue.Before(earliest)) {
		earliest = s.recoveryDue
	}
	if earliest.IsZero() {
		return nil, nil
	}
	delay := earliest.Sub(s.manager.clock.Now())
	if delay < 0 {
		delay = 0
	}
	timer := s.manager.clock.NewTimer(delay)
	return timer, timer.C
}

func (s *runtimeState) degradeAll(reason string) {
	for _, definition := range s.configs {
		s.manager.setHealth(definition, domain.WatcherDegraded, reason)
	}
	s.closeObserver()
	s.pending = map[string]pendingCandidate{}
	s.configs = map[string]domain.FilesystemWatcher{}
	s.scheduleRecovery()
}

func (s *runtimeState) degradeAffectedRoots(path, reason string) bool {
	affected := false
	for id, definition := range s.configs {
		root := observationRoot(definition)
		if !samePath(root, path) && !within(root, path) {
			continue
		}
		affected = true
		s.manager.setHealth(definition, domain.WatcherDegraded, reason)
		delete(s.configs, id)
		s.discardPending(id)
	}
	if affected {
		s.scheduleRecovery()
	}
	return affected
}

func (s *runtimeState) forgetObserved(path string) {
	for observed := range s.observed {
		if samePath(observed, path) || within(observed, path) {
			delete(s.observed, observed)
		}
	}
}

func (s *runtimeState) scheduleRecovery() {
	s.needsRecovery = true
	s.recoveryDue = s.manager.clock.Now().Add(recoveryInterval)
}

func (s *runtimeState) closeObserver() {
	if s.observer != nil {
		_ = s.observer.Close()
		s.observer = nil
	}
}

func observationRoot(definition domain.FilesystemWatcher) string {
	if definition.Kind == domain.WatcherFile {
		return filepath.Dir(definition.Path)
	}
	return definition.Path
}

func (m *Manager) setHealth(definition domain.FilesystemWatcher, state domain.WatcherHealthState, reason string) {
	reason = boundedReason(reason)
	m.mu.Lock()
	previous, exists := m.health[definition.ID]
	if exists && previous.State == state && previous.Reason == reason {
		m.mu.Unlock()
		return
	}
	health := domain.WatcherHealth{State: state, Reason: reason, ChangedAt: m.clock.Now().UTC()}
	m.health[definition.ID] = health
	reporter := m.reporter
	m.mu.Unlock()
	if reporter != nil {
		reporter(definition, health)
	}
	m.log.Info("watcher: health transition", "watcher", definition.ID, "name", definition.Name, "state", state, "reason", reason)
}

func (m *Manager) replaceConfigs(configs map[string]domain.FilesystemWatcher) {
	m.mu.Lock()
	for id := range m.health {
		if _, ok := configs[id]; !ok {
			delete(m.health, id)
		}
	}
	m.configs = configs
	m.mu.Unlock()
}

func matches(definition domain.FilesystemWatcher, candidate string) bool {
	if definition.Kind == domain.WatcherFile {
		return samePath(definition.Path, candidate)
	}
	if !within(candidate, definition.Path) || samePath(candidate, definition.Path) {
		return false
	}
	rel, err := filepath.Rel(definition.Path, candidate)
	if err != nil || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || rel == ".." {
		return false
	}
	if !definition.Recursive && filepath.Dir(rel) != "." {
		return false
	}
	matched, err := filepath.Match(definition.Pattern, filepath.Base(candidate))
	return err == nil && matched
}

func within(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func samePath(left, right string) bool {
	left, right = normalizedPath(left), normalizedPath(right)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(left, right)
	}
	return left == right
}

func normalizedPath(path string) string { return filepath.Clean(path) }

func hasEnabled(definitions []domain.FilesystemWatcher) bool {
	for _, definition := range definitions {
		if definition.Enabled {
			return true
		}
	}
	return false
}

func registrationReason(err error) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "configured observation root is missing"
	case errors.Is(err, fs.ErrPermission):
		return "configured observation root is not accessible"
	case errors.Is(err, errUnsupportedLink):
		return "configured observation root uses an unsupported link"
	case errors.Is(err, errNotDirectory):
		return "configured observation root is not a directory"
	default:
		if detail := platformErrorDetail(err); detail != "" {
			return boundedReason("configured observation root could not be watched: " + detail)
		}
		return "configured observation root could not be watched"
	}
}

func platformErrorDetail(err error) string {
	var pathError *fs.PathError
	if errors.As(err, &pathError) {
		return platformErrorDetail(pathError.Err)
	}
	var linkError *os.LinkError
	if errors.As(err, &linkError) {
		return platformErrorDetail(linkError.Err)
	}
	var syscallError *os.SyscallError
	if errors.As(err, &syscallError) {
		return platformErrorDetail(syscallError.Err)
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return strings.TrimSpace(errno.Error())
	}
	return ""
}

func boundedReason(reason string) string {
	reason = strings.TrimSpace(reason)
	if len(reason) <= maxHealthReason {
		return reason
	}
	return fmt.Sprintf("%s...", reason[:maxHealthReason-3])
}
