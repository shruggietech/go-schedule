//go:build !windows

package executor

import (
	"os/exec"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

func TestCredentialForUserAcceptsUint32Boundaries(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		gid     string
		wantUID uint32
		wantGID uint32
	}{
		{name: "zero", uid: "0", gid: "0", wantUID: 0, wantGID: 0},
		{name: "maximum", uid: "4294967295", gid: "4294967295", wantUID: ^uint32(0), wantGID: ^uint32(0)},
		{name: "leading zeroes", uid: "00042", gid: "00084", wantUID: 42, wantGID: 84},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := credentialForUser("synthetic", &user.User{Uid: tt.uid, Gid: tt.gid})
			if err != nil {
				t.Fatalf("credentialForUser() error = %v", err)
			}
			if got.Uid != tt.wantUID || got.Gid != tt.wantGID {
				t.Fatalf("credentialForUser() = {%d, %d}, want {%d, %d}", got.Uid, got.Gid, tt.wantUID, tt.wantGID)
			}
		})
	}
}

func TestCredentialForUserRejectsInvalidIdentifiers(t *testing.T) {
	invalid := []string{"", "-1", " 1", "1 ", "1.0", "0x10", "not-a-number", "4294967296"}
	for _, field := range []string{"uid", "gid"} {
		for _, value := range invalid {
			t.Run(field+"_"+value, func(t *testing.T) {
				u := &user.User{Uid: "42", Gid: "84"}
				if field == "uid" {
					u.Uid = value
				} else {
					u.Gid = value
				}
				_, err := credentialForUser("synthetic", u)
				if err == nil {
					t.Fatal("credentialForUser() error = nil, want rejection")
				}
				if !strings.Contains(err.Error(), "invalid "+field) || !strings.Contains(err.Error(), "synthetic") {
					t.Fatalf("credentialForUser() error = %q, want field and account context", err)
				}
			})
		}
	}
}

func TestApplyResolvedRunAsPreservesCommandOnCredentialFailure(t *testing.T) {
	tests := []struct {
		name string
		uid  string
		gid  string
	}{
		{name: "invalid uid", uid: "-1", gid: "84"},
		{name: "valid uid invalid gid", uid: "42", gid: "4294967296"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCredential := &syscall.Credential{Uid: 7, Gid: 8}
			originalAttr := &syscall.SysProcAttr{Credential: originalCredential}
			originalEnv := []string{"HOME=/original", "USER=original", "LOGNAME=original"}
			cmd := exec.Command("true")
			cmd.SysProcAttr = originalAttr
			cmd.Env = append([]string(nil), originalEnv...)

			err := applyResolvedRunAs(cmd, "synthetic", &user.User{
				Username: "resolved",
				HomeDir:  "/resolved",
				Uid:      tt.uid,
				Gid:      tt.gid,
			}, false)
			if err == nil {
				t.Fatal("applyResolvedRunAs() error = nil, want rejection")
			}
			if cmd.SysProcAttr != originalAttr || cmd.SysProcAttr.Credential != originalCredential {
				t.Fatal("applyResolvedRunAs() replaced process state after validation failure")
			}
			if !reflect.DeepEqual(cmd.Env, originalEnv) {
				t.Fatalf("applyResolvedRunAs() Env = %#v, want %#v", cmd.Env, originalEnv)
			}
		})
	}
}

func TestApplyRunAsSetsIdentityAndPreservesExplicitHome(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("true")
	cmd.Env = []string{"HOME=/explicit", "LOGNAME=wrong", "USER=wrong"}
	if err := applyRunAs(cmd, current.Username, true); err != nil {
		t.Fatal(err)
	}
	wantUID, err := strconv.ParseUint(current.Uid, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	wantGID, err := strconv.ParseUint(current.Gid, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	if cmd.SysProcAttr == nil || cmd.SysProcAttr.Credential == nil {
		t.Fatal("applyRunAs() did not assign process credentials")
	}
	if got := cmd.SysProcAttr.Credential; got.Uid != uint32(wantUID) || got.Gid != uint32(wantGID) {
		t.Fatalf("applyRunAs() credential = {%d, %d}, want {%d, %d}", got.Uid, got.Gid, wantUID, wantGID)
	}
	want := map[string]string{"HOME": "/explicit", "LOGNAME": current.Username, "USER": current.Username}
	for key, value := range want {
		found := false
		for _, item := range cmd.Env {
			if item == key+"="+value {
				found = true
			}
			if strings.HasPrefix(item, key+"=") && item != key+"="+value {
				t.Fatalf("%s has stale value %q", key, item)
			}
		}
		if !found {
			t.Fatalf("%s=%q missing from %#v", key, value, cmd.Env)
		}
	}
}

func TestApplyRunAsSupportsNumericAccountLookup(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("true")
	if err := applyRunAs(cmd, current.Uid, true); err != nil {
		t.Fatal(err)
	}
	wantUID, err := strconv.ParseUint(current.Uid, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	if cmd.SysProcAttr == nil || cmd.SysProcAttr.Credential == nil || cmd.SysProcAttr.Credential.Uid != uint32(wantUID) {
		t.Fatalf("applyRunAs() did not retain numeric UID %d", wantUID)
	}
}

func TestApplyRunAsEmptyPreservesCommand(t *testing.T) {
	originalCredential := &syscall.Credential{Uid: 7, Gid: 8}
	originalAttr := &syscall.SysProcAttr{Credential: originalCredential}
	originalEnv := []string{"HOME=/original", "USER=original", "LOGNAME=original"}
	cmd := exec.Command("true")
	cmd.SysProcAttr = originalAttr
	cmd.Env = append([]string(nil), originalEnv...)

	if err := applyRunAs(cmd, "", false); err != nil {
		t.Fatal(err)
	}
	if cmd.SysProcAttr != originalAttr || cmd.SysProcAttr.Credential != originalCredential {
		t.Fatal("applyRunAs() changed process state for empty run_as")
	}
	if !reflect.DeepEqual(cmd.Env, originalEnv) {
		t.Fatalf("applyRunAs() Env = %#v, want %#v", cmd.Env, originalEnv)
	}
}

func TestApplyRunAsRemovesDuplicateIdentityAndDefaultsHome(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("true")
	cmd.Env = []string{
		"USER=inherited", "LOGNAME=inherited", "HOME=/daemon",
		"USER=task", "LOGNAME=task",
	}
	if err := applyRunAs(cmd, current.Username, false); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"HOME": current.HomeDir, "LOGNAME": current.Username, "USER": current.Username}
	for key, value := range want {
		matches := 0
		for _, item := range cmd.Env {
			if strings.HasPrefix(item, key+"=") {
				matches++
				if item != key+"="+value {
					t.Fatalf("%s has stale value %q", key, item)
				}
			}
		}
		if matches != 1 {
			t.Fatalf("%s has %d entries in %#v", key, matches, cmd.Env)
		}
	}
}
