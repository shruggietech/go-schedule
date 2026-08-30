package main

import (
	"path/filepath"
	"testing"
)

func TestCheckRepositoryTextAndSVGFailures(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"bom", append([]byte{0xef, 0xbb, 0xbf}, []byte("# Brand\n")...), "UTF-8 BOM"},
		{"invalid UTF-8", []byte{0xff, 0xfe}, "invalid UTF-8"},
		{"mojibake", []byte("Fran\u00c3\u00a7ais\n"), "mojibake"},
		{"missing title", []byte(`<svg xmlns="http://www.w3.org/2000/svg"><path/></svg>`), "accessible title"},
		{"live text", []byte(`<svg xmlns="http://www.w3.org/2000/svg"><title>X</title><text>X</text></svg>`), "live text"},
		{"font family", []byte(`<svg xmlns="http://www.w3.org/2000/svg"><title>X</title><path style="font-family:sans-serif"/></svg>`), "font-family"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := newFixture(t)
			path := "README.md"
			if test.name == "missing title" || test.name == "live text" || test.name == "font family" {
				path = "logos/mark.svg"
			}
			f.write(t, filepath.ToSlash(filepath.Join("brand", path)), test.data)
			_, failures := checkRepository(f.root)
			requireFailure(t, failures, test.want)
		})
	}
}
