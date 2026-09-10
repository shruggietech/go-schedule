package client

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func serverCertificatePEM(t *testing.T, server *httptest.Server) string {
	t.Helper()
	certificate := server.Certificate()
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Raw}))
}

func TestSwitchableClientDoesNotRetargetInflightRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	first := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		once.Do(func() { close(started) })
		<-release
		_, _ = w.Write([]byte(`{"status":"ok","version":"first"}`))
	}))
	first.TLS = first.Config.TLSConfig
	first.StartTLS()
	defer first.Close()
	second := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"status":"ok","version":"second"}`))
	}))
	second.TLS = second.Config.TLSConfig
	second.StartTLS()
	defer second.Close()
	firstClient, _ := NewRemote(first.URL, serverCertificatePEM(t, first), "first-token", "first-id")
	secondClient, _ := NewRemote(second.URL, serverCertificatePEM(t, second), "second-token", "second-id")
	router := NewSwitchable(firstClient)
	result := make(chan string, 1)
	go func() { health, _ := router.Health(context.Background()); result <- health.Version }()
	<-started
	router.Use(secondClient)
	close(release)
	if version := <-result; version != "first" {
		t.Fatalf("in-flight version=%q", version)
	}
	health, err := router.Health(context.Background())
	if err != nil || health.Version != "second" {
		t.Fatalf("next health=%+v err=%v", health, err)
	}
}

func TestRemoteClientPinsIdentityMapsPathAndAddsBearer(t *testing.T) {
	var taskAuthorization string
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/v1/manifest":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"installation_id":"daemon-1","display_name":"Remote","product_version":"v1.0.0","remote_api_versions":["v1"],"capabilities":["tasks"],"platform":{"os":"linux","architecture":"amd64"}}`))
		case "/api/v1/tasks":
			taskAuthorization = request.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tasks":[]}`))
		default:
			http.NotFound(w, request)
		}
	}))
	server.TLS = server.Config.TLSConfig
	server.StartTLS()
	defer server.Close()

	remote, err := NewRemote(server.URL, serverCertificatePEM(t, server), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := remote.ListTaskObservations(context.Background(), "", "", false, 0, 10, 100); err == nil || !strings.Contains(err.Error(), "identity") {
		t.Fatalf("unverified err = %v", err)
	}
	manifest, err := remote.VerifyIdentity(context.Background())
	if err != nil || manifest.InstallationID != "daemon-1" {
		t.Fatalf("manifest = %#v, err = %v", manifest, err)
	}
	if _, err := remote.ListTaskObservations(context.Background(), "", "", false, 0, 10, 100); err != nil {
		t.Fatal(err)
	}
	if taskAuthorization != "Bearer bearer-canary" {
		t.Fatalf("authorization = %q", taskAuthorization)
	}
}

func TestRemoteVerifyAccessReturnsCurrentServerAuthority(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/api/v1/manifest":
			_, _ = w.Write([]byte(`{"installation_id":"daemon-1","display_name":"Remote","product_version":"v1.0.0","remote_api_versions":["v1"],"capabilities":["tasks"],"platform":{"os":"linux","architecture":"amd64"}}`))
		case "/api/v1/access/current":
			_, _ = w.Write([]byte(`{"id":"actor-1","kind":"desktop","display_name":"Remote desktop","capability":"observe","state":"active","builtin":false,"created_at":"2026-09-09T12:00:00Z","updated_at":"2026-09-09T12:00:00Z"}`))
		default:
			http.NotFound(w, request)
		}
	}))
	server.TLS = server.Config.TLSConfig
	server.StartTLS()
	defer server.Close()
	remote, err := NewRemote(server.URL, serverCertificatePEM(t, server), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := remote.VerifyIdentity(context.Background()); err != nil {
		t.Fatal(err)
	}
	capability, err := remote.VerifyAccess(context.Background())
	if err != nil || capability != domain.CapabilityObserve {
		t.Fatalf("capability=%q err=%v", capability, err)
	}
}

func TestRemoteClientRejectsIdentityMismatchAndRedirect(t *testing.T) {
	redirected := false
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/api/v1/manifest" {
			http.Redirect(w, request, "/captured", http.StatusTemporaryRedirect)
			return
		}
		redirected = true
	}))
	server.TLS = server.Config.TLSConfig
	server.StartTLS()
	defer server.Close()
	remote, err := NewRemote(server.URL, serverCertificatePEM(t, server), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := remote.VerifyIdentity(context.Background()); err == nil {
		t.Fatal("redirect was accepted")
	}
	if redirected {
		t.Fatal("redirect target was contacted")
	}

	mismatch := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"installation_id":"daemon-2","display_name":"Remote","product_version":"v1.0.0","remote_api_versions":["v1"],"capabilities":[],"platform":{"os":"linux","architecture":"amd64"}}`))
	}))
	mismatch.TLS = mismatch.Config.TLSConfig
	mismatch.StartTLS()
	defer mismatch.Close()
	remote, err = NewRemote(mismatch.URL, serverCertificatePEM(t, mismatch), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := remote.VerifyIdentity(context.Background()); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("err = %v", err)
	}
}

func TestRemoteTaskResponsesDecodeObservationProjection(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/api/v1/manifest":
			_, _ = w.Write([]byte(`{"installation_id":"daemon-1","display_name":"Remote","product_version":"v1.0.0","remote_api_versions":["v1"],"capabilities":["tasks"],"platform":{"os":"linux","architecture":"amd64"}}`))
		case "/api/v1/tasks":
			_, _ = w.Write([]byte(`{"tasks":[{"id":"task-1","name":"Backup","enabled":true,"state":"active","timezone":"UTC","readiness":"ready","readiness_reason":"Ready.","has_schedule":true,"schedule_summary":"Every day","policy_summary":"Queue one","next_runs":[],"updated_at":"2026-09-09T12:00:00Z"}]}`))
		case "/api/v1/tasks/task-1":
			_, _ = w.Write([]byte(`{"id":"task-1","name":"Backup","enabled":true,"state":"active","timezone":"UTC","readiness":"ready","readiness_reason":"Ready.","has_schedule":true,"schedule_summary":"Every day","policy_summary":"Queue one","next_runs":[],"updated_at":"2026-09-09T12:00:00Z"}`))
		default:
			http.NotFound(w, request)
		}
	}))
	server.TLS = server.Config.TLSConfig
	server.StartTLS()
	defer server.Close()
	remote, err := NewRemote(server.URL, serverCertificatePEM(t, server), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := remote.VerifyIdentity(context.Background()); err != nil {
		t.Fatal(err)
	}
	details, err := remote.ListTaskDetails(context.Background(), "", "")
	if err != nil || len(details) != 1 || details[0].Task.ID != "task-1" || details[0].Task.Name != "Backup" || details[0].Schedule == nil || details[0].Schedule.HumanSummary != "Every day" || !details[0].Readiness.CommandReady {
		t.Fatalf("details=%+v err=%v", details, err)
	}
	detail, err := remote.GetTask(context.Background(), "task-1")
	if err != nil || detail.Task.ID != "task-1" || detail.PolicySummary != "Queue one" {
		t.Fatalf("detail=%+v err=%v", detail, err)
	}
}

func TestRemoteIdentityVerificationHasTransportAndOperationBounds(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	server.TLS = server.Config.TLSConfig
	server.StartTLS()
	defer server.Close()
	remote, err := NewRemote(server.URL, serverCertificatePEM(t, server), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	transport := remote.http.Transport.(*http.Transport)
	if transport.TLSHandshakeTimeout <= 0 || transport.ResponseHeaderTimeout <= 0 {
		t.Fatalf("transport timeouts are incomplete: %#v", transport)
	}
	remote.identityTimeout = 25 * time.Millisecond
	started := time.Now()
	if _, err := remote.VerifyIdentity(context.Background()); err == nil || time.Since(started) > time.Second {
		t.Fatalf("verification err=%v elapsed=%s", err, time.Since(started))
	}
}

func TestRemoteTaskDetailsReadEveryObservationPage(t *testing.T) {
	const taskCount = 205
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if request.URL.Path == "/api/v1/manifest" {
			_, _ = w.Write([]byte(`{"installation_id":"daemon-1","display_name":"Remote","product_version":"v1.0.0","remote_api_versions":["v1"],"capabilities":["tasks"],"platform":{"os":"linux","architecture":"amd64"}}`))
			return
		}
		offset, _ := strconv.Atoi(request.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
		end := min(offset+limit, taskCount)
		observations := make([]server.TaskObservationResponse, 0, end-offset)
		for i := offset; i < end; i++ {
			observations = append(observations, server.TaskObservationResponse{ID: "task-" + strconv.Itoa(i), Name: "Task", Enabled: true, State: "active", Timezone: "UTC", Readiness: "ready", NextRuns: []time.Time{}})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"tasks": observations})
	}))
	server.TLS = server.Config.TLSConfig
	server.StartTLS()
	defer server.Close()
	remote, err := NewRemote(server.URL, serverCertificatePEM(t, server), "bearer-canary", "daemon-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := remote.VerifyIdentity(context.Background()); err != nil {
		t.Fatal(err)
	}
	details, err := remote.ListTaskDetails(context.Background(), "", "")
	if err != nil || len(details) != taskCount {
		t.Fatalf("details=%d err=%v", len(details), err)
	}
	if details[0].Task.ID != "task-0" || details[taskCount-1].Task.ID != "task-204" {
		t.Fatalf("first=%+v last=%+v", details[0].Task, details[taskCount-1].Task)
	}
}

func TestRemoteMutationTransportFailureIsUncertainAndSingleAttempt(t *testing.T) {
	attempts := 0
	remote := &Client{http: &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		return nil, errors.New("private bearer-canary transport detail")
	})}, baseURL: "https://example.test", pathPrefix: "/api", remote: true, bearer: "bearer-canary"}
	remote.verified.Store(true)
	err := remote.RunNow(context.Background(), "task-1")
	var uncertain *MutationUncertainError
	if !errors.As(err, &uncertain) || attempts != 1 {
		t.Fatalf("err=%T %v attempts=%d", err, err, attempts)
	}
	if strings.Contains(err.Error(), "bearer-canary") || !strings.Contains(err.Error(), "may have completed") {
		t.Fatalf("unsafe or unactionable error=%q", err)
	}
}

func TestRemoteMutationAmbiguousServerResponseIsUncertain(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
	}{{name: "invalid success body", status: http.StatusOK, body: "{"}, {name: "server failure", status: http.StatusInternalServerError, body: "not-json"}} {
		t.Run(test.name, func(t *testing.T) {
			attempts := 0
			remote := &Client{http: &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
				attempts++
				return &http.Response{StatusCode: test.status, Body: io.NopCloser(strings.NewReader(test.body)), Header: make(http.Header)}, nil
			})}, baseURL: "https://example.test", pathPrefix: "/api", remote: true}
			remote.verified.Store(true)
			_, err := remote.CreateTask(context.Background(), server.TaskCreateRequest{Name: "task"})
			var uncertain *MutationUncertainError
			if !errors.As(err, &uncertain) || attempts != 1 {
				t.Fatalf("err=%T %v attempts=%d", err, err, attempts)
			}
		})
	}
}

func TestConnectionErrorClassifiesTrustFailureWithoutCauseDisclosure(t *testing.T) {
	err := NewConnectionError("GET /v1/manifest", x509.UnknownAuthorityError{})
	if err.Kind != ConnectionTrustFailure || strings.Contains(err.Error(), "certificate signed") {
		t.Fatalf("error=%+v message=%q", err, err.Error())
	}
}
