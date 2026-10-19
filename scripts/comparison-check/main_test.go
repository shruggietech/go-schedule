package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	data, err := os.ReadFile("../../site/data/comparisons/v1.5.1.json")
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "site/data/comparisons/v1.5.1.json", data)
	writeFixture(t, root, "CHANGELOG.md", []byte("# Changelog\n\n## [1.5.1] - 2026-09-27\n\nPublished capability.\n\n## [1.5.0] - 2026-09-23\n\nOlder.\n"))
	writeFixture(t, root, ".github/release-notes/v1.5.1.md", []byte("# v1.5.1\n\nPublished capability.\n"))
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	decisions := map[string]string{}
	for _, row := range snapshot.Capabilities {
		decisions[row.ID] = "no-change"
	}
	r := Review{Tag: "v1.5.1", SourceRef: "refs/tags/v1.5.1", ReviewedAt: "2026-10-19", DatasetSHA256: digest(data), ChangelogSHA256: digest([]byte("## [1.5.1] - 2026-09-27\n\nPublished capability.\n\n")), NotesSHA256: digest([]byte("# v1.5.1\n\nPublished capability.\n")), Decisions: decisions}
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "site/data/reviews/v1.5.1.json", raw)
	return root
}

func writeFixture(t *testing.T, root, relative string, data []byte) {
	t.Helper()
	file := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(file), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func TestReleaseReviewBindsTagCommitAndContent(t *testing.T) {
	root := fixture(t)
	if err := checkRelease(root, "v1.5.1", strings.Repeat("a", 40)); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"site/data/comparisons/v1.5.1.json", "CHANGELOG.md", ".github/release-notes/v1.5.1.md"} {
		t.Run(file, func(t *testing.T) {
			root := fixture(t)
			p := filepath.Join(root, filepath.FromSlash(file))
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			// A valid semantic edit, rather than only malformed input, must invalidate review.
			data = []byte(strings.ReplaceAll(string(data), "Published capability.", "Changed capability."))
			if strings.HasSuffix(file, ".json") {
				data = append(data, '\n')
			}
			writeFixture(t, root, file, data)
			if err := checkRelease(root, "v1.5.1", strings.Repeat("a", 40)); err == nil {
				t.Fatal("changed reviewed content accepted")
			}
		})
	}
}

func TestRejectIncompleteOrMismatchedReview(t *testing.T) {
	for _, name := range []string{"missing", "tag", "source", "date", "decisions", "unknown-decision"} {
		t.Run(name, func(t *testing.T) {
			root := fixture(t)
			file := filepath.Join(root, "site/data/reviews/v1.5.1.json")
			if name == "missing" {
				if err := os.Remove(file); err != nil {
					t.Fatal(err)
				}
			} else {
				data, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				var r Review
				if err := json.Unmarshal(data, &r); err != nil {
					t.Fatal(err)
				}
				switch name {
				case "tag":
					r.Tag = "v1.5.2"
				case "source":
					r.SourceRef = "refs/heads/main"
				case "date":
					r.ReviewedAt = "not-a-date"
				case "decisions":
					delete(r.Decisions, "mcp")
				case "unknown-decision":
					r.Decisions["mcp"] = "skip"
				}
				data, err = json.Marshal(r)
				if err != nil {
					t.Fatal(err)
				}
				writeFixture(t, root, "site/data/reviews/v1.5.1.json", data)
			}
			if err := checkRelease(root, "v1.5.1", strings.Repeat("a", 40)); err == nil {
				t.Fatal("invalid review accepted")
			}
		})
	}
	root := fixture(t)
	for _, tag := range []string{"../../escape", "v1.5", "main", "v1.5.1\n"} {
		if err := checkRelease(root, tag, strings.Repeat("a", 40)); err == nil {
			t.Fatalf("invalid tag %q accepted", tag)
		}
	}
	if err := checkRelease(root, "v1.5.1", "main"); err == nil {
		t.Fatal("unresolved source accepted")
	}
}

func TestSnapshotRejectsMissingEvidenceAndDuplicateIdentity(t *testing.T) {
	for _, mode := range []string{"source", "status", "duplicate", "missing-cell", "context"} {
		t.Run(mode, func(t *testing.T) {
			root := fixture(t)
			file := filepath.Join(root, "site/data/comparisons/v1.5.1.json")
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var s Snapshot
			if err := json.Unmarshal(data, &s); err != nil {
				t.Fatal(err)
			}
			c := s.Capabilities[0].Cells["goschedule"]
			switch mode {
			case "source":
				c.Source = ""
			case "status":
				c.Status = "best"
			case "duplicate":
				s.Capabilities = append(s.Capabilities, s.Capabilities[0])
			case "missing-cell":
				delete(s.Capabilities[0].Cells, "jams")
			case "context":
				s.Products[0].Context = ""
			}
			s.Capabilities[0].Cells["goschedule"] = c
			if err := validateSnapshot(s); err == nil {
				t.Fatal("invalid snapshot accepted")
			}
		})
	}
}
