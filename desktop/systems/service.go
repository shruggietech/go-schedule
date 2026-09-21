package systems

import (
	"context"
	"errors"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	maxConcurrentTargets = 4
	targetTimeout        = 5 * time.Second
)

type profileStore interface {
	Load() (clientprofile.Collection, error)
}

type secretStore interface {
	Load(string, string) (string, error)
}

type summaryClient interface {
	SystemSummary(context.Context) (domain.SystemSummary, error)
}

type target struct {
	registration Registration
	backend      connection.Backend
	client       summaryClient
}

// Service refreshes independent daemon observations without changing the selected connection.
type Service struct {
	profiles profileStore
	secrets  secretStore
	local    *client.Client

	mu          sync.Mutex
	generation  uint64
	cancel      context.CancelFunc
	cache       map[string]cachedSummary
	now         func() time.Time
	loadTargets func() ([]target, []Observation)
	timeout     time.Duration
	concurrency int
}

// New creates an All Systems service over existing profile, credential, and local client stores.
func New(profiles profileStore, secrets secretStore, local *client.Client) *Service {
	return &Service{profiles: profiles, secrets: secrets, local: local, cache: make(map[string]cachedSummary), now: time.Now, timeout: targetTimeout, concurrency: maxConcurrentTargets}
}

// Refresh returns one observation for every registration present at completion and optionally publishes generation-scoped partial snapshots.
func (s *Service) Refresh(ctx context.Context, publish ...func(Snapshot)) Snapshot {
	started := s.now().UTC()
	refreshCtx, generation := s.begin(ctx)
	var targets []target
	var setupFailures []Observation
	if s.loadTargets != nil {
		targets, setupFailures = s.loadTargets()
	} else {
		targets, setupFailures = s.targets()
	}
	results := make(chan Observation, len(targets))
	jobs := make(chan target)

	var workers sync.WaitGroup
	workerCount := min(s.concurrency, len(targets))
	for range workerCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for value := range jobs {
				results <- s.observe(refreshCtx, generation, value)
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, value := range targets {
			select {
			case jobs <- value:
			case <-refreshCtx.Done():
				return
			}
		}
	}()
	go func() {
		workers.Wait()
		close(results)
	}()

	observations := append([]Observation(nil), setupFailures...)
	if len(observations) > 0 && len(publish) > 0 {
		publish[0](s.snapshot(generation, started, observations, false, false))
	}
	for value := range results {
		observations = append(observations, value)
		if len(publish) > 0 {
			publish[0](s.snapshot(generation, started, observations, false, false))
		}
	}
	final := s.snapshot(generation, started, observations, true, true)
	if len(publish) > 0 {
		publish[0](final)
	}
	return final
}

func (s *Service) snapshot(generation uint64, started time.Time, observations []Observation, complete, filterCurrent bool) Snapshot {
	filtered := append([]Observation(nil), observations...)
	if filterCurrent {
		current := s.currentKeys(observations)
		kept := filtered[:0]
		for _, value := range filtered {
			if current[value.Registration.Key] {
				kept = append(kept, value)
			}
		}
		filtered = kept
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Registration.Kind != filtered[j].Registration.Kind {
			return filtered[i].Registration.Kind == "local"
		}
		if filtered[i].Registration.Label != filtered[j].Registration.Label {
			return filtered[i].Registration.Label < filtered[j].Registration.Label
		}
		return filtered[i].Registration.Key < filtered[j].Registration.Key
	})
	completedAt := ""
	if complete {
		completedAt = s.now().UTC().Format(time.RFC3339)
	}
	return Snapshot{Generation: generation, StartedAt: started.Format(time.RFC3339), CompletedAt: completedAt, Complete: complete, Observations: filtered}
}

func (s *Service) begin(parent context.Context) (context.Context, uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
	}
	s.generation++
	ctx, cancel := context.WithCancel(parent)
	s.cancel = cancel
	return ctx, s.generation
}

func (s *Service) targets() ([]target, []Observation) {
	localRegistration := Registration{Key: "local", Kind: "local", Label: "This computer", DaemonID: "local", ShortDaemonID: "local", Platform: runtime.GOOS}
	result := []target{{registration: localRegistration, backend: connection.NewLocalBackend(s.local), client: s.local}}
	collection, err := s.profiles.Load()
	if err != nil {
		return result, nil
	}
	failures := make([]Observation, 0)
	for _, profile := range collection.Profiles {
		registration := registrationOf(profile)
		token, loadErr := s.secrets.Load(profile.DaemonID, profile.CredentialID)
		if loadErr != nil {
			failures = append(failures, s.failed(registration, &connection.Failure{State: connection.StateUnavailable, Message: "The saved credential is unavailable.", Action: "Repair this connection before opening it.", Cause: loadErr}))
			continue
		}
		remote, remoteErr := client.NewRemote(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
		if remoteErr != nil {
			failures = append(failures, s.failed(registration, &connection.Failure{State: connection.StateIncompatible, Message: "The saved connection profile is invalid.", Action: "Repair or remove this connection.", Cause: remoteErr}))
			continue
		}
		result = append(result, target{registration: registration, backend: connection.NewRemoteBackend(remote, profile.Capability), client: remote})
	}
	return result, failures
}

func (s *Service) observe(parent context.Context, generation uint64, value target) Observation {
	ctx, cancel := context.WithTimeout(parent, s.timeout)
	defer cancel()
	health, err := value.backend.Health(ctx)
	if err != nil {
		return s.failed(value.registration, err)
	}
	value.registration.DaemonID = health.ID
	value.registration.ShortDaemonID = connection.ShortID(health.ID)
	value.registration.Platform = health.Platform
	value.registration.Architecture = health.Architecture
	value.registration.Version = health.Version
	summary, err := value.client.SystemSummary(ctx)
	if err != nil {
		return s.failed(value.registration, classifySummaryFailure(value.registration.Kind, err))
	}
	if summary.Schema != domain.SystemSummarySchema || summary.ObservedAt.IsZero() {
		return s.failed(value.registration, &connection.Failure{State: connection.StateIncompatible, Message: "This scheduler returned an unsupported operational summary.", Action: "Update the scheduler service."})
	}
	s.mu.Lock()
	if s.generation == generation {
		s.cache[value.registration.Key] = cachedSummary{summary: summary, at: summary.ObservedAt}
	}
	s.mu.Unlock()
	copy := summary
	return Observation{Registration: value.registration, State: connection.StateConnected, ObservedAt: summary.ObservedAt.UTC().Format(time.RFC3339), Summary: &copy}
}

func (s *Service) failed(registration Registration, err error) Observation {
	failure := &connection.Failure{State: connection.StateUnavailable, Message: "This scheduler could not be observed.", Action: "Check the connection and try again.", Cause: err}
	var typed *connection.Failure
	if errors.As(err, &typed) {
		failure = typed
	}
	result := Observation{Registration: registration, State: failure.State, Failure: &Failure{State: failure.State, Message: failure.Message, Action: failure.Action}}
	s.mu.Lock()
	cached, ok := s.cache[registration.Key]
	s.mu.Unlock()
	if ok {
		copy := cached.summary
		result.Stale = true
		result.Summary = &copy
		result.ObservedAt = cached.at.UTC().Format(time.RFC3339)
	}
	return result
}

func (s *Service) currentKeys(fallback []Observation) map[string]bool {
	result := map[string]bool{"local": true}
	collection, err := s.profiles.Load()
	if err != nil {
		for _, observation := range fallback {
			result[observation.Registration.Key] = true
		}
		return result
	}
	for _, profile := range collection.Profiles {
		result[profile.ID] = true
	}
	return result
}

func registrationOf(profile clientprofile.Profile) Registration {
	return Registration{Key: profile.ID, ProfileID: profile.ID, Kind: "remote", Label: profile.Label, Endpoint: profile.Endpoint, DaemonID: profile.DaemonID, ShortDaemonID: connection.ShortID(profile.DaemonID), Platform: profile.Platform, Architecture: profile.Architecture, Version: profile.ProductVersion}
}

func classifySummaryFailure(kind string, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return &connection.Failure{State: connection.StateTimedOut, Message: "This scheduler did not return its operational summary in time.", Action: "Check the connection and try again.", Cause: err}
	}
	var connectionErr *client.ConnectionError
	if errors.As(err, &connectionErr) {
		if connectionErr.Kind == client.ConnectionTrustFailure {
			return &connection.Failure{State: connection.StateTrustChanged, Message: "The remote certificate is no longer trusted.", Action: "Inspect the certificate change, then repair this connection explicitly.", Cause: err}
		}
		return &connection.Failure{State: connection.StateUnavailable, Message: "This scheduler became unavailable during refresh.", Action: "Check the connection and try again.", Cause: err}
	}
	if kind == "local" {
		return &connection.Failure{State: connection.StateUnavailable, Message: "The local scheduler summary is unavailable.", Action: "Check the service, then try again.", Cause: err}
	}
	var status *client.StatusError
	if errors.As(err, &status) {
		switch status.Code {
		case "not_found":
			return &connection.Failure{State: connection.StateIncompatible, Message: "This scheduler does not support the operational overview.", Action: "Update the remote scheduler.", Cause: err}
		case "credential_revoked":
			return &connection.Failure{State: connection.StateRevoked, Message: "The remote credential was revoked during refresh.", Action: "Repair the connection with a new pairing phrase.", Cause: err}
		case "unauthorized", "authentication_failed":
			return &connection.Failure{State: connection.StateUnauthorized, Message: "The remote credential was rejected during refresh.", Action: "Repair the connection or verify its credential.", Cause: err}
		case "forbidden":
			return &connection.Failure{State: connection.StateForbidden, Message: "The remote credential cannot read operational summaries.", Action: "Ask an administrator to restore Observe authority.", Cause: err}
		}
	}
	return &connection.Failure{State: connection.StateUnavailable, Message: "The remote scheduler summary is unavailable.", Action: "Check the connection and try again.", Cause: err}
}
