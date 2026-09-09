package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"runtime"

	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

// ManifestPlatform is the minimum platform information needed for compatibility.
type ManifestPlatform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

// ManifestResponse is the daemon-owned local discovery contract.
type ManifestResponse struct {
	InstallationID    string           `json:"installation_id"`
	DisplayName       string           `json:"display_name"`
	ProductVersion    string           `json:"product_version"`
	LocalAPIVersions  []string         `json:"local_api_versions"`
	RemoteAPIVersions []string         `json:"remote_api_versions"`
	OperatingMode     string           `json:"operating_mode"`
	Capabilities      []string         `json:"capabilities"`
	Platform          ManifestPlatform `json:"platform"`
}

type manifestRenameRequest struct {
	DisplayName string `json:"display_name"`
}

type manifestResetRequest struct {
	ConfirmInstallationID string `json:"confirm_installation_id"`
}

var daemonCapabilities = []string{"activity", "actor-authorization", "agent-access", "chains", "groups", "management-audit", "notifications", "schedule", "tasks", "triggers", "watchers"}

func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	s.writeManifest(w, r)
}

func (s *Server) handleRenameManifest(w http.ResponseWriter, r *http.Request) {
	var request manifestRenameRequest
	if err := decodeSingleJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if _, err := s.store.RenameDaemon(request.DisplayName); err != nil {
		if errors.Is(err, domain.ErrInvalidDaemonDisplayName) {
			writeError(w, http.StatusBadRequest, CodeValidation, "display_name", "display name must contain 1 through 80 characters and no control characters")
			return
		}
		s.internal(w, err)
		return
	}
	s.writeManifest(w, r)
}

func (s *Server) handleResetManifest(w http.ResponseWriter, r *http.Request) {
	var request manifestResetRequest
	if err := decodeSingleJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	if request.ConfirmInstallationID == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "confirm_installation_id", "current installation ID is required")
		return
	}
	if _, err := s.store.ResetDaemonIdentity(request.ConfirmInstallationID); err != nil {
		if errors.Is(err, store.ErrIdentityMismatch) {
			writeError(w, http.StatusConflict, CodeConflict, "confirm_installation_id", "installation identity changed or confirmation does not match")
			return
		}
		s.internal(w, err)
		return
	}
	s.writeManifest(w, r)
}

func (s *Server) writeManifest(w http.ResponseWriter, r *http.Request) {
	identity, err := s.store.DaemonIdentity()
	if err != nil {
		s.internal(w, err)
		return
	}
	platform := runtime.GOOS
	if platform == "darwin" {
		platform = "macos"
	}
	remoteVersions := []string{}
	capabilities := append([]string(nil), daemonCapabilities...)
	if r.Header.Get("X-Go-Schedule-Transport") == "remote" {
		remoteVersions = []string{"v1"}
		capabilities = append(capabilities, "remote-json")
	}
	writeJSON(w, http.StatusOK, ManifestResponse{
		InstallationID:    identity.InstallationID,
		DisplayName:       identity.DisplayName,
		ProductVersion:    buildinfo.Version,
		LocalAPIVersions:  []string{"v1"},
		RemoteAPIVersions: remoteVersions,
		OperatingMode:     "local_only",
		Capabilities:      capabilities,
		Platform:          ManifestPlatform{OS: platform, Architecture: runtime.GOARCH},
	})
}

func decodeSingleJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 64*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON value")
	}
	return nil
}
