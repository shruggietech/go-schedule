package remotepairing

import (
	"context"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/remoteenroll"
)

type Draft struct {
	ProfileLabel    string            `json:"profile_label"`
	RepairProfileID string            `json:"repair_profile_id,omitempty"`
	Address         string            `json:"address"`
	DaemonID        string            `json:"daemon_id"`
	PairingID       string            `json:"pairing_id"`
	Phrase          string            `json:"phrase"`
	CertificatePEM  string            `json:"certificate_pem"`
	DisplayName     string            `json:"display_name"`
	Capability      domain.Capability `json:"capability"`
}

type Result struct {
	Action       string `json:"action"`
	Outcome      string `json:"outcome"`
	Message      string `json:"message"`
	CredentialID string `json:"credential_id,omitempty"`
	ProfileID    string `json:"profile_id,omitempty"`
}

type secretStore interface {
	Probe() error
	Save(string, string, string) error
	Delete(string, string) error
}
type profileStore interface {
	Load() (clientprofile.Collection, error)
	Add(clientprofile.Profile) (clientprofile.Profile, error)
	Replace(clientprofile.Profile) (clientprofile.Profile, error)
}
type Service struct {
	secrets  secretStore
	profiles profileStore
	now      func() time.Time
}

func New() *Service {
	return &Service{secrets: clientsecret.New(), profiles: clientprofile.NewStore(""), now: time.Now}
}
func NewWithStores(secrets secretStore, profiles profileStore) *Service {
	return &Service{secrets: secrets, profiles: profiles, now: time.Now}
}

func (s *Service) Pair(ctx context.Context, draft Draft) Result {
	if s.secrets.Probe() != nil {
		return rejected("Native credential storage is unavailable. No pairing attempt was sent.")
	}
	issued, manifest, err := remoteenroll.Exchange(ctx, remoteenroll.Draft{Address: draft.Address, DaemonID: draft.DaemonID, PairingID: draft.PairingID, Phrase: draft.Phrase, CertificatePEM: draft.CertificatePEM, DisplayName: draft.DisplayName, Kind: domain.ActorKindDesktop, Capability: draft.Capability})
	if err != nil {
		return rejected("Pairing failed. Verify the HTTPS endpoint, daemon identity, certificate, and one-time phrase.")
	}
	if err := s.secrets.Save(issued.DaemonID, issued.ID, issued.Token); err != nil {
		return rejected("The credential could not be stored securely. Revoke it from the daemon.")
	}
	now := s.now().UTC()
	label := strings.TrimSpace(draft.ProfileLabel)
	if label == "" {
		label = manifest.DisplayName
	}
	profileID, err := clientprofile.NewID()
	if err != nil {
		_ = s.secrets.Delete(issued.DaemonID, issued.ID)
		return rejected("The connection profile could not be created.")
	}
	profile := clientprofile.Profile{ID: profileID, Label: label, Endpoint: draft.Address, DaemonID: issued.DaemonID, CredentialID: issued.ID, CertificatePEM: draft.CertificatePEM, ClientKind: string(issued.Actor.Kind), Capability: string(issued.Actor.Capability), DaemonDisplayName: manifest.DisplayName, Platform: manifest.Platform.OS, Architecture: manifest.Platform.Architecture, ProductVersion: manifest.ProductVersion, CreatedAt: now, UpdatedAt: now, LastSuccessfulAt: now}
	profile, err = clientprofile.Normalize(profile)
	if err != nil {
		_ = s.secrets.Delete(issued.DaemonID, issued.ID)
		return rejected("The connection profile could not be validated safely.")
	}
	oldCredentialID := ""
	if draft.RepairProfileID != "" {
		collection, loadErr := s.profiles.Load()
		if loadErr != nil {
			_ = s.secrets.Delete(issued.DaemonID, issued.ID)
			return rejected("The existing connection profile could not be loaded.")
		}
		existing, findErr := collection.Find(draft.RepairProfileID)
		if findErr != nil || existing.DaemonID != issued.DaemonID {
			_ = s.secrets.Delete(issued.DaemonID, issued.ID)
			return rejected("Repair was rejected because the daemon identity changed.")
		}
		if strings.TrimSpace(draft.ProfileLabel) == "" {
			profile.Label = existing.Label
		}
		profile.ID, profile.CreatedAt, oldCredentialID = existing.ID, existing.CreatedAt, existing.CredentialID
		if _, err = s.profiles.Replace(profile); err == nil {
			if deleteErr := s.secrets.Delete(existing.DaemonID, oldCredentialID); deleteErr != nil {
				if _, rollbackErr := s.profiles.Replace(existing); rollbackErr == nil {
					_ = s.secrets.Delete(issued.DaemonID, issued.ID)
					return rejected("The old native credential could not be deleted, so the repair was rolled back.")
				}
				issued.Token = ""
				return Result{Action: "pair_remote", Outcome: "accepted", Message: "The connection was repaired, but the superseded native credential must be removed manually.", CredentialID: issued.ID, ProfileID: profile.ID}
			}
		}
	} else {
		_, err = s.profiles.Add(profile)
	}
	if err != nil {
		_ = s.secrets.Delete(issued.DaemonID, issued.ID)
		return rejected("The connection profile could not be saved safely.")
	}
	issued.Token = ""
	return Result{Action: "pair_remote", Outcome: "accepted", Message: "Remote connection paired and saved.", CredentialID: issued.ID, ProfileID: profile.ID}
}

func rejected(message string) Result {
	return Result{Action: "pair_remote", Outcome: "rejected", Message: message}
}
