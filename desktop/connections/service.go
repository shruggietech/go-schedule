package connections

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

type profileStore interface {
	Load() (clientprofile.Collection, error)
	Rename(string, string) (clientprofile.Profile, error)
	SetActive(string) error
	RemoveWith(string, func(clientprofile.Profile) error) (clientprofile.Profile, error)
}
type secretStore interface {
	Load(string, string) (string, error)
	Delete(string, string) error
}
type manager interface {
	Switch(connection.Backend, connection.Target) bool
}

type Service struct {
	profiles profileStore
	secrets  secretStore
	local    *client.Client
	router   *client.Client
	manager  manager
}

func New(profiles profileStore, secrets secretStore, local, router *client.Client, manager manager) *Service {
	return &Service{profiles: profiles, secrets: secrets, local: local, router: router, manager: manager}
}

func (s *Service) Workspace() Result {
	collection, err := s.profiles.Load()
	if err != nil {
		return rejected("load_connections", "Connection profiles could not be loaded.")
	}
	workspace := workspaceOf(collection)
	return Result{Action: "load_connections", Outcome: "accepted", Message: "Connection profiles loaded.", Workspace: &workspace}
}

func (s *Service) RestoreSelection(ctx context.Context) Result {
	collection, err := s.profiles.Load()
	if err != nil || collection.ActiveDesktopProfileID == "" {
		return s.selectLocal(false)
	}
	return s.selectProfile(ctx, collection.ActiveDesktopProfileID, false)
}

func (s *Service) Select(ctx context.Context, id string) Result {
	if strings.TrimSpace(id) == "" {
		return s.selectLocal(true)
	}
	return s.selectProfile(ctx, id, true)
}

func (s *Service) selectLocal(persist bool) Result {
	if persist {
		if err := s.profiles.SetActive(""); err != nil {
			return rejected("select_connection", "This computer could not be saved as the active connection.")
		}
	}
	if !s.manager.Switch(connection.NewLocalBackend(s.local), connection.LocalTarget()) {
		return rejected("select_connection", "The connection manager is closing.")
	}
	s.router.Use(s.local)
	workspace := s.Workspace()
	return Result{Action: "select_connection", Outcome: "accepted", Message: "This computer is selected.", Workspace: workspace.Workspace}
}

func (s *Service) selectProfile(ctx context.Context, id string, persist bool) Result {
	collection, err := s.profiles.Load()
	if err != nil {
		return rejected("select_connection", "Connection profiles could not be loaded.")
	}
	profile, err := collection.Find(id)
	if err != nil {
		return rejected("select_connection", "The selected connection profile does not exist or is ambiguous.")
	}
	if persist {
		if err := s.profiles.SetActive(profile.ID); err != nil {
			return rejected("select_connection", "The selected connection could not be saved.")
		}
	}
	target := connection.RemoteTarget(profile.ID, profile.DaemonID, profile.Label, profile.Endpoint, profile.CertificateFingerprint, profile.Platform, profile.Architecture, profile.ProductVersion)
	if !s.manager.Switch(connection.NewUnavailableRemoteBackend("The selected remote connection cannot be used until its credential and trust settings are available.", "Repair the connection or select another target."), target) {
		return rejected("select_connection", "The connection manager is closing.")
	}
	s.router.Use(client.NewUnavailableRemote(profile.Endpoint))
	token, err := s.secrets.Load(profile.DaemonID, profile.CredentialID)
	if err != nil {
		return rejected("select_connection", "The selected connection credential is unavailable. Repair the connection.")
	}
	remote, err := client.NewRemote(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
	if err != nil {
		return rejected("select_connection", "The selected connection profile is invalid.")
	}
	if !s.manager.Switch(connection.NewRemoteBackend(remote, profile.Capability), target) {
		return rejected("select_connection", "The connection manager is closing.")
	}
	s.router.Use(remote)
	workspace := s.Workspace()
	return Result{Action: "select_connection", Outcome: "accepted", Message: fmt.Sprintf("%s is selected.", profile.Label), Workspace: workspace.Workspace}
}

func (s *Service) Rename(id, label string) Result {
	if _, err := s.profiles.Rename(id, label); err != nil {
		return rejected("rename_connection", "The connection label could not be changed.")
	}
	workspace := s.Workspace()
	return Result{Action: "rename_connection", Outcome: "accepted", Message: "Connection label changed.", Workspace: workspace.Workspace}
}

func (s *Service) Remove(id string) Result {
	collection, err := s.profiles.Load()
	if err != nil {
		return rejected("remove_connection", "Connection profiles could not be loaded.")
	}
	profile, err := collection.Find(id)
	if err != nil {
		return rejected("remove_connection", "The selected connection profile does not exist or is ambiguous.")
	}
	if collection.ActiveDesktopProfileID == profile.ID {
		if result := s.selectLocal(true); result.Outcome != "accepted" {
			return rejected("remove_connection", "The active connection could not be changed to This computer.")
		}
	}
	if _, err := s.profiles.RemoveWith(profile.ID, func(current clientprofile.Profile) error {
		return s.secrets.Delete(current.DaemonID, current.CredentialID)
	}); err != nil {
		return rejected("remove_connection", "The native credential and profile could not be removed together. The profile was kept.")
	}
	workspace := s.Workspace()
	return Result{Action: "remove_connection", Outcome: "accepted", Message: "Connection removed.", Workspace: workspace.Workspace}
}

func workspaceOf(collection clientprofile.Collection) Workspace {
	workspace := Workspace{ActiveProfileID: collection.ActiveDesktopProfileID, Profiles: make([]Profile, 0, len(collection.Profiles))}
	for _, value := range collection.Profiles {
		last := ""
		if !value.LastSuccessfulAt.IsZero() {
			last = value.LastSuccessfulAt.Format(time.RFC3339)
		}
		workspace.Profiles = append(workspace.Profiles, Profile{ID: value.ID, Label: value.Label, Endpoint: value.Endpoint, DaemonID: value.DaemonID, ShortDaemonID: connection.ShortID(value.DaemonID), Fingerprint: value.CertificateFingerprint, Capability: value.Capability, Platform: value.Platform, Architecture: value.Architecture, ProductVersion: value.ProductVersion, LastSuccessfulAt: last, Active: value.ID == collection.ActiveDesktopProfileID})
	}
	return workspace
}

func rejected(action, message string) Result {
	return Result{Action: action, Outcome: "rejected", Message: message}
}
