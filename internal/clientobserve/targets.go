package clientobserve

import (
	"context"
	"fmt"
	"unicode"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

type ProfileStore interface {
	Load() (clientprofile.Collection, error)
}
type SecretStore interface {
	LoadContext(context.Context, string, string) (string, error)
}
type RemoteFactory func(endpoint, certificate, bearer, daemonID string) (Reader, error)

// Targets enumerates the explicit local target and saved registrations without changing the selected connection.
// A profile-store failure remains visible alongside the local observation rather than being silently omitted.
func Targets(ctx context.Context, local Target, profiles ProfileStore, secrets SecretStore, factory RemoteFactory) []Target {
	result := []Target{local}
	if ctx.Err() != nil {
		return result
	}
	if profiles == nil {
		return result
	}
	collection, err := profiles.Load()
	if err != nil {
		return append(result, Target{Registration: Registration{Key: "profiles", Label: "Saved connections"}, SetupFailure: &Failure{Code: "profile_store_unavailable", Message: "Saved connections could not be read. Check the profile store and try again."}})
	}
	if collection.Version != clientprofile.CurrentVersion || len(collection.Profiles) > clientprofile.MaxProfiles {
		return append(result, invalidInventory())
	}
	if factory == nil {
		factory = func(endpoint, certificate, bearer, daemonID string) (Reader, error) {
			return client.NewRemote(endpoint, certificate, bearer, daemonID)
		}
	}
	seen := map[string]bool{"local": true}
	for index, profile := range collection.Profiles {
		profile, normalizeErr := clientprofile.Normalize(profile)
		if normalizeErr != nil || !plainIdentity(profile.ID) || !plainIdentity(profile.DaemonID) || seen[profile.ID] {
			result = append(result, Target{Registration: Registration{Key: fmt.Sprintf("invalid-profile-%d", index+1), Label: fmt.Sprintf("Invalid saved connection %d", index+1)}, SetupFailure: &Failure{Code: "invalid_profile", Message: "This saved connection contains invalid metadata. Repair or remove it in the profile store."}})
			continue
		}
		seen[profile.ID] = true
		target := Target{Registration: Registration{Key: profile.ID, Label: profile.Label, Endpoint: profile.Endpoint, ExpectedDaemonID: profile.DaemonID}}
		target.Resolve = func(ctx context.Context) (Reader, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if secrets == nil {
				return nil, &setupError{Failure{Code: "credential_unavailable", Message: "The saved credential is unavailable. Repair this connection."}}
			}
			token, loadErr := secrets.LoadContext(ctx, profile.DaemonID, profile.CredentialID)
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if loadErr != nil {
				return nil, &setupError{Failure{Code: "credential_unavailable", Message: "The saved credential is unavailable. Repair this connection."}}
			}
			remote, remoteErr := factory(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
			if remoteErr != nil {
				return nil, &setupError{Failure{Code: "invalid_profile", Message: "The saved connection is invalid. Repair or remove this connection."}}
			}
			return remote, nil
		}
		result = append(result, target)
	}
	return result
}

type setupError struct{ failure Failure }

func (e *setupError) Error() string { return e.failure.Message }
func plainIdentity(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func invalidInventory() Target {
	return Target{Registration: Registration{Key: "profiles", Label: "Saved connections"}, SetupFailure: &Failure{Code: "invalid_profile_store", Message: "Saved connections have an unsupported version or exceed the profile limit. Repair the profile store."}}
}
