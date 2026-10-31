// Package clientobserve performs bounded, identified operational reads without a UI dependency.
package clientobserve

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

var (
	ErrPartial = errors.New("one or more schedulers could not be observed; see their individual outcomes")
	ErrInvalid = errors.New("invalid observation request")
)

// Registration identifies a source without retaining a credential or trust bundle.
type Registration struct {
	Key              string `json:"key"`
	Label            string `json:"label"`
	Endpoint         string `json:"endpoint"`
	ExpectedDaemonID string `json:"expected_daemon_id,omitempty"`
	ObservedDaemonID string `json:"observed_daemon_id,omitempty"`
}

// Failure is stable, credential-free guidance for one unsuccessful source.
type Failure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Reader is satisfied by the shared API client and injectable fixture clients.
type Reader interface {
	Health(context.Context) (server.HealthResponse, error)
	VerifyIdentity(context.Context) (server.ManifestResponse, error)
	VerifyAccess(context.Context) (domain.Capability, error)
	SystemSummary(context.Context) (domain.SystemSummary, error)
	Search(context.Context, string, []domain.SearchKind, int) (domain.DaemonSearch, error)
}

type Target struct {
	Registration Registration
	Client       Reader
	SetupFailure *Failure
	// Resolve defers credential access and transport creation until this target's bounded worker runs.
	Resolve func(context.Context) (Reader, error)
}

type Observation struct {
	Registration Registration          `json:"registration"`
	ObservedAt   time.Time             `json:"observed_at"`
	Summary      *domain.SystemSummary `json:"summary,omitempty"`
	Search       *domain.DaemonSearch  `json:"search,omitempty"`
	Failure      *Failure              `json:"failure,omitempty"`
}

type Snapshot struct {
	Operation    string        `json:"operation"`
	CompletedAt  time.Time     `json:"completed_at"`
	Observations []Observation `json:"observations"`
}

type SearchRequest struct {
	Query string
	Kinds []domain.SearchKind
	Limit int
}

// Observer has a maximum of four concurrent targets and an independent deadline per target.
// Configure its clock and timeout before starting an observation.
type Observer struct {
	Now     func() time.Time
	Timeout time.Duration
}

func New() *Observer { return &Observer{Now: time.Now, Timeout: 5 * time.Second} }

func (o *Observer) Summary(ctx context.Context, targets []Target) (Snapshot, error) {
	return o.observe(ctx, targets, "summary", func(ctx context.Context, reader Reader, result *Observation) error {
		summary, err := reader.SystemSummary(ctx)
		if err != nil {
			return err
		}
		if summary.Schema != domain.SystemSummarySchema || summary.ObservedAt.IsZero() {
			return incompatible()
		}
		result.Summary = &summary
		result.ObservedAt = summary.ObservedAt.UTC()
		return nil
	})
}

func (o *Observer) Search(ctx context.Context, targets []Target, request SearchRequest) (Snapshot, error) {
	request.Query = strings.TrimSpace(request.Query)
	if request.Query == "" || utf8.RuneCountInString(request.Query) > 200 || request.Limit < 0 || request.Limit > 50 {
		return Snapshot{}, ErrInvalid
	}
	for _, kind := range request.Kinds {
		if !kind.Valid() {
			return Snapshot{}, ErrInvalid
		}
	}
	if request.Limit == 0 {
		request.Limit = 50
	}
	return o.observe(ctx, targets, "search", func(ctx context.Context, reader Reader, result *Observation) error {
		search, err := reader.Search(ctx, request.Query, request.Kinds, request.Limit)
		if err != nil {
			return err
		}
		if search.Schema != domain.DaemonSearchSchema || search.ObservedAt.IsZero() || search.Query != request.Query || len(search.Results) > request.Limit {
			return incompatible()
		}
		for _, match := range search.Results {
			if !match.Kind.Valid() || strings.TrimSpace(match.ObjectID) == "" || strings.TrimSpace(match.Name) == "" {
				return incompatible()
			}
		}
		if search.Results == nil {
			search.Results = []domain.SearchMatch{}
		}
		result.Search = &search
		result.ObservedAt = search.ObservedAt.UTC()
		return nil
	})
}

func (o *Observer) observe(ctx context.Context, targets []Target, operation string, read func(context.Context, Reader, *Observation) error) (Snapshot, error) {
	now := o.Now
	if now == nil {
		now = time.Now
	}
	timeout := o.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	results := make([]Observation, len(targets))
	jobs := make(chan int)
	var workers sync.WaitGroup
	for range min(4, len(targets)) {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for index := range jobs {
				target := targets[index]
				result := Observation{Registration: target.Registration, ObservedAt: now().UTC()}
				if target.SetupFailure != nil {
					failure := *target.SetupFailure
					result.Failure = &failure
					results[index] = result
					continue
				}
				callCtx, cancel := context.WithTimeout(ctx, timeout)
				err := callCtx.Err()
				reader := target.Client
				if err == nil && target.Resolve != nil {
					reader, err = target.Resolve(callCtx)
				}
				if err == nil {
					err = callCtx.Err()
				}
				if err == nil && reader == nil {
					err = errors.New("missing target client")
				}
				if err == nil {
					err = verify(callCtx, reader, &result.Registration)
				}
				if err == nil {
					err = callCtx.Err()
				}
				if err == nil {
					err = read(callCtx, reader, &result)
				}
				if err != nil {
					result.Failure = Classify(err)
				}
				cancel()
				results[index] = result
			}
		}()
	}
	for index := range targets {
		jobs <- index
	}
	close(jobs)
	workers.Wait()
	sort.SliceStable(results, func(i, j int) bool { return results[i].Registration.Key < results[j].Registration.Key })
	snapshot := Snapshot{Operation: operation, CompletedAt: now().UTC(), Observations: results}
	if err := ctx.Err(); err != nil {
		return snapshot, err
	}
	for _, result := range results {
		if result.Failure != nil {
			return snapshot, ErrPartial
		}
	}
	return snapshot, nil
}

func verify(ctx context.Context, reader Reader, registration *Registration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	health, err := reader.Health(ctx)
	if err != nil {
		return err
	}
	if health.Status != "ok" {
		return incompatible()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	manifest, err := reader.VerifyIdentity(ctx)
	if err != nil {
		return err
	}
	registration.ObservedDaemonID = manifest.InstallationID
	if manifest.InstallationID == "" || manifest.ProductVersion != health.Version {
		return incompatible()
	}
	if registration.ExpectedDaemonID != "" && registration.ExpectedDaemonID != manifest.InstallationID {
		return &client.StatusError{Code: "identity_changed"}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	capability, err := reader.VerifyAccess(ctx)
	if err != nil {
		return err
	}
	if !capability.Valid() || !capability.Allows(domain.CapabilityObserve) {
		return &client.StatusError{Code: "forbidden"}
	}
	return nil
}

func incompatible() error { return &client.StatusError{Code: "incompatible"} }

// Classify discards raw transport/API text, which can contain credentials or sensitive payloads.
func Classify(err error) *Failure {
	if errors.Is(err, context.Canceled) {
		return &Failure{Code: "canceled", Message: "The observation was canceled."}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &Failure{Code: "timed_out", Message: "This scheduler did not respond in time. Check the connection and try again."}
	}
	var setup *setupError
	if errors.As(err, &setup) {
		result := setup.failure
		return &result
	}
	var status *client.StatusError
	if errors.As(err, &status) {
		switch status.Code {
		case "credential_revoked":
			return &Failure{Code: status.Code, Message: "The saved credential was revoked. Repair this connection with a new pairing phrase."}
		case "unauthorized", "authentication_failed":
			return &Failure{Code: "unauthorized", Message: "The saved credential was rejected. Repair this connection."}
		case "forbidden":
			return &Failure{Code: status.Code, Message: "This credential cannot read operational data. Restore Observe authority or select another connection."}
		case "identity_changed", "conflict":
			return &Failure{Code: "identity_changed", Message: "The scheduler identity differs from the saved target. Verify the target before repairing this connection."}
		case "incompatible", "not_found":
			return &Failure{Code: "incompatible", Message: "This scheduler returned an unsupported operational response. Update the scheduler service."}
		}
	}
	var connection *client.ConnectionError
	if errors.As(err, &connection) {
		switch connection.Kind {
		case client.ConnectionTrustFailure:
			return &Failure{Code: "trust_changed", Message: "The remote certificate is no longer trusted. Verify its change before repairing this connection."}
		case client.ConnectionAccessDenied:
			return &Failure{Code: "access_denied", Message: "Access to the local scheduler was denied. Check service permissions for this account."}
		case client.ConnectionTimeout:
			return &Failure{Code: "timed_out", Message: "This scheduler did not respond in time. Check the connection and try again."}
		}
	}
	return &Failure{Code: "unavailable", Message: "This scheduler is unavailable. Check its service or connection and try again."}
}
