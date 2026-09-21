package config

import (
	"path/filepath"
	"testing"

	"github.com/shruggietech/go-schedule/internal/platform"
)

func TestRemoteDefaultsDisabledAndRequiresSafeExplicitConfiguration(t *testing.T) {
	if Default().Remote.Enabled {
		t.Fatal("remote access must default disabled")
	}
	tests := []struct {
		name   string
		config RemoteConfig
		valid  bool
	}{
		{"loopback", RemoteConfig{Enabled: true, BindAddress: "127.0.0.1:8443", CertificateFile: "cert.pem", PrivateKeyFile: "key.pem"}, true},
		{"private", RemoteConfig{Enabled: true, BindAddress: "10.0.0.2:8443", CertificateFile: "cert.pem", PrivateKeyFile: "key.pem"}, true},
		{"wildcard unacknowledged", RemoteConfig{Enabled: true, BindAddress: "0.0.0.0:8443", CertificateFile: "cert.pem", PrivateKeyFile: "key.pem"}, false},
		{"wildcard acknowledged", RemoteConfig{Enabled: true, BindAddress: "0.0.0.0:8443", CertificateFile: "cert.pem", PrivateKeyFile: "key.pem", AcknowledgePublicExposure: true}, true},
		{"hostname", RemoteConfig{Enabled: true, BindAddress: "localhost:8443", CertificateFile: "cert.pem", PrivateKeyFile: "key.pem"}, false},
		{"missing certificate", RemoteConfig{Enabled: true, BindAddress: "127.0.0.1:8443", PrivateKeyFile: "key.pem"}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.config.Validate() == nil; got != test.valid {
				t.Fatalf("valid=%v, want %v", got, test.valid)
			}
		})
	}
}

func TestRemoteMCPRequiresCanonicalHTTPSResource(t *testing.T) {
	base := RemoteConfig{Enabled: true, BindAddress: "127.0.0.1:8443", CertificateFile: "cert.pem", PrivateKeyFile: "key.pem"}
	for _, test := range []struct {
		name  string
		mcp   RemoteMCPConfig
		valid bool
	}{
		{"disabled", RemoteMCPConfig{}, true},
		{"canonical", RemoteMCPConfig{Enabled: true, ResourceURL: "https://scheduler.example/mcp"}, true},
		{"bounded lifetime", RemoteMCPConfig{Enabled: true, ResourceURL: "https://scheduler.example/mcp", AccessTokenLifetimeSeconds: 3600}, true},
		{"http", RemoteMCPConfig{Enabled: true, ResourceURL: "http://scheduler.example/mcp"}, false},
		{"wrong path", RemoteMCPConfig{Enabled: true, ResourceURL: "https://scheduler.example/api/mcp"}, false},
		{"query", RemoteMCPConfig{Enabled: true, ResourceURL: "https://scheduler.example/mcp?x=1"}, false},
		{"excessive lifetime", RemoteMCPConfig{Enabled: true, ResourceURL: "https://scheduler.example/mcp", AccessTokenLifetimeSeconds: 3601}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate := base
			candidate.MCP = test.mcp
			if got := candidate.Validate() == nil; got != test.valid {
				t.Fatalf("valid=%v, want %v", got, test.valid)
			}
		})
	}
	if err := (RemoteConfig{MCP: RemoteMCPConfig{Enabled: true, ResourceURL: "https://scheduler.example/mcp"}}).Validate(); err == nil {
		t.Fatal("MCP must not be enabled without the remote listener")
	}
}

func TestDefaultPathLivesUnderPlatformDataDirectory(t *testing.T) {
	path := DefaultPath()
	if !filepath.IsAbs(path) || filepath.Base(path) != "config.json" || filepath.Dir(path) != platform.DataDir() {
		t.Fatalf("default path=%q", path)
	}
}
