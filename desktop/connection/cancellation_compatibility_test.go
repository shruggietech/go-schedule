package connection

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestProductVersionDoesNotReplaceAdvertisedProtocol(t *testing.T) {
	for _, version := range []string{"v2.0.3", "v99.0.0"} {
		health := server.HealthResponse{Status: "ok", Version: version}
		manifest := server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Owned", ProductVersion: version, LocalAPIVersions: []string{"v1"}, RemoteAPIVersions: []string{"v1"}}
		if _, err := NewLocalBackend(localDaemonFake{health: health, manifest: manifest}).Health(context.Background()); err != nil {
			t.Errorf("local product %s with v1: %v", version, err)
		}
		if _, err := NewRemoteBackend(remoteDaemonStub{health: health, manifest: manifest, capability: domain.CapabilityOperate}, "operate").Health(context.Background()); err != nil {
			t.Errorf("remote product %s with v1: %v", version, err)
		}
	}
	for _, versions := range [][]string{nil, {"v2"}, {"v10"}} {
		health := server.HealthResponse{Status: "ok", Version: "v1.5.0"}
		manifest := server.ManifestResponse{InstallationID: "daemon-1", DisplayName: "Owned", ProductVersion: health.Version, LocalAPIVersions: versions, RemoteAPIVersions: versions}
		if _, err := NewLocalBackend(localDaemonFake{health: health, manifest: manifest}).Health(context.Background()); err == nil {
			t.Errorf("local accepted unsupported APIs %v", versions)
		}
		if _, err := NewRemoteBackend(remoteDaemonStub{health: health, manifest: manifest}, "observe").Health(context.Background()); err == nil {
			t.Errorf("remote accepted unsupported APIs %v", versions)
		}
	}
}
