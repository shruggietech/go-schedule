//go:build !windows

package commandexample

import "os/exec"

func configureTestCommand(*exec.Cmd) {}
