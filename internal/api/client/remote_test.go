package client

import (
	"context"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

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
