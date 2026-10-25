// Command release-source-check validates local release metadata before a tag is
// created and in the tagged Release workflow, without mutating release state.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var tagPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
var repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*/[A-Za-z0-9][A-Za-z0-9._-]*$`)
var datePattern = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)

func main() {
	root := flag.String("root", ".", "repository root containing README.md and CHANGELOG.md")
	repository := flag.String("repository", "", "GitHub owner/repository for canonical release links")
	tag := flag.String("tag", "", "selected vMAJOR.MINOR.PATCH source version (does not verify a Git tag)")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "release-source-check: unexpected operands; use -root, -repository and -tag")
		os.Exit(1)
	}
	if err := checkSource(*root, *repository, *tag); err != nil {
		fmt.Fprintf(os.Stderr, "release-source-check: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("release-source-check: source metadata OK repository=%s tag=%s (no tag or package identity verified)\n", *repository, *tag)
}

func readLines(root, name string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		return nil, fmt.Errorf("%s: read source file: %w; restore the required release-preparation file", name, err)
	}
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n"), nil
}

func countContaining(lines []string, marker string) int {
	count := 0
	for _, line := range lines {
		if strings.Contains(line, marker) {
			count++
		}
	}
	return count
}

func checkSource(root, repository, tag string) error {
	if !tagPattern.MatchString(tag) || !repositoryPattern.MatchString(repository) {
		return errors.New("use -repository owner/repository and -tag vMAJOR.MINOR.PATCH; source paths and links require these exact identities")
	}
	readme, err := readLines(root, "README.md")
	if err != nil {
		return err
	}
	badge := fmt.Sprintf(`<a href="https://github.com/%s/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/%s?color=58A6FF"></a>`, repository, repository)
	if countContaining(readme, `<img alt="Release"`) != 1 {
		return errors.New("README.md: require exactly one Release badge line; add the canonical badge if missing and remove duplicates")
	}
	badgeCount, healthCount := 0, 0
	plain := strings.TrimPrefix(tag, "v")
	for _, line := range readme {
		if strings.TrimSpace(line) == badge {
			badgeCount++
		}
		if line == "daemon ok (version "+plain+")" {
			healthCount++
		}
	}
	if badgeCount != 1 {
		return errors.New("README.md: Release badge must use the canonical latest published GitHub release URL")
	}
	if countContaining(readme, "daemon ok (version ") != 1 || healthCount != 1 {
		return fmt.Errorf("README.md: require exactly one health example: daemon ok (version %s)", plain)
	}
	changelog, err := readLines(root, "CHANGELOG.md")
	if err != nil {
		return err
	}
	prefix := "## [" + plain + "] - "
	if countContaining(changelog, "## ["+plain+"]") != 1 {
		return fmt.Errorf("CHANGELOG.md: require exactly one dated %s release heading", tag)
	}
	var heading string
	for _, line := range changelog {
		if strings.Contains(line, prefix) {
			heading = line
			break
		}
	}
	date := strings.TrimPrefix(heading, prefix)
	if !strings.HasPrefix(heading, prefix) || !datePattern.MatchString(date) {
		return fmt.Errorf("CHANGELOG.md: use the exact heading %sYYYY-MM-DD with a calendar-valid date", prefix)
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("CHANGELOG.md: invalid release date %q: %w; use a calendar-valid YYYY-MM-DD date", date, err)
	}
	notesPath := filepath.ToSlash(filepath.Join(".github/release-notes", tag+".md"))
	notes, err := readLines(root, notesPath)
	if err != nil {
		return err
	}
	anchor := strings.ReplaceAll(plain, ".", "") + "---" + date
	link := fmt.Sprintf("Read the [full changelog](https://github.com/%s/blob/%s/CHANGELOG.md#%s) for every change.", repository, tag, anchor)
	count := 0
	for _, line := range notes {
		if line == link {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("%s: require exactly one whole-line canonical changelog link: %s", notesPath, link)
	}
	return nil
}
