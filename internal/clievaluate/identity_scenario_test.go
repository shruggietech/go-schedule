package clievaluate

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/cli"
	"github.com/shruggietech/go-schedule/v2/internal/clientsecret"
	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

func TestWrongIdentityScenarioReachesRemoteComparison(t *testing.T) {
	const token = "owned-fixture-credential"
	var requests atomic.Int32
	remote := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Header.Get("Authorization") != "" {
			t.Error("credential was sent before remote identity verification")
		}
		_ = json.NewEncoder(w).Encode(server.ManifestResponse{InstallationID: "actual-daemon"})
	}))
	defer remote.Close()
	certificate := filepath.Join(t.TempDir(), "trust.pem")
	if err := os.WriteFile(certificate, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: remote.Certificate().Raw}), 0o600); err != nil {
		t.Fatal(err)
	}
	s := &session{recorder: &clifixture.Recorder{}, vars: map[string]string{"remote": remote.URL, "remote_credential": "owned-credential", "remote_token": token, "certificate": certificate}}
	secrets := clientsecret.NewWithBackend(s.recorder)
	if err := secrets.Save("actual-daemon", "owned-credential", token); err != nil {
		t.Fatal(err)
	}
	var definition step
	for _, candidate := range additionalSteps() {
		if candidate.Spec.ID == "selector.wrong-identity" {
			definition = candidate
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.before(ctx, definition); err != nil {
		t.Fatal(err)
	}
	if retrieved, err := secrets.Load(wrongDaemonID, "owned-credential"); err != nil || retrieved != token {
		t.Fatal("scenario did not make its owned credential retrievable before identity comparison")
	}
	args, err := s.arguments(definition.Spec.Argv)
	if err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	code := cli.ExecuteWithDependencies(ctx, args, cli.Dependencies{Secrets: secrets, Out: &out, Err: &diagnostic})
	if code != 1 || requests.Load() == 0 || !strings.Contains(diagnostic.String(), "remote daemon identity does not match the selected target") || strings.Contains(diagnostic.String(), token) {
		t.Fatalf("exit=%d requests=%d out=%q diagnostic=%q", code, requests.Load(), out.String(), diagnostic.String())
	}
	if definition.OutputContains == "" || !strings.Contains(diagnostic.String(), definition.OutputContains) || strings.Contains("load native credential: not found", definition.OutputContains) {
		t.Fatal("scenario accepts credential lookup failure instead of requiring the remote identity refusal")
	}
}
