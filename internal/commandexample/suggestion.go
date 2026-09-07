// Package commandexample defines the one safe first-task example for each supported execution host.
package commandexample

import "strings"

// Suggestion is an exact allowlisted direct process invocation, not shell input.
type Suggestion struct {
	Platform    string   `json:"platform"`
	Display     string   `json:"display"`
	Program     string   `json:"program"`
	Args        []string `json:"args"`
	Explanation string   `json:"explanation"`
}

// ForPlatform returns a safe example only for a supported packaged host.
func ForPlatform(platform string) (Suggestion, bool) {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "windows":
		return Suggestion{Platform: "windows", Display: "cmd.exe /d /c ver", Program: "cmd.exe", Args: []string{"/d", "/c", "ver"}, Explanation: "Prints the Windows version, then exits. Its captured output appears in Activity."}, true
	case "darwin", "macos":
		return Suggestion{Platform: "macos", Display: "/usr/bin/sw_vers", Program: "/usr/bin/sw_vers", Explanation: "Prints the macOS product and version, then exits. Its captured output appears in Activity."}, true
	case "linux":
		return Suggestion{Platform: "linux", Display: "uname -a", Program: "uname", Args: []string{"-a"}, Explanation: "Prints Linux system information, then exits. Its captured output appears in Activity."}, true
	default:
		return Suggestion{}, false
	}
}

// Safe verifies the complete allowlist so callers cannot construct variants.
func (s Suggestion) Safe() bool {
	want, ok := ForPlatform(s.Platform)
	return ok && s.Display == want.Display && s.Program == want.Program && equalStrings(s.Args, want.Args) && s.Explanation == want.Explanation
}

// Recognizes reports whether captured output matches the platform example.
func (s Suggestion) Recognizes(output string) bool {
	if !s.Safe() {
		return false
	}
	switch s.Platform {
	case "windows":
		return strings.Contains(strings.ToLower(output), "windows")
	case "macos":
		return strings.Contains(output, "ProductName") && strings.Contains(output, "ProductVersion")
	case "linux":
		return strings.TrimSpace(output) != ""
	default:
		return false
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
