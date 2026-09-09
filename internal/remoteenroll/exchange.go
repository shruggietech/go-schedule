// Package remoteenroll performs one redirect-free, certificate-pinned pairing exchange.
package remoteenroll

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type Draft struct {
	Address, DaemonID, PairingID, Phrase, CertificatePEM, DisplayName string
	Kind                                                              domain.ActorKind
	Capability                                                        domain.Capability
}

func Exchange(ctx context.Context, draft Draft) (domain.IssuedCredential, server.ManifestResponse, error) {
	base, err := url.Parse(strings.TrimRight(draft.Address, "/"))
	if err != nil || base.Scheme != "https" || base.Host == "" || draft.DaemonID == "" || draft.PairingID == "" || draft.Phrase == "" || draft.DisplayName == "" || !draft.Capability.Valid() {
		return domain.IssuedCredential{}, server.ManifestResponse{}, errors.New("pairing input is invalid")
	}
	pool := x509.NewCertPool()
	if draft.CertificatePEM == "" || !pool.AppendCertsFromPEM([]byte(draft.CertificatePEM)) {
		return domain.IssuedCredential{}, server.ManifestResponse{}, errors.New("trusted certificate is invalid")
	}
	payload, _ := json.Marshal(map[string]any{"daemon_id": draft.DaemonID, "pairing_id": draft.PairingID, "phrase": draft.Phrase, "display_name": draft.DisplayName, "kind": draft.Kind, "capability": draft.Capability})
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, base.String()+"/api/v1/enroll", bytes.NewReader(payload))
	if err != nil {
		return domain.IssuedCredential{}, server.ManifestResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13, RootCAs: pool}}}
	response, err := httpClient.Do(request)
	if err != nil {
		return domain.IssuedCredential{}, server.ManifestResponse{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
		return domain.IssuedCredential{}, server.ManifestResponse{}, errors.New("pairing was rejected")
	}
	var issued domain.IssuedCredential
	if json.NewDecoder(io.LimitReader(response.Body, 64<<10)).Decode(&issued) != nil || issued.DaemonID != draft.DaemonID || issued.Token == "" || issued.ID == "" {
		return domain.IssuedCredential{}, server.ManifestResponse{}, errors.New("pairing response is invalid")
	}
	remote, err := client.NewRemote(base.String(), draft.CertificatePEM, issued.Token, draft.DaemonID)
	if err != nil {
		return domain.IssuedCredential{}, server.ManifestResponse{}, err
	}
	manifest, err := remote.VerifyIdentity(ctx)
	if err != nil {
		return domain.IssuedCredential{}, server.ManifestResponse{}, err
	}
	return issued, manifest, nil
}
