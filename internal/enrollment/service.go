package enrollment

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

const (
	PhraseWords     = 10
	PairingLifetime = 10 * time.Minute
	PairingAttempts = 5
)

var ErrRejected = errors.New("enrollment rejected")

var words = [...]string{
	"amber", "anchor", "apple", "april", "arrow", "atlas", "baker", "beacon",
	"birch", "bloom", "bravo", "brook", "cabin", "cactus", "candle", "cedar",
	"charm", "cobalt", "comet", "coral", "delta", "dove", "ember", "falcon",
	"fern", "fjord", "flame", "flora", "frost", "glade", "globe", "harbor",
	"hazel", "heron", "indigo", "iris", "island", "jade", "juniper", "lagoon",
	"lark", "lemon", "lilac", "lotus", "maple", "mesa", "mint", "nova",
	"oasis", "olive", "onyx", "opal", "orbit", "pearl", "pine", "quartz",
	"raven", "river", "sable", "solar", "spruce", "stone", "tulip", "willow",
}

type persistence interface {
	CreatePairing(domain.PairingSession, []byte, []byte) error
	ListPairings() ([]domain.PairingSession, error)
	CancelPairing(string) (domain.PairingSession, error)
	ExchangePairing(string, []byte, domain.Actor, domain.ClientCredential, []byte, time.Time) (domain.Actor, domain.ClientCredential, error)
	ListCredentials() ([]domain.ClientCredential, error)
	RotateCredential(string, []byte, string, time.Time) (domain.ClientCredential, error)
	RevokeCredential(string, time.Time) (domain.ClientCredential, error)
	AuthenticateCredential([]byte, time.Time) (domain.ClientCredential, domain.Actor, error)
	DaemonIdentity() (domain.DaemonIdentity, error)
}

type Service struct {
	store  persistence
	random io.Reader
	now    func() time.Time
	memory uint32
}

func New(st persistence) *Service {
	return &Service{store: st, random: rand.Reader, now: func() time.Time { return time.Now().UTC() }, memory: 64 * 1024}
}

func (s *Service) Create(name string, kind domain.ActorKind, capability domain.Capability) (domain.PairingSecret, error) {
	name, err := domain.NormalizeActorDisplayName(name)
	if err != nil || !kind.Valid() || kind == domain.ActorKindLocalOS || !capability.Valid() {
		return domain.PairingSecret{}, domain.ErrInvalidActor
	}
	phrase, err := s.phrase()
	if err != nil {
		return domain.PairingSecret{}, err
	}
	salt := make([]byte, 16)
	if _, err := io.ReadFull(s.random, salt); err != nil {
		return domain.PairingSecret{}, err
	}
	now := s.now()
	session := domain.PairingSession{ID: uuid.NewString(), DisplayName: name, Kind: kind, Capability: capability, AttemptsRemaining: PairingAttempts, State: domain.PairingActive, CreatedAt: now, ExpiresAt: now.Add(PairingLifetime)}
	if err := s.store.CreatePairing(session, salt, s.verify(phrase, salt)); err != nil {
		return domain.PairingSecret{}, err
	}
	identity, err := s.store.DaemonIdentity()
	if err != nil {
		return domain.PairingSecret{}, err
	}
	return domain.PairingSecret{PairingSession: session, Phrase: phrase, DaemonID: identity.InstallationID}, nil
}

func (s *Service) List() ([]domain.PairingSession, error)          { return s.store.ListPairings() }
func (s *Service) Cancel(id string) (domain.PairingSession, error) { return s.store.CancelPairing(id) }
func (s *Service) Credentials() ([]domain.ClientCredential, error) { return s.store.ListCredentials() }

func (s *Service) Exchange(id, phrase, expectedDaemonID, displayName string, kind domain.ActorKind, capability domain.Capability) (domain.IssuedCredential, error) {
	identity, err := s.store.DaemonIdentity()
	displayName, nameErr := domain.NormalizeActorDisplayName(displayName)
	if err != nil || nameErr != nil || expectedDaemonID != identity.InstallationID || !validPhrase(phrase) {
		return domain.IssuedCredential{}, ErrRejected
	}
	session, salt, err := s.pairingForExchange(id)
	if err != nil || session.DisplayName != displayName || session.Kind != kind || session.Capability != capability {
		return domain.IssuedCredential{}, ErrRejected
	}
	now := s.now()
	tokenBytes := make([]byte, 32)
	if _, err := io.ReadFull(s.random, tokenBytes); err != nil {
		return domain.IssuedCredential{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	digest := sha256.Sum256(tokenBytes)
	fingerprint := base64.RawURLEncoding.EncodeToString(digest[:])[:12]
	actor := domain.Actor{ID: uuid.NewString(), State: domain.ActorStateActive, CreatedAt: now, UpdatedAt: now}
	credential := domain.ClientCredential{ID: uuid.NewString(), Fingerprint: fingerprint, State: domain.CredentialActive, CreatedAt: now, UpdatedAt: now}
	actor, credential, err = s.store.ExchangePairing(id, s.verify(phrase, salt), actor, credential, digest[:], now)
	if err != nil {
		if errors.Is(err, store.ErrEnrollmentRejected) {
			return domain.IssuedCredential{}, ErrRejected
		}
		return domain.IssuedCredential{}, err
	}
	return domain.IssuedCredential{ClientCredential: credential, Actor: actor, Token: token, DaemonID: identity.InstallationID}, nil
}

type pairingReader interface {
	PairingForExchange(string) (domain.PairingSession, []byte, error)
}

func (s *Service) pairingForExchange(id string) (domain.PairingSession, []byte, error) {
	reader, ok := s.store.(pairingReader)
	if !ok {
		return domain.PairingSession{}, nil, ErrRejected
	}
	return reader.PairingForExchange(id)
}

func (s *Service) Authenticate(token string) (domain.ClientCredential, domain.Actor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != 32 {
		return domain.ClientCredential{}, domain.Actor{}, ErrRejected
	}
	digest := sha256.Sum256(raw)
	credential, actor, err := s.store.AuthenticateCredential(digest[:], s.now())
	if err != nil {
		return domain.ClientCredential{}, domain.Actor{}, ErrRejected
	}
	return credential, actor, nil
}

func (s *Service) Rotate(id string) (domain.IssuedCredential, error) {
	raw := make([]byte, 32)
	if _, err := io.ReadFull(s.random, raw); err != nil {
		return domain.IssuedCredential{}, err
	}
	digest := sha256.Sum256(raw)
	fingerprint := base64.RawURLEncoding.EncodeToString(digest[:])[:12]
	credential, err := s.store.RotateCredential(id, digest[:], fingerprint, s.now())
	if err != nil {
		return domain.IssuedCredential{}, err
	}
	return domain.IssuedCredential{ClientCredential: credential, Token: base64.RawURLEncoding.EncodeToString(raw)}, nil
}

func (s *Service) Revoke(id string) (domain.ClientCredential, error) {
	return s.store.RevokeCredential(id, s.now())
}

func (s *Service) phrase() (string, error) {
	raw := make([]byte, PhraseWords)
	if _, err := io.ReadFull(s.random, raw); err != nil {
		return "", err
	}
	result := make([]string, len(raw))
	for i, value := range raw {
		result[i] = words[int(value)&63]
	}
	return strings.Join(result, "-"), nil
}

func validPhrase(value string) bool {
	parts := strings.Split(value, "-")
	if len(parts) != PhraseWords {
		return false
	}
	for _, part := range parts {
		found := false
		for _, word := range words {
			if part == word {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *Service) verify(phrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(phrase), salt, 1, s.memory, 4, 32)
}
