// Command gosched is the scheduler CLI: a thin client that manages tasks,
// groups, runs, alerts, and the system service through the daemon's local API.
package main

import (
	"os"

	"github.com/shruggietech/go-schedule/internal/cli"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
)

func main() {
	if handled, code := clientsecret.RunReadHelper(os.Args[1:]); handled {
		os.Exit(code)
	}
	os.Exit(cli.Execute())
}
