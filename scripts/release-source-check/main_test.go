package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testRepo = "shruggietech/go-schedule"
const testTag = "v1.6.0"
const testHeading = "## [1.6.0] - 2026-10-22"
const testLink = "Read the [full changelog](https://github.com/shruggietech/go-schedule/blob/v1.6.0/CHANGELOG.md#160---2026-10-22) for every change."
const testBadge = "<a href=\"https://github.com/shruggietech/go-schedule/releases/latest\"><img alt=\"Release\" src=\"https://img.shields.io/github/v/release/shruggietech/go-schedule?color=58A6FF\"></a>"

func fixture(t *testing.T, tag, readme, changelog, notes string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".github/release-notes"), 0755); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{"README.md": readme, "CHANGELOG.md": changelog, ".github/release-notes/" + tag + ".md": notes} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestSourceMetadata(t *testing.T) {
	goodReadme := testBadge + "\ndaemon ok (version v1.6.0)\n"
	cases := []struct{ name, readme, heading, notes, want string }{
		{"valid", goodReadme, testHeading, testLink, ""},
		{"badge whitespace", "  " + testBadge + " \ndaemon ok (version v1.6.0)\n", testHeading, testLink, ""},
		{"CRLF", strings.ReplaceAll(goodReadme, "\n", "\r\n"), testHeading + "\r", testLink + "\r\n", ""},
		{"no badge", "daemon ok (version v1.6.0)", testHeading, testLink, "README.md"},
		{"duplicate badge", testBadge + "\n" + goodReadme, testHeading, testLink, "README.md"},
		{"two badges same line", testBadge + goodReadme, testHeading, testLink, "README.md"},
		{"static badge", strings.ReplaceAll(goodReadme, "github/v/release/", "badge/release/"), testHeading, testLink, "README.md"},
		{"wrong health", strings.ReplaceAll(goodReadme, "version v1.6.0", "version v1.5.1"), testHeading, testLink, "README.md"},
		{"duplicate health", goodReadme + "daemon ok (version v1.6.0)\n", testHeading, testLink, "README.md"},
		{"no heading", goodReadme, "## [1.5.1] - 2026-10-22", testLink, "CHANGELOG.md"},
		{"duplicate heading", goodReadme, testHeading + "\n" + testHeading, testLink, "CHANGELOG.md"},
		{"missing date", goodReadme, "## [1.6.0]", testLink, "CHANGELOG.md"},
		{"duplicate malformed heading", goodReadme, testHeading + "\n## [1.6.0]", testLink, "CHANGELOG.md"},
		{"not date", goodReadme, "## [1.6.0] - not-a-date", strings.ReplaceAll(testLink, "2026-10-22", "not-a-date"), "CHANGELOG.md"},
		{"impossible date", goodReadme, "## [1.6.0] - 2026-02-30", strings.ReplaceAll(testLink, "2026-10-22", "2026-02-30"), "CHANGELOG.md"},
		{"single digit month", goodReadme, "## [1.6.0] - 2026-2-03", testLink, "CHANGELOG.md"},
		{"embedded heading", goodReadme, "prefix " + testHeading, testLink, "CHANGELOG.md"},
		{"extra heading text", goodReadme, testHeading + " extra", testLink, "CHANGELOG.md"},
		{"wrong tag", goodReadme, testHeading, strings.ReplaceAll(testLink, "blob/v1.6.0", "blob/v1.5.1"), "release-notes"},
		{"wrong repository", goodReadme, testHeading, strings.ReplaceAll(testLink, "shruggietech/go-schedule", "other/repo"), "release-notes"},
		{"wrong anchor", goodReadme, testHeading, strings.ReplaceAll(testLink, "160---2026", "160--2026"), "release-notes"},
		{"duplicate link", goodReadme, testHeading, testLink + "\n" + testLink, "release-notes"},
		{"embedded link", goodReadme, testHeading, "prefix " + testLink, "release-notes"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := fixture(t, testTag, c.readme, c.heading+"\n", c.notes+"\n")
			err := checkSource(root, testRepo, testTag)
			if c.want == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("error=%v, want %s", err, c.want)
			}
		})
	}
}

func TestMissingInputsAndArguments(t *testing.T) {
	for _, name := range []string{"README.md", "CHANGELOG.md", ".github/release-notes/v1.6.0.md"} {
		t.Run(name, func(t *testing.T) {
			root := fixture(t, testTag, testBadge+"\ndaemon ok (version v1.6.0)\n", testHeading+"\n", testLink+"\n")
			if err := os.Remove(filepath.Join(root, name)); err != nil {
				t.Fatal(err)
			}
			if err := checkSource(root, testRepo, testTag); err == nil || !strings.Contains(err.Error(), name) {
				t.Fatalf("missing file accepted: %v", err)
			}
		})
	}
	for _, c := range []struct{ repo, tag string }{{testRepo, ""}, {testRepo, "../../bad"}, {testRepo, "v1.6.0/../../bad"}, {testRepo, "v1.6.0-rc1"}, {"", testTag}, {"https://github.com/other/repo", testTag}, {"../other", testTag}, {"owner/repo/extra", testTag}} {
		if err := checkSource(t.TempDir(), c.repo, c.tag); err == nil {
			t.Fatalf("invalid inputs accepted: %+v", c)
		}
	}
}

func TestCandidateTag(t *testing.T) {
	for _, c := range []struct {
		name, changelog, want string
	}{
		{"later minor", "## [Unreleased]\n\n## [1.7.0] - 2026-10-27\n## [1.6.0] - 2026-10-22\n", "v1.7.0"},
		{"future version", "## [9.8.7] - 2027-01-01\n", "v9.8.7"},
		{"missing candidate", "## [Unreleased]\n", ""},
		{"malformed first candidate", "## [broken] - 2026-10-27\n## [1.6.0] - 2026-10-22\n", ""},
		{"unclosed first candidate", "## [1.7.0\n## [1.6.0] - 2026-10-22\n", ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := fixture(t, testTag, "", c.changelog, "")
			tag, err := candidateTag(root)
			if c.want == "" {
				if err == nil || !strings.Contains(err.Error(), "CHANGELOG.md") {
					t.Fatalf("tag=%q error=%v, want attributable candidate error", tag, err)
				}
				return
			}
			if err != nil || tag != c.want {
				t.Fatalf("tag=%q error=%v, want %q", tag, err, c.want)
			}
		})
	}
}

// candidateTag follows the first declared version while leaving exact metadata
// checks to checkSource. A malformed new heading must not fall back to history.
func candidateTag(root string) (string, error) {
	lines, err := readLines(root, "CHANGELOG.md")
	if err != nil {
		return "", err
	}
	for _, line := range lines {
		version, ok := strings.CutPrefix(line, "## [")
		if !ok {
			continue
		}
		version, _, closed := strings.Cut(version, "]")
		if closed && version == "Unreleased" {
			continue
		}
		tag := "v" + version
		if !closed || !tagPattern.MatchString(tag) {
			return "", fmt.Errorf("CHANGELOG.md: invalid current candidate heading %q", line)
		}
		return tag, nil
	}
	return "", fmt.Errorf("CHANGELOG.md: no current candidate heading")
}

func TestCurrentCandidate(t *testing.T) {
	tag, err := candidateTag("../..")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("current prepared source: %s", tag)
	if err := checkSource("../..", testRepo, tag); err != nil {
		t.Fatal(err)
	}
}

func TestDeclaredCandidateKeepsVersionChecks(t *testing.T) {
	for _, tag := range []string{"v1.7.0", "v9.8.7"} {
		for _, wrongHealth := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/wrong-health=%t", tag, wrongHealth), func(t *testing.T) {
				version := strings.TrimPrefix(tag, "v")
				health := tag
				if wrongHealth {
					health = "v1.6.0"
				}
				heading := "## [Unreleased]\n\n## [" + version + "] - 2026-10-27\n\n" + testHeading + "\n"
				link := "Read the [full changelog](https://github.com/" + testRepo + "/blob/" + tag + "/CHANGELOG.md#" + strings.ReplaceAll(version, ".", "") + "---2026-10-27) for every change.\n"
				root := fixture(t, tag, testBadge+"\ndaemon ok (version "+health+")\n", heading, link)
				selected, err := candidateTag(root)
				if err != nil || selected != tag {
					t.Fatalf("selected=%q error=%v, want %q", selected, err, tag)
				}
				err = checkSource(root, testRepo, selected)
				if wrongHealth {
					if err == nil || !strings.Contains(err.Error(), "README.md") {
						t.Fatalf("mismatched current health accepted: %v", err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

// Release workflows inject the complete tag, and daemon health returns that
// value unchanged. The source example must match those actual output bytes.
func TestHealthExampleMatchesReleaseInjectedVersion(t *testing.T) {
	const tag = "v2.0.0"
	heading := "## [2.0.0] - 2026-10-31\n"
	link := "Read the [full changelog](https://github.com/" + testRepo + "/blob/v2.0.0/CHANGELOG.md#200---2026-10-31) for every change.\n"
	for _, version := range []string{tag, "2.0.0"} {
		t.Run(version, func(t *testing.T) {
			root := fixture(t, tag, testBadge+"\ndaemon ok (version "+version+")\n", heading, link)
			err := checkSource(root, testRepo, tag)
			if version == tag {
				if err != nil {
					t.Fatalf("release-injected health rejected: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "README.md") {
				t.Fatalf("unprefixed health accepted: %v", err)
			}
		})
	}
}
