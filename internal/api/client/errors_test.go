package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestClassifyConnectionError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want ConnectionFailureKind
	}{
		{name: "access denied", err: fmt.Errorf("open pipe: %w", os.ErrPermission), want: ConnectionAccessDenied},
		{name: "daemon unavailable", err: fmt.Errorf("dial pipe: %w", os.ErrNotExist), want: ConnectionUnavailable},
		{name: "connection refused", err: syscall.ECONNREFUSED, want: ConnectionUnavailable},
		{name: "timeout", err: context.DeadlineExceeded, want: ConnectionTimeout},
		{name: "other transport", err: syscall.ECONNRESET, want: ConnectionOtherTransport},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewConnectionError("GET /v1/tasks", tt.err)
			if got.Kind != tt.want {
				t.Fatalf("kind = %q, want %q", got.Kind, tt.want)
			}
			if !errors.Is(got, tt.err) {
				t.Fatalf("connection error does not unwrap to %v", tt.err)
			}
		})
	}
}

func TestAccessDeniedCopyDoesNotClaimDaemonIsAbsent(t *testing.T) {
	t.Parallel()
	err := NewConnectionError("GET /v1/tasks", os.ErrPermission)
	if got := err.Error(); got == "" || containsFold(got, "is the daemon running") {
		t.Fatalf("access-denied copy = %q", got)
	}
}

func TestStatusErrorIsNotAConnectionError(t *testing.T) {
	t.Parallel()
	var connection *ConnectionError
	if errors.As(&StatusError{Code: "conflict", Message: "conflict"}, &connection) {
		t.Fatal("API status error classified as connection error")
	}
}

func containsFold(value, fragment string) bool {
	for i := 0; i+len(fragment) <= len(value); i++ {
		if equalFoldASCII(value[i:i+len(fragment)], fragment) {
			return true
		}
	}
	return false
}

func equalFoldASCII(left, right string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		l, r := left[i], right[i]
		if 'A' <= l && l <= 'Z' {
			l += 'a' - 'A'
		}
		if 'A' <= r && r <= 'Z' {
			r += 'a' - 'A'
		}
		if l != r {
			return false
		}
	}
	return true
}
