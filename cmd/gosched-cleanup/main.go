package main

import (
	"os"

	"github.com/shruggietech/go-schedule/internal/winuninstall"
)

const (
	exitSuccess       = 0
	exitUsage         = 64
	exitIncomplete    = 2
	exitInternalError = 3
)

func main() {
	os.Exit(run(os.Args[1:], winuninstall.Wipe))
}

func run(args []string, wipe func() winuninstall.Result) int {
	if len(args) != 1 || args[0] != "wipe" {
		return exitUsage
	}
	switch wipe().State {
	case winuninstall.StateComplete:
		return exitSuccess
	case winuninstall.StateRefused, winuninstall.StatePartial:
		return exitIncomplete
	default:
		return exitInternalError
	}
}
