package releasegate

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestIssueDispositionMappingsAreExactAndCopied(t *testing.T) {
	t.Parallel()

	mappings := IssueDispositionMappings()
	wantIssues := []int{96, 98, 101, 104, 105, 106, 109, 111, 112, 113}
	if len(mappings) != len(wantIssues) {
		t.Fatalf("IssueDispositionMappings() count = %d, want %d", len(mappings), len(wantIssues))
	}
	for i, issue := range wantIssues {
		if mappings[i].Issue != issue {
			t.Fatalf("mapping[%d].Issue = %d, want %d", i, mappings[i].Issue, issue)
		}
	}
	if got := len(mappings[0].ObservationIDs); got != 36 {
		t.Fatalf("#96 observation count = %d, want 36", got)
	}
	if got := len(mappings[1].ObservationIDs); got != 16 {
		t.Fatalf("#98 observation count = %d, want 16", got)
	}
	want98 := append([]string{}, RequiredScenarioIDs()[20:36]...)
	if !reflect.DeepEqual(mappings[1].ObservationIDs, want98) {
		t.Fatalf("#98 observations = %v, want %v", mappings[1].ObservationIDs, want98)
	}
	if !reflect.DeepEqual(mappings[0].ChildIssues, []int{97, 98, 94}) {
		t.Fatalf("#96 child issues = %v", mappings[0].ChildIssues)
	}
	if !reflect.DeepEqual(mappings[0].PrerequisiteIssues, []int{89, 90}) {
		t.Fatalf("#96 prerequisite issues = %v", mappings[0].PrerequisiteIssues)
	}

	mappings[0].ObservationIDs[0] = "mutated"
	mappings[0].ChildIssues[0] = 999
	if IssueDispositionMappings()[0].ObservationIDs[0] == "mutated" {
		t.Fatal("IssueDispositionMappings returned mutable package state")
	}
	if IssueDispositionMappings()[0].ChildIssues[0] == 999 {
		t.Fatal("IssueDispositionMappings returned mutable relationship state")
	}
}

func TestRenderDispositionPacketIsCompleteDeterministicAndSafe(t *testing.T) {
	t.Parallel()

	_, _, evidence := passingEvidence(t)
	evidence.Environments[0].Snapshot = "snapshot | <unsafe> @octocat `tick`\r\nnext"
	observation := findObservation(&evidence, "access.intended-user")
	observation.Summary = "summary | <unsafe> @octocat `tick`\nnext"
	observation.AttachmentPaths = append(observation.AttachmentPaths, "attachments/@octocat`proof`.png")

	first, err := RenderDispositionPacket(evidence)
	if err != nil {
		t.Fatal(err)
	}
	second, err := RenderDispositionPacket(evidence)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("identical evidence rendered different packet bytes")
	}
	parent := t.TempDir()
	firstDir := filepath.Join(parent, "first")
	secondDir := filepath.Join(parent, "second")
	if err := WriteDispositionPacket(firstDir, evidence); err != nil {
		t.Fatal(err)
	}
	if err := WriteDispositionPacket(secondDir, evidence); err != nil {
		t.Fatal(err)
	}

	wantNames := []string{
		"issue-096.md", "issue-098.md", "issue-101.md", "issue-104.md",
		"issue-105.md", "issue-106.md", "issue-109.md", "issue-111.md",
		"issue-112.md", "issue-113.md", "packet.json",
	}
	gotNames := make([]string, 0, len(first))
	byName := make(map[string][]byte, len(first))
	for _, file := range first {
		gotNames = append(gotNames, file.Name)
		byName[file.Name] = file.Data
	}
	sort.Strings(gotNames)
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("packet files = %v, want %v", gotNames, wantNames)
	}
	for _, name := range wantNames {
		firstBytes, err := os.ReadFile(filepath.Join(firstDir, name))
		if err != nil {
			t.Fatal(err)
		}
		secondBytes, err := os.ReadFile(filepath.Join(secondDir, name))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(firstBytes, secondBytes) {
			t.Errorf("%s differs between separate packet destinations", name)
		}
	}

	issue96 := string(byName["issue-096.md"])
	for _, required := range []string{
		"Formal v1.0.0 candidate evidence for #96",
		"https://github.com/shruggietech/go-schedule/actions/runs/1234/attempts/1",
		evidence.Candidate.SHA256,
		"production candidate, archive, attachment, and manifest validation passed",
		"access.intended-user",
		"setup.shortcut-defaults",
		"remove.reinstall-after-wipe",
		"Coordinator child issues: #97, #98, #94",
		"Completed implementation prerequisites: #89, #90",
		"Independent closure boundary: #98",
		"&#64;octocat",
		"&lt;unsafe&gt;",
		"summary \\|",
		"attachments/&#64;octocat&#96;proof&#96;.png",
		"<br>",
	} {
		if !strings.Contains(issue96, required) {
			t.Errorf("issue-096.md missing %q", required)
		}
	}
	for _, forbidden := range []string{"@octocat", "<unsafe>"} {
		if strings.Contains(issue96, forbidden) {
			t.Errorf("issue-096.md contains unsafe/local text %q", forbidden)
		}
	}
	if strings.Contains(string(byName["issue-101.md"]), "access.intended-user") {
		t.Fatal("issue-101.md contains an unmapped observation")
	}

	var index struct {
		SchemaVersion int `json:"schema_version"`
		Issues        []struct {
			Issue        int      `json:"issue"`
			File         string   `json:"file"`
			Observations []string `json:"observations"`
		} `json:"issues"`
	}
	if err := json.Unmarshal(byName["packet.json"], &index); err != nil {
		t.Fatal(err)
	}
	if index.SchemaVersion != 1 || len(index.Issues) != 10 {
		t.Fatalf("packet index schema=%d issues=%d", index.SchemaVersion, len(index.Issues))
	}
	if index.Issues[0].Issue != 96 || len(index.Issues[0].Observations) != 36 {
		t.Fatalf("coordinator index = %+v", index.Issues[0])
	}
}

func TestWriteDispositionPacketCommitsAtomicallyAndRejectsExistingTarget(t *testing.T) {
	t.Parallel()

	_, _, evidence := passingEvidence(t)
	parent := t.TempDir()
	target := filepath.Join(parent, "packet")
	if err := WriteDispositionPacket(target, evidence); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 11 {
		t.Fatalf("packet entry count = %d, want 11", len(entries))
	}
	if err := WriteDispositionPacket(target, evidence); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("existing target error = %v", err)
	}
	assertNoDispositionStaging(t, parent)
}

func TestWriteDispositionPacketCleansUpWriteFailure(t *testing.T) {
	t.Parallel()

	_, _, evidence := passingEvidence(t)
	parent := t.TempDir()
	target := filepath.Join(parent, "packet")
	writes := 0
	err := writeDispositionPacket(target, evidence, func(name string, data []byte, mode fs.FileMode) error {
		writes++
		if writes == 2 {
			return errors.New("injected write failure")
		}
		return os.WriteFile(name, data, mode)
	}, renameDispositionNoReplace)
	if err == nil || !strings.Contains(err.Error(), "injected write failure") {
		t.Fatalf("write failure = %v", err)
	}
	if _, err := os.Lstat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("target exists after failure: %v", err)
	}
	assertNoDispositionStaging(t, parent)
}

func TestWriteDispositionPacketDoesNotReplaceTargetCreatedBeforeCommit(t *testing.T) {
	t.Parallel()

	_, _, evidence := passingEvidence(t)
	parent := t.TempDir()
	target := filepath.Join(parent, "packet")
	err := writeDispositionPacket(target, evidence, os.WriteFile, func(staging, destination string) error {
		if err := os.Mkdir(destination, 0o700); err != nil {
			return err
		}
		return renameDispositionNoReplace(staging, destination)
	})
	if err == nil {
		t.Fatal("commit replaced a concurrently created empty target directory")
	}
	entries, readErr := os.ReadDir(target)
	if readErr != nil {
		t.Fatalf("concurrent target was removed: %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("concurrent target was replaced with %d packet entries", len(entries))
	}
	assertNoDispositionStaging(t, parent)
}

func TestWriteDispositionPacketRejectsLinkedParent(t *testing.T) {
	t.Parallel()

	_, _, evidence := passingEvidence(t)
	root := t.TempDir()
	realParent := filepath.Join(root, "real")
	if err := os.Mkdir(realParent, 0o700); err != nil {
		t.Fatal(err)
	}
	linkedParent := filepath.Join(root, "linked")
	if err := os.Symlink(realParent, linkedParent); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf("symbolic links unavailable: %v", err)
		}
		t.Fatal(err)
	}
	err := WriteDispositionPacket(filepath.Join(linkedParent, "packet"), evidence)
	if err == nil || !strings.Contains(err.Error(), "symbolic link") {
		t.Fatalf("linked parent error = %v", err)
	}
}

func TestWriteDispositionPacketAllowsLinkedAncestorOfRealParent(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX ancestor-link regression; Windows reparse semantics are covered by direct-parent rejection")
	}

	_, _, evidence := passingEvidence(t)
	root := t.TempDir()
	realRoot := filepath.Join(root, "real")
	realParent := filepath.Join(realRoot, "parent")
	if err := os.MkdirAll(realParent, 0o700); err != nil {
		t.Fatal(err)
	}
	linkedRoot := filepath.Join(root, "linked")
	if err := os.Symlink(realRoot, linkedRoot); err != nil {
		t.Fatal(err)
	}
	if err := WriteDispositionPacket(filepath.Join(linkedRoot, "parent", "packet"), evidence); err != nil {
		t.Fatalf("real parent beneath a linked ancestor should remain usable: %v", err)
	}
}

func assertNoDispositionStaging(t *testing.T, parent string) {
	t.Helper()
	entries, err := os.ReadDir(parent)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".go-schedule-dispositions-") {
			t.Fatalf("staging directory remains: %s", entry.Name())
		}
	}
}
