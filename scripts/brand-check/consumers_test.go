package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRepositoryConsumerFailures(t *testing.T) {
	tests := []struct {
		name string
		edit func(*testing.T, *fixture)
		want string
	}{
		{"missing", func(t *testing.T, f *fixture) {
			f.mapFile.Mappings = []consumerMapping{{Source: "README.md", Targets: []string{"docs/logo.md"}, Purpose: "test"}}
			f.flush(t)
		}, "missing consumer"},
		{"different", func(t *testing.T, f *fixture) {
			f.write(t, "docs/logo.md", []byte("different\n"))
			f.mapFile.Mappings = []consumerMapping{{Source: "README.md", Targets: []string{"docs/logo.md"}, Purpose: "test"}}
			f.flush(t)
		}, "differs from canonical"},
		{"duplicate target", func(t *testing.T, f *fixture) {
			f.write(t, "docs/logo.md", []byte("# Brand\n"))
			f.mapFile.Mappings = []consumerMapping{
				{Source: "README.md", Targets: []string{"docs/logo.md"}, Purpose: "one"},
				{Source: "README.md", Targets: []string{"docs/logo.md"}, Purpose: "two"},
			}
			f.flush(t)
		}, "duplicate consumer target"},
		{"malformed", func(t *testing.T, f *fixture) { f.write(t, "brand/repository-consumers.json", []byte("{")) }, "parse"},
		{"missing control", func(t *testing.T, f *fixture) { os.Remove(filepath.Join(f.root, "brand", "REPOSITORY.md")) }, "missing required control file"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := newFixture(t)
			test.edit(t, f)
			_, failures := checkRepository(f.root)
			requireFailure(t, failures, test.want)
		})
	}
}

func TestCheckRepositoryConsumerSuccess(t *testing.T) {
	f := newFixture(t)
	f.write(t, "docs/logo.md", []byte("# Brand\n"))
	f.mapFile.Mappings = []consumerMapping{{Source: "README.md", Targets: []string{"docs/logo.md"}, Purpose: "test"}}
	f.flush(t)
	stats, failures := checkRepository(f.root)
	if len(failures) != 0 {
		t.Fatalf("unexpected failures: %v", failures)
	}
	if stats.Consumers != 1 {
		t.Fatalf("consumers = %d, want 1", stats.Consumers)
	}
}
