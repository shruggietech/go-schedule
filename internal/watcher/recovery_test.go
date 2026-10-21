package watcher

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type recoveryDispatcher struct {
	calls []string
	err   error
}

func (d *recoveryDispatcher) FireFilesystemWatcher(id string) error {
	d.calls = append(d.calls, id)
	return d.err
}

type recoveryFixture struct {
	state           *runtimeState
	manager         *Manager
	clock           *clock.FakeClock
	store           *watcherStoreStub
	dispatcher      *recoveryDispatcher
	definition      domain.FilesystemWatcher
	observer        *fakeObserver
	factoryFailures int
	addFailures     int
	replacements    int
	logs            bytes.Buffer
}

func newRecoveryFixture(t *testing.T, present bool) *recoveryFixture {
	t.Helper()
	path := filepath.Join(canonicalTempDir(t), "ready.txt")
	if present {
		if err := os.WriteFile(path, []byte("ready"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	f := &recoveryFixture{clock: clock.NewFake(time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)), dispatcher: &recoveryDispatcher{}}
	f.definition = domain.FilesystemWatcher{ID: "watcher", Name: "ready", Kind: domain.WatcherFile, Path: path, Enabled: true, Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: "task"}
	f.store = &watcherStoreStub{items: []domain.FilesystemWatcher{f.definition}}
	f.manager = newManager(f.store, f.dispatcher, f.clock, slog.New(slog.NewTextHandler(&f.logs, nil)), func() (observer, error) {
		if f.factoryFailures > 0 {
			f.factoryFailures--
			return nil, os.ErrPermission
		}
		f.observer = newFakeObserver()
		f.observer.addErr = func(string) error {
			if f.addFailures > 0 {
				f.addFailures--
				return os.ErrPermission
			}
			return nil
		}
		return f.observer, nil
	})
	f.state = &runtimeState{manager: f.manager, pending: map[string]pendingCandidate{}}
	f.state.rebuild()
	t.Cleanup(f.state.closeObserver)
	return f
}

func (f *recoveryFixture) key() string             { return candidateKey(f.definition.ID, f.definition.Path) }
func (f *recoveryFixture) advance(d time.Duration) { f.clock.Advance(d); f.state.processDue() }
func (f *recoveryFixture) finishCandidate() {
	f.advance(f.definition.Debounce)
	f.advance(f.definition.Stability)
}
func (f *recoveryFixture) hint() {
	f.state.handleEvent(fsnotify.Event{Name: f.definition.Path, Op: fsnotify.Write})
}
func (f *recoveryFixture) replace(t *testing.T) {
	t.Helper()
	prior := inspectFile(f.definition.Path)
	f.replacements++
	temporary := filepath.Join(filepath.Dir(f.definition.Path), fmt.Sprintf("replacement-%d.txt", f.replacements))
	if err := os.WriteFile(temporary, []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	if prior.present {
		if err := os.Chtimes(temporary, prior.version.modTime, prior.version.modTime); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(f.definition.Path); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Rename(temporary, f.definition.Path); err != nil {
		t.Fatal(err)
	}
	if prior.present && sameVersion(prior.version, inspectFile(f.definition.Path).version) {
		t.Fatal("fixture did not replace file identity")
	}
}
func (f *recoveryFixture) degrade(t *testing.T, mode string) {
	t.Helper()
	switch mode {
	case "observer":
		f.state.degradeAll("filesystem observer failed")
	case "root":
		root := filepath.Dir(f.definition.Path)
		f.state.forgetObserved(root)
		if !f.state.degradeAffectedRoots(root, "configured observation root was removed or replaced") {
			t.Fatal("root was not degraded")
		}
	default:
		t.Fatalf("unknown degradation mode %q", mode)
	}
	f.assertHealth(t, domain.WatcherDegraded, f.definition)
}
func (f *recoveryFixture) assertHealth(t *testing.T, expected domain.WatcherHealthState, definition domain.FilesystemWatcher) {
	t.Helper()
	got := f.manager.Health(definition.ID)
	if got.State != expected || got.ObservedDefinitionRevision != domain.FilesystemWatcherRevision(definition) || !got.ObservedAt.Equal(f.clock.Now()) {
		t.Fatalf("health=%+v, want state %s, current revision and observation time", got, expected)
	}
}
func (f *recoveryFixture) assertCalls(t *testing.T, expected int) {
	t.Helper()
	if len(f.dispatcher.calls) != expected {
		t.Fatalf("replacement dispatches=%d, want %d", len(f.dispatcher.calls), expected)
	}
	for _, id := range f.dispatcher.calls {
		if id != f.definition.ID {
			t.Fatalf("unexpected dispatched watcher %q", id)
		}
	}
}

func TestRecoveryReconcilesReplacementAcrossObserverFailure(t *testing.T) {
	for _, mode := range []string{"observer", "root"} {
		t.Run(mode, func(t *testing.T) {
			f := newRecoveryFixture(t, true)
			f.degrade(t, mode)
			f.replace(t)
			f.advance(recoveryInterval)
			f.assertHealth(t, domain.WatcherActive, f.definition)
			f.finishCandidate()
			f.assertCalls(t, 1)
			f.advance(reconciliationInterval)
			f.hint()
			f.finishCandidate()
			f.assertCalls(t, 1)
			f.replace(t)
			f.hint()
			f.finishCandidate()
			f.assertCalls(t, 2)
		})
	}
}

func TestRecoveryPreservesPendingDebounceAndStability(t *testing.T) {
	for _, mode := range []string{"observer", "root"} {
		for _, stable := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/snapshot=%t", mode, stable), func(t *testing.T) {
				f := newRecoveryFixture(t, true)
				f.replace(t)
				f.hint()
				if stable {
					f.advance(f.definition.Debounce)
				}
				before := f.state.pending[f.key()]
				f.degrade(t, mode)
				f.advance(recoveryInterval)
				after, ok := f.state.pending[f.key()]
				if !ok || after.generation != before.generation || !after.dueAt.Equal(before.dueAt) || !after.stableSince.Equal(before.stableSince) || after.snapshot != before.snapshot || !after.reconciled {
					t.Fatalf("pending candidate changed across recovery: before=%+v after=%+v exists=%t", before, after, ok)
				}
				f.finishCandidate()
				f.assertCalls(t, 1)
				f.hint()
				f.finishCandidate()
				f.assertCalls(t, 1)
			})
		}
	}
}

func TestRecoveryRetainsStateAcrossRepeatedRegistrationFailures(t *testing.T) {
	for _, mode := range []string{"factory", "registration"} {
		t.Run(mode, func(t *testing.T) {
			f := newRecoveryFixture(t, true)
			f.replace(t)
			f.hint()
			f.advance(f.definition.Debounce)
			before := f.state.pending[f.key()]
			if mode == "factory" {
				f.factoryFailures = 2
				f.degrade(t, "observer")
			} else {
				f.addFailures = 2
				f.degrade(t, "root")
			}
			for attempt := 0; attempt < 2; attempt++ {
				f.advance(recoveryInterval)
				if health := f.manager.Health(f.definition.ID); health.State != domain.WatcherDegraded {
					t.Fatalf("failed attempt health=%+v", health)
				}
				f.state.processDue()
				if after, exists := f.state.pending[f.key()]; !exists || !after.dueAt.Equal(before.dueAt) || after.snapshot != before.snapshot {
					t.Fatal("failed recovery lost or reset pending stability")
				}
				if !f.state.fileStates[f.definition.ID].initialized {
					t.Fatal("failed recovery lost prior identity")
				}
				timer, ch := f.state.nextTimer()
				if timer == nil {
					t.Fatal("missing recovery timer")
				}
				f.clock.Advance(0)
				select {
				case <-ch:
					t.Fatal("inactive pending candidate spun an already-due timer")
				default:
				}
				timer.Stop()
			}
			f.advance(recoveryInterval)
			f.assertHealth(t, domain.WatcherActive, f.definition)
			f.finishCandidate()
			f.assertCalls(t, 1)
		})
	}
}

func TestRecoveryDoesNotReplayPreviouslyDispatchedVersion(t *testing.T) {
	for _, mode := range []string{"observer", "root"} {
		for _, native := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/native=%t", mode, native), func(t *testing.T) {
				f := newRecoveryFixture(t, true)
				f.replace(t)
				if native {
					f.hint()
				} else {
					f.state.reconcileFiles(f.clock.Now())
				}
				f.finishCandidate()
				f.assertCalls(t, 1)
				f.advance(2 * duplicateHintWindow)
				f.degrade(t, mode)
				f.advance(recoveryInterval)
				f.hint()
				f.finishCandidate()
				f.assertCalls(t, 1)
				f.advance(duplicateHintWindow)
				f.hint()
				f.finishCandidate()
				f.assertCalls(t, 2)
				f.replace(t)
				f.hint()
				f.finishCandidate()
				f.assertCalls(t, 3)
				if len(f.state.dispatched) != 1 {
					t.Fatalf("dispatch history=%d, want one latest identity", len(f.state.dispatched))
				}
			})
		}
	}
}

func TestRecoveryIsolatesChangedDisabledDeletedAndReloadedDefinitions(t *testing.T) {
	for _, change := range []string{"target", "path", "disabled", "deleted", "reload"} {
		t.Run(change, func(t *testing.T) {
			f := newRecoveryFixture(t, true)
			f.replace(t)
			f.hint()
			f.advance(f.definition.Debounce)
			previousGeneration := f.state.generation
			f.degrade(t, "observer")
			current := f.definition
			switch change {
			case "target":
				current.TargetTaskID = "another-task"
			case "path":
				current.Path = filepath.Join(filepath.Dir(current.Path), "another.txt")
				if err := os.WriteFile(current.Path, []byte("ready"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "disabled":
				current.Enabled = false
			}
			if change == "deleted" {
				f.store.items = nil
			} else {
				f.store.items = []domain.FilesystemWatcher{current}
			}
			if change == "reload" {
				f.state.rebuild()
			} else {
				f.advance(recoveryInterval)
			}
			f.advance(reconciliationInterval)
			f.finishCandidate()
			f.assertCalls(t, 0)
			if len(f.state.pending) != 0 || len(f.state.dispatched) != 0 {
				t.Fatal("obsolete candidate or dispatch identity survived definition change")
			}
			if change == "reload" && f.state.generation != previousGeneration+1 {
				t.Fatal("explicit reload did not replace generation")
			}
			if change != "reload" && f.state.generation != previousGeneration {
				t.Fatal("recovery replaced generation")
			}
			switch change {
			case "disabled":
				health := f.manager.Health(current.ID)
				if health.State != domain.WatcherDisabled || health.ObservedDefinitionRevision != domain.FilesystemWatcherRevision(current) {
					t.Fatalf("disabled health=%+v", health)
				}
				if len(f.state.fileStates) != 0 {
					t.Fatal("disabled definition retained file identity")
				}
			case "deleted":
				if len(f.state.definitions) != 0 || len(f.state.fileStates) != 0 || len(f.manager.health) != 0 {
					t.Fatal("deleted definition retained runtime state")
				}
			default:
				health := f.manager.Health(current.ID)
				if health.State != domain.WatcherActive || health.ObservedDefinitionRevision != domain.FilesystemWatcherRevision(current) {
					t.Fatalf("changed definition health=%+v", health)
				}
				if !sameVersion(f.state.fileStates[current.ID].version, inspectFile(current.Path).version) {
					t.Fatal("current revision did not get a fresh baseline")
				}
			}
		})
	}
}

func TestRecoveryHandlesAbsentFilesAndLaterCreation(t *testing.T) {
	for _, initiallyPresent := range []bool{false, true} {
		for _, createdDuringFailure := range []bool{false, true} {
			t.Run(fmt.Sprintf("present=%t/createDuringFailure=%t", initiallyPresent, createdDuringFailure), func(t *testing.T) {
				f := newRecoveryFixture(t, initiallyPresent)
				f.degrade(t, "observer")
				if initiallyPresent {
					if err := os.Remove(f.definition.Path); err != nil {
						t.Fatal(err)
					}
				}
				if createdDuringFailure {
					f.replace(t)
				}
				f.advance(recoveryInterval)
				f.assertHealth(t, domain.WatcherActive, f.definition)
				f.finishCandidate()
				if !createdDuringFailure {
					f.assertCalls(t, 0)
					f.replace(t)
					f.advance(reconciliationInterval)
					f.finishCandidate()
				}
				f.assertCalls(t, 1)
				f.advance(reconciliationInterval)
				f.hint()
				f.finishCandidate()
				f.assertCalls(t, 1)
			})
		}
	}
}

func TestRecoveryDoesNotRetryRejectedDispatch(t *testing.T) {
	f := newRecoveryFixture(t, true)
	f.dispatcher.err = errors.New("target task is disabled")
	f.replace(t)
	f.hint()
	f.finishCandidate()
	f.assertCalls(t, 1)
	if !strings.Contains(f.logs.String(), "task dispatch rejected") || !strings.Contains(f.logs.String(), "target task is disabled") {
		t.Fatalf("missing contextual rejection: %s", f.logs.String())
	}
	f.degrade(t, "observer")
	f.advance(recoveryInterval)
	f.assertHealth(t, domain.WatcherActive, f.definition)
	f.hint()
	f.finishCandidate()
	f.assertCalls(t, 1)
}

func TestRecoveryWaitsForRemovedRootAndDispatchesRecreatedFile(t *testing.T) {
	f := newRecoveryFixture(t, true)
	root := filepath.Dir(f.definition.Path)
	if err := os.Remove(f.definition.Path); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(root); err != nil {
		t.Fatal(err)
	}
	f.state.handleEvent(fsnotify.Event{Name: root, Op: fsnotify.Remove})
	f.assertHealth(t, domain.WatcherDegraded, f.definition)
	f.advance(recoveryInterval)
	if health := f.manager.Health(f.definition.ID); health.State != domain.WatcherDegraded || health.Reason != "configured observation root is missing" {
		t.Fatalf("missing-root health=%+v", health)
	}
	if !f.state.fileStates[f.definition.ID].present {
		t.Fatal("missing-root recovery discarded prior meaningful identity")
	}
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	f.replace(t)
	f.advance(recoveryInterval)
	f.assertHealth(t, domain.WatcherActive, f.definition)
	f.finishCandidate()
	f.assertCalls(t, 1)
}

func TestRecoveryPreservesUnaffectedPendingWorkWhenDefinitionChanges(t *testing.T) {
	f := newRecoveryFixture(t, true)
	other := f.definition
	other.ID = "other"
	other.Path = filepath.Join(filepath.Dir(other.Path), "other.txt")
	if err := os.WriteFile(other.Path, []byte("ready"), 0o600); err != nil {
		t.Fatal(err)
	}
	f.store.items = append(f.store.items, other)
	f.state.rebuild()
	f.replace(t)
	f.hint()
	f.advance(f.definition.Debounce)
	before := f.state.pending[f.key()]
	f.degrade(t, "observer")
	other.TargetTaskID = "changed-task"
	f.store.items[1] = other
	f.advance(recoveryInterval)
	if after, exists := f.state.pending[f.key()]; !exists || after.snapshot != before.snapshot || !after.stableSince.Equal(before.stableSince) {
		t.Fatal("another definition change discarded unchanged pending work")
	}
	f.finishCandidate()
	f.assertCalls(t, 1)
	if health := f.manager.Health(other.ID); health.State != domain.WatcherActive || health.ObservedDefinitionRevision != domain.FilesystemWatcherRevision(other) {
		t.Fatalf("other revision health=%+v", health)
	}
}
