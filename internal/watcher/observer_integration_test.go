package watcher

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestNativeObserverDetectsOneHundredAtomicReplacements(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{{ID: "watcher-1", Name: "target", Kind: domain.WatcherFile, Path: target, Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}}
	dispatcher := &dispatcherStub{ch: make(chan string, 2)}
	manager := New(store, dispatcher, clock.NewReal(), slog.Default())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	for trial := 0; trial < 100; trial++ {
		temporary := filepath.Join(root, "temporary.txt")
		if err := os.WriteFile(temporary, []byte("complete"), 0o600); err != nil {
			t.Fatal(err)
		}
		_ = os.Remove(target)
		if err := os.Rename(temporary, target); err != nil {
			t.Fatal(err)
		}
		select {
		case id := <-dispatcher.ch:
			if id != "watcher-1" {
				t.Fatalf("trial %d dispatch id = %q", trial, id)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("native observer did not dispatch atomic replacement trial %d", trial)
		}
		select {
		case id := <-dispatcher.ch:
			t.Fatalf("trial %d duplicate dispatch %q", trial, id)
		default:
		}
	}
	cancel()
	<-done
}

func TestManagerRejectsLinkedObservationRoot(t *testing.T) {
	realRoot := t.TempDir()
	link := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(realRoot, link); err != nil {
		t.Skipf("symbolic links unavailable: %v", err)
	}
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{{ID: "watcher-1", Name: "linked", Kind: domain.WatcherDirectory, Path: link, Pattern: "*", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}}
	manager := New(store, &dispatcherStub{ch: make(chan string, 1)}, clock.NewReal(), slog.Default())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	if health := manager.Health("watcher-1"); health.State != domain.WatcherDegraded {
		t.Fatalf("health = %+v", health)
	}
	cancel()
	<-done
}

func TestNativeObserverDetectsDirectWriteInExistingNestedDirectory(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	store := &watcherStoreStub{items: []domain.FilesystemWatcher{{ID: "watcher-1", Name: "nested", Kind: domain.WatcherDirectory, Path: root, Pattern: "*.json", Recursive: true, Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task-1", Enabled: true}}}
	dispatcher := &dispatcherStub{ch: make(chan string, 2)}
	manager := New(store, dispatcher, clock.NewReal(), slog.Default())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- manager.Run(ctx) }()
	<-manager.Ready()
	if err := os.WriteFile(filepath.Join(nested, "ready.json"), []byte("complete"), 0o600); err != nil {
		t.Fatal(err)
	}
	select {
	case id := <-dispatcher.ch:
		if id != "watcher-1" {
			t.Fatalf("dispatch id = %q", id)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("native observer did not dispatch nested direct write")
	}
	select {
	case id := <-dispatcher.ch:
		t.Fatalf("duplicate dispatch %q", id)
	default:
	}
	cancel()
	<-done
}
