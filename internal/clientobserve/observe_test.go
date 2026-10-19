package clientobserve

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type fakeReader struct {
	id           string
	accessErr    error
	active, peak *atomic.Int32
	block        bool
	reads        atomic.Int32
	searchResult domain.DaemonSearch
}

func (f *fakeReader) Health(context.Context) (server.HealthResponse, error) {
	return server.HealthResponse{Status: "ok", Version: "1.5.0"}, nil
}
func (f *fakeReader) VerifyIdentity(context.Context) (server.ManifestResponse, error) {
	return server.ManifestResponse{InstallationID: f.id, ProductVersion: "1.5.0"}, nil
}
func (f *fakeReader) VerifyAccess(context.Context) (domain.Capability, error) {
	return domain.CapabilityObserve, f.accessErr
}
func (f *fakeReader) SystemSummary(ctx context.Context) (domain.SystemSummary, error) {
	f.reads.Add(1)
	if f.active != nil {
		n := f.active.Add(1)
		defer f.active.Add(-1)
		for {
			old := f.peak.Load()
			if n <= old || f.peak.CompareAndSwap(old, n) {
				break
			}
		}
	}
	if f.block {
		<-ctx.Done()
		return domain.SystemSummary{}, ctx.Err()
	}
	time.Sleep(time.Millisecond)
	return domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: time.Now()}, nil
}
func (f *fakeReader) Search(context.Context, string, []domain.SearchKind, int) (domain.DaemonSearch, error) {
	f.reads.Add(1)
	return f.searchResult, nil
}

func TestSummaryPreservesIdentityAndIndependentAccessFailures(t *testing.T) {
	good := &fakeReader{id: "same-object-daemon-a"}
	bad := &fakeReader{id: "daemon-b", accessErr: &client.StatusError{Code: "credential_revoked", Message: "SECRET must never appear"}}
	result, err := New().Summary(context.Background(), []Target{{Registration: Registration{Key: "a", Label: "Alpha", ExpectedDaemonID: good.id}, Client: good}, {Registration: Registration{Key: "b", Label: "Beta", ExpectedDaemonID: bad.id}, Client: bad}})
	if !errors.Is(err, ErrPartial) || len(result.Observations) != 2 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if result.Observations[0].Registration.ObservedDaemonID != good.id || result.Observations[0].Summary == nil {
		t.Fatalf("lost source %+v", result.Observations[0])
	}
	if result.Observations[1].Failure.Code != "credential_revoked" || bad.reads.Load() != 0 {
		t.Fatalf("access failure=%+v reads=%d", result.Observations[1], bad.reads.Load())
	}
	encoded, _ := json.Marshal(result)
	if strings.Contains(string(encoded), "SECRET") {
		t.Fatal("raw server error leaked")
	}
}

func TestPerTargetDeadlinePreservesSuccessfulPeer(t *testing.T) {
	observer := New()
	observer.Timeout = 20 * time.Millisecond
	result, err := observer.Summary(context.Background(), []Target{{Registration: Registration{Key: "slow"}, Client: &fakeReader{id: "slow", block: true}}, {Registration: Registration{Key: "good"}, Client: &fakeReader{id: "good"}}})
	if !errors.Is(err, ErrPartial) || result.Observations[0].Summary == nil || result.Observations[1].Failure.Code != "timed_out" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestSearchChecksEnvelopeAndKeepsSameObjectOnDistinctSources(t *testing.T) {
	valid := domain.DaemonSearch{Schema: domain.DaemonSearchSchema, ObservedAt: time.Now(), Query: "backup", Results: []domain.SearchMatch{{Kind: domain.SearchKindTask, ObjectID: "same-task", Name: "Backup"}}}
	result, err := New().Search(context.Background(), []Target{{Registration: Registration{Key: "a"}, Client: &fakeReader{id: "daemon-a", searchResult: valid}}, {Registration: Registration{Key: "b"}, Client: &fakeReader{id: "daemon-b", searchResult: valid}}}, SearchRequest{Query: " backup ", Limit: 10})
	if err != nil || len(result.Observations) != 2 || result.Observations[0].Registration.ObservedDaemonID == result.Observations[1].Registration.ObservedDaemonID {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	valid.Query = "wrong"
	result, err = New().Search(context.Background(), []Target{{Client: &fakeReader{id: "daemon", searchResult: valid}}}, SearchRequest{Query: "backup", Limit: 10})
	if !errors.Is(err, ErrPartial) || result.Observations[0].Failure.Code != "incompatible" {
		t.Fatalf("malformed result=%+v err=%v", result, err)
	}
}

type fakeProfiles struct {
	value clientprofile.Collection
	err   error
}

func (p fakeProfiles) Load() (clientprofile.Collection, error) { return p.value, p.err }

type fakeSecrets struct {
	err   error
	loads *atomic.Int32
}

func (s fakeSecrets) LoadContext(context.Context, string, string) (string, error) {
	if s.loads != nil {
		s.loads.Add(1)
	}
	return "SECRET", s.err
}

func validObservedProfile(t *testing.T) clientprofile.Profile {
	t.Helper()
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	certificate := &x509.Certificate{SerialNumber: big.NewInt(1), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, public, private)
	if err != nil {
		t.Fatal(err)
	}
	return clientprofile.Profile{ID: "remote", Label: "Remote", Endpoint: "https://remote.example", DaemonID: "remote-daemon", CredentialID: "credential-id", CertificatePEM: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})), ClientKind: "cli", Capability: "observe", CreatedAt: time.Now(), UpdatedAt: time.Now()}
}

func TestTargetSetupFailuresRemainVisibleWithoutCredentialMaterial(t *testing.T) {
	local := Target{Registration: Registration{Key: "local"}, Client: &fakeReader{id: "local"}}
	profiles := fakeProfiles{value: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{validObservedProfile(t)}}}
	targets := Targets(context.Background(), local, profiles, fakeSecrets{err: errors.New("SECRET")}, nil)
	if len(targets) != 2 || targets[1].Resolve == nil {
		t.Fatalf("targets=%+v", targets)
	}
	result, err := New().Summary(context.Background(), targets)
	encoded, _ := json.Marshal(result)
	if !errors.Is(err, ErrPartial) || strings.Contains(string(encoded), "SECRET") || strings.Contains(string(encoded), "CERTIFICATE") || strings.Contains(string(encoded), "credential-id") {
		t.Fatalf("result=%s err=%v", encoded, err)
	}
	targets = Targets(context.Background(), local, fakeProfiles{err: errors.New("broken store")}, nil, nil)
	if len(targets) != 2 || targets[1].SetupFailure.Code != "profile_store_unavailable" {
		t.Fatalf("profile load disappeared: %+v", targets)
	}
}

func TestCanceledDiscoveryAndPendingTargetsNeverLoadCredentials(t *testing.T) {
	var loads atomic.Int32
	profiles := fakeProfiles{value: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{validObservedProfile(t)}}}
	local := Target{Registration: Registration{Key: "local"}, Client: &fakeReader{id: "local"}}
	targets := Targets(context.Background(), local, profiles, fakeSecrets{loads: &loads}, func(string, string, string, string) (Reader, error) { return &fakeReader{id: "remote-daemon"}, nil })
	if loads.Load() != 0 {
		t.Fatal("discovery read credentials before a bounded worker")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := New().Summary(ctx, targets)
	if !errors.Is(err, context.Canceled) || loads.Load() != 0 || len(result.Observations) != 2 {
		t.Fatalf("result=%+v err=%v loads=%d", result, err, loads.Load())
	}
	if len(Targets(ctx, local, profiles, fakeSecrets{loads: &loads}, nil)) != 1 {
		t.Fatal("canceled discovery continued reading profiles")
	}
}

func TestMalformedDiscoveryDoesNotEchoUnsafeMetadata(t *testing.T) {
	profile := validObservedProfile(t)
	profile.Label = "unsafe\nlabel"
	profile.Endpoint = "https://SECRET@remote.example"
	profiles := fakeProfiles{value: clientprofile.Collection{Version: clientprofile.CurrentVersion, Profiles: []clientprofile.Profile{profile}}}
	result, err := New().Summary(context.Background(), Targets(context.Background(), Target{Registration: Registration{Key: "local"}, Client: &fakeReader{id: "local"}}, profiles, fakeSecrets{}, nil))
	encoded, _ := json.Marshal(result)
	if !errors.Is(err, ErrPartial) || strings.Contains(string(encoded), "SECRET") || strings.Contains(string(encoded), "unsafe") || !strings.Contains(string(encoded), "invalid_profile") {
		t.Fatalf("result=%s err=%v", encoded, err)
	}
}

func TestIdentityMismatchFailsBeforeFeatureRead(t *testing.T) {
	f := &fakeReader{id: "replacement"}
	result, err := New().Summary(context.Background(), []Target{{Registration: Registration{Key: "saved", ExpectedDaemonID: "original"}, Client: f}})
	if !errors.Is(err, ErrPartial) || result.Observations[0].Failure.Code != "identity_changed" || f.reads.Load() != 0 {
		t.Fatalf("result=%+v err=%v reads=%d", result, err, f.reads.Load())
	}
}

func TestSummaryBoundsWorkersAndCancelsEveryTarget(t *testing.T) {
	var active, peak atomic.Int32
	targets := make([]Target, 12)
	for i := range targets {
		targets[i] = Target{Registration: Registration{Key: string(rune('a' + i))}, Client: &fakeReader{id: "daemon", active: &active, peak: &peak, block: true}}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	result, err := New().Summary(ctx, targets)
	if !errors.Is(err, context.DeadlineExceeded) || len(result.Observations) != len(targets) || peak.Load() > 4 || active.Load() != 0 {
		t.Fatalf("count=%d peak=%d active=%d err=%v", len(result.Observations), peak.Load(), active.Load(), err)
	}
	for _, o := range result.Observations {
		if o.Failure == nil {
			t.Fatalf("missing cancellation outcome %+v", o)
		}
	}
}

func TestSearchRejectsInvalidRequestBeforeAnyRead(t *testing.T) {
	f := &fakeReader{id: "daemon"}
	_, err := New().Search(context.Background(), []Target{{Client: f}}, SearchRequest{Query: "x", Limit: 51})
	if !errors.Is(err, ErrInvalid) || f.reads.Load() != 0 {
		t.Fatalf("err=%v", err)
	}
}

// blockedSecrets models the pre-context boundary and its cancellation-aware replacement.
type blockedSecrets struct {
	release chan struct{}
	active  atomic.Int32
}

func (s *blockedSecrets) Load(string, string) (string, error) {
	s.active.Add(1)
	defer s.active.Add(-1)
	<-s.release
	return "", context.Canceled
}
func (s *blockedSecrets) LoadContext(ctx context.Context, _, _ string) (string, error) {
	s.active.Add(1)
	defer s.active.Add(-1)
	<-ctx.Done()
	return "", ctx.Err()
}
func TestCredentialLookupDeadlinesReleaseAllWorkersAndKeepLocalResult(t *testing.T) {
	for _, operation := range []string{"summary", "search", "cancel"} {
		t.Run(operation, func(t *testing.T) {
			secrets := &blockedSecrets{release: make(chan struct{})}
			profile := validObservedProfile(t)
			collection := clientprofile.Collection{Version: clientprofile.CurrentVersion}
			for i := 0; i < 8; i++ {
				p := profile
				p.ID = string(rune('a' + i))
				collection.Profiles = append(collection.Profiles, p)
			}
			localReader := &fakeReader{id: "local", searchResult: domain.DaemonSearch{Schema: domain.DaemonSearchSchema, ObservedAt: time.Now(), Query: "backup"}}
			targets := Targets(context.Background(), Target{Registration: Registration{Key: "local"}, Client: localReader}, fakeProfiles{value: collection}, secrets, nil)
			observer := New()
			observer.Timeout = 20 * time.Millisecond
			ctx := context.Background()
			if operation == "cancel" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 20*time.Millisecond)
				defer cancel()
				observer.Timeout = time.Second
			}
			type response struct {
				snapshot Snapshot
				err      error
			}
			done := make(chan response, 1)
			go func() {
				var result Snapshot
				var err error
				if operation == "search" {
					result, err = observer.Search(ctx, targets, SearchRequest{Query: "backup", Limit: 10})
				} else {
					result, err = observer.Summary(ctx, targets)
				}
				done <- response{result, err}
			}()
			select {
			case got := <-done:
				close(secrets.release)
				expected := ErrPartial
				if operation == "cancel" {
					expected = context.DeadlineExceeded
				}
				if !errors.Is(got.err, expected) || len(got.snapshot.Observations) != 9 || secrets.active.Load() != 0 {
					t.Fatalf("err=%v observations=%d active=%d", got.err, len(got.snapshot.Observations), secrets.active.Load())
				}
				good := 0
				for _, o := range got.snapshot.Observations {
					if o.Registration.Key == "local" {
						if o.Failure != nil {
							t.Fatal(o.Failure)
						}
						good++
					} else if o.Failure == nil || o.Failure.Code != "timed_out" {
						t.Fatalf("remote outcome %+v", o)
					}
				}
				if good != 1 {
					t.Fatal("successful local observation was lost")
				}
			case <-time.After(500 * time.Millisecond):
				close(secrets.release)
				<-done
				t.Fatal("credential lookup prevented bounded partial results and worker shutdown")
			}
		})
	}
}
