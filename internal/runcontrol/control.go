// Package runcontrol defines the shared execution-cancellation cause and feedback contract.
package runcontrol

import (
	"context"
	"errors"
	"time"
)

// ErrRequested identifies an explicit operator request, separately from daemon shutdown.
var ErrRequested = errors.New("run cancellation requested")

// ErrNotActive identifies a known history marker that is not an executing process.
var ErrNotActive = errors.New("run is not an active execution")

// CleanupTimeout bounds one cancellation attempt, not the lifetime of unreaped work.
const CleanupTimeout = 5 * time.Second

type reporterKey struct{}

// WithFailureReporter attaches non-blocking ownership feedback to one execution.
func WithFailureReporter(ctx context.Context, report func(error)) context.Context {
	return context.WithValue(ctx, reporterKey{}, report)
}

// ReportFailure reports unconfirmed cancellation without fabricating a terminal result.
func ReportFailure(ctx context.Context, err error) {
	if err == nil {
		return
	}
	if report, ok := ctx.Value(reporterKey{}).(func(error)); ok && report != nil {
		report(err)
	}
}
