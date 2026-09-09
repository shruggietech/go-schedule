// Package clientprofile owns secret-free remote connection profile metadata.
package clientprofile

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	CurrentVersion = 1
	MaxProfiles    = 100
)

var (
	ErrInvalid   = errors.New("invalid connection profile")
	ErrNotFound  = errors.New("connection profile not found")
	ErrAmbiguous = errors.New("connection profile label is ambiguous")
	ErrBusy      = errors.New("connection profile store is busy")
	ErrFuture    = errors.New("connection profile version is newer than this application")
)

// Profile pins one user-owned client relationship without retaining its bearer credential.
type Profile struct {
	ID                     string    `json:"id"`
	Label                  string    `json:"label"`
	Endpoint               string    `json:"endpoint"`
	DaemonID               string    `json:"daemon_id"`
	CredentialID           string    `json:"credential_id"`
	CertificatePEM         string    `json:"certificate_pem"`
	CertificateFingerprint string    `json:"certificate_fingerprint"`
	ClientKind             string    `json:"client_kind"`
	Capability             string    `json:"capability"`
	DaemonDisplayName      string    `json:"daemon_display_name"`
	Platform               string    `json:"platform,omitempty"`
	Architecture           string    `json:"architecture,omitempty"`
	ProductVersion         string    `json:"product_version,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	LastSuccessfulAt       time.Time `json:"last_successful_at,omitempty"`
}

// Collection is the complete versioned user profile document.
type Collection struct {
	Version                int       `json:"version"`
	ActiveDesktopProfileID string    `json:"active_desktop_profile_id,omitempty"`
	Profiles               []Profile `json:"profiles"`
}

// NewID creates a random local profile identity.
func NewID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

// Normalize validates and canonicalizes a profile before persistence.
func Normalize(profile Profile) (Profile, error) {
	profile.ID = strings.TrimSpace(profile.ID)
	profile.Label = strings.TrimSpace(profile.Label)
	profile.DaemonID = strings.TrimSpace(profile.DaemonID)
	profile.CredentialID = strings.TrimSpace(profile.CredentialID)
	profile.ClientKind = strings.TrimSpace(profile.ClientKind)
	profile.Capability = strings.TrimSpace(profile.Capability)
	profile.DaemonDisplayName = strings.TrimSpace(profile.DaemonDisplayName)
	if profile.ID == "" || profile.DaemonID == "" || profile.CredentialID == "" || !validLabel(profile.Label) {
		return Profile{}, ErrInvalid
	}
	endpoint, err := NormalizeEndpoint(profile.Endpoint)
	if err != nil {
		return Profile{}, err
	}
	fingerprint, err := Fingerprint(profile.CertificatePEM)
	if err != nil {
		return Profile{}, err
	}
	if profile.CertificateFingerprint != "" && !strings.EqualFold(profile.CertificateFingerprint, fingerprint) {
		return Profile{}, ErrInvalid
	}
	profile.Endpoint = endpoint
	profile.CertificateFingerprint = fingerprint
	if profile.ClientKind != "desktop" && profile.ClientKind != "cli" && profile.ClientKind != "json" {
		return Profile{}, ErrInvalid
	}
	if profile.Capability != "observe" && profile.Capability != "operate" && profile.Capability != "manage" && profile.Capability != "enroll" {
		return Profile{}, ErrInvalid
	}
	if profile.CreatedAt.IsZero() || profile.UpdatedAt.IsZero() {
		return Profile{}, ErrInvalid
	}
	return profile, nil
}

// NormalizeEndpoint accepts one HTTPS origin and removes its optional trailing slash.
func NormalizeEndpoint(value string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || (parsed.Path != "" && parsed.Path != "/") {
		return "", ErrInvalid
	}
	parsed.Path = ""
	return strings.TrimRight(parsed.String(), "/"), nil
}

// Fingerprint validates certificate-only PEM and returns the leaf SHA-256 fingerprint.
func Fingerprint(value string) (string, error) {
	rest := []byte(value)
	var first *x509.Certificate
	for len(rest) > 0 {
		block, remaining := pem.Decode(rest)
		if block == nil || block.Type != "CERTIFICATE" {
			return "", ErrInvalid
		}
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return "", ErrInvalid
		}
		if first == nil {
			first = certificate
		}
		rest = remaining
	}
	if first == nil {
		return "", ErrInvalid
	}
	sum := sha256.Sum256(first.Raw)
	return hex.EncodeToString(sum[:]), nil
}

func validLabel(value string) bool {
	if !utf8.ValidString(value) || value == "" || utf8.RuneCountInString(value) > 80 {
		return false
	}
	for _, character := range value {
		if unicode.IsControl(character) {
			return false
		}
	}
	return true
}

// Find resolves one exact ID or one unambiguous case-insensitive label.
func (collection Collection) Find(reference string) (Profile, error) {
	reference = strings.TrimSpace(reference)
	for _, profile := range collection.Profiles {
		if profile.ID == reference {
			return profile, nil
		}
	}
	var found *Profile
	for i := range collection.Profiles {
		if strings.EqualFold(collection.Profiles[i].Label, reference) {
			if found != nil {
				return Profile{}, ErrAmbiguous
			}
			value := collection.Profiles[i]
			found = &value
		}
	}
	if found == nil {
		return Profile{}, ErrNotFound
	}
	return *found, nil
}
