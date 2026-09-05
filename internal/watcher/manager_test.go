package watcher

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type watcherStoreStub struct {
	mu    sync.Mutex
	items []domain.FilesystemWatcher
}

func (s *watcherStoreStub) ListFilesystemWatchers() ([]domain.FilesystemWatcher, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]domain.FilesystemWatcher(nil), s.items...), nil
}

type dispatcherStub struct{ ch chan string }

func (d *dispatcherStub) FireFilesystemWatcher(id string) error { d.ch <- id; return nil }

type fakeObserver struct {
	events chan fsnotify.Event
	errors chan error
	mu     sync.Mutex
	added  []string
	closed bool
	addErr func(string) error
}

func newFakeObserver() *fakeObserver {
	return &fakeObserver{events: make(chan fsnotify.Event, 256), errors: make(chan error, 16)}
}
func (o *fakeObserver) Add(path string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.added = append(o.added, path)
	if o.addErr != nil {
		return o.addErr(path)
	}
	return nil
}
func (o *fakeObserver) Close() error                  { o.mu.Lock(); o.closed = true; o.mu.Unlock(); return nil }
func (o *fakeObserver) Events() <-chan fsnotify.Event { return o.events }
func (o *fakeObserver) Errors() <-chan error          { return o.errors }

func TestManagerCoalescesWriteStormAndRecordsNoPath(t *testing.T) {
	root := canonicalTempDir(t)
	path := filepath.Join(root, "ready.txt")
	if err := os.WriteFile(path, []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	clk := clock.NewFake(now)
	watcher := domain.FilesystemWatcher{ID: "watcher-1", Name: "ready", Kind: domain.WatcherDirectory, Path: root, Pattern: "*.txt", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{watcher}}
	dispatcher := &dispatcherStub{ch: make(chan string, 2)}
	fake := newFakeObserver()
	manager := newManager(store, dispatcher, clk, slog.Default(), func() (observer, error) { return fake, nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	for i := 0; i < 100; i++ {
		fake.events <- fsnotify.Event{Name: path, Op: fsnotify.Write}
	}
	waitFor(t, func() bool { return len(fake.events) == 0 })
	time.Sleep(10 * time.Millisecond)
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	clk.Advance(25 * time.Millisecond)
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	clk.Advance(25 * time.Millisecond)
	select {
	case id := <-dispatcher.ch:
		if id != watcher.ID {
			t.Fatalf("dispatch id = %q", id)
		}
	case <-time.After(time.Second):
		t.Fatal("watcher did not dispatch")
	}
	select {
	case id := <-dispatcher.ch:
		t.Fatalf("unexpected duplicate dispatch %q", id)
	default:
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestManagerReloadCancelsPriorGeneration(t *testing.T) {
	root := canonicalTempDir(t)
	path := filepath.Join(root, "ready.txt")
	if err := os.WriteFile(path, []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	clk := clock.NewFake(time.Now())
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{{ID: "watcher-1", Name: "ready", Kind: domain.WatcherFile, Path: path, Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}}
	dispatcher := &dispatcherStub{ch: make(chan string, 1)}
	first := newFakeObserver()
	second := newFakeObserver()
	var calls atomic.Int32
	manager := newManager(store, dispatcher, clk, slog.Default(), func() (observer, error) {
		if calls.Add(1) == 1 {
			return first, nil
		}
		return second, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	first.events <- fsnotify.Event{Name: path, Op: fsnotify.Write}
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	manager.Reload()
	waitFor(t, func() bool { return calls.Load() == 2 })
	clk.Advance(time.Second)
	select {
	case <-dispatcher.ch:
		t.Fatal("obsolete generation dispatched")
	default:
	}
	cancel()
	<-done
}

func TestManagerRecoversMissingRootWithoutReplay(t *testing.T) {
	root := filepath.Join(canonicalTempDir(t), "later")
	clk := clock.NewFake(time.Now())
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{{ID: "watcher-1", Name: "later", Kind: domain.WatcherDirectory, Path: root, Pattern: "*", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}}
	dispatcher := &dispatcherStub{ch: make(chan string, 1)}
	manager := newManager(store, dispatcher, clk, slog.Default(), func() (observer, error) { return newFakeObserver(), nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	if health := manager.Health("watcher-1"); health.State != domain.WatcherDegraded {
		t.Fatalf("health = %s", health.State)
	}
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	clk.Advance(recoveryInterval)
	waitFor(t, func() bool { return manager.Health("watcher-1").State == domain.WatcherActive })
	select {
	case <-dispatcher.ch:
		t.Fatal("recovery replayed an event")
	default:
	}
	cancel()
	<-done
}

func TestMatchesExactDirectoryPatternAndRecursion(t *testing.T) {
	root := canonicalTempDir(t)
	direct := filepath.Join(root, "ready.json")
	nested := filepath.Join(root, "nested", "ready.json")
	file := domain.FilesystemWatcher{Kind: domain.WatcherFile, Path: direct}
	if !matches(file, direct) || matches(file, nested) {
		t.Fatal("exact file selection failed")
	}
	directory := domain.FilesystemWatcher{Kind: domain.WatcherDirectory, Path: root, Pattern: "*.json"}
	if !matches(directory, direct) || matches(directory, nested) || matches(directory, filepath.Join(root, "ready.txt")) {
		t.Fatal("non-recursive directory selection failed")
	}
	directory.Recursive = true
	if !matches(directory, nested) {
		t.Fatal("recursive directory selection failed")
	}
}

func TestObserverOverflowDegradesOnceAndRecovers(t *testing.T) {
	root := canonicalTempDir(t)
	clk := clock.NewFake(time.Now())
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{{ID: "watcher-1", Name: "incoming", Kind: domain.WatcherDirectory, Path: root, Pattern: "*", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}}
	dispatcher := &dispatcherStub{ch: make(chan string, 1)}
	first := newFakeObserver()
	second := newFakeObserver()
	var factories atomic.Int32
	manager := newManager(store, dispatcher, clk, slog.Default(), func() (observer, error) {
		if factories.Add(1) == 1 {
			return first, nil
		}
		return second, nil
	})
	var transitions atomic.Int32
	manager.SetHealthReporter(func(domain.FilesystemWatcher, domain.WatcherHealth) { transitions.Add(1) })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	first.errors <- fsnotify.ErrEventOverflow
	waitFor(t, func() bool { return manager.Health("watcher-1").State == domain.WatcherDegraded })
	if got := transitions.Load(); got != 2 {
		t.Fatalf("transitions after overflow = %d, want active plus degraded", got)
	}
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	clk.Advance(recoveryInterval)
	waitFor(t, func() bool { return manager.Health("watcher-1").State == domain.WatcherActive })
	if got := transitions.Load(); got != 3 {
		t.Fatalf("transitions after recovery = %d", got)
	}
	cancel()
	<-done
	second.mu.Lock()
	closed := second.closed
	second.mu.Unlock()
	if !closed {
		t.Fatal("observer handle remained open after shutdown")
	}
}

func TestMissingWatcherRecoveryPreservesHealthyLongStabilityCandidate(t *testing.T) {
	root := canonicalTempDir(t)
	missing := filepath.Join(canonicalTempDir(t), "missing")
	path := filepath.Join(root, "ready.txt")
	if err := os.WriteFile(path, []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	clk := clock.NewFake(time.Now())
	definitions := []domain.FilesystemWatcher{{ID: "healthy", Name: "healthy", Kind: domain.WatcherDirectory, Path: root, Pattern: "*", Debounce: 25 * time.Millisecond, Stability: 3 * time.Second, TargetTaskID: "task-1", Enabled: true}, {ID: "missing", Name: "missing", Kind: domain.WatcherDirectory, Path: missing, Pattern: "*", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}
	store := &watcherStoreStub{items: definitions}
	dispatcher := &dispatcherStub{ch: make(chan string, 1)}
	fake := newFakeObserver()
	manager := newManager(store, dispatcher, clk, slog.Default(), func() (observer, error) { return fake, nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	fake.events <- fsnotify.Event{Name: path, Op: fsnotify.Write}
	waitFor(t, func() bool { return len(fake.events) == 0 })
	time.Sleep(10 * time.Millisecond)
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	clk.Advance(25 * time.Millisecond)
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	clk.Advance(recoveryInterval)
	waitFor(t, func() bool { return clk.Waiters() > 0 })
	clk.Advance(time.Second)
	select {
	case id := <-dispatcher.ch:
		if id != "healthy" {
			t.Fatalf("dispatch id = %q", id)
		}
	case <-time.After(time.Second):
		t.Fatal("healthy candidate was lost during degraded watcher recovery")
	}
	cancel()
	<-done
}

func TestNewSubdirectoryRegistrationFailureDegradesWatcher(t *testing.T) {
	root := canonicalTempDir(t)
	clk := clock.NewFake(time.Now())
	definition := domain.FilesystemWatcher{ID: "watcher-1", Name: "recursive", Kind: domain.WatcherDirectory, Path: root, Pattern: "*", Recursive: true, Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{definition}}
	fake := newFakeObserver()
	manager := newManager(store, &dispatcherStub{ch: make(chan string, 1)}, clk, slog.Default(), func() (observer, error) { return fake, nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	fake.addErr = func(path string) error {
		if samePath(path, nested) {
			return os.ErrPermission
		}
		return nil
	}
	fake.events <- fsnotify.Event{Name: nested, Op: fsnotify.Create}
	waitFor(t, func() bool { return manager.Health("watcher-1").State == domain.WatcherDegraded })
	cancel()
	<-done
}

func TestSymlinkCandidateDoesNotDispatch(t *testing.T) {
	root := canonicalTempDir(t)
	target := filepath.Join(root, "target.txt")
	link := filepath.Join(root, "link.txt")
	if err := os.WriteFile(target, []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links unavailable: %v", err)
	}
	clk := clock.NewFake(time.Now())
	definition := domain.FilesystemWatcher{ID: "watcher-1", Name: "links", Kind: domain.WatcherDirectory, Path: root, Pattern: "*.txt", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{definition}}
	dispatcher := &dispatcherStub{ch: make(chan string, 1)}
	fake := newFakeObserver()
	manager := newManager(store, dispatcher, clk, slog.Default(), func() (observer, error) { return fake, nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	fake.events <- fsnotify.Event{Name: link, Op: fsnotify.Create}
	waitFor(t, func() bool { return len(fake.events) == 0 })
	time.Sleep(5 * time.Millisecond)
	clk.Advance(time.Second)
	select {
	case id := <-dispatcher.ch:
		t.Fatalf("symlink candidate dispatched %q", id)
	default:
	}
	cancel()
	<-done
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("condition was not met")
		}
		time.Sleep(time.Millisecond)
	}
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()
	directory := t.TempDir()
	resolved, err := filepath.EvalSymlinks(directory)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}
