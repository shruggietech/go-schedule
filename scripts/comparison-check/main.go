// Command comparison-check validates public comparison evidence and binds release
// reviews to exact candidate content before release workflows mutate assets.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Cell records a qualified capability claim with its official evidence.
type Cell struct {
	Status     string `json:"status"`
	Note       string `json:"note"`
	Source     string `json:"source"`
	ReviewedAt string `json:"reviewed_at"`
}

// Product identifies the edition and deployment context used by a comparison.
type Product struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Context string `json:"context"`
	URL     string `json:"url"`
}

// Capability groups the same bounded question across all products.
type Capability struct {
	ID       string          `json:"id"`
	Category string          `json:"category"`
	Name     string          `json:"name"`
	Cells    map[string]Cell `json:"cells"`
}

// Opportunity links an observed gap to an existing planning decision.
type Opportunity struct {
	Capability string `json:"capability"`
	Issue      string `json:"issue"`
	Note       string `json:"note"`
}

// Snapshot is a versioned comparison, selected separately from candidate review.
type Snapshot struct {
	SchemaVersion int           `json:"schema_version"`
	Tag           string        `json:"tag"`
	SourceCommit  string        `json:"source_commit,omitempty"`
	ReviewedAt    string        `json:"reviewed_at"`
	Products      []Product     `json:"products"`
	Capabilities  []Capability  `json:"capabilities"`
	Opportunities []Opportunity `json:"opportunities"`
}

// Review binds explicit capability decisions to exact reviewed release content.
type Review struct {
	Tag             string            `json:"tag"`
	SourceRef       string            `json:"source_ref"`
	ReviewedAt      string            `json:"reviewed_at"`
	DatasetSHA256   string            `json:"dataset_sha256"`
	ChangelogSHA256 string            `json:"changelog_sha256"`
	NotesSHA256     string            `json:"notes_sha256"`
	Decisions       map[string]string `json:"decisions"`
}

var tagPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
var commitPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

func main() {
	root := flag.String("root", ".", "repository root")
	tag := flag.String("tag", "", "exact release tag for candidate review")
	commit := flag.String("commit", "", "exact resolved tagged source commit")
	flag.Parse()
	var err error
	if *tag != "" || *commit != "" {
		err = checkRelease(*root, *tag, *commit)
	} else {
		err = checkSnapshots(*root)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "comparison-check: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("comparison-check: OK tag=%s commit=%s\n", *tag, *commit)
}

func decode[T any](data []byte, target *T) error {
	d := json.NewDecoder(strings.NewReader(string(data)))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("expected exactly one JSON document")
	}
	return nil
}

func validDate(date string) bool { _, err := time.Parse("2006-01-02", date); return err == nil }
func validURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host != "" && u.User == nil
}
func contentDigest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func validateSnapshot(s Snapshot) error {
	if s.SchemaVersion != 1 || !tagPattern.MatchString(s.Tag) || !validDate(s.ReviewedAt) {
		return errors.New("snapshot needs schema_version 1, a vMAJOR.MINOR.PATCH tag and YYYY-MM-DD review date")
	}
	if s.SourceCommit != "" && !commitPattern.MatchString(s.SourceCommit) {
		return errors.New("snapshot source_commit must be a resolved 40-character commit")
	}
	products := map[string]bool{}
	for _, p := range s.Products {
		if p.ID == "" || products[p.ID] || p.Name == "" || p.Context == "" || !validURL(p.URL) {
			return fmt.Errorf("product %q needs unique identity, edition/deployment context and HTTPS source", p.ID)
		}
		products[p.ID] = true
	}
	for _, id := range []string{"goschedule", "jams", "activebatch", "scriptrunner", "controlm", "logicapps", "resolve", "opcon"} {
		if !products[id] {
			return fmt.Errorf("required product %s is missing", id)
		}
	}
	ids := map[string]bool{}
	statuses := map[string]bool{"supported": true, "partial": true, "planned": true, "unavailable": true, "unknown": true}
	for _, c := range s.Capabilities {
		if c.ID == "" || ids[c.ID] || c.Name == "" || c.Category == "" {
			return fmt.Errorf("capability %q needs unique identity, name and category", c.ID)
		}
		ids[c.ID] = true
		if len(c.Cells) != len(products) {
			return fmt.Errorf("capability %s must contain one cell per product", c.ID)
		}
		for id := range products {
			cell, ok := c.Cells[id]
			if !ok || !statuses[cell.Status] || cell.Note == "" || !validURL(cell.Source) || !validDate(cell.ReviewedAt) {
				return fmt.Errorf("capability %s product %s needs a bounded status, explanation, official HTTPS source and date", c.ID, id)
			}
		}
	}
	if len(ids) == 0 {
		return errors.New("snapshot must contain capability rows")
	}
	for _, o := range s.Opportunities {
		if !ids[o.Capability] || !strings.HasPrefix(o.Issue, "https://github.com/shruggietech/go-schedule/issues/") || o.Note == "" {
			return errors.New("opportunity must identify an existing capability and project issue with rationale")
		}
	}
	return nil
}

func readSnapshot(root, tag string) (Snapshot, []byte, error) {
	var s Snapshot
	data, err := os.ReadFile(filepath.Join(root, "site/data/comparisons", tag+".json"))
	if err != nil {
		return s, nil, fmt.Errorf("read snapshot for %s: %w; add its reviewed comparison before releasing", tag, err)
	}
	if err := decode(data, &s); err != nil {
		return s, nil, fmt.Errorf("decode snapshot %s: %w", tag, err)
	}
	if s.Tag != tag {
		return s, nil, errors.New("snapshot tag differs from its filename")
	}
	return s, data, validateSnapshot(s)
}

func checkSnapshots(root string) error {
	files, err := filepath.Glob(filepath.Join(root, "site/data/comparisons/*.json"))
	if err != nil {
		return fmt.Errorf("find snapshots: %w", err)
	}
	if len(files) == 0 {
		return errors.New("no comparison snapshots found")
	}
	for _, file := range files {
		tag := strings.TrimSuffix(filepath.Base(file), ".json")
		if !tagPattern.MatchString(tag) {
			return fmt.Errorf("invalid snapshot filename %s", file)
		}
		if _, _, err := readSnapshot(root, tag); err != nil {
			return err
		}
	}
	return nil
}

func changelogSection(data []byte, tag string) ([]byte, error) {
	startMarker := "## [" + strings.TrimPrefix(tag, "v") + "]"
	lines := strings.SplitAfter(string(data), "\n")
	var section strings.Builder
	started := false
	for _, line := range lines {
		if strings.HasPrefix(line, startMarker) {
			if started {
				return nil, errors.New("duplicate changelog release section")
			}
			started = true
		} else if started && strings.HasPrefix(line, "## [") {
			break
		}
		if started {
			section.WriteString(line)
		}
	}
	if !started {
		return nil, fmt.Errorf("CHANGELOG.md needs a section for %s", tag)
	}
	return []byte(section.String()), nil
}

func checkRelease(root, tag, commit string) error {
	if !tagPattern.MatchString(tag) || !commitPattern.MatchString(commit) {
		return errors.New("release review requires -tag vMAJOR.MINOR.PATCH and -commit with the resolved 40-character tagged source")
	}
	s, data, err := readSnapshot(root, tag)
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(filepath.Join(root, "site/data/reviews", tag+".json"))
	if err != nil {
		return fmt.Errorf("read review for %s: %w; complete the comparison review before staging", tag, err)
	}
	var r Review
	if err := decode(raw, &r); err != nil {
		return fmt.Errorf("decode release review: %w", err)
	}
	if r.Tag != tag || r.SourceRef != "refs/tags/"+tag || !validDate(r.ReviewedAt) {
		return errors.New("review tag/source_ref/date does not identify this release; review the exact candidate")
	}
	if r.DatasetSHA256 != contentDigest(data) {
		return errors.New("comparison changed after review; update dataset_sha256 after reviewing the candidate")
	}
	change, err := os.ReadFile(filepath.Join(root, "CHANGELOG.md"))
	if err != nil {
		return fmt.Errorf("read changelog: %w", err)
	}
	section, err := changelogSection(change, tag)
	if err != nil {
		return err
	}
	notes, err := os.ReadFile(filepath.Join(root, ".github/release-notes", tag+".md"))
	if err != nil {
		return fmt.Errorf("read release notes: %w", err)
	}
	if r.ChangelogSHA256 != contentDigest(section) || r.NotesSHA256 != contentDigest(notes) {
		return errors.New("release changelog/notes changed after comparison review; refresh reviewed content digests")
	}
	if len(r.Decisions) != len(s.Capabilities) {
		return errors.New("review must record changed or no-change for every capability")
	}
	for _, row := range s.Capabilities {
		if decision := r.Decisions[row.ID]; decision != "changed" && decision != "no-change" {
			return fmt.Errorf("review capability %s needs changed or no-change", row.ID)
		}
	}
	fmt.Printf("comparison review: tag=%s source_ref=%s commit=%s dataset_sha256=%s notes_sha256=%s\n", tag, r.SourceRef, commit, r.DatasetSHA256, r.NotesSHA256)
	return nil
}
