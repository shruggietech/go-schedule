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

func TestDefaultPathLivesUnderPlatformDataDirectory(t *testing.T) {
	path := DefaultPath()
	if !filepath.IsAbs(path) || filepath.Base(path) != "config.json" || filepath.Dir(path) != platform.DataDir() {
		t.Fatalf("default path=%q", path)
	}
}
