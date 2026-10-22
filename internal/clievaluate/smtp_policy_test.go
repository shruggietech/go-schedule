package clievaluate

import (
	"context"
	"testing"
	"time"
)

func TestOwnedSMTPTransportEvaluation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := smtpTransportPolicy(ctx, NewRedactor())
	if result.Code != 0 || result.Detail != "" {
		t.Fatalf("policy: %+v", result)
	}
	for _, check := range result.Checks {
		if !check {
			t.Fatalf("policy checks: %+v", result)
		}
	}
}
