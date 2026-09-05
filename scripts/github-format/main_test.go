package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatMarkdownUnwrapsProseAndListItems(t *testing.T) {
	input := "# Title\n\nA paragraph that was\nwrapped by an author.\n\n- A list item that was\n  wrapped as well.\n"
	want := "# Title\n\nA paragraph that was wrapped by an author.\n\n- A list item that was wrapped as well.\n"
	if got := formatMarkdown(input); got != want {
		t.Fatalf("formatMarkdown() = %q, want %q", got, want)
	}
}

func TestFormatMarkdownPreservesNestedBlocks(t *testing.T) {
	input := "- Item text that\n  continues here.\n\n  ## Nested heading\n\n  ```text\n  literal\n  ```\n\n  - Nested item\n"
	want := "- Item text that continues here.\n\n  ## Nested heading\n\n  ```text\n  literal\n  ```\n\n  - Nested item\n"
	if got := formatMarkdown(input); got != want {
		t.Fatalf("nested blocks changed:\n%s", got)
	}
}

func TestFormatMarkdownPreservesLiteralStructures(t *testing.T) {
	input := "---\ntitle: Example\n---\n\n| A | B |\n|---|---|\n| 1 | 2 |\n\n```text\nline one\nline two\n```\n"
	if got := formatMarkdown(input); got != input {
		t.Fatalf("literal structures changed:\n%s", got)
	}
}

func TestFormatMarkdownPreservesNestedFences(t *testing.T) {
	input := "````markdown\n```mermaid\nflowchart TB\n  A --> B\n```\n````\n"
	if got := formatMarkdown(input); got != input {
		t.Fatalf("nested fences changed:\n%s", got)
	}
}

func TestFormatMarkdownPreservesAlertMarkerAndUnwrapsBody(t *testing.T) {
	input := "> [!CAUTION]\n> A warning that was\n> wrapped by an author.\n"
	want := "> [!CAUTION]\n> A warning that was wrapped by an author.\n"
	if got := formatMarkdown(input); got != want {
		t.Fatalf("alert changed:\n%s", got)
	}
}

func TestFormatMarkdownUnwrapsIssueReferenceContinuation(t *testing.T) {
	input := "- A list item referencing\n  #123 without treating it as a heading.\n"
	want := "- A list item referencing #123 without treating it as a heading.\n"
	if got := formatMarkdown(input); got != want {
		t.Fatalf("issue reference continuation changed:\n%s", got)
	}
}

func TestFormatMarkdownPreservesSemanticMetadataLines(t *testing.T) {
	input := "**Feature Branch**: `example`\n**Created**: 2026-09-05\n**Status**: Draft\n\n**Label:** First value.  \n**Other:** Second value.\n"
	if got := formatMarkdown(input); got != input {
		t.Fatalf("semantic metadata changed:\n%s", got)
	}
}

func TestFormatMarkdownUnwrapsMetadataProse(t *testing.T) {
	input := "**Storage**: Embedded database for tasks, groups,\nschedules, triggers, and history.\n"
	want := "**Storage**: Embedded database for tasks, groups, schedules, triggers, and history.\n"
	if got := formatMarkdown(input); got != want {
		t.Fatalf("metadata prose changed:\n%s", got)
	}
}

func TestReplaceEmDashesUsesPlainPunctuationAndTableFallback(t *testing.T) {
	input := "A clause " + emDash + " another clause.\n" + emDash + " continued.\n| Value | " + emDash + " |\n| Result | PASS " + emDash + " verified |\n"
	want := "A clause, another clause.\ncontinued.\n| Value | N/A |\n| Result | PASS, verified |\n"
	if got := replaceEmDashes(input); got != want {
		t.Fatalf("replaceEmDashes() = %q, want %q", got, want)
	}
}

func TestRepositorySourcesContainNoForbiddenFormatting(t *testing.T) {
	problems, err := checkRepository("../..", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("run go run ./scripts/github-format -fix; formatting violations: %v", problems)
	}
}

func TestRepositoryCheckRejectsInvalidUTF8AndBOM(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "invalid.md"), []byte{0xff}, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bom.md"), []byte("\ufeff# Title\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	problems, err := checkRepository(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 2 {
		t.Fatalf("problems = %v, want invalid UTF-8 and BOM", problems)
	}
}
