//go:build !windows

package integration

import "os/exec"

func configureScriptProcess(_ *exec.Cmd) {}
