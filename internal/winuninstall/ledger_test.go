package winuninstall

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupResultPath(t *testing.T) {
	base := filepath.Join(t.TempDir(), "ProgramData")
	want := filepath.Join(base, "ShruggieTech", "go-schedule-uninstall", "b6f3c2e1-7a4d-4c9e-9b2a-1f6d8e5a0c34", "cleanup-result.json")
	if got := CleanupResultPath(base); got != want {
		t.Fatalf("CleanupResultPath() = %q, want %q", got, want)
	}
}

func TestWriteLedgerIsUTF8WithoutBOMAndReplacesAtomically(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cleanup-result.json")
	first := Result{Schema: resultSchema, State: StateRunning}
	second := Result{Schema: resultSchema, State: StatePartial, Error: "accès refusé"}

	if err := writeLedger(path, first); err != nil {
		t.Fatal(err)
	}
	if err := writeLedger(path, second); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		t.Fatal("ledger contains a UTF-8 BOM")
	}
	var got Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid ledger JSON: %v", err)
	}
	if got.State != StatePartial || got.Error != second.Error {
		t.Fatalf("ledger = %#v", got)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temporary ledger remains: %v", err)
	}
}
