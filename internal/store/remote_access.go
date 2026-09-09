package store

import (
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

var (
	ErrEnrollmentRejected = errors.New("store: enrollment rejected")
	ErrCredentialRejected = errors.New("store: credential rejected")
)

func (s *Store) CreatePairing(session domain.PairingSession, salt, verifier []byte) error {
	if len(salt) != 16 || len(verifier) != 32 || session.State != domain.PairingActive || session.AttemptsRemaining != 5 || !session.ExpiresAt.After(session.CreatedAt) {
		return domain.ErrInvalidActor
	}
	if _, err := domain.NormalizeActorDisplayName(session.DisplayName); err != nil || !session.Kind.Valid() || session.Kind == domain.ActorKindLocalOS || !session.Capability.Valid() {
		return domain.ErrInvalidActor
	}
	_, err := s.db.Exec(`INSERT INTO pairing_sessions(id,display_name,kind,capability,salt,verifier,attempts_remaining,state,created_at,expires_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, session.ID, session.DisplayName, session.Kind, session.Capability, salt, verifier, session.AttemptsRemaining, session.State, fmtTime(session.CreatedAt), fmtTime(session.ExpiresAt))
	if err != nil {
		return fmt.Errorf("store: create pairing: %w", err)
	}
	return nil
}

func (s *Store) PairingSalt(id string) ([]byte, error) {
	var salt []byte
	if err := s.db.QueryRow(`SELECT salt FROM pairing_sessions WHERE id=?`, id).Scan(&salt); err != nil {
		return nil, ErrEnrollmentRejected
	}
	return append([]byte(nil), salt...), nil
}

func (s *Store) PairingForExchange(id string) (domain.PairingSession, []byte, error) {
	session, err := s.getPairing(id)
	if err != nil {
		return domain.PairingSession{}, nil, ErrEnrollmentRejected
	}
	salt, err := s.PairingSalt(id)
	if err != nil {
		return domain.PairingSession{}, nil, ErrEnrollmentRejected
	}
	return session, salt, nil
}

func (s *Store) ListPairings() ([]domain.PairingSession, error) {
	rows, err := s.db.Query(`SELECT id,display_name,kind,capability,attempts_remaining,state,created_at,expires_at,completed_at FROM pairing_sessions ORDER BY created_at DESC,id DESC`)
	if err != nil {
		return nil, fmt.Errorf("store: list pairings: %w", err)
	}
	defer rows.Close()
	var result []domain.PairingSession
	for rows.Next() {
		item, err := scanPairing(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) CancelPairing(id string) (domain.PairingSession, error) {
	now := time.Now().UTC()
	result, err := s.db.Exec(`UPDATE pairing_sessions SET state='cancelled',completed_at=? WHERE id=? AND state='active'`, fmtTime(now), id)
	if err != nil {
		return domain.PairingSession{}, fmt.Errorf("store: cancel pairing: %w", err)
	}
	if count, _ := result.RowsAffected(); count != 1 {
		return domain.PairingSession{}, ErrNotFound
	}
	return s.getPairing(id)
}

func (s *Store) getPairing(id string) (domain.PairingSession, error) {
	item, err := scanPairing(s.db.QueryRow(`SELECT id,display_name,kind,capability,attempts_remaining,state,created_at,expires_at,completed_at FROM pairing_sessions WHERE id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PairingSession{}, ErrNotFound
	}
	return item, err
}

func scanPairing(row rowScanner) (domain.PairingSession, error) {
	var item domain.PairingSession
	var created, expires string
	var completed sql.NullString
	if err := row.Scan(&item.ID, &item.DisplayName, &item.Kind, &item.Capability, &item.AttemptsRemaining, &item.State, &created, &expires, &completed); err != nil {
		return item, err
	}
	var err error
	if item.CreatedAt, err = parseTime(created); err != nil {
		return item, err
	}
	if item.ExpiresAt, err = parseTime(expires); err != nil {
		return item, err
	}
	if item.CompletedAt, err = parseTimePtr(completed); err != nil {
		return item, err
	}
	if item.State == domain.PairingActive && !item.ExpiresAt.After(time.Now()) {
		item.State = domain.PairingExpired
	}
	return item, nil
}

// ExchangePairing atomically consumes one pairing verifier and creates its actor and credential.
func (s *Store) ExchangePairing(id string, candidateVerifier []byte, actor domain.Actor, credential domain.ClientCredential, digest []byte, now time.Time) (domain.Actor, domain.ClientCredential, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return actor, credential, err
	}
	defer func() { _ = tx.Rollback() }()
	var name string
	var kind domain.ActorKind
	var capability domain.Capability
	var verifier []byte
	var attempts int
	var state domain.PairingState
	var expires string
	if err := tx.QueryRow(`SELECT display_name,kind,capability,verifier,attempts_remaining,state,expires_at FROM pairing_sessions WHERE id=?`, id).Scan(&name, &kind, &capability, &verifier, &attempts, &state, &expires); err != nil {
		return actor, credential, ErrEnrollmentRejected
	}
	expiresAt, err := parseTime(expires)
	if err != nil || state != domain.PairingActive || attempts <= 0 || !expiresAt.After(now) {
		if state == domain.PairingActive && !expiresAt.After(now) {
			if _, updateErr := tx.Exec(`UPDATE pairing_sessions SET state='expired',completed_at=? WHERE id=? AND state='active'`, fmtTime(now), id); updateErr != nil {
				return actor, credential, updateErr
			}
			if commitErr := tx.Commit(); commitErr != nil {
				return actor, credential, commitErr
			}
		}
		return actor, credential, ErrEnrollmentRejected
	}
	if len(candidateVerifier) != len(verifier) || subtle.ConstantTimeCompare(candidateVerifier, verifier) != 1 {
		attempts--
		newState := domain.PairingActive
		var completed any
		if attempts == 0 {
			newState = domain.PairingExhausted
			completed = fmtTime(now)
		}
		if _, updateErr := tx.Exec(`UPDATE pairing_sessions SET attempts_remaining=?,state=?,completed_at=? WHERE id=? AND state='active'`, attempts, newState, completed, id); updateErr != nil {
			return actor, credential, updateErr
		}
		if commitErr := tx.Commit(); commitErr != nil {
			return actor, credential, commitErr
		}
		return actor, credential, ErrEnrollmentRejected
	}
	actor.DisplayName, actor.Kind, actor.Capability = name, kind, capability
	if _, err := tx.Exec(`INSERT INTO actors(id,kind,display_name,capability,state,builtin,created_at,updated_at,expires_at) VALUES(?,?,?,?,?,0,?,?,?)`, actor.ID, actor.Kind, actor.DisplayName, actor.Capability, actor.State, fmtTime(actor.CreatedAt), fmtTime(actor.UpdatedAt), fmtTimePtr(actor.ExpiresAt)); err != nil {
		return actor, credential, fmt.Errorf("store: enroll actor: %w", err)
	}
	credential.ActorID = actor.ID
	if _, err := tx.Exec(`INSERT INTO client_credentials(id,actor_id,digest,fingerprint,state,created_at,updated_at,expires_at) VALUES(?,?,?,?,?,?,?,?)`, credential.ID, credential.ActorID, digest, credential.Fingerprint, credential.State, fmtTime(credential.CreatedAt), fmtTime(credential.UpdatedAt), fmtTimePtr(credential.ExpiresAt)); err != nil {
		return actor, credential, fmt.Errorf("store: enroll credential: %w", err)
	}
	result, err := tx.Exec(`UPDATE pairing_sessions SET state='consumed',completed_at=? WHERE id=? AND state='active'`, fmtTime(now), id)
	if err != nil {
		return actor, credential, err
	}
	if count, _ := result.RowsAffected(); count != 1 {
		return actor, credential, ErrEnrollmentRejected
	}
	if err := tx.Commit(); err != nil {
		return actor, credential, err
	}
	return actor, credential, nil
}

func (s *Store) ListCredentials() ([]domain.ClientCredential, error) {
	rows, err := s.db.Query(`SELECT id,actor_id,fingerprint,state,created_at,updated_at,last_used_at,expires_at,revoked_at FROM client_credentials ORDER BY created_at,id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.ClientCredential
	for rows.Next() {
		item, err := scanCredential(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func scanCredential(row rowScanner) (domain.ClientCredential, error) {
	var item domain.ClientCredential
	var created, updated string
	var lastUsed, expires, revoked sql.NullString
	if err := row.Scan(&item.ID, &item.ActorID, &item.Fingerprint, &item.State, &created, &updated, &lastUsed, &expires, &revoked); err != nil {
		return item, err
	}
	var err error
	if item.CreatedAt, err = parseTime(created); err != nil {
		return item, err
	}
	if item.UpdatedAt, err = parseTime(updated); err != nil {
		return item, err
	}
	if item.LastUsedAt, err = parseTimePtr(lastUsed); err != nil {
		return item, err
	}
	if item.ExpiresAt, err = parseTimePtr(expires); err != nil {
		return item, err
	}
	if item.RevokedAt, err = parseTimePtr(revoked); err != nil {
		return item, err
	}
	return item, nil
}

func (s *Store) AuthenticateCredential(digest []byte, now time.Time) (domain.ClientCredential, domain.Actor, error) {
	var stored []byte
	row := s.db.QueryRow(`SELECT c.id,c.actor_id,c.fingerprint,c.state,c.created_at,c.updated_at,c.last_used_at,c.expires_at,c.revoked_at,c.digest,a.id,a.kind,a.display_name,a.capability,a.state,a.builtin,a.created_at,a.updated_at,a.expires_at FROM client_credentials c JOIN actors a ON a.id=c.actor_id WHERE c.digest=?`, digest)
	var credential domain.ClientCredential
	var actor domain.Actor
	var cCreated, cUpdated, aCreated, aUpdated string
	var cLast, cExpires, cRevoked, aExpires sql.NullString
	var builtin int
	if err := row.Scan(&credential.ID, &credential.ActorID, &credential.Fingerprint, &credential.State, &cCreated, &cUpdated, &cLast, &cExpires, &cRevoked, &stored, &actor.ID, &actor.Kind, &actor.DisplayName, &actor.Capability, &actor.State, &builtin, &aCreated, &aUpdated, &aExpires); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if len(stored) != len(digest) || subtle.ConstantTimeCompare(stored, digest) != 1 || credential.State != domain.CredentialActive {
		return credential, actor, ErrCredentialRejected
	}
	var err error
	if credential.CreatedAt, err = parseTime(cCreated); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if credential.UpdatedAt, err = parseTime(cUpdated); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if credential.LastUsedAt, err = parseTimePtr(cLast); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if credential.ExpiresAt, err = parseTimePtr(cExpires); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if credential.RevokedAt, err = parseTimePtr(cRevoked); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	actor.Builtin = builtin == 1
	if actor.CreatedAt, err = parseTime(aCreated); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if actor.UpdatedAt, err = parseTime(aUpdated); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if actor.ExpiresAt, err = parseTimePtr(aExpires); err != nil {
		return credential, actor, ErrCredentialRejected
	}
	if (credential.ExpiresAt != nil && !credential.ExpiresAt.After(now)) || !actor.ActiveAt(now) {
		return credential, actor, ErrCredentialRejected
	}
	_, _ = s.db.Exec(`UPDATE client_credentials SET last_used_at=? WHERE id=?`, fmtTime(now), credential.ID)
	return credential, actor, nil
}

func (s *Store) RotateCredential(id string, digest []byte, fingerprint string, now time.Time) (domain.ClientCredential, error) {
	result, err := s.db.Exec(`UPDATE client_credentials SET digest=?,fingerprint=?,updated_at=?,last_used_at=NULL WHERE id=? AND state='active'`, digest, fingerprint, fmtTime(now), id)
	if err != nil {
		return domain.ClientCredential{}, fmt.Errorf("store: rotate credential: %w", err)
	}
	if count, _ := result.RowsAffected(); count != 1 {
		return domain.ClientCredential{}, ErrNotFound
	}
	return s.getCredential(id)
}

func (s *Store) RevokeCredential(id string, now time.Time) (domain.ClientCredential, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return domain.ClientCredential{}, err
	}
	defer func() { _ = tx.Rollback() }()
	var actorID string
	if err := tx.QueryRow(`SELECT actor_id FROM client_credentials WHERE id=? AND state='active'`, id).Scan(&actorID); err != nil {
		return domain.ClientCredential{}, ErrNotFound
	}
	if _, err := tx.Exec(`UPDATE client_credentials SET state='revoked',updated_at=?,revoked_at=? WHERE id=?`, fmtTime(now), fmtTime(now), id); err != nil {
		return domain.ClientCredential{}, err
	}
	if _, err := tx.Exec(`UPDATE actors SET state='revoked',updated_at=? WHERE id=? AND builtin=0`, fmtTime(now), actorID); err != nil {
		return domain.ClientCredential{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.ClientCredential{}, err
	}
	return s.getCredential(id)
}

func (s *Store) getCredential(id string) (domain.ClientCredential, error) {
	item, err := scanCredential(s.db.QueryRow(`SELECT id,actor_id,fingerprint,state,created_at,updated_at,last_used_at,expires_at,revoked_at FROM client_credentials WHERE id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return item, ErrNotFound
	}
	return item, err
}
