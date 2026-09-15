// Command windows-qualification-session prepares disposable Windows sessions.
// It never executes supplied packages or publishes release material.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

type input struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256"`
	Source string `json:"source"`
}

type manifest struct {
	Schema     int     `json:"schema_version"`
	Repository string  `json:"repository"`
	Tag        string  `json:"tag"`
	Commit     string  `json:"commit"`
	RunID      int64   `json:"run_id"`
	RunAttempt int     `json:"run_attempt"`
	Inputs     []input `json:"inputs"`
}

var filenames = map[string]string{"candidate": "go-schedule_v1.4.0_windows_amd64.msi", "baseline": "go-schedule_v1.1.1_windows_amd64.msi", "powershell": "powershell.zip", "webview2": "webview2.exe", "bootstrap": "Start-QualificationGuest.ps1", "collector": "Invoke-ReleaseCandidateAttended.ps1"}

func main() {
	path := flag.String("manifest", "", "absolute input manifest JSON")
	out := flag.String("output", "", "absolute new output directory")
	flag.Parse()
	if err := run(*path, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(path, out string) error {
	if !filepath.IsAbs(path) {
		return errors.New("manifest path must be absolute")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var m manifest
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	if err := d.Decode(&m); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("trailing manifest content")
	}
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	return prepare(m, out, filepath.Join(root, "test", "windows"))
}

func rejectLinks(path string) error {
	for p := filepath.Clean(path); ; p = filepath.Dir(p) {
		info, err := os.Lstat(p)
		if err != nil {
			return err
		}
		linked, err := isLinkedPath(p, info)
		if err != nil {
			return err
		}
		if linked {
			return errors.New("linked input/output path refused")
		}
		if filepath.Dir(p) == p {
			return nil
		}
	}
}

func fileIdentity(path string) (int64, string, error) {
	if !filepath.IsAbs(path) {
		return 0, "", errors.New("input path must be absolute")
	}
	if err := rejectLinks(path); err != nil {
		return 0, "", err
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0, "", err
	}
	if !info.Mode().IsRegular() {
		return 0, "", errors.New("input must be a regular file")
	}
	h := sha256.New()
	n, err := io.Copy(h, f)
	return n, hex.EncodeToString(h.Sum(nil)), err
}

func prepare(m manifest, out, helpers string) error {
	if m.Schema != 1 || m.Repository != "shruggietech/go-schedule" || m.Tag != "v1.4.0" || !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(m.Commit) || m.RunID <= 0 || m.RunAttempt <= 0 {
		return errors.New("invalid candidate provenance")
	}
	if !filepath.IsAbs(out) || filepath.Dir(filepath.Clean(out)) == filepath.Clean(out) {
		return errors.New("output must be an absolute new non-root directory")
	}
	if _, err := os.Lstat(out); !os.IsNotExist(err) {
		return errors.New("output already exists or cannot be inspected")
	}
	if err := rejectLinks(filepath.Dir(out)); err != nil {
		return errors.New("output parent must exist and must not be linked")
	}
	if len(m.Inputs) != 4 {
		return errors.New("exactly four package inputs required")
	}
	seen := map[string]bool{}
	var identities []os.FileInfo
	for _, in := range m.Inputs {
		if _, ok := filenames[in.Role]; !ok || in.Role == "bootstrap" || in.Role == "collector" || seen[in.Role] {
			return errors.New("unknown or duplicate input role")
		}
		seen[in.Role] = true
		if in.Source == "" || in.Bytes <= 0 || !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(in.SHA256) {
			return fmt.Errorf("missing identity for %s", in.Role)
		}
		u, err := url.Parse(in.Source)
		if err != nil || u.Scheme != "https" || u.User != nil {
			return errors.New("input source must be an HTTPS identity")
		}
		switch in.Role {
		case "candidate":
			if u.Host != "github.com" || u.Path != fmt.Sprintf("/shruggietech/go-schedule/actions/runs/%d", m.RunID) {
				return errors.New("candidate source must identify its staging run")
			}
		case "baseline":
			if u.Host != "github.com" || u.Path != "/shruggietech/go-schedule/releases/download/v1.1.1/go-schedule_v1.1.1_windows_amd64.msi" {
				return errors.New("baseline must identify the public v1.1.1 MSI")
			}
		case "powershell":
			if u.Host != "github.com" || !regexp.MustCompile(`^/PowerShell/PowerShell/releases/download/v[^/]+/[^/]+win-x64\.zip$`).MatchString(u.Path) {
				return errors.New("PowerShell must identify an upstream x64 portable ZIP")
			}
		case "webview2":
			if u.Host != "developer.microsoft.com" && u.Host != "msedge.sf.dl.delivery.mp.microsoft.com" {
				return errors.New("WebView2 must identify a Microsoft source")
			}
		}
		// Reject either direction of containment, including output over inputs.
		for _, pair := range [][2]string{{out, in.Path}, {in.Path, out}} {
			rel, err := filepath.Rel(pair[0], pair[1])
			if err == nil && rel != ".." && !regexp.MustCompile(`^\.\.[/\\]`).MatchString(rel) {
				return errors.New("input/output overlap")
			}
		}
		n, digest, err := fileIdentity(in.Path)
		if err != nil {
			return err
		}
		if n != in.Bytes || digest != in.SHA256 {
			return fmt.Errorf("input identity mismatch: %s", in.Role)
		}
		info, err := os.Stat(in.Path)
		if err != nil {
			return err
		}
		for _, prior := range identities {
			if os.SameFile(prior, info) {
				return errors.New("role input file aliases refused")
			}
		}
		identities = append(identities, info)
	}
	for _, role := range []string{"bootstrap", "collector"} {
		path := filepath.Join(helpers, filenames[role])
		n, digest, err := fileIdentity(path)
		if err != nil {
			return err
		}
		m.Inputs = append(m.Inputs, input{Role: role, Path: path, Bytes: n, SHA256: digest, Source: "reviewed repository helper"})
	}
	// Create-only output deliberately preserves partial packages for diagnosis.
	if err := os.Mkdir(out, 0700); err != nil {
		return err
	}
	inputs := filepath.Join(out, "inputs")
	if err := os.Mkdir(inputs, 0700); err != nil {
		return err
	}
	for i, in := range m.Inputs {
		data, err := os.ReadFile(in.Path)
		if err != nil {
			return err
		}
		h := sha256.Sum256(data)
		if int64(len(data)) != in.Bytes || hex.EncodeToString(h[:]) != in.SHA256 {
			return fmt.Errorf("input changed while packaging: %s", in.Role)
		}
		m.Inputs[i].Path = filenames[in.Role]
		if err := os.WriteFile(filepath.Join(inputs, filenames[in.Role]), data, 0600); err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(inputs, "manifest.json"), append(data, '\n'), 0600); err != nil {
		return err
	}
	for _, scenario := range []string{"fresh", "upgrade"} {
		exports := filepath.Join(out, scenario+"-exports")
		if err := os.Mkdir(exports, 0700); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(out, scenario+".wsb"), []byte(configuration(inputs, exports, scenario)), 0600); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(out, "README.txt"), []byte("Preparation only, NOT release qualification. Launch each .wsb in a separate fresh session. Diagnostics export automatically. Complete test/windows/README.md native release matrix and issues #229-#233. Normal-user and high/mixed-DPI checks require suitable separate Windows 11 environments. No tags or releases were changed.\n"), 0600)
}

func escaped(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func configuration(inputs, exports, scenario string) string {
	return fmt.Sprintf("<Configuration>\n<Networking>Disable</Networking>\n<MappedFolders>\n<MappedFolder><HostFolder>%s</HostFolder><SandboxFolder>C:\\qualification-inputs</SandboxFolder><ReadOnly>true</ReadOnly></MappedFolder>\n<MappedFolder><HostFolder>%s</HostFolder><SandboxFolder>C:\\qualification-exports</SandboxFolder><ReadOnly>false</ReadOnly></MappedFolder>\n</MappedFolders>\n<LogonCommand><Command>powershell.exe -NoProfile -NonInteractive -WindowStyle Hidden -ExecutionPolicy Bypass -File C:\\qualification-inputs\\Start-QualificationGuest.ps1 -Scenario %s</Command></LogonCommand>\n</Configuration>\n", escaped(inputs), escaped(exports), scenario)
}
