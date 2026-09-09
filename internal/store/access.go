package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	AuditRetentionCount = 10_000
	AuditRetentionAge   = 90 * 24 * time.Hour
)

var ErrBuiltinActor = errors.New("store: built-in actor is protected")

func (s *Store) initializeLocalActor() error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM actors WHERE builtin=1`).Scan(&count); err != nil {
		return fmt.Errorf("store: inspect built-in actor: %w", err)
	}
	if count == 1 {
		_, err := s.LocalActor()
		return err
	}
	if count != 0 {
		return fmt.Errorf("store: invalid built-in actor count %d", count)
	}
	now := fmtTime(time.Now())
	if _, err := s.db.Exec(`INSERT INTO actors(id,kind,display_name,capability,state,builtin,created_at,updated_at) VALUES(?,?,?,?,?,1,?,?)`, uuid.NewString(), domain.ActorKindLocalOS, "Local operating system", domain.CapabilityEnroll, domain.ActorStateActive, now, now); err != nil {
		return fmt.Errorf("store: initialize local actor: %w", err)
	}
	_, err := s.LocalActor()
	return err
}

func (s *Store) LocalActor() (domain.Actor, error) {
	actor, err := s.scanActor(s.db.QueryRow(`SELECT id,kind,display_name,capability,state,builtin,created_at,updated_at,expires_at FROM actors WHERE builtin=1`))
	if err != nil {
		return domain.Actor{}, err
	}
	if !actor.Builtin || actor.Kind != domain.ActorKindLocalOS || actor.Capability != domain.CapabilityEnroll || actor.State != domain.ActorStateActive || actor.ExpiresAt != nil {
		return domain.Actor{}, fmt.Errorf("store: validate local actor: %w", domain.ErrInvalidActor)
	}
	return actor, nil
}

func (s *Store) GetActor(id string) (domain.Actor, error) {
	actor, err := s.scanActor(s.db.QueryRow(`SELECT id,kind,display_name,capability,state,builtin,created_at,updated_at,expires_at FROM actors WHERE id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Actor{}, ErrNotFound
	}
	return actor, err
}

func (s *Store) ListActors() ([]domain.Actor, error) {
	rows, err := s.db.Query(`SELECT id,kind,display_name,capability,state,builtin,created_at,updated_at,expires_at FROM actors ORDER BY builtin DESC,created_at,id`)
	if err != nil {
		return nil, fmt.Errorf("store: list actors: %w", err)
	}
	defer rows.Close()
	var actors []domain.Actor
	for rows.Next() {
		actor, err := s.scanActor(rows)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, rows.Err()
}

func (s *Store) CreateActor(kind domain.ActorKind, name string, capability domain.Capability, expiresAt *time.Time) (domain.Actor, error) {
	name, err := domain.NormalizeActorDisplayName(name)
	if err != nil || !kind.Valid() || kind == domain.ActorKindLocalOS || !capability.Valid() {
		return domain.Actor{}, domain.ErrInvalidActor
	}
	if expiresAt != nil && !expiresAt.After(time.Now()) {
		return domain.Actor{}, domain.ErrInvalidActor
	}
	now := time.Now().UTC()
	id := uuid.NewString()
	if _, err := s.db.Exec(`INSERT INTO actors(id,kind,display_name,capability,state,builtin,created_at,updated_at,expires_at) VALUES(?,?,?,?,?,0,?,?,?)`, id, kind, name, capability, domain.ActorStateActive, fmtTime(now), fmtTime(now), fmtTimePtr(expiresAt)); err != nil {
		return domain.Actor{}, fmt.Errorf("store: create actor: %w", err)
	}
	return s.GetActor(id)
}

type ActorUpdate struct {
	DisplayName *string
	Capability  *domain.Capability
	State       *domain.ActorState
	ExpiresAt   **time.Time
}

func (s *Store) UpdateActor(id string, update ActorUpdate) (domain.Actor, error) {
	actor, err := s.GetActor(id)
	if err != nil {
		return domain.Actor{}, err
	}
	if actor.Builtin {
		return domain.Actor{}, ErrBuiltinActor
	}
	if update.DisplayName != nil {
		actor.DisplayName, err = domain.NormalizeActorDisplayName(*update.DisplayName)
		if err != nil {
			return domain.Actor{}, err
		}
	}
	if update.Capability != nil {
		if !update.Capability.Valid() {
			return domain.Actor{}, domain.ErrInvalidActor
		}
		actor.Capability = *update.Capability
	}
	if update.State != nil {
		if !update.State.Valid() || actor.State == domain.ActorStateRevoked && *update.State != domain.ActorStateRevoked {
			return domain.Actor{}, domain.ErrInvalidActor
		}
		actor.State = *update.State
	}
	if update.ExpiresAt != nil {
		actor.ExpiresAt = *update.ExpiresAt
	}
	actor.UpdatedAt = time.Now().UTC()
	if _, err := s.db.Exec(`UPDATE actors SET display_name=?,capability=?,state=?,updated_at=?,expires_at=? WHERE id=?`, actor.DisplayName, actor.Capability, actor.State, fmtTime(actor.UpdatedAt), fmtTimePtr(actor.ExpiresAt), id); err != nil {
		return domain.Actor{}, fmt.Errorf("store: update actor: %w", err)
	}
	return s.GetActor(id)
}

func (s *Store) RevokeActor(id string) (domain.Actor, error) {
	state := domain.ActorStateRevoked
	return s.UpdateActor(id, ActorUpdate{State: &state})
}

type rowScanner interface{ Scan(...any) error }

func (s *Store) scanActor(row rowScanner) (domain.Actor, error) {
	var actor domain.Actor
	var builtin int
	var created, updated string
	var expires sql.NullString
	if err := row.Scan(&actor.ID, &actor.Kind, &actor.DisplayName, &actor.Capability, &actor.State, &builtin, &created, &updated, &expires); err != nil {
		return domain.Actor{}, err
	}
	actor.Builtin = builtin == 1
	var err error
	if actor.CreatedAt, err = parseTime(created); err != nil {
		return domain.Actor{}, fmt.Errorf("store: parse actor creation time: %w", err)
	}
	if actor.UpdatedAt, err = parseTime(updated); err != nil {
		return domain.Actor{}, fmt.Errorf("store: parse actor update time: %w", err)
	}
	if actor.ExpiresAt, err = parseTimePtr(expires); err != nil {
		return domain.Actor{}, fmt.Errorf("store: parse actor expiration time: %w", err)
	}
	if actor.State == domain.ActorStateActive && actor.ExpiresAt != nil && !actor.ExpiresAt.After(time.Now()) {
		actor.State = domain.ActorStateExpired
	}
	if _, err := uuid.Parse(actor.ID); err != nil || !actor.Kind.Valid() || !actor.Capability.Valid() || !actor.State.Valid() {
		return domain.Actor{}, fmt.Errorf("store: validate actor: %w", domain.ErrInvalidActor)
	}
	name, err := domain.NormalizeActorDisplayName(actor.DisplayName)
	if err != nil || name != actor.DisplayName {
		return domain.Actor{}, fmt.Errorf("store: validate actor name: %w", domain.ErrInvalidActor)
	}
	return actor, nil
}

func (s *Store) BeginAudit(event domain.AuditEvent) error {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	if event.CorrelationID == "" {
		event.CorrelationID = uuid.NewString()
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	event.Result = domain.AuditResultUncertain
	return s.insertAudit(event)
}

func (s *Store) RecordDeniedAudit(event domain.AuditEvent) error {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	if event.CorrelationID == "" {
		event.CorrelationID = uuid.NewString()
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	event.Result = domain.AuditResultDenied
	event.CompletedAt = &event.OccurredAt
	return s.insertAudit(event)
}

func (s *Store) insertAudit(event domain.AuditEvent) error {
	if err := validateAuditEvent(event); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin audit record: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec(`INSERT INTO audit_events(id,actor_id,daemon_id,operation,target_kind,target_id,result,correlation_id,occurred_at,completed_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, event.ID, nullableString(event.ActorID), event.DaemonID, event.Operation, event.TargetKind, event.TargetID, event.Result, event.CorrelationID, fmtTime(event.OccurredAt), fmtTimePtr(event.CompletedAt)); err != nil {
		return fmt.Errorf("store: insert audit record: %w", err)
	}
	cutoff := fmtTime(event.OccurredAt.Add(-AuditRetentionAge))
	if _, err := tx.Exec(`DELETE FROM audit_events WHERE occurred_at < ?`, cutoff); err != nil {
		return fmt.Errorf("store: prune audit age: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM audit_events WHERE id IN (SELECT id FROM audit_events ORDER BY occurred_at DESC,id DESC LIMIT -1 OFFSET ?)`, AuditRetentionCount); err != nil {
		return fmt.Errorf("store: prune audit count: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit audit record: %w", err)
	}
	return nil
}

func (s *Store) CompleteAudit(id string, result domain.AuditResult) error {
	if result != domain.AuditResultSucceeded && result != domain.AuditResultFailed {
		return fmt.Errorf("store: complete audit: invalid result")
	}
	now := fmtTime(time.Now())
	changed, err := s.db.Exec(`UPDATE audit_events SET result=?,completed_at=? WHERE id=? AND result='uncertain'`, result, now, id)
	if err != nil {
		return fmt.Errorf("store: complete audit: %w", err)
	}
	count, err := changed.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("store: complete audit: event not uncertain")
	}
	return nil
}

func (s *Store) ListAudit(query domain.AuditQuery) ([]domain.AuditEvent, error) {
	if query.Limit == 0 {
		query.Limit = 100
	}
	if query.Limit < 1 || query.Limit > 1000 || (query.Result != "" && !query.Result.Valid()) || (query.Since != nil && query.Until != nil && query.Since.After(*query.Until)) {
		return nil, fmt.Errorf("store: invalid audit query")
	}
	clauses := []string{"1=1"}
	args := []any{}
	if query.ActorID != "" {
		clauses = append(clauses, "actor_id=?")
		args = append(args, query.ActorID)
	}
	if query.Operation != "" {
		clauses = append(clauses, "operation=?")
		args = append(args, query.Operation)
	}
	if query.Result != "" {
		clauses = append(clauses, "result=?")
		args = append(args, query.Result)
	}
	if query.Since != nil {
		clauses = append(clauses, "occurred_at>=?")
		args = append(args, fmtTime(*query.Since))
	}
	if query.Until != nil {
		clauses = append(clauses, "occurred_at<=?")
		args = append(args, fmtTime(*query.Until))
	}
	args = append(args, query.Limit)
	rows, err := s.db.Query(`SELECT id,COALESCE(actor_id,''),daemon_id,operation,target_kind,target_id,result,correlation_id,occurred_at,completed_at FROM audit_events WHERE `+strings.Join(clauses, " AND ")+` ORDER BY occurred_at,id LIMIT ?`, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list audit: %w", err)
	}
	defer rows.Close()
	var events []domain.AuditEvent
	for rows.Next() {
		var event domain.AuditEvent
		var occurred string
		var completed sql.NullString
		if err := rows.Scan(&event.ID, &event.ActorID, &event.DaemonID, &event.Operation, &event.TargetKind, &event.TargetID, &event.Result, &event.CorrelationID, &occurred, &completed); err != nil {
			return nil, err
		}
		event.OccurredAt, err = parseTime(occurred)
		if err != nil {
			return nil, err
		}
		event.CompletedAt, err = parseTimePtr(completed)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func validateAuditEvent(event domain.AuditEvent) error {
	if _, err := uuid.Parse(event.ID); err != nil {
		return fmt.Errorf("store: invalid audit event")
	}
	if _, err := uuid.Parse(event.DaemonID); err != nil {
		return fmt.Errorf("store: invalid audit daemon")
	}
	if _, err := uuid.Parse(event.CorrelationID); err != nil {
		return fmt.Errorf("store: invalid audit correlation")
	}
	if event.ActorID != "" {
		if _, err := uuid.Parse(event.ActorID); err != nil {
			return fmt.Errorf("store: invalid audit actor")
		}
	}
	if event.Operation == "" || event.TargetKind == "" || !event.Result.Valid() || event.OccurredAt.IsZero() {
		return fmt.Errorf("store: invalid audit event")
	}
	return nil
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
