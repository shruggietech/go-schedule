package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestCredentialRotationResponseIncludesActualIdentitiesAndOneTimeToken(t *testing.T) {
	s := newTestServer(t)
	pairing, err := s.enrollment.Create("Rotation API client", domain.ActorKindCLI, domain.CapabilityManage)
	if err != nil {
		t.Fatal(err)
	}
	issued, err := s.enrollment.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
	if err != nil {
		t.Fatal(err)
	}
	response := doJSON(t, s, http.MethodPost, "/v1/access/credentials/"+issued.ID+"/rotate", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("rotation status = %d", response.Code)
	}
	var rotated domain.IssuedCredential
	if err := json.Unmarshal(response.Body.Bytes(), &rotated); err != nil {
		t.Fatal(err)
	}
	if rotated.Actor.ID != issued.Actor.ID || rotated.Actor.ID != rotated.ActorID || rotated.Actor.DisplayName != issued.Actor.DisplayName || rotated.Actor.Kind != issued.Actor.Kind || rotated.Actor.Capability != issued.Actor.Capability || rotated.Actor.CreatedAt.IsZero() || rotated.Actor.UpdatedAt.IsZero() || rotated.DaemonID != issued.DaemonID {
		t.Fatal("rotation response omitted actual identity context")
	}
	if rotated.Token == "" || rotated.Token == issued.Token || rotated.LastUsedAt != nil {
		t.Fatal("rotation response omitted replacement token or recorded authentication prematurely")
	}
	if strings.Contains(response.Body.String(), "0001-01-01") {
		t.Fatal("rotation response contains zero timestamps")
	}
	credentials := doJSON(t, s, http.MethodGet, "/v1/access/credentials", nil)
	if credentials.Code != http.StatusOK {
		t.Fatalf("credential list status = %d", credentials.Code)
	}
	var listed map[string]any
	if err := json.Unmarshal(credentials.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(credentials.Body.String(), rotated.Token) || strings.Contains(credentials.Body.String(), issued.Token) || strings.Contains(credentials.Body.String(), "\"token\"") {
		t.Fatal("credential list disclosed a one-time token")
	}
	if _, _, err := s.enrollment.Authenticate(issued.Token); err == nil {
		t.Fatal("old token remained usable")
	}
	if _, _, err := s.enrollment.Authenticate(rotated.Token); err != nil {
		t.Fatal("replacement token was not usable")
	}
}
