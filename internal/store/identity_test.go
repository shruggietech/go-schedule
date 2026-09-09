package store

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestFreshStoresHaveUniqueStableIdentities(t *testing.T) {
	seen := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		path := filepath.Join(t.TempDir(), "daemon.db")
		st, err := Open(path)
		if err != nil {
			t.Fatal(err)
		}
		identity, err := st.DaemonIdentity()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := uuid.Parse(identity.InstallationID); err != nil {
			t.Fatalf("invalid installation ID %q: %v", identity.InstallationID, err)
		}
		if identity.DisplayName != domain.DefaultDaemonDisplayName {
			t.Fatalf("display name = %q", identity.DisplayName)
		}
		if _, duplicate := seen[identity.InstallationID]; duplicate {
			t.Fatalf("duplicate installation ID %q", identity.InstallationID)
		}
		seen[identity.InstallationID] = struct{}{}
		if err := st.Close(); err != nil {
			t.Fatal(err)
		}
		st, err = Open(path)
		if err != nil {
			t.Fatal(err)
		}
		reopened, err := st.DaemonIdentity()
		if err != nil {
			t.Fatal(err)
		}
		if reopened != identity {
			t.Fatalf("reopened identity = %+v, want %+v", reopened, identity)
		}
		_ = st.Close()
	}
}

func TestDaemonIdentityConcurrentReadsKeepSingleton(t *testing.T) {
	st := openMem(t)
	var wg sync.WaitGroup
	ids := make(chan string, 20)
	for i := 0; i < cap(ids); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			identity, err := st.DaemonIdentity()
			if err != nil {
				t.Errorf("read identity: %v", err)
				return
			}
			ids <- identity.InstallationID
		}()
	}
	wg.Wait()
	close(ids)
	var expected string
	for id := range ids {
		if expected == "" {
			expected = id
		}
		if id != expected {
			t.Fatalf("identity = %q, want %q", id, expected)
		}
	}
	var count int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM daemon_identity`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("identity rows = %d, err=%v", count, err)
	}
}

func TestRenameDaemonValidatesWithoutMutation(t *testing.T) {
	st := openMem(t)
	original, err := st.DaemonIdentity()
	if err != nil {
		t.Fatal(err)
	}
	renamed, err := st.RenameDaemon("  Workshop scheduler  ")
	if err != nil {
		t.Fatal(err)
	}
	if renamed.DisplayName != "Workshop scheduler" || renamed.InstallationID != original.InstallationID {
		t.Fatalf("renamed identity = %+v", renamed)
	}
	invalid := []string{"", "   ", "line\nbreak", string(make([]rune, 81))}
	invalid[3] = string([]rune("x"))
	for utf8.RuneCountInString(invalid[3]) < 81 {
		invalid[3] += "x"
	}
	for _, name := range invalid {
		if _, err := st.RenameDaemon(name); !errors.Is(err, domain.ErrInvalidDaemonDisplayName) {
			t.Fatalf("RenameDaemon(%q) error = %v", name, err)
		}
		current, err := st.DaemonIdentity()
		if err != nil {
			t.Fatal(err)
		}
		if current != renamed {
			t.Fatalf("invalid rename mutated identity: %+v", current)
		}
	}
}

func TestResetDaemonIdentityRequiresExactCurrentID(t *testing.T) {
	st := openMem(t)
	if _, err := st.db.Exec(`INSERT INTO groups(id,name,enabled,created_at,updated_at) VALUES('group-1','Preserved',1,'2026-09-09T00:00:00Z','2026-09-09T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	before, err := st.RenameDaemon("Workshop scheduler")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.ResetDaemonIdentity("wrong"); !errors.Is(err, ErrIdentityMismatch) {
		t.Fatalf("mismatch error = %v", err)
	}
	unchanged, _ := st.DaemonIdentity()
	if unchanged != before {
		t.Fatalf("mismatch mutated identity: %+v", unchanged)
	}
	after, err := st.ResetDaemonIdentity(before.InstallationID)
	if err != nil {
		t.Fatal(err)
	}
	if after.InstallationID == before.InstallationID || after.DisplayName != before.DisplayName || after.CreatedAt != before.CreatedAt {
		t.Fatalf("reset identity = %+v, before %+v", after, before)
	}
	var groups int
	if err := st.db.QueryRow(`SELECT COUNT(*) FROM groups WHERE id='group-1' AND name='Preserved'`).Scan(&groups); err != nil || groups != 1 {
		t.Fatalf("scheduler data after reset = %d, err=%v", groups, err)
	}
	if _, err := st.ResetDaemonIdentity(before.InstallationID); !errors.Is(err, ErrIdentityMismatch) {
		t.Fatalf("stale confirmation error = %v", err)
	}
}

func TestDaemonIdentityRejectsCorruptPersistentValues(t *testing.T) {
	st := openMem(t)
	if _, err := st.db.Exec(`UPDATE daemon_identity SET installation_id='not-a-uuid' WHERE singleton=1`); err != nil {
		t.Fatal(err)
	}
	if _, err := st.DaemonIdentity(); err == nil {
		t.Fatal("corrupt persistent identity was accepted")
	}
}

func TestReopenExistingIdentityDoesNotRequireRandomGeneration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "daemon.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	want, err := st.DaemonIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	original := installationIDGenerator
	installationIDGenerator = func() (string, error) { return "", errors.New("random unavailable") }
	t.Cleanup(func() { installationIDGenerator = original })
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen existing identity: %v", err)
	}
	defer reopened.Close()
	got, err := reopened.DaemonIdentity()
	if err != nil || got != want {
		t.Fatalf("identity=%+v err=%v, want %+v", got, err, want)
	}
}

func TestDatabaseRestoreAndClonePreserveIdentityUntilReset(t *testing.T) {
	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "source.db")
	clonePath := filepath.Join(directory, "clone.db")
	source, err := Open(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	want, err := source.RenameDaemon("Restored scheduler")
	if err != nil {
		t.Fatal(err)
	}
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(clonePath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	clone, err := Open(clonePath)
	if err != nil {
		t.Fatal(err)
	}
	copied, err := clone.DaemonIdentity()
	if err != nil || copied != want {
		t.Fatalf("copied identity=%+v err=%v, want %+v", copied, err, want)
	}
	reset, err := clone.ResetDaemonIdentity(copied.InstallationID)
	if err != nil {
		t.Fatal(err)
	}
	if reset.InstallationID == copied.InstallationID || reset.DisplayName != copied.DisplayName {
		t.Fatalf("reset clone=%+v", reset)
	}
	if err := clone.Close(); err != nil {
		t.Fatal(err)
	}
	restored, err := Open(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()
	got, err := restored.DaemonIdentity()
	if err != nil || got != want {
		t.Fatalf("source identity=%+v err=%v, want %+v", got, err, want)
	}
}
