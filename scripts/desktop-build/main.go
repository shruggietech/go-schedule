// Command desktop-build binds Wails resources and Go identity to one selected version.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/shruggietech/go-schedule/v2/internal/buildinfo"
	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

var releasePattern = regexp.MustCompile("^v?(0|[1-9][0-9]*)[.](0|[1-9][0-9]*)[.](0|[1-9][0-9]*)$")

func numericVersion(v string) (string, error) {
	if v == "0.0.1-dev" {
		return "0.0.1", nil
	}
	m := releasePattern.FindStringSubmatch(v)
	if m == nil {
		return "", fmt.Errorf("version %q: use MAJOR.MINOR.PATCH (optional v prefix), or 0.0.1-dev", v)
	}
	for i, max := range []int{255, 255, 65535} {
		n, e := strconv.Atoi(m[i+1])
		if e != nil || n > max {
			return "", fmt.Errorf("version component %q exceeds Windows Installer limit %d", m[i+1], max)
		}
	}
	return strings.TrimPrefix(v, "v"), nil
}
func build(root, v, extra string, args []string, run func([]string) error) (retErr error) {
	numeric, err := numericVersion(v)
	if err != nil {
		return err
	}
	if strings.Contains(extra, "github.com/shruggietech/go-schedule/v2/internal/buildinfo.Version") {
		return errors.New("supply version with -version, not a linker Version override")
	}
	for _, arg := range args {
		if strings.HasPrefix(strings.TrimLeft(arg, "-"), "ldflags") {
			return errors.New("supply linker flags with desktop-build -ldflags before --")
		}
	}
	path := filepath.Join(root, "desktop", "wails.json")
	lockPath := filepath.Join(root, "desktop", ".version-build.lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("acquire desktop build lock (remove stale lock only after confirming no build is running): %w", err)
	}
	defer func() { retErr = errors.Join(retErr, lock.Close(), os.Remove(lockPath)) }()
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	var cfg map[string]any
	if err = json.Unmarshal(original, &cfg); err != nil {
		return fmt.Errorf("read Wails configuration: %w", err)
	}
	if cfg == nil {
		return errors.New("wails configuration must be an object")
	}
	info, ok := cfg["info"].(map[string]any)
	if !ok {
		info = map[string]any{}
		cfg["info"] = info
	}
	info["productVersion"] = numeric
	if v == "0.0.1-dev" {
		info["comments"] = "Development build (0.0.1-dev)"
	} else {
		info["comments"] = "Release " + v
	}
	updated, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, os.WriteFile(path, original, stat.Mode().Perm())) }()
	if err = os.WriteFile(path, append(updated, '\n'), stat.Mode().Perm()); err != nil {
		return err
	}
	flags := strings.TrimSpace(extra + " -X github.com/shruggietech/go-schedule/v2/internal/buildinfo.Version=" + v)
	command := []string{"run", "github.com/wailsapp/wails/v2/cmd/wails@v2.17.1", "build", "-ldflags", flags}
	return run(append(command, args...))
}
func main() {
	v := flag.String("version", buildinfo.Version, "selected release version, default development")
	extra := flag.String("ldflags", "", "additional linker flags")
	flag.Parse()
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = build(root, *v, *extra, flag.Args(), func(args []string) error {
		exe := os.Getenv("GO")
		if exe == "" {
			exe = "go"
		}
		cmd := exec.Command(exe, args...)
		cmd.Dir = filepath.Join(root, "desktop")
		platform.HideConsole(cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), "CI=true", "npm_config_yes=true", "GIT_TERMINAL_PROMPT=0", "GCM_INTERACTIVE=never")
		if e := cmd.Run(); e != nil {
			return fmt.Errorf("wails build: %w", e)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
