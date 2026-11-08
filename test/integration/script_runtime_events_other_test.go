//go:build !windows

package integration

import (
	"context"
	"errors"

	"github.com/shruggietech/go-schedule/v2/internal/incident"
)

func scriptRuntimeEvents(context.Context, string) incident.RuntimeEventQueryResult {
	return incident.RuntimeEventQueryResult{Error: errors.New("windows event readback is unavailable on this platform"), Unavailable: true}
}
