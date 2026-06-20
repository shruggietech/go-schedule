// Command gosched is the scheduler CLI: a thin client that manages tasks,
// groups, runs, alerts, and the system service through the daemon's local API.
package main

import (
	"os"

	"github.com/shruggietech/go-scheduler/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
