package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPowerShellCompletionHelpExplainsProfileInvocation(t *testing.T) {
	var out, diagnostic bytes.Buffer
	code := ExecuteWithDependencies(context.Background(), []string{"completion", "powershell", "--help"}, Dependencies{Out: &out, Err: &diagnostic})
	if code != 0 || diagnostic.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q", code, diagnostic.String())
	}
	const invocation = "gosched completion powershell | Out-String | Invoke-Expression"
	help := out.String()
	if strings.Count(help, invocation) != 2 || !strings.Contains(help, "add the following invocation as a line in your PowerShell profile ($PROFILE)") || strings.Contains(help, "add the output") {
		t.Fatalf("PowerShell help does not distinguish current loading from a persisted invocation: %s", help)
	}
}

func TestCompletionGuidancePreservesGeneratorAndOtherShells(t *testing.T) {
	for _, noDescriptions := range []bool{false, true} {
		var out, diagnostic, expected bytes.Buffer
		args := []string{"completion", "powershell"}
		if noDescriptions {
			args = append(args, "--no-descriptions")
		}
		code := ExecuteWithDependencies(context.Background(), args, Dependencies{Out: &out, Err: &diagnostic})
		root := newRoot()
		var err error
		if noDescriptions {
			err = root.GenPowerShellCompletion(&expected)
		} else {
			err = root.GenPowerShellCompletionWithDesc(&expected)
		}
		if err != nil || code != 0 || diagnostic.Len() != 0 || !bytes.Equal(out.Bytes(), expected.Bytes()) {
			t.Fatalf("noDescriptions=%v exit=%d stderr=%q generator error=%v", noDescriptions, code, diagnostic.String(), err)
		}
	}
	standard := &cobra.Command{Use: "gosched"}
	standard.AddCommand(&cobra.Command{Use: "fixture"})
	standard.InitDefaultCompletionCmd()
	actual := newRoot()
	for _, shell := range []string{"bash", "fish", "zsh"} {
		path := []string{"completion", shell}
		got, _, err := actual.Find(path)
		if err != nil {
			t.Fatal(err)
		}
		want, _, err := standard.Find(path)
		if err != nil {
			t.Fatal(err)
		}
		if got.Long != want.Long {
			t.Errorf("%s completion guidance unexpectedly changed", shell)
		}
	}
}
