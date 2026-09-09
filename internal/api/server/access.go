package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/authorization"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

type ActorCreateRequest struct {
	Kind        domain.ActorKind  `json:"kind"`
	DisplayName string            `json:"display_name"`
	Capability  domain.Capability `json:"capability"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
}

type ActorUpdateRequest struct {
	DisplayName     *string            `json:"display_name,omitempty"`
	Capability      *domain.Capability `json:"capability,omitempty"`
	State           *domain.ActorState `json:"state,omitempty"`
	ExpiresAt       *time.Time         `json:"expires_at,omitempty"`
	ClearExpiration bool               `json:"clear_expiration,omitempty"`
}

type auditResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (w *auditResponseWriter) Header() http.Header { return w.header }
func (w *auditResponseWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
}
func (w *auditResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.body.Write(data)
}

func (s *Server) authorizeAndAudit(w http.ResponseWriter, r *http.Request) {
	operation, known := authorization.Lookup(r.Method, r.URL.Path)
	if !known {
		s.mux.ServeHTTP(w, r)
		return
	}
	actorID, resolveErr := s.resolveActor(r)
	actor, actorErr := s.store.GetActor(actorID)
	authorized := resolveErr == nil && actorErr == nil && authorization.Allows(actor, operation.ID, time.Now())
	identity, identityErr := s.store.DaemonIdentity()
	if identityErr != nil {
		writeError(w, http.StatusServiceUnavailable, CodeAuditUnavailable, "", "authorization evidence is unavailable")
		return
	}
	auditActorID := actorID
	if _, err := uuid.Parse(auditActorID); resolveErr != nil || err != nil {
		auditActorID = ""
	}
	event := domain.AuditEvent{ID: uuid.NewString(), ActorID: auditActorID, DaemonID: identity.InstallationID, Operation: operation.ID, TargetKind: operation.TargetKind, TargetID: authorization.TargetID(operation, r.URL.Path), CorrelationID: uuid.NewString(), OccurredAt: time.Now().UTC()}
	if !authorized {
		if err := s.store.RecordDeniedAudit(event); err != nil {
			writeError(w, http.StatusServiceUnavailable, CodeAuditUnavailable, "", "authorization evidence is unavailable")
			return
		}
		writeError(w, http.StatusForbidden, CodeForbidden, "", "actor is not authorized for this operation")
		return
	}
	if operation.Audit == authorization.AuditNone {
		s.mux.ServeHTTP(w, r)
		return
	}
	if err := s.store.BeginAudit(event); err != nil {
		writeError(w, http.StatusServiceUnavailable, CodeAuditUnavailable, "", "audit intent could not be persisted")
		return
	}
	capture := &auditResponseWriter{header: make(http.Header)}
	s.mux.ServeHTTP(capture, r)
	if capture.status == 0 {
		capture.status = http.StatusOK
	}
	result := domain.AuditResultSucceeded
	if capture.status >= 400 {
		result = domain.AuditResultFailed
	}
	if err := s.store.CompleteAudit(event.ID, result); err != nil {
		writeError(w, http.StatusInternalServerError, CodeAuditUnavailable, "", "audit result could not be persisted")
		return
	}
	for key, values := range capture.header {
		w.Header()[key] = append([]string(nil), values...)
	}
	w.WriteHeader(capture.status)
	_, _ = w.Write(capture.body.Bytes())
}

func (s *Server) handleListActors(w http.ResponseWriter, _ *http.Request) {
	actors, err := s.store.ListActors()
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Actors []domain.Actor `json:"actors"`
	}{Actors: actors})
}

func (s *Server) handleCreateActor(w http.ResponseWriter, r *http.Request) {
	var request ActorCreateRequest
	if err := decodeSingleJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	actor, err := s.store.CreateActor(request.Kind, request.DisplayName, request.Capability, request.ExpiresAt)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidActor) {
			writeError(w, http.StatusBadRequest, CodeValidation, "actor", "actor fields are invalid")
			return
		}
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, actor)
}

func (s *Server) handleUpdateActor(w http.ResponseWriter, r *http.Request) {
	var request ActorUpdateRequest
	if err := decodeSingleJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if request.DisplayName == nil && request.Capability == nil && request.State == nil && request.ExpiresAt == nil && !request.ClearExpiration {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "provide at least one actor field")
		return
	}
	var expires **time.Time
	if request.ExpiresAt != nil {
		value := request.ExpiresAt
		expires = &value
	}
	if request.ClearExpiration {
		var value *time.Time
		expires = &value
	}
	actor, err := s.store.UpdateActor(r.PathValue("id"), store.ActorUpdate{DisplayName: request.DisplayName, Capability: request.Capability, State: request.State, ExpiresAt: expires})
	s.writeActorResult(w, actor, err)
}

func (s *Server) handleRevokeActor(w http.ResponseWriter, r *http.Request) {
	actor, err := s.store.RevokeActor(r.PathValue("id"))
	s.writeActorResult(w, actor, err)
}

func (s *Server) writeActorResult(w http.ResponseWriter, actor domain.Actor, err error) {
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, actor)
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, CodeNotFound, "actor", "actor not found")
	case errors.Is(err, store.ErrBuiltinActor):
		writeError(w, http.StatusConflict, CodeConflict, "actor", "built-in local actor is protected")
	case errors.Is(err, domain.ErrInvalidActor):
		writeError(w, http.StatusBadRequest, CodeValidation, "actor", "actor fields are invalid")
	default:
		s.internal(w, err)
	}
}

func parseAuditQuery(r *http.Request) (domain.AuditQuery, error) {
	query := domain.AuditQuery{ActorID: r.URL.Query().Get("actor_id"), Operation: r.URL.Query().Get("operation"), Result: domain.AuditResult(r.URL.Query().Get("result"))}
	if value := r.URL.Query().Get("limit"); value != "" {
		limit, err := strconv.Atoi(value)
		if err != nil {
			return query, err
		}
		query.Limit = limit
	}
	if value := r.URL.Query().Get("since"); value != "" {
		parsed, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return query, err
		}
		query.Since = &parsed
	}
	if value := r.URL.Query().Get("until"); value != "" {
		parsed, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return query, err
		}
		query.Until = &parsed
	}
	if query.Limit < 0 || query.Limit > 1000 || (query.Result != "" && !query.Result.Valid()) || (query.Since != nil && query.Until != nil && query.Since.After(*query.Until)) {
		return query, fmt.Errorf("invalid audit query")
	}
	if query.ActorID != "" {
		if _, err := uuid.Parse(query.ActorID); err != nil {
			return query, fmt.Errorf("invalid audit actor")
		}
	}
	if query.Operation != "" {
		if _, ok := authorization.LookupID(query.Operation); !ok {
			return query, fmt.Errorf("invalid audit operation")
		}
	}
	return query, nil
}

func (s *Server) handleListAudit(w http.ResponseWriter, r *http.Request) {
	query, err := parseAuditQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "query", "audit filters are invalid")
		return
	}
	events, err := s.store.ListAudit(query)
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Events []domain.AuditEvent `json:"events"`
	}{Events: events})
}

func (s *Server) handleExportAudit(w http.ResponseWriter, r *http.Request) {
	query, err := parseAuditQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "query", "audit filters are invalid")
		return
	}
	if query.Limit == 0 {
		query.Limit = 1000
	}
	events, err := s.store.ListAudit(query)
	if err != nil {
		s.internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	for _, event := range events {
		if err := encoder.Encode(event); err != nil {
			return
		}
	}
}
