//go:build !windows

package executor

import (
	"os/exec"
	"os/user"
	"strings"
	"testing"
)

func TestApplyRunAsSetsIdentityAndPreservesExplicitHome(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("true")
	cmd.Env = []string{"HOME=/explicit", "LOGNAME=wrong", "USER=wrong"}
	if err := applyRunAs(cmd, current.Username); err != nil {
		t.Fatal(err)
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
