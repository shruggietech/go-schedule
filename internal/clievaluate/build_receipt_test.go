package clievaluate

import (
	"strings"
	"testing"
)

func TestBuildReceiptRejectsDifferentSourceAndExecutableBytes(t *testing.T) {
	subject := Subject{Repository: "https://github.com/shruggietech/go-schedule", Commit: strings.Repeat("a", 40), Dirty: true, SourceSHA256: strings.Repeat("b", 64)}
	for _, platform := range []string{"windows", "linux"} {
		for _, role := range []string{"cli", "daemon", "runner"} {
			subject.Binaries = append(subject.Binaries, Binary{role, platform, "amd64", preparedName(role, platform), strings.Repeat("c", 64)})
		}
	}
	receipt := BuildReceipt{Version, "go-schedule.cli-qa-build-receipt", subject.Repository, subject.Commit, subject.Dirty, subject.SourceSHA256, append([]Binary{}, subject.Binaries...)}
	directory := t.TempDir()
	if err := writeBuildReceipt(directory, receipt); err != nil {
		t.Fatal(err)
	}
	if err := requireBuildReceipt(directory, subject); err != nil {
		t.Fatal(err)
	}
	for _, change := range []func(*Subject){
		func(s *Subject) { s.SourceSHA256 = strings.Repeat("d", 64) },
		func(s *Subject) { s.Dirty = false },
		func(s *Subject) { s.Commit = strings.Repeat("e", 40) },
		func(s *Subject) { s.Binaries[0].SHA256 = strings.Repeat("f", 64) },
	} {
		current := subject
		current.Binaries = append([]Binary{}, subject.Binaries...)
		change(&current)
		if err := requireBuildReceipt(directory, current); err == nil {
			t.Fatal("different build/source identity accepted")
		}
	}
	if err := requireBuildReceipt(t.TempDir(), subject); err == nil {
		t.Fatal("missing receipt accepted")
	}
	duplicate := receipt
	duplicate.Binaries = append([]Binary{}, receipt.Binaries...)
	duplicate.Binaries[1] = duplicate.Binaries[0]
	if err := validateBuildReceipt(duplicate, subject); err == nil {
		t.Fatal("duplicate receipt binary accepted")
	}
	malformed := t.TempDir()
	store, err := OpenArtifacts(malformed, 1<<20, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = store.Close() }()
	if err := store.Replace("build-receipt.json", []byte(`{"schema_version":"1","kind":"go-schedule.cli-qa-build-receipt"}`)); err != nil {
		t.Fatal(err)
	}
	if err := requireBuildReceipt(malformed, subject); err == nil {
		t.Fatal("incomplete receipt accepted")
	}
}
