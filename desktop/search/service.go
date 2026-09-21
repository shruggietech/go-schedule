package search

import (
	"context"
	"errors"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/systems"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const maxConcurrentTargets = 8
const targetTimeout = 3 * time.Second

type profileStore interface {
	Load() (clientprofile.Collection, error)
}
type secretStore interface {
	Load(string, string) (string, error)
}

type daemonClient interface {
	Search(context.Context, string, []domain.SearchKind, int) (domain.DaemonSearch, error)
	GetTask(context.Context, string) (server.TaskResponse, error)
	SetTaskEnabled(context.Context, string, bool) error
	RunNow(context.Context, string) error
	ListAlertsPage(context.Context, bool, int, int, int) ([]domain.Alert, error)
	AckAlert(context.Context, string) error
}

type target struct {
	registration systems.Registration
	backend      connection.Backend
	client       daemonClient
}

// Service fans one query out without changing the desktop's selected connection.
type Service struct {
	profiles      profileStore
	secrets       secretStore
	local         *client.Client
	mu            sync.Mutex
	generation    uint64
	cancel        context.CancelFunc
	now           func() time.Time
	timeout       time.Duration
	concurrency   int
	loadTargets   func() ([]target, []Observation)
	resolveTarget func(string) (target, error)
}

// New creates a cross-daemon search service over current registrations.
func New(profiles profileStore, secrets secretStore, local *client.Client) *Service {
	return &Service{profiles: profiles, secrets: secrets, local: local, now: time.Now, timeout: targetTimeout, concurrency: maxConcurrentTargets}
}

// Search returns progressive source-labeled results for the current registrations.
func (s *Service) Search(ctx context.Context, request Request, publish ...func(Snapshot)) Snapshot {
	query := strings.TrimSpace(request.Query)
	started := s.now().UTC()
	searchCtx, generation := s.begin(ctx)
	if query == "" || utf8.RuneCountInString(query) > 200 || request.Limit < 0 || request.Limit > 50 || !validKinds(request.Kinds) {
		return Snapshot{Generation: generation, Query: query, StartedAt: started.Format(time.RFC3339), CompletedAt: s.now().UTC().Format(time.RFC3339), Complete: true, Observations: []Observation{}}
	}
	limit := request.Limit
	if limit == 0 {
		limit = 50
	}
	var targets []target
	var failures []Observation
	if s.loadTargets != nil {
		targets, failures = s.loadTargets()
	} else {
		targets, failures = s.targets()
	}
	observations := append([]Observation(nil), failures...)
	positions := map[string]int{}
	for i := range observations {
		positions[observations[i].Registration.Key] = i
	}
	for _, value := range targets {
		positions[value.registration.Key] = len(observations)
		observations = append(observations, Observation{Registration: value.registration, State: connection.StateConnecting, Matches: []Match{}})
	}
	snapshot := func(complete bool) Snapshot {
		values := append([]Observation(nil), observations...)
		if complete {
			current := s.currentKeys(observations)
			kept := values[:0]
			for _, value := range values {
				if current[value.Registration.Key] {
					kept = append(kept, value)
				}
			}
			values = kept
		}
		sort.SliceStable(values, func(i, j int) bool { return registrationLess(values[i].Registration, values[j].Registration) })
		completed := ""
		if complete {
			completed = s.now().UTC().Format(time.RFC3339)
		}
		return Snapshot{Generation: generation, Query: query, StartedAt: started.Format(time.RFC3339), CompletedAt: completed, Complete: complete, Observations: values}
	}
	if len(publish) > 0 {
		publish[0](snapshot(false))
	}
	jobs := make(chan target)
	results := make(chan Observation, len(targets))
	var workers sync.WaitGroup
	for range min(s.concurrency, len(targets)) {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for value := range jobs {
				results <- s.observe(searchCtx, value, query, request.Kinds, limit)
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, value := range targets {
			select {
			case jobs <- value:
			case <-searchCtx.Done():
				return
			}
		}
	}()
	go func() { workers.Wait(); close(results) }()
	for value := range results {
		if index, ok := positions[value.Registration.Key]; ok {
			observations[index] = value
		}
		if len(publish) > 0 {
			publish[0](snapshot(false))
		}
	}
	final := snapshot(true)
	if len(publish) > 0 {
		publish[0](final)
	}
	return final
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

func (s *Service) observe(parent context.Context, value target, query string, kinds []domain.SearchKind, limit int) Observation {
	ctx, cancel := context.WithTimeout(parent, s.timeout)
	defer cancel()
	health, err := value.backend.Health(ctx)
	if err != nil {
		return failed(value.registration, err)
	}
	value.registration.DaemonID = health.ID
	value.registration.ShortDaemonID = connection.ShortID(health.ID)
	value.registration.Platform, value.registration.Architecture, value.registration.Version = health.Platform, health.Architecture, health.Version
	result, err := value.client.Search(ctx, query, kinds, limit)
	if err != nil {
		return failed(value.registration, err)
	}
	if err := validateResult(result, query, limit); err != nil {
		return failed(value.registration, &connection.Failure{State: connection.StateIncompatible, Message: "This scheduler returned an unsupported search response.", Action: "Update the scheduler service.", Cause: err})
	}
	matches := make([]Match, 0, len(result.Results))
	for _, item := range result.Results {
		available, reason := availableActions(item.ActionHints, health.Permissions)
		matches = append(matches, Match{RegistrationKey: value.registration.Key, ExpectedDaemonID: health.ID, SourceLabel: value.registration.Label, SourceShortID: connection.ShortID(health.ID), Result: item, AvailableActions: available, DisabledReason: reason})
	}
	return Observation{Registration: value.registration, State: connection.StateConnected, ObservedAt: result.ObservedAt.UTC().Format(time.RFC3339), Truncated: result.Truncated, Matches: matches}
}

func (s *Service) currentKeys(fallback []Observation) map[string]bool {
	result := map[string]bool{"local": true}
	if s.profiles == nil {
		for _, observation := range fallback {
			result[observation.Registration.Key] = true
		}
		return result
	}
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

func validateResult(result domain.DaemonSearch, query string, limit int) error {
	if result.Schema != domain.DaemonSearchSchema || result.ObservedAt.IsZero() || result.Query != query || len(result.Results) > limit {
		return errors.New("invalid search envelope")
	}
	for _, item := range result.Results {
		if !item.Kind.Valid() || strings.TrimSpace(item.ObjectID) == "" || strings.TrimSpace(item.Name) == "" {
			return errors.New("invalid search match")
		}
		for _, action := range item.ActionHints {
			if !action.Valid() {
				return errors.New("invalid search action")
			}
		}
	}
	return nil
}

func (s *Service) targets() ([]target, []Observation) {
	localRegistration := systems.Registration{Key: "local", Kind: "local", Label: "This computer", DaemonID: "local", ShortDaemonID: "local", Platform: runtime.GOOS}
	result := []target{{registration: localRegistration, backend: connection.NewLocalBackend(s.local), client: s.local}}
	collection, err := s.profiles.Load()
	if err != nil {
		return result, nil
	}
	failures := []Observation{}
	for _, profile := range collection.Profiles {
		registration := registrationOf(profile)
		token, loadErr := s.secrets.Load(profile.DaemonID, profile.CredentialID)
		if loadErr != nil {
			failures = append(failures, failed(registration, loadErr))
			continue
		}
		remote, remoteErr := client.NewRemote(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
		if remoteErr != nil {
			failures = append(failures, failed(registration, remoteErr))
			continue
		}
		result = append(result, target{registration: registration, backend: connection.NewRemoteBackend(remote, profile.Capability), client: remote})
	}
	return result, failures
}

func (s *Service) targetFor(key string) (target, error) {
	if s.resolveTarget != nil {
		return s.resolveTarget(key)
	}
	if key == "local" {
		return target{registration: systems.Registration{Key: "local", Kind: "local", Label: "This computer"}, backend: connection.NewLocalBackend(s.local), client: s.local}, nil
	}
	collection, err := s.profiles.Load()
	if err != nil {
		return target{}, err
	}
	profile, err := collection.Find(key)
	if err != nil || profile.ID != key {
		return target{}, clientprofile.ErrNotFound
	}
	token, err := s.secrets.Load(profile.DaemonID, profile.CredentialID)
	if err != nil {
		return target{}, err
	}
	remote, err := client.NewRemote(profile.Endpoint, profile.CertificatePEM, token, profile.DaemonID)
	if err != nil {
		return target{}, err
	}
	return target{registration: registrationOf(profile), backend: connection.NewRemoteBackend(remote, profile.Capability), client: remote}, nil
}

func registrationOf(profile clientprofile.Profile) systems.Registration {
	return systems.Registration{Key: profile.ID, ProfileID: profile.ID, Kind: "remote", Label: profile.Label, Endpoint: profile.Endpoint, DaemonID: profile.DaemonID, ShortDaemonID: connection.ShortID(profile.DaemonID), Platform: profile.Platform, Architecture: profile.Architecture, Version: profile.ProductVersion}
}

func failed(registration systems.Registration, err error) Observation {
	state := connection.StateUnavailable
	message, action := "This scheduler could not be searched.", "Check the connection and try again."
	var typed *connection.Failure
	if errors.As(err, &typed) {
		state, message, action = typed.State, typed.Message, typed.Action
	}
	if errors.Is(err, context.DeadlineExceeded) {
		state, message = connection.StateTimedOut, "This scheduler did not return search results in time."
	}
	return Observation{Registration: registration, State: state, Matches: []Match{}, Failure: &systems.Failure{State: state, Message: message, Action: action}}
}

func availableActions(hints []domain.SearchAction, permissions []string) ([]domain.SearchAction, string) {
	canOperate := false
	for _, permission := range permissions {
		if permission == "operate" || permission == "manage" || permission == "enroll" {
			canOperate = true
		}
	}
	result := []domain.SearchAction{domain.SearchActionOpen}
	if canOperate {
		result = append(result, hints...)
	} else if len(hints) > 0 {
		return result, "This credential has Observe authority only."
	}
	return result, ""
}

func validKinds(kinds []domain.SearchKind) bool {
	for _, kind := range kinds {
		if !kind.Valid() {
			return false
		}
	}
	return true
}
func registrationLess(a, b systems.Registration) bool {
	if a.Kind != b.Kind {
		return a.Kind == "local"
	}
	if a.Label != b.Label {
		return a.Label < b.Label
	}
	return a.Key < b.Key
}
