package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestRunCancellationCommandAcknowledgesExactIdentity(t *testing.T) {
	for _, args := range [][]string{{"runs", "cancel", "run-1"}, {"--json", "runs", "cancel", "run-1"}} {
		calls := 0
		api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) {
			calls++
			if r.Method != "POST" || r.URL.Path != "/v1/runs/run-1/cancel" {
				t.Errorf("request %s %s", r.Method, r.URL.Path)
			}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(domain.RunCancellation{RunID: "run-1", TaskID: "task-1", Status: "requested"})
		})
		var out, diagnostics bytes.Buffer
		root := newRootWithDependencies(Dependencies{Client: api, Out: &out, Err: &diagnostics})
		root.SetArgs(args)
		if err := root.ExecuteContext(context.Background()); err != nil || calls != 1 || !strings.Contains(out.String(), "run-1") || !strings.Contains(out.String(), "requested") || diagnostics.Len() != 0 {
			t.Fatalf("args=%v out=%s stderr=%s calls=%d err=%v", args, out.String(), diagnostics.String(), calls, err)
		}
	}
}
func TestRunCancellationOutputFailureDoesNotReplay(t *testing.T) {
	calls := 0
	api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		json.NewEncoder(w).Encode(domain.RunCancellation{RunID: "run-1", Status: "requested"})
	})
	root := newRootWithDependencies(Dependencies{Client: api, Out: cancellationBrokenWriter{}, Err: &bytes.Buffer{}})
	root.SetArgs([]string{"runs", "cancel", "run-1"})
	if err := root.Execute(); err == nil || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}

type cancellationBrokenWriter struct{}

func (cancellationBrokenWriter) Write([]byte) (int, error) {
	return 0, errors.New("owned output failure")
}

func TestHumanActiveRunDetailReportsUnconfirmedCancellation(t *testing.T) {
	var out bytes.Buffer
	if err := writeRunDetail(&out, domain.Run{ID: "run-1", CancellationStatus: "failed", CancellationError: "owned cleanup did not finish"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Cancellation status: failed") || !strings.Contains(out.String(), "owned cleanup did not finish") {
		t.Fatalf("out=%s", out.String())
	}
}
