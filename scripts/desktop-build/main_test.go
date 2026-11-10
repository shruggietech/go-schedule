package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	for _, v := range []string{"v2.0.3", "3.7.19", "0.0.1-dev"} {
		if _, err := numericVersion(v); err != nil {
			t.Fatal(err)
		}
	}
	for _, v := range []string{"", "v2", "v02.0.3", "256.0.0", "1.256.0", "1.0.65536", "1.2.3-beta", "1.2.3;echo"} {
		if _, err := numericVersion(v); err == nil {
			t.Fatalf("accepted %q", v)
		}
	}
}
func TestBuildRestoresConfig(t *testing.T) {
	for _, v := range []string{"v2.0.3", "3.7.19", "0.0.1-dev"} {
		t.Run(v, func(t *testing.T) {
			dir := t.TempDir()
			desktop := filepath.Join(dir, "desktop")
			if err := os.Mkdir(desktop, 0700); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(desktop, "wails.json")
			original := []byte("{\"name\":\"go-schedule\",\"other\":true}\n")
			if err := os.WriteFile(path, original, 0600); err != nil {
				t.Fatal(err)
			}
			failure := errors.New("child failed")
			for _, childErr := range []error{nil, failure} {
				err := build(dir, v, "-s -w", []string{"-clean"}, func(args []string) error {
					b, e := os.ReadFile(path)
					if e != nil {
						t.Fatal(e)
					}
					var cfg map[string]any
					if e = json.Unmarshal(b, &cfg); e != nil {
						t.Fatal(e)
					}
					info := cfg["info"].(map[string]any)
					want, _ := numericVersion(v)
					if info["productVersion"] != want || cfg["other"] != true {
						t.Fatalf("config %s", b)
					}
					if !strings.Contains(strings.Join(args, " "), "Version="+v) {
						t.Fatal(args)
					}
					if v == "0.0.1-dev" && !strings.Contains(info["comments"].(string), "Development") {
						t.Fatal(info)
					}
					if e := build(dir, v, "", nil, func([]string) error { return nil }); e == nil {
						t.Fatal("parallel build accepted")
					}
					return childErr
				})
				if !errors.Is(err, childErr) {
					t.Fatalf("got %v want %v", err, childErr)
				}
				b, e := os.ReadFile(path)
				if e != nil || string(b) != string(original) {
					t.Fatal("configuration not restored", e)
				}
			}
			for _, arg := range []string{"-ldflags", "--ldflags", "--ldflags=wrong"} {
				if e := build(dir, v, "", []string{arg}, func([]string) error { return nil }); e == nil {
					t.Errorf("accepted forwarded %s", arg)
				}
			}
			if e := build(dir, v, "-X github.com/shruggietech/go-schedule/v2/internal/buildinfo.Version=wrong", nil, func([]string) error { t.Fatal("ran conflicting flags"); return nil }); e == nil {
				t.Fatal("accepted override")
			}
		})
	}
}

func TestInvalidConfigurationDoesNotRun(t *testing.T) {
	for _, content := range []string{"null", "{broken"} {
		dir := t.TempDir()
		if e := os.Mkdir(filepath.Join(dir, "desktop"), 0700); e != nil {
			t.Fatal(e)
		}
		path := filepath.Join(dir, "desktop", "wails.json")
		if e := os.WriteFile(path, []byte(content), 0600); e != nil {
			t.Fatal(e)
		}
		if e := build(dir, "v2.0.3", "", nil, func([]string) error { t.Fatal("ran malformed config"); return nil }); e == nil {
			t.Fatal("accepted malformed config")
		}
		b, e := os.ReadFile(path)
		if e != nil || string(b) != content {
			t.Fatal("changed original config")
		}
	}
}
