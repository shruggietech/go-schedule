package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestRequestObservationExcludesSecretsAndRecordsDecodeFailure(t *testing.T) {
	c := NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("invalid secret-response")) }))
	var observations []RequestObservation
	ctx := WithRequestObserver(context.Background(), func(o RequestObservation) { observations = append(observations, o) })
	var out any
	if err := c.do(ctx, http.MethodGet, "/v1/tasks/private-id?token=private-query", nil, &out); err == nil {
		t.Fatal("malformed response accepted")
	}
	if len(observations) != 1 || observations[0].Route != "/v1/tasks/{id}" || observations[0].Outcome != "failed" {
		t.Fatalf("bad observation: %+v", observations)
	}
	data, err := json.Marshal(observations)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "private") || strings.Contains(string(data), "secret") {
		t.Fatal("request content leaked")
	}
	if _, err := c.Health(context.Background()); err == nil {
		t.Fatal("malformed health accepted")
	}
	if len(observations) != 1 {
		t.Fatal("observer escaped its context")
	}
}

func TestRequestObservationPreservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var observed []RequestObservation
	ctx = WithRequestObserver(ctx, func(o RequestObservation) { observed = append(observed, o) })
	c := New("unused-owned-endpoint")
	if _, err := c.Health(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("get cancellation: %v", err)
	}
	if _, err := c.ListCredentials(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("do cancellation: %v", err)
	}
	if len(observed) != 2 || observed[0].Outcome != "canceled" || observed[1].Outcome != "canceled" || observed[1].Route != "/v1/access/credentials" {
		t.Fatalf("cancellation evidence: %+v", observed)
	}
}
