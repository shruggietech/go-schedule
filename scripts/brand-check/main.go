// Command brand-check validates the canonical brand kit and every declared
// repository consumer without requiring the optional asset-generation stack.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

type artifact struct {
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256"`
}

type manifestFile struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	GeneratedAt string     `json:"generated_at"`
	Files       []artifact `json:"files"`
}

type consumerMapping struct {
	Source  string   `json:"source"`
	Targets []string `json:"targets"`
	Purpose string   `json:"purpose"`
}

type consumerFile struct {
	Version  int               `json:"version"`
	Mappings []consumerMapping `json:"mappings"`
}

type checkStats struct {
	Artifacts int
	SVGs      int
	Consumers int
}

var (
	titlePattern      = regexp.MustCompile(`(?i)<title(?:\s|>)`)
	liveTextPattern   = regexp.MustCompile(`(?i)<text(?:\s|>)`)
	fontFamilyPattern = regexp.MustCompile(`(?i)font-family\s*(?:=|:)`)
	drivePathPattern  = regexp.MustCompile(`^[A-Za-z]:`)
)

var textExtensions = map[string]bool{
	".css": true, ".desktop": true, ".html": true, ".js": true,
	".json": true, ".jsx": true, ".md": true, ".py": true, ".svg": true,
}

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: brand-check [repository-root]")
		os.Exit(2)
	}
	root := "."
	if len(os.Args) == 2 {
		root = os.Args[1]
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "brand-check: resolve repository root: %v\n", err)
		os.Exit(2)
	}
	stats, failures := checkRepository(absolute)
	if len(failures) != 0 {
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		fmt.Fprintf(os.Stderr, "brand-check: FAILED with %d issue(s)\n", len(failures))
		os.Exit(1)
	}
	fmt.Printf("brand-check: OK - %d artifacts, %d SVGs, %d consumers\n", stats.Artifacts, stats.SVGs, stats.Consumers)
}

func checkRepository(root string) (checkStats, []string) {
	var stats checkStats
	var failures []string
	brandRoot := filepath.Join(root, "brand")

	required := []string{"README.md", "REPOSITORY.md", "VERIFY.md", "brand-guide.pdf", "manifest.json", "repository-consumers.json"}
	for _, relative := range required {
		if info, err := os.Stat(filepath.Join(brandRoot, filepath.FromSlash(relative))); err != nil || !info.Mode().IsRegular() {
			failures = append(failures, fmt.Sprintf("brand/%s: missing required control file", relative))
		}
	}

	manifest, manifestFailures := readManifest(filepath.Join(brandRoot, "manifest.json"))
	failures = append(failures, manifestFailures...)
	known := map[string]bool{
		"manifest.json": true, "VERIFY.md": true,
		"REPOSITORY.md": true, "repository-consumers.json": true,
	}

	for _, item := range manifest.Files {
		if failure := validateRelativePath(item.Path); failure != "" {
			failures = append(failures, fmt.Sprintf("brand/manifest.json: unsafe artifact path %q: %s", item.Path, failure))
			continue
		}
		if known[item.Path] {
			failures = append(failures, fmt.Sprintf("brand/manifest.json: duplicate artifact path: %s", item.Path))
			continue
		}
		known[item.Path] = true
		stats.Artifacts++
		full := filepath.Join(brandRoot, filepath.FromSlash(item.Path))
		actualBytes, digest, err := hashFile(full)
		if err != nil {
			if os.IsNotExist(err) {
				failures = append(failures, fmt.Sprintf("brand/%s: missing artifact", item.Path))
			} else {
				failures = append(failures, fmt.Sprintf("brand/%s: read artifact: %v", item.Path, err))
			}
			continue
		}
		if actualBytes != item.Bytes {
			failures = append(failures, fmt.Sprintf("brand/%s: byte length %d, manifest requires %d", item.Path, actualBytes, item.Bytes))
		}
		if !strings.EqualFold(digest, item.SHA256) {
			failures = append(failures, fmt.Sprintf("brand/%s: SHA-256 %s, manifest requires %s", item.Path, digest, item.SHA256))
		}
		if textExtensions[strings.ToLower(filepath.Ext(item.Path))] {
			data, err := os.ReadFile(full)
			if err == nil {
				failures = append(failures, validateText("brand/"+item.Path, data)...)
				if strings.EqualFold(filepath.Ext(item.Path), ".svg") {
					stats.SVGs++
					failures = append(failures, validateSVG("brand/"+item.Path, data)...)
				}
			}
		}
	}

	walkErr := filepath.WalkDir(brandRoot, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(brandRoot, file)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if !known[relative] {
			failures = append(failures, fmt.Sprintf("brand/%s: unlisted artifact", relative))
		}
		return nil
	})
	if walkErr != nil {
		failures = append(failures, fmt.Sprintf("brand: inventory walk: %v", walkErr))
	}

	consumers, consumerFailures := readConsumers(filepath.Join(brandRoot, "repository-consumers.json"))
	failures = append(failures, consumerFailures...)
	seenSources := make(map[string]bool)
	seenTargets := make(map[string]bool)
	for _, mapping := range consumers.Mappings {
		if failure := validateRelativePath(mapping.Source); failure != "" {
			failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: unsafe source path %q: %s", mapping.Source, failure))
			continue
		}
		if seenSources[mapping.Source] {
			failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: duplicate source mapping: %s", mapping.Source))
		}
		seenSources[mapping.Source] = true
		if !known[mapping.Source] {
			failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: source is not a canonical artifact: %s", mapping.Source))
		}
		if strings.TrimSpace(mapping.Purpose) == "" {
			failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: mapping %s has no purpose", mapping.Source))
		}
		if len(mapping.Targets) == 0 {
			failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: mapping %s has no targets", mapping.Source))
		}
		sourcePath := filepath.Join(brandRoot, filepath.FromSlash(mapping.Source))
		for _, target := range mapping.Targets {
			if failure := validateRelativePath(target); failure != "" {
				failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: unsafe target path %q: %s", target, failure))
				continue
			}
			if seenTargets[target] {
				failures = append(failures, fmt.Sprintf("brand/repository-consumers.json: duplicate consumer target: %s", target))
				continue
			}
			seenTargets[target] = true
			stats.Consumers++
			targetPath := filepath.Join(root, filepath.FromSlash(target))
			same, err := equalFiles(sourcePath, targetPath)
			if err != nil {
				if os.IsNotExist(err) {
					failures = append(failures, fmt.Sprintf("%s: missing consumer for canonical brand/%s", target, mapping.Source))
				} else {
					failures = append(failures, fmt.Sprintf("%s: compare with canonical brand/%s: %v", target, mapping.Source, err))
				}
			} else if !same {
				failures = append(failures, fmt.Sprintf("%s: differs from canonical brand/%s", target, mapping.Source))
			}
		}
	}

	sort.Strings(failures)
	return stats, failures
}

func readManifest(file string) (manifestFile, []string) {
	var result manifestFile
	data, err := os.ReadFile(file)
	if err != nil {
		return result, []string{fmt.Sprintf("brand/manifest.json: read: %v", err)}
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, []string{fmt.Sprintf("brand/manifest.json: parse: %v", err)}
	}
	if strings.TrimSpace(result.Name) == "" || strings.TrimSpace(result.Version) == "" || len(result.Files) == 0 {
		return result, []string{"brand/manifest.json: name, version, and files are required"}
	}
	return result, nil
}

func readConsumers(file string) (consumerFile, []string) {
	var result consumerFile
	data, err := os.ReadFile(file)
	if err != nil {
		return result, []string{fmt.Sprintf("brand/repository-consumers.json: read: %v", err)}
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, []string{fmt.Sprintf("brand/repository-consumers.json: parse: %v", err)}
	}
	if result.Version != 1 {
		return result, []string{fmt.Sprintf("brand/repository-consumers.json: unsupported version %d", result.Version)}
	}
	return result, nil
}

func validateRelativePath(value string) string {
	if value == "" {
		return "empty"
	}
	if strings.Contains(value, `\`) {
		return "must use forward slashes"
	}
	if path.IsAbs(value) || drivePathPattern.MatchString(value) || path.Clean(value) != value || value == "." || strings.HasPrefix(value, "../") {
		return "must be a clean relative path"
	}
	return ""
}

func hashFile(file string) (int64, string, error) {
	stream, err := os.Open(file)
	if err != nil {
		return 0, "", err
	}
	defer stream.Close()
	digest := sha256.New()
	length, err := io.Copy(digest, stream)
	if err != nil {
		return 0, "", err
	}
	return length, hex.EncodeToString(digest.Sum(nil)), nil
}

func equalFiles(left, right string) (bool, error) {
	leftData, err := os.ReadFile(left)
	if err != nil {
		return false, err
	}
	rightData, err := os.ReadFile(right)
	if err != nil {
		return false, err
	}
	return bytes.Equal(leftData, rightData), nil
}

func validateText(name string, data []byte) []string {
	var failures []string
	if bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		failures = append(failures, fmt.Sprintf("%s: UTF-8 BOM is not allowed", name))
	}
	if !utf8.Valid(data) {
		return append(failures, fmt.Sprintf("%s: invalid UTF-8", name))
	}
	text := string(data)
	for _, marker := range []string{"\u00c3", "\u00c2", "\ufffd", "\u00e2\u20ac"} {
		if strings.Contains(text, marker) {
			failures = append(failures, fmt.Sprintf("%s: possible mojibake marker %q", name, marker))
			break
		}
	}
	return failures
}

func validateSVG(name string, data []byte) []string {
	var failures []string
	if !titlePattern.Match(data) {
		failures = append(failures, fmt.Sprintf("%s: SVG lacks an accessible title", name))
	}
	if liveTextPattern.Match(data) {
		failures = append(failures, fmt.Sprintf("%s: SVG contains live text", name))
	}
	if fontFamilyPattern.Match(data) {
		failures = append(failures, fmt.Sprintf("%s: SVG contains a font-family dependency", name))
	}
	return failures
}
