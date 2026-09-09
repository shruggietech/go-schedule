package clientprofile

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func testCertificate(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "daemon.test"}, DNSNames: []string{"daemon.test"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

func validProfile(t *testing.T) Profile {
	now := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	return Profile{ID: "profile-1", Label: "Production", Endpoint: "https://daemon.test:8443/", DaemonID: "daemon-1", CredentialID: "credential-1", CertificatePEM: testCertificate(t), ClientKind: "cli", Capability: "manage", DaemonDisplayName: "Scheduler", CreatedAt: now, UpdatedAt: now}
}

func TestNormalizeCanonicalizesSafeProfile(t *testing.T) {
	profile, err := Normalize(validProfile(t))
	if err != nil {
		t.Fatal(err)
	}
	if profile.Endpoint != "https://daemon.test:8443" || len(profile.CertificateFingerprint) != 64 {
		t.Fatalf("normalized = %#v", profile)
	}
}

func TestNormalizeRejectsUnsafeEndpointAndPrivateKey(t *testing.T) {
	profile := validProfile(t)
	for _, endpoint := range []string{"http://daemon.test", "https://user@daemon.test", "https://daemon.test/api", "https://daemon.test?q=1"} {
		profile.Endpoint = endpoint
		if _, err := Normalize(profile); err == nil {
			t.Fatalf("accepted %q", endpoint)
		}
	}
	profile = validProfile(t)
	profile.CertificatePEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("secret")}))
	if _, err := Normalize(profile); err == nil {
		t.Fatal("accepted private key")
	}
}

func TestFindRejectsAmbiguousLabels(t *testing.T) {
	a, b := validProfile(t), validProfile(t)
	b.ID = "profile-2"
	collection := Collection{Version: CurrentVersion, Profiles: []Profile{a, b}}
	if _, err := collection.Find("Production"); err != ErrAmbiguous {
		t.Fatalf("err = %v", err)
	}
	if profile, err := collection.Find("profile-2"); err != nil || profile.ID != "profile-2" {
		t.Fatalf("profile = %#v, err = %v", profile, err)
	}
}
