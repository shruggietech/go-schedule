package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeMCPHTTPManager struct {
	status      MCPHTTPStatusResponse
	enableReq   MCPHTTPEnableRequest
	credential  string
	disableCall int
	err         error
}

func (m *fakeMCPHTTPManager) Status() MCPHTTPStatusResponse { return m.status }
func (m *fakeMCPHTTPManager) Enable(_ context.Context, req MCPHTTPEnableRequest) (MCPHTTPCredentialResponse, error) {
	m.enableReq = req
	return MCPHTTPCredentialResponse{MCPHTTPStatusResponse: m.status, Credential: m.credential}, m.err
}
func (m *fakeMCPHTTPManager) Rotate(context.Context) (MCPHTTPCredentialResponse, error) {
	return MCPHTTPCredentialResponse{MCPHTTPStatusResponse: m.status, Credential: m.credential}, m.err
}
func (m *fakeMCPHTTPManager) Disable(context.Context) (MCPHTTPStatusResponse, error) {
	m.disableCall++
	return MCPHTTPStatusResponse{AllowedOrigins: []string{}}, m.err
}

func TestMCPHTTPLifecycleRoutesAndSecretBoundary(t *testing.T) {
	s := newTestServer(t)
	manager := &fakeMCPHTTPManager{status: MCPHTTPStatusResponse{Enabled: true, Endpoint: "http://127.0.0.1:43123/mcp", AllowedOrigins: []string{}}, credential: "one-time-secret"}
	s.SetMCPHTTPManager(manager)

	recorder := httptest.NewRecorder()
	s.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/mcp/http", nil))
	if recorder.Code != http.StatusOK || strings.Contains(recorder.Body.String(), manager.credential) {
		t.Fatalf("status response=%d %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	s.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/v1/mcp/http/enable", strings.NewReader(`{"port":43123,"allowed_origins":["http://127.0.0.1:3000"]}`)))
	if recorder.Code != http.StatusCreated || !strings.Contains(recorder.Body.String(), manager.credential) || manager.enableReq.Port != 43123 || len(manager.enableReq.AllowedOrigins) != 1 {
		t.Fatalf("enable response=%d %s request=%+v", recorder.Code, recorder.Body.String(), manager.enableReq)
	}

	recorder = httptest.NewRecorder()
	s.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/v1/mcp/http/rotate", nil))
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), manager.credential) {
		t.Fatalf("rotate response=%d %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	s.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/v1/mcp/http/disable", nil))
	if recorder.Code != http.StatusOK || manager.disableCall != 1 || strings.Contains(recorder.Body.String(), manager.credential) {
		t.Fatalf("disable response=%d %s calls=%d", recorder.Code, recorder.Body.String(), manager.disableCall)
	}
}

func TestMCPHTTPRouteErrorsAreSafeAndTyped(t *testing.T) {
	s := newTestServer(t)
	manager := &fakeMCPHTTPManager{}
	s.SetMCPHTTPManager(manager)
	tests := []struct {
		name   string
		body   string
		err    error
		status int
		code   string
	}{
		{name: "bad json", body: "{", status: http.StatusBadRequest, code: CodeValidation},
		{name: "validation", body: `{"port":1}`, err: &MCPHTTPControlError{Code: CodeValidation, Field: "allowed_origins", Message: "invalid origin"}, status: http.StatusBadRequest, code: CodeValidation},
		{name: "conflict", body: `{"port":1}`, err: &MCPHTTPControlError{Code: CodeConflict, Message: "port unavailable"}, status: http.StatusConflict, code: CodeConflict},
		{name: "internal", body: `{"port":1}`, err: errors.New("secret internal detail"), status: http.StatusInternalServerError, code: CodeInternal},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager.err = tc.err
			recorder := httptest.NewRecorder()
			s.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/v1/mcp/http/enable", strings.NewReader(tc.body)))
			var envelope APIError
			if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
				t.Fatal(err)
			}
			if recorder.Code != tc.status || envelope.Error.Code != tc.code || strings.Contains(recorder.Body.String(), "secret internal detail") {
				t.Fatalf("response=%d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
