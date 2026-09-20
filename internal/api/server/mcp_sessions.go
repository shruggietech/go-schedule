package server

import (
	"errors"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/mcpsession"
	"github.com/shruggietech/go-schedule/internal/store"
)

const MCPSessionHeader = "X-Go-Schedule-MCP-Session"
const ExpectedDaemonHeader = "X-Go-Schedule-Expected-Daemon"

type MCPSessionCreateRequest struct {
	ClientName string            `json:"client_name"`
	Capability domain.Capability `json:"capability"`
}

type MCPSessionCredentialResponse struct {
	mcpsession.Session
	Credential string `json:"credential"`
}

func (s *Server) handleCreateMCPSession(w http.ResponseWriter, r *http.Request) {
	var request MCPSessionCreateRequest
	if err := decodeSingleJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	session, credential, err := s.mcpSessions.Create(request.ClientName, request.Capability)
	if err != nil {
		if errors.Is(err, mcpsession.ErrInvalidSession) || errors.Is(err, domain.ErrInvalidActor) {
			writeError(w, http.StatusBadRequest, CodeValidation, "session", "client name and capability must define a valid Observe, Operate, or Manage session")
			return
		}
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, MCPSessionCredentialResponse{Session: session, Credential: credential})
}

func (s *Server) handleRevokeMCPSession(w http.ResponseWriter, r *http.Request) {
	if _, err := s.mcpSessions.Revoke(r.PathValue("id")); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			s.notFoundOr(w, err)
			return
		}
		s.internal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
