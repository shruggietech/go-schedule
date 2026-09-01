package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

// ConnectionFailureKind identifies why the local daemon transport failed.
type ConnectionFailureKind string

const (
	// ConnectionUnavailable means the local endpoint is absent or refused.
	ConnectionUnavailable ConnectionFailureKind = "unavailable"
	// ConnectionAccessDenied means the OS rejected the caller's authorization.
	ConnectionAccessDenied ConnectionFailureKind = "access_denied"
	// ConnectionTimeout means the operation exceeded its deadline.
	ConnectionTimeout ConnectionFailureKind = "timeout"
	// ConnectionOtherTransport means the transport failed for another reason.
	ConnectionOtherTransport ConnectionFailureKind = "other_transport"
)

// ConnectionError is a classified local-transport failure. Cause is retained
// so callers can use errors.Is and errors.As for platform-specific diagnosis.
type ConnectionError struct {
	Kind      ConnectionFailureKind
	Operation string
	Cause     error
}

// NewConnectionError classifies cause and wraps it with actionable context.
func NewConnectionError(operation string, cause error) *ConnectionError {
	return &ConnectionError{Kind: classifyConnectionError(cause), Operation: operation, Cause: cause}
}

func classifyConnectionError(err error) ConnectionFailureKind {
	if errors.Is(err, os.ErrPermission) {
		return ConnectionAccessDenied
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ConnectionTimeout
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ConnectionTimeout
	}
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ECONNREFUSED) {
		return ConnectionUnavailable
	}
	return ConnectionOtherTransport
}

func (e *ConnectionError) Error() string {
	if e == nil {
		return "connection failed"
	}
	switch e.Kind {
	case ConnectionAccessDenied:
		return fmt.Sprintf("api: %s: daemon IPC access was denied for the current session: %v", e.Operation, e.Cause)
	case ConnectionUnavailable:
		return fmt.Sprintf("api: %s: daemon unavailable: %v (check that the goschedd service is running)", e.Operation, e.Cause)
	case ConnectionTimeout:
		return fmt.Sprintf("api: %s: daemon connection timed out: %v", e.Operation, e.Cause)
	default:
		return fmt.Sprintf("api: %s: daemon connection failed: %v", e.Operation, e.Cause)
	}
}

// Unwrap returns the original transport failure.
func (e *ConnectionError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
