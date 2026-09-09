package server

import (
	"errors"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

type PairingCreateRequest struct {
	DisplayName string            `json:"display_name"`
	Kind        domain.ActorKind  `json:"kind"`
	Capability  domain.Capability `json:"capability"`
}

func (s *Server) handleCreatePairing(w http.ResponseWriter, r *http.Request) {
	var request PairingCreateRequest
	if decodeSingleJSON(r, &request) != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid pairing request")
		return
	}
	result, err := s.enrollment.Create(request.DisplayName, request.Kind, request.Capability)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid pairing request")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (s *Server) handleListPairings(w http.ResponseWriter, _ *http.Request) {
	items, err := s.enrollment.List()
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Pairings []domain.PairingSession `json:"pairings"`
	}{items})
}

func (s *Server) handleCancelPairing(w http.ResponseWriter, r *http.Request) {
	item, err := s.enrollment.Cancel(r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, CodeNotFound, "pairing", "pairing not found")
		return
	}
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleListCredentials(w http.ResponseWriter, _ *http.Request) {
	items, err := s.enrollment.Credentials()
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Credentials []domain.ClientCredential `json:"credentials"`
	}{items})
}

func (s *Server) handleRotateCredential(w http.ResponseWriter, r *http.Request) {
	item, err := s.enrollment.Rotate(r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, CodeNotFound, "credential", "credential not found")
		return
	}
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleRevokeCredential(w http.ResponseWriter, r *http.Request) {
	item, err := s.enrollment.Revoke(r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, CodeNotFound, "credential", "credential not found")
		return
	}
	if err != nil {
		s.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
