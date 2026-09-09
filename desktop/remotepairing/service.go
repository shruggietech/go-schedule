package remotepairing

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type Draft struct {
	Address        string            `json:"address"`
	DaemonID       string            `json:"daemon_id"`
	PairingID      string            `json:"pairing_id"`
	Phrase         string            `json:"phrase"`
	CertificatePEM string            `json:"certificate_pem"`
	DisplayName    string            `json:"display_name"`
	Capability     domain.Capability `json:"capability"`
}

type Result struct {
	Action       string `json:"action"`
	Outcome      string `json:"outcome"`
	Message      string `json:"message"`
	CredentialID string `json:"credential_id,omitempty"`
}

type secretStore interface {
	Probe() error
	Save(string, string, string) error
}
type Service struct{ secrets secretStore }

func New() *Service                           { return &Service{secrets: clientsecret.New()} }
func NewWithStore(store secretStore) *Service { return &Service{secrets: store} }

func (s *Service) Pair(ctx context.Context, draft Draft) Result {
	if s.secrets.Probe() != nil {
		return rejected("Native credential storage is unavailable. No pairing attempt was sent.")
	}
	base, err := url.Parse(strings.TrimRight(draft.Address, "/"))
	if err != nil || base.Scheme != "https" || base.Host == "" || draft.DaemonID == "" || draft.PairingID == "" || draft.Phrase == "" || draft.DisplayName == "" || !draft.Capability.Valid() {
		return rejected("Provide a valid HTTPS address, daemon ID, pairing ID, and phrase.")
	}
	pool := x509.NewCertPool()
	if draft.CertificatePEM == "" || !pool.AppendCertsFromPEM([]byte(draft.CertificatePEM)) {
		return rejected("Provide the trusted daemon certificate in PEM format.")
	}
	payload, _ := json.Marshal(map[string]any{"daemon_id": draft.DaemonID, "pairing_id": draft.PairingID, "phrase": draft.Phrase, "display_name": draft.DisplayName, "kind": domain.ActorKindDesktop, "capability": draft.Capability})
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, base.String()+"/api/v1/enroll", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13, RootCAs: pool}}}
	response, err := client.Do(request)
	if err != nil {
		return rejected("The trusted HTTPS daemon could not be reached.")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
		return rejected("Pairing was rejected. Request a new phrase and verify the daemon identity.")
	}
	var issued domain.IssuedCredential
	decoder := json.NewDecoder(io.LimitReader(response.Body, 64<<10))
	if decoder.Decode(&issued) != nil || issued.DaemonID != draft.DaemonID || issued.Token == "" || issued.ID == "" {
		return rejected("The daemon returned an invalid pairing response.")
	}
	if err := s.secrets.Save(issued.DaemonID, issued.ID, issued.Token); err != nil {
		return rejected("The credential could not be stored securely. Revoke it from the daemon.")
	}
	issued.Token = ""
	return Result{Action: "pair_remote", Outcome: "accepted", Message: "Remote client paired and stored in native credential storage.", CredentialID: issued.ID}
}

func rejected(message string) Result {
	return Result{Action: "pair_remote", Outcome: "rejected", Message: message}
}
