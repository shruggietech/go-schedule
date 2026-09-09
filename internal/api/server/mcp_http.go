package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// MCPHTTPStatusResponse describes the optional runtime-only localhost endpoint.
// It never contains the bearer credential.
type MCPHTTPStatusResponse struct {
	Enabled               bool       `json:"enabled"`
	Endpoint              string     `json:"endpoint,omitempty"`
	AllowedOrigins        []string   `json:"allowed_origins"`
	CredentialFingerprint string     `json:"credential_fingerprint,omitempty"`
	EnabledAt             *time.Time `json:"enabled_at,omitempty"`
}

// MCPHTTPEnableRequest contains the only caller-selectable listener settings.
type MCPHTTPEnableRequest struct {
	Port           int      `json:"port"`
	AllowedOrigins []string `json:"allowed_origins"`
}

// MCPHTTPCredentialResponse returns a newly issued credential once.
type MCPHTTPCredentialResponse struct {
	MCPHTTPStatusResponse
	Credential string `json:"credential"`
}

// MCPHTTPManager controls one daemon-owned runtime listener.
type MCPHTTPManager interface {
	Status() MCPHTTPStatusResponse
	Enable(context.Context, MCPHTTPEnableRequest) (MCPHTTPCredentialResponse, error)
	Rotate(context.Context) (MCPHTTPCredentialResponse, error)
	Disable(context.Context) (MCPHTTPStatusResponse, error)
}

// MCPHTTPControlError classifies safe lifecycle failures for the local API.
type MCPHTTPControlError struct {
	Code    string
	Field   string
	Message string
}

func (e *MCPHTTPControlError) Error() string { return e.Message }

func mcpHTTPError(w http.ResponseWriter, err error) {
	var control *MCPHTTPControlError
	if errors.As(err, &control) {
		status := http.StatusBadRequest
		if control.Code == CodeConflict {
			status = http.StatusConflict
		}
		writeError(w, status, control.Code, control.Field, control.Message)
		return
	}
	writeError(w, http.StatusInternalServerError, CodeInternal, "", "localhost MCP operation failed")
}

func (s *Server) handleMCPHTTPStatus(w http.ResponseWriter, _ *http.Request) {
	if s.mcpHTTP == nil {
		writeJSON(w, http.StatusOK, MCPHTTPStatusResponse{AllowedOrigins: []string{}})
		return
	}
	writeJSON(w, http.StatusOK, s.mcpHTTP.Status())
}

func (s *Server) handleMCPHTTPEnable(w http.ResponseWriter, r *http.Request) {
	if s.mcpHTTP == nil {
		writeError(w, http.StatusServiceUnavailable, CodeInternal, "", "localhost MCP control is unavailable")
		return
	}
	var req MCPHTTPEnableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON body")
		return
	}
	result, err := s.mcpHTTP.Enable(r.Context(), req)
	if err != nil {
		mcpHTTPError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (s *Server) handleMCPHTTPRotate(w http.ResponseWriter, r *http.Request) {
	if s.mcpHTTP == nil {
		writeError(w, http.StatusServiceUnavailable, CodeInternal, "", "localhost MCP control is unavailable")
		return
	}
	result, err := s.mcpHTTP.Rotate(r.Context())
	if err != nil {
		mcpHTTPError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleMCPHTTPDisable(w http.ResponseWriter, r *http.Request) {
	if s.mcpHTTP == nil {
		writeJSON(w, http.StatusOK, MCPHTTPStatusResponse{AllowedOrigins: []string{}})
		return
	}
	result, err := s.mcpHTTP.Disable(r.Context())
	if err != nil {
		mcpHTTPError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
