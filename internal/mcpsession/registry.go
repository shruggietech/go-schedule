// Package mcpsession owns runtime-only local MCP identities.
package mcpsession

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

const Lifetime = 24 * time.Hour

var ErrInvalidSession = errors.New("invalid MCP session")

type actorStore interface {
	CreateActor(domain.ActorKind, string, domain.Capability, *time.Time) (domain.Actor, error)
	RevokeActor(string) (domain.Actor, error)
}

// Session is the non-secret runtime session projection.
type Session struct {
	ID         string            `json:"id"`
	ActorID    string            `json:"actor_id"`
	Capability domain.Capability `json:"capability"`
	CreatedAt  time.Time         `json:"created_at"`
	ExpiresAt  time.Time         `json:"expires_at"`
}

type entry struct {
	session Session
	digest  [sha256.Size]byte
}

// Registry resolves random session secrets to persisted MCP actors. Entries
// remain until process exit so requests made with a revoked secret can still
// be attributed to the revoked actor and denied by the shared authorization
// path.
type Registry struct {
	mu       sync.RWMutex
	store    actorStore
	now      func() time.Time
	sessions map[string]entry
}

func New(st actorStore) *Registry {
	return &Registry{store: st, now: time.Now, sessions: make(map[string]entry)}
}

// Create returns a session and its secret exactly once.
func (r *Registry) Create(displayName string, capability domain.Capability) (Session, string, error) {
	if capability != domain.CapabilityObserve && capability != domain.CapabilityOperate && capability != domain.CapabilityManage {
		return Session{}, "", ErrInvalidSession
	}
	created := r.now().UTC()
	expires := created.Add(Lifetime)
	actor, err := r.store.CreateActor(domain.ActorKindMCP, displayName, capability, &expires)
	if err != nil {
		return Session{}, "", err
	}
	secret, digest, err := newSecret()
	if err != nil {
		_, _ = r.store.RevokeActor(actor.ID)
		return Session{}, "", err
	}
	session := Session{ID: uuid.NewString(), ActorID: actor.ID, Capability: capability, CreatedAt: created, ExpiresAt: expires}
	r.mu.Lock()
	r.sessions[session.ID] = entry{session: session, digest: digest}
	r.mu.Unlock()
	return session, secret, nil
}

// Resolve returns the actor even after revocation. Authorization reads the
// actor state from SQLite on every request and therefore records an attributed
// denial rather than treating a revoked client as anonymous.
func (r *Registry) Resolve(secret string) (string, error) {
	if secret == "" {
		return "", ErrInvalidSession
	}
	digest := sha256.Sum256([]byte(secret))
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, candidate := range r.sessions {
		if subtle.ConstantTimeCompare(candidate.digest[:], digest[:]) == 1 {
			return candidate.session.ActorID, nil
		}
	}
	return "", ErrInvalidSession
}

func (r *Registry) Revoke(id string) (Session, error) {
	r.mu.RLock()
	candidate, ok := r.sessions[id]
	r.mu.RUnlock()
	if !ok {
		return Session{}, store.ErrNotFound
	}
	if _, err := r.store.RevokeActor(candidate.session.ActorID); err != nil && !errors.Is(err, store.ErrBuiltinActor) {
		return Session{}, err
	}
	return candidate.session, nil
}

// List returns the secret-free runtime sessions known to this daemon process.
func (r *Registry) List() []Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sessions := make([]Session, 0, len(r.sessions))
	for _, candidate := range r.sessions {
		sessions = append(sessions, candidate.session)
	}
	return sessions
}

func newSecret() (string, [sha256.Size]byte, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", [sha256.Size]byte{}, err
	}
	secret := base64.RawURLEncoding.EncodeToString(raw[:])
	return secret, sha256.Sum256([]byte(secret)), nil
}
