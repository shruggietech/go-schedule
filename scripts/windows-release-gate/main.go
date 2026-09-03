// Command windows-release-gate validates exact Windows candidate evidence.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shruggietech/go-schedule/internal/releasegate"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	switch args[0] {
	case "verify-candidate":
		return runVerifyCandidate(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "verify-bundle":
		return runVerifyBundle(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "windows-release-gate: unknown command %q\n", args[0])
		usage(stderr)
		return 2
	}
}

func runVerifyCandidate(args []string, stdout, stderr io.Writer) int {
	set, options := newFlagSet("verify-candidate", stderr)
	manifestPath := set.String("candidate-manifest", "", "path to the staged candidate manifest")
	artifactPath := set.String("artifact", "", "path to the exact candidate MSI")
	if err := set.Parse(args); err != nil {
		return 2
	}
	if set.NArg() != 0 || *manifestPath == "" || *artifactPath == "" {
		fmt.Fprintln(stderr, "windows-release-gate: verify-candidate requires --candidate-manifest and --artifact")
		return 2
	}
	manifest, err := loadCandidate(*manifestPath)
	if err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: %v\n", err)
		return 2
	}
	failures := releasegate.ValidateCandidate(manifest, *artifactPath, *options)
	if len(failures) != 0 {
		for _, failure := range failures {
			fmt.Fprintf(stderr, "windows-release-gate: %s\n", failure)
		}
		fmt.Fprintf(stderr, "windows-release-gate: FAILED with %d issue(s)\n", len(failures))
		return 1
	}
	fmt.Fprintf(stdout, "windows-release-gate: candidate OK - %s %s (%s)\n", manifest.Tag, manifest.ProductCode, manifest.SHA256)
	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	set, options := newFlagSet("validate", stderr)
	evidencePath := set.String("evidence", "", "path to evidence.json")
	artifactPath := set.String("artifact", "", "path to the exact candidate MSI")
	manifestPath := set.String("candidate-manifest", "", "path to the staged candidate manifest")
	if err := set.Parse(args); err != nil {
		return 2
	}
	if set.NArg() != 0 || *evidencePath == "" || *artifactPath == "" {
		fmt.Fprintln(stderr, "windows-release-gate: validate requires --evidence and --artifact")
		return 2
	}
	evidence, err := loadEvidence(*evidencePath)
	if err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: %v\n", err)
		return 2
	}
	if _, err := os.Stat(*artifactPath); err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: candidate artifact: %v\n", err)
		return 2
	}
	failures := releasegate.Validate(evidence, filepath.Dir(*evidencePath), *artifactPath, *options)
	if *manifestPath != "" {
		manifest, err := loadCandidate(*manifestPath)
		if err != nil {
			fmt.Fprintf(stderr, "windows-release-gate: %v\n", err)
			return 2
		}
		failures = append(failures, releasegate.ValidateCandidateManifest(evidence.Candidate, manifest)...)
	}
	return report(failures, evidence, stdout, stderr)
}

func runVerifyBundle(args []string, stdout, stderr io.Writer) int {
	set, options := newFlagSet("verify-bundle", stderr)
	bundlePath := set.String("bundle", "", "path to the evidence ZIP")
	artifactPath := set.String("artifact", "", "path to the exact candidate MSI")
	manifestPath := set.String("candidate-manifest", "", "path to the staged candidate manifest")
	if err := set.Parse(args); err != nil {
		return 2
	}
	if set.NArg() != 0 || *bundlePath == "" || *artifactPath == "" {
		fmt.Fprintln(stderr, "windows-release-gate: verify-bundle requires --bundle and --artifact")
		return 2
	}
	if _, err := os.Stat(*bundlePath); err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: evidence bundle: %v\n", err)
		return 2
	}
	if _, err := os.Stat(*artifactPath); err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: candidate artifact: %v\n", err)
		return 2
	}
	root, cleanup, err := releasegate.ExtractBundle(*bundlePath)
	if err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: FAILED: %v\n", err)
		return 1
	}
	defer cleanup()
	evidence, err := loadEvidence(filepath.Join(root, "evidence.json"))
	if err != nil {
		fmt.Fprintf(stderr, "windows-release-gate: FAILED: %v\n", err)
		return 1
	}
	failures := releasegate.Validate(evidence, root, *artifactPath, *options)
	failures = append(failures, releasegate.ValidateBundleContents(root, evidence)...)
	if *manifestPath != "" {
		manifest, err := loadCandidate(*manifestPath)
		if err != nil {
			fmt.Fprintf(stderr, "windows-release-gate: %v\n", err)
			return 2
		}
		failures = append(failures, releasegate.ValidateCandidateManifest(evidence.Candidate, manifest)...)
	}
	return report(failures, evidence, stdout, stderr)
}

func newFlagSet(name string, stderr io.Writer) (*flag.FlagSet, *releasegate.ExpectedIdentity) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(stderr)
	options := &releasegate.ExpectedIdentity{}
	set.StringVar(&options.Repository, "repository", "", "expected owner/repository")
	set.StringVar(&options.Tag, "tag", "", "expected vMAJOR.MINOR.PATCH tag")
	set.StringVar(&options.Commit, "commit", "", "expected 40-character tag commit")
	return set, options
}

func loadEvidence(name string) (releasegate.Evidence, error) {
	file, err := os.Open(name)
	if err != nil {
		return releasegate.Evidence{}, fmt.Errorf("open evidence: %w", err)
	}
	defer file.Close()
	evidence, err := releasegate.DecodeEvidence(file)
	if err != nil {
		return releasegate.Evidence{}, err
	}
	return evidence, nil
}

func loadCandidate(name string) (releasegate.Candidate, error) {
	file, err := os.Open(name)
	if err != nil {
		return releasegate.Candidate{}, fmt.Errorf("open candidate manifest: %w", err)
	}
	defer file.Close()
	candidate, err := releasegate.DecodeCandidate(file)
	if err != nil {
		return releasegate.Candidate{}, err
	}
	return candidate, nil
}

func report(failures []string, evidence releasegate.Evidence, stdout, stderr io.Writer) int {
	if len(failures) != 0 {
		for _, failure := range failures {
			fmt.Fprintf(stderr, "windows-release-gate: %s\n", failure)
		}
		fmt.Fprintf(stderr, "windows-release-gate: FAILED with %d issue(s)\n", len(failures))
		return 1
	}
	fmt.Fprintf(stdout, "windows-release-gate: OK - %s %s (%s) with %d required observations\n",
		evidence.Candidate.Tag, evidence.Candidate.ProductCode, evidence.Candidate.SHA256, len(evidence.Observations))
	return 0
}

func usage(writer io.Writer) {
	fmt.Fprintln(writer, "usage:")
	fmt.Fprintln(writer, "  windows-release-gate verify-candidate --candidate-manifest JSON --artifact MSI [--repository OWNER/REPO --tag TAG --commit SHA]")
	fmt.Fprintln(writer, "  windows-release-gate validate --evidence FILE --artifact MSI [--candidate-manifest JSON --repository OWNER/REPO --tag TAG --commit SHA]")
	fmt.Fprintln(writer, "  windows-release-gate verify-bundle --bundle ZIP --artifact MSI [--candidate-manifest JSON --repository OWNER/REPO --tag TAG --commit SHA]")
}
