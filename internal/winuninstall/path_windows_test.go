//go:build windows

package winuninstall

import "testing"

func TestValidateCandidateRejectsUnsafeWindowsPaths(t *testing.T) {
	tests := []string{"", `relative\goschedule`, `C:\`, `\\server\share\goschedule`, `\\?\C:\ProgramData\goschedule`, `C:\ProgramData\goschedule:stream`, `C:\ProgramData\goschedule.`, `C:\ProgramData\CON\goschedule`}
	for _, path := range tests {
		if err := validateLexicalPath(path); err == nil {
			t.Errorf("validateLexicalPath(%q) succeeded", path)
		}
	}
}

func TestValidateCandidateAcceptsOwnedLeaves(t *testing.T) {
	for _, path := range []string{`C:\ProgramData\goschedule`, `D:\Profiles\Ada\AppData\Roaming\fyne\tech.shruggie.goschedule`} {
		if err := validateLexicalPath(path); err != nil {
			t.Errorf("validateLexicalPath(%q): %v", path, err)
		}
	}
}
