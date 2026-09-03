package winuninstall

import (
	"errors"
	"testing"
)

type fakeBackend struct {
	targets       []Target
	discoverErr   error
	preflightErrs map[string]error
	removeErrs    map[string]error
	preflighted   []string
	removed       []string
	written       []Result
	writeErrors   []error
	cleared       bool
}

func (f *fakeBackend) Discover() ([]Target, error) { return f.targets, f.discoverErr }
func (f *fakeBackend) Preflight(target Target) (bool, error) {
	f.preflighted = append(f.preflighted, target.Path)
	return true, f.preflightErrs[target.Path]
}
func (f *fakeBackend) Remove(target Target) error {
	f.removed = append(f.removed, target.Path)
	return f.removeErrs[target.Path]
}
func (f *fakeBackend) WriteResult(result Result) error {
	f.written = append(f.written, result)
	if len(f.writeErrors) != 0 {
		err := f.writeErrors[0]
		f.writeErrors = f.writeErrors[1:]
		return err
	}
	return nil
}
func (f *fakeBackend) ClearResult() error { f.cleared = true; return nil }

func TestRunPreflightsEveryTargetBeforeDeleting(t *testing.T) {
	b := &fakeBackend{
		targets:       []Target{{Kind: TargetMachine, Path: `C:\ProgramData\goschedule`}, {Kind: TargetProfile, Path: `D:\Profiles\Ada\AppData\Roaming\fyne\tech.shruggie.goschedule`}},
		preflightErrs: map[string]error{`D:\Profiles\Ada\AppData\Roaming\fyne\tech.shruggie.goschedule`: errors.New("reparse ancestor")},
	}

	result := Run(b)

	if result.State != StateRefused {
		t.Fatalf("state = %q, want %q", result.State, StateRefused)
	}
	if len(b.preflighted) != 2 {
		t.Fatalf("preflighted %d targets, want 2", len(b.preflighted))
	}
	if len(b.removed) != 0 {
		t.Fatalf("removed %v before all preflights passed", b.removed)
	}
	if len(result.Entries) != 2 || result.Entries[1].Outcome != OutcomeRefused {
		t.Fatalf("entries = %#v", result.Entries)
	}
}

func TestRunContinuesAfterRemovalFailureAndReportsPartial(t *testing.T) {
	b := &fakeBackend{
		targets:    []Target{{Kind: TargetMachine, Path: "machine"}, {Kind: TargetProfile, SID: "S-1-5-21-1", Path: "profile"}},
		removeErrs: map[string]error{"machine": errors.New("locked")},
	}

	result := Run(b)

	if result.State != StatePartial || result.Remaining != 1 {
		t.Fatalf("result = %#v, want partial with one remaining", result)
	}
	if len(b.removed) != 2 {
		t.Fatalf("removed = %v, want both safe targets attempted", b.removed)
	}
	if b.cleared {
		t.Fatal("partial cleanup cleared its retained evidence")
	}
}

func TestRunCompleteClearsResultEvidence(t *testing.T) {
	b := &fakeBackend{targets: []Target{{Kind: TargetMachine, Path: "machine"}}}

	result := Run(b)

	if result.State != StateComplete || !b.cleared {
		t.Fatalf("result = %#v, cleared = %v", result, b.cleared)
	}
}

func TestRunDiscoveryFailureIsInternalError(t *testing.T) {
	b := &fakeBackend{discoverErr: errors.New("profile registry unavailable")}

	result := Run(b)

	if result.State != StateInternalError || result.Error == "" {
		t.Fatalf("result = %#v", result)
	}
}

func TestRunRewritesTerminalStateAfterInitialResultFailure(t *testing.T) {
	b := &fakeBackend{
		targets:     []Target{{Kind: TargetMachine, Path: "machine"}},
		writeErrors: []error{errors.New("registry unavailable"), nil},
	}

	result := Run(b)

	if result.State != StateInternalError || result.CompletedAt.IsZero() {
		t.Fatalf("result = %#v, want completed internal error", result)
	}
	if len(b.written) != 2 {
		t.Fatalf("writes = %d, want running plus terminal retry", len(b.written))
	}
	if b.written[1].State != StateInternalError || b.written[1].CompletedAt.IsZero() {
		t.Fatalf("terminal retry = %#v", b.written[1])
	}
	if len(b.preflighted) != 0 || len(b.removed) != 0 {
		t.Fatalf("cleanup proceeded after evidence failure: preflight=%v remove=%v", b.preflighted, b.removed)
	}
}
