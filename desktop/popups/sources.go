package popups

import (
	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
)

type profileStore interface {
	Load() (clientprofile.Collection, error)
}
type secretStore interface {
	Load(string, string) (string, error)
}

// RegisteredSources constructs identity-pinned clients and skips broken registrations.
// It does not alter the desktop's selected target or probe a daemon for write authority.
func RegisteredSources(profiles profileStore, secrets secretStore, local *client.Client) SourceLoader {
	return func() []Source {
		result := make([]Source, 0, maxSources)
		if local != nil {
			result = append(result, Source{Label: "This computer", Client: local, Backend: connection.NewLocalBackend(local)})
		}
		if profiles == nil || secrets == nil {
			return result
		}
		collection, err := profiles.Load()
		if err != nil {
			return result
		}
		for _, profile := range collection.Profiles {
			if len(result) >= maxSources {
				break
			}
			token, loadErr := secrets.Load(profile.DaemonID, profile.CredentialID)
			if loadErr != nil {
				continue
			}
			remote, remoteErr := client.NewRemote(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
			if remoteErr != nil {
				continue
			}
			result = append(result, Source{ProfileID: profile.ID, DaemonID: profile.DaemonID, Label: profile.Label, Client: remote, Backend: connection.NewRemoteBackend(remote, profile.Capability)})
		}
		return result
	}
}
